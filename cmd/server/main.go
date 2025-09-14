package main

import (
	"fmt"
	"log"
	"net/http"
	"sommelierr/internal/application"
	"sommelierr/internal/config"
	"sommelierr/internal/infrastructure/radarr"
	"sommelierr/internal/infrastructure/sonarr"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	moviesRepo := radarr.New(cfg.RadarrHost, cfg.RadarrAPIKey)
	seriesRepo := sonarr.New(cfg.SonarrHost, cfg.SonarrAPIKey)

	getRandomMovie := &application.GetRandomMovie{Repo: moviesRepo}
	getRandomSeries := &application.GetRandomSeries{Repo: seriesRepo}

	apiHandler := &application.APIHandler{GetRandomMovie: getRandomMovie, GetRandomSeries: getRandomSeries}
	mux := http.NewServeMux()
	application.RegisterRoutes(mux, apiHandler)

	addr := ":" + cfg.ServerPort
	fmt.Printf("sommelierr listening on %s\n", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("listen error: %v", err)
	}
}
