package main

import (
	"fmt"
	"log"
	"net/http"
	"sommelierr/internal/application"
	"sommelierr/internal/config"
	"sommelierr/internal/infrastructure/radarr"
)

func main() {
	// 1️⃣ Load configuration (reads .env + real env)
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	// 2️⃣ Build the Radarr repository implementation
	repo := radarr.New(cfg.RadarrHost, cfg.APIKey)

	// 3️⃣ Build the application use‑case
	getRandom := &application.GetRandomMovie{Repo: repo}

	// 4️⃣ Wire HTTP handlers
	apiHandler := &application.APIHandler{GetRandom: getRandom}
	mux := http.NewServeMux()
	application.RegisterRoutes(mux, apiHandler)

	// 5️⃣ Run the HTTP server
	addr := ":" + cfg.Port
	fmt.Printf("🚀 server listening on %s (Radarr host: %s)\n", addr, cfg.RadarrHost)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("listen error: %v", err)
	}
}