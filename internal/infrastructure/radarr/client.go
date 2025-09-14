package radarr

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sommelierr/internal/domain"
	"time"
)

type client struct {
	baseURL *url.URL
	apiKey  string
	httpCli *http.Client
}

func New(baseURL, apiKey string) domain.MovieRepository {
	parsed, _ := url.Parse(baseURL)
	return &client{
		baseURL: parsed,
		apiKey:  apiKey,
		httpCli: &http.Client{
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

	resp, err := c.httpCli.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("radarr returned %d", resp.StatusCode)
	}

	var raw []struct {
		ID            int            `json:"id"`
		Title         string         `json:"title"`
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

	var result []*domain.Movie
	for _, r := range raw {
		if !r.HasFile {
			continue
		}
		added, _ := time.Parse(time.RFC3339, r.Added)

		m := &domain.Movie{
			ID:            r.ID,
			Title:         r.Title,
			OriginalTitle: r.OriginalTitle,
			Year:          r.Year,
			Overview:      r.Overview,
			Images:        r.Images,
			Added:         added,
		}

		for _, img := range r.Images {
			if img.CoverType == "poster" {
				if img.URL != "" {
					m.PosterURL = c.baseURL.String() + img.URL
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
