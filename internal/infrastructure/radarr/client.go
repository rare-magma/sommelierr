package radarr

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sommelierr/internal/domain"
	"strings"
	"time"
)

type client struct {
	baseURL      *url.URL
	apiKey       string
	httpClient   *http.Client
	excludeLabel string
}

type Tag struct {
	Label string `json:"label"`
	Id    int    `json:"id"`
}

func New(baseURL *url.URL, apiKey string, excludeLabel string) domain.MovieRepository {
	return &client{
		baseURL:      baseURL,
		apiKey:       apiKey,
		excludeLabel: excludeLabel,
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

func (c *client) ListAvailable() ([]*domain.Movie, error) {
	rel := &url.URL{Path: "/api/v3/movie"}
	q := rel.Query()
	q.Set("apikey", c.apiKey)
	rel.RawQuery = q.Encode()

	u := c.baseURL.ResolveReference(rel)

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("Radarr returned %d", resp.StatusCode)
	}

	var raw []struct {
		Tags          []int          `json:"tags"`
		Title         string         `json:"title"`
		TitleSlug     string         `json:"titleSlug"`
		OriginalTitle string         `json:"originalTitle"`
		Year          int            `json:"year"`
		Overview      string         `json:"overview"`
		Images        []domain.Image `json:"images"`
		Added         string         `json:"added"`
		HasFile       bool           `json:"hasFile"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, err
	}

	filterId, err := c.findTagByLabel(c.excludeLabel)
	if err != nil {
		filterId = -1
	}

	var result []*domain.Movie

movies:
	for _, r := range raw {
		if !r.HasFile {
			continue
		}
		if filterId != -1 {
			for _, tag := range r.Tags {
				if tag == filterId {
					continue movies
				}
			}
		}
		added, _ := time.Parse(time.RFC3339, r.Added)
		sourceUrl := fmt.Sprintf("%s/movie/%s", c.baseURL.String(), r.TitleSlug)

		m := &domain.Movie{
			Title:         r.Title,
			OriginalTitle: r.OriginalTitle,
			Year:          r.Year,
			Overview:      r.Overview,
			Images:        r.Images,
			Added:         added,
			SourceURL:     sourceUrl,
		}

		for _, img := range r.Images {
			if img.CoverType == "poster" {
				if img.URL != "" {
					path := url.QueryEscape(strings.SplitAfter(img.URL, ".jpg")[0])
					m.PosterURL = "/image" + fmt.Sprintf("?redirectTo=radarr&path=%s", path)
				} else if img.RemoteURL != "" {
					m.PosterURL = img.RemoteURL
				}
				break
			}
		}
		result = append(result, m)
	}
	return result, nil
}

func (c *client) findTagByLabel(label string) (int, error) {
	if label == "" {
		return 0, fmt.Errorf("No tag label supplied")
	}

	rel := &url.URL{Path: "/api/v3/tag"}
	q := rel.Query()
	q.Set("apikey", c.apiKey)
	rel.RawQuery = q.Encode()

	u := c.baseURL.ResolveReference(rel)

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return 0, err
	}
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return 0, fmt.Errorf("Radarr returned %d", resp.StatusCode)
	}

	var raw []Tag
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return 0, err
	}

	for _, tag := range raw {
		if tag.Label == label {
			return tag.Id, nil
		}
	}
	return 0, fmt.Errorf("Tag %s not found", label)
}
