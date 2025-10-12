package application

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"html/template"
	"net/http"
	"net/url"
	"sommelierr/internal/domain"
	"strings"
)

//go:embed ui/index.html
var html string

//go:embed ui/styles.css
var css string

type Model struct {
	Style string
}

type APIHandler struct {
	GetRandomMovie  *GetRandomMovie
	GetRandomSeries *GetRandomSeries
}

func processTemplate() []byte {
	model := Model{
		Style: css,
	}

	funcs := template.FuncMap{
		"safeCss": func(s string) template.CSS {
			return template.CSS(s)
		},
	}
	template, err := template.New("index").Funcs(funcs).Parse(html)
	if err != nil {
		panic(err)
	}

	var output bytes.Buffer
	err = template.Execute(&output, model)
	if err != nil {
		panic(err)
	}
	return output.Bytes()
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
		PosterB64: movie.Poster,
		PosterURL: movie.RemotePosterURL,
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
		PosterB64: series.Poster,
		PosterURL: series.RemotePosterURL,
		SourceURL: series.SourceURL,
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

func UIHandler() http.Handler {
	html := processTemplate()

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
		_, _ = w.Write([]byte(html))
	})
}

func RegisterRoutes(mux *http.ServeMux, api *APIHandler, radarrHost url.URL, sonarrHost url.URL) {
	mux.HandleFunc("/movie", api.RandomMovieHandler)
	mux.HandleFunc("/series", api.RandomSeriesHandler)
	mux.Handle("/", UIHandler())
}
