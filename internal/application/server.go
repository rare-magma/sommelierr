package application

import (
	_ "embed"
	"encoding/json"
	"net/http"
	"strings"
)

//go:embed ui/index.html
var indexHTML string

// APIHandler bundles the use‑case(s) we expose.
type APIHandler struct {
	GetRandom *GetRandomMovie
}

// RandomMovieHandler returns JSON for the UI.
func (h *APIHandler) RandomMovieHandler(w http.ResponseWriter, r *http.Request) {
	movie, err := h.GetRandom.Execute()
	if err != nil {
		if err == ErrNoMovies {
			http.Error(w, "no movies available", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	resp := struct {
		Title     string `json:"title"`
		Year      int    `json:"year"`
		Overview  string `json:"overview,omitempty"`
		PosterURL string `json:"posterUrl,omitempty"`
	}{
		Title:     movie.Title,
		Year:      movie.Year,
		Overview:  movie.Overview,
		PosterURL: movie.PosterURL,
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

// UIHandler serves the embedded index.html content
func UIHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip serving index.html for the API endpoint
		if strings.HasPrefix(r.URL.Path, "/random-movie") {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(indexHTML))
	})
}

// RegisterRoutes attaches the handlers to a ServeMux.
func RegisterRoutes(mux *http.ServeMux, api *APIHandler) {
	mux.HandleFunc("/random-movie", api.RandomMovieHandler)

	// Serve embedded index.html for all other routes
	mux.Handle("/", UIHandler())
}