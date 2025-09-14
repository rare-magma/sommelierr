package sonarr

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

type statistics struct {
	PercentOfEpisodes float64 `json:"percentOfEpisodes"`
}

func New(baseURL *url.URL, apiKey string) domain.SeriesRepository {
	return &client{
		baseURL: baseURL,
		apiKey:  apiKey,
		httpCli: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

func (c *client) ListAvailable() ([]*domain.Series, error) {
	rel := &url.URL{Path: "/api/v3/series"}
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
		return nil, fmt.Errorf("Sonarr returned %d", resp.StatusCode)
	}

	var raw []struct {
		Title      string         `json:"title"`
		TitleSlug  string         `json:"titleSlug"`
		Year       int            `json:"year"`
		Overview   string         `json:"overview"`
		Images     []domain.Image `json:"images"`
		Added      string         `json:"added"`
		Statistics statistics     `json:"statistics"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, err
	}

	var result []*domain.Series
	for _, r := range raw {

		if r.Statistics.PercentOfEpisodes == 0 {
			continue
		}
		added, _ := time.Parse(time.RFC3339, r.Added)
		sourceUrl := fmt.Sprintf("%s/series/%s", c.baseURL.String(), r.TitleSlug)

		m := &domain.Series{
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
