package sonarr

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
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

type Statistics struct {
	PercentOfEpisodes float64 `json:"percentOfEpisodes"`
}

func New(baseURL *url.URL, apiKey string, excludeLabel string) domain.SeriesRepository {
	return &client{
		baseURL:      baseURL,
		apiKey:       apiKey,
		excludeLabel: excludeLabel,
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

func (c *client) ListAvailable() ([]*domain.Series, error) {
	rel := &url.URL{Path: "/api/v3/series"}
	u := c.baseURL.ResolveReference(rel)

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Api-Key", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("sonarr returned %d", resp.StatusCode)
	}

	var raw []struct {
		Id         int            `json:"id"`
		Tags       []int          `json:"tags"`
		Title      string         `json:"title"`
		TitleSlug  string         `json:"titleSlug"`
		Year       int            `json:"year"`
		Overview   string         `json:"overview"`
		Images     []domain.Image `json:"images"`
		Added      string         `json:"added"`
		Statistics Statistics     `json:"statistics"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, err
	}

	filterId, err := c.findTagByLabel(c.excludeLabel)
	if err != nil {
		filterId = -1
	}

	var result []*domain.Series
series:
	for _, r := range raw {

		if r.Statistics.PercentOfEpisodes == 0 {
			continue
		}
		if filterId != -1 {
			for _, tag := range r.Tags {
				if tag == filterId {
					continue series
				}
			}
		}
		added, _ := time.Parse(time.RFC3339, r.Added)
		sourceUrl := fmt.Sprintf("%s/series/%s", c.baseURL.String(), r.TitleSlug)

		s := &domain.Series{
			Id:        r.Id,
			Title:     r.Title,
			Year:      r.Year,
			Overview:  r.Overview,
			Images:    r.Images,
			Added:     added,
			SourceURL: sourceUrl,
		}

		for _, img := range r.Images {
			if img.CoverType == "poster" {
				if img.URL != "" {
					s.PosterURL = img.URL
				} else if img.RemoteURL != "" {
					s.RemotePosterURL = img.RemoteURL
				}
				break
			}
		}
		result = append(result, s)
	}
	return result, nil
}

func (c *client) GetPoster(id int, imageUrl string) (string, error) {
	imageFilename := strings.Split(imageUrl, "/")
	imageName := strings.Split(imageFilename[len(imageFilename)-1], "?")[0]
	rel := &url.URL{Path: fmt.Sprintf("/api/v3/mediacover/%d/%s", id, imageName)}
	u := c.baseURL.ResolveReference(rel)

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "*/*")
	req.Header.Set("X-Api-Key", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("sonarr returned %d", resp.StatusCode)
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	imageType := resp.Header.Get("Content-Type")
	return fmt.Sprintf("data:%s;base64, %s", imageType, base64.StdEncoding.EncodeToString(data)), nil
}

func (c *client) findTagByLabel(label string) (int, error) {
	if label == "" {
		return 0, fmt.Errorf("no tag label supplied")
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
		return 0, fmt.Errorf("sonarr returned %d", resp.StatusCode)
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
	return 0, fmt.Errorf("tag %s not found", label)
}
