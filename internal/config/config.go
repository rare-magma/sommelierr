package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Config holds the values we need at runtime.
type Config struct {
	RadarrHost string // e.g. http://localhost:7878
	APIKey     string // Radarr API key
	Port       string // HTTP server port (default 8080)
}

// Load reads a .env file (if it exists) and then the process environment.
// Values found in the real environment win over those from .env.
func Load() (*Config, error) {
	// 1️⃣ Load .env file (if present) – this is a tiny parser, no external lib.
	if err := loadDotEnv(); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("reading .env: %w", err)
	}

	// 2️⃣ Pull values from the environment (real env overrides .env)
	host := strings.TrimRight(os.Getenv("RADARR_HOST"), "/")
	if host == "" {
		host = "http://localhost:7878"
	}
	key := os.Getenv("RADARR_API_KEY")
	if key == "" {
		return nil, fmt.Errorf("RADARR_API_KEY is required")
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return &Config{
		RadarrHost: host,
		APIKey:     key,
		Port:       port,
	}, nil
}

// loadDotEnv parses a file named ".env" in the current working directory.
// It follows the same rules as the popular godotenv package:
//   * ignore empty lines and lines that start with '#'
//   * split on the first '=', trim surrounding whitespace and optional quotes
func loadDotEnv() error {
	f, err := os.Open(".env")
	if err != nil {
		return err // caller decides whether missing is an error
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// skip comments / empty lines
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		kv := strings.SplitN(line, "=", 2)
		if len(kv) != 2 {
			continue // malformed line – ignore
		}
		key := strings.TrimSpace(kv[0])
		val := strings.TrimSpace(kv[1])

		// Remove optional surrounding quotes (single or double)
		if len(val) >= 2 {
			if (val[0] == '"' && val[len(val)-1] == '"') ||
				(val[0] == '\'' && val[len(val)-1] == '\'') {
				val = val[1 : len(val)-1]
			}
		}
		// Set only if not already defined in the real environment.
		if _, exists := os.LookupEnv(key); !exists {
			_ = os.Setenv(key, val)
		}
	}
	return scanner.Err()
}