package application

import (
	_ "embed"
	"encoding/json"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"sommelierr/internal/domain"
	"strings"
)

//go:embed ui/index.html
var indexHTML string

type APIHandler struct {
	GetRandomMovie  *GetRandomMovie
	GetRandomSeries *GetRandomSeries
}

func (h *APIHandler) RandomMovieHandler(w http.ResponseWriter, r *http.Request) {
	movie, err := h.GetRandomMovie.Execute()
	if err != nil {
		if err == ErrNoMovies {
			http.Error(w, "no movies available", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	resp := domain.Media{
		Title:     movie.Title,
		Year:      movie.Year,
		Overview:  movie.Overview,
		PosterURL: movie.PosterURL,
		SourceURL: movie.SourceURL,
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

func (h *APIHandler) RandomSeriesHandler(w http.ResponseWriter, r *http.Request) {
	series, err := h.GetRandomSeries.Execute()
	if err != nil {
		if err == ErrNoSeries {
			http.Error(w, "no series available", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	resp := domain.Media{
		Title:     series.Title,
		Year:      series.Year,
		Overview:  series.Overview,
		PosterURL: series.PosterURL,
		SourceURL: series.SourceURL,
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

func UIHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip serving index.html for the API endpoints
		if strings.HasPrefix(r.URL.Path, "/movie") {
			http.NotFound(w, r)
			return
		}
		if strings.HasPrefix(r.URL.Path, "/series") {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(indexHTML))
	})
}

func ImageHandler(radarrHost url.URL, sonarrHost url.URL) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		redirectHost := r.URL.Query().Get("redirectTo")
		redirectPath := r.URL.Query().Get("path")
		var target url.URL

		switch redirectHost {
		case "radarr":
			target = radarrHost
		case "sonarr":
			target = sonarrHost
		default:
			http.Error(w, "Invalid redirectTo", http.StatusInternalServerError)
			return
		}
		pattern := `^\/MediaCover\/[0-9]+\/poster\.jpg$`
		matched, err := regexp.MatchString(pattern, redirectPath)
		if err != nil {
			http.Error(w, "Invalid path", http.StatusInternalServerError)
			return
		}
		if !matched {
			http.Error(w, "Invalid path", http.StatusInternalServerError)
			return
		}

		proxy := httputil.NewSingleHostReverseProxy(&target)
		director := proxy.Director
		proxy.Director = func(r *http.Request) {
			director(r)
			r.Host = target.Host
			r.URL.Path = redirectPath
		}
		proxy.ServeHTTP(w, r)
	}
}

func RegisterRoutes(mux *http.ServeMux, api *APIHandler, radarrHost url.URL, sonarrHost url.URL) {
	mux.HandleFunc("/movie", api.RandomMovieHandler)
	mux.HandleFunc("/series", api.RandomSeriesHandler)
	mux.HandleFunc("/image", ImageHandler(radarrHost, sonarrHost))
	mux.Handle("/", UIHandler())
}
