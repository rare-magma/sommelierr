package application

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	RadarrHost   url.URL
	RadarrAPIKey string
	SonarrHost   url.URL
	SonarrAPIKey string
	ServerPort   int64
}

func LoadConfig() (*Config, error) {
	if err := loadDotEnv(); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("Failed to load .env: %w", err)
	}

	radarrHostString := os.Getenv("RADARR_HOST")
	if radarrHostString == "" {
		return nil, fmt.Errorf("RADARR_HOST is required")
	}
	radarrHost, err := url.Parse(radarrHostString)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse Radarr host: %w", err)
	}
	radarrKey := os.Getenv("RADARR_API_KEY")
	if radarrKey == "" {
		return nil, fmt.Errorf("RADARR_API_KEY is required")
	}

	sonarrHostString := os.Getenv("SONARR_HOST")
	if sonarrHostString == "" {
		return nil, fmt.Errorf("SONARR_HOST is required")
	}
	sonarrHost, err := url.Parse(sonarrHostString)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse Sonarr host: %w", err)
	}
	sonarrKey := os.Getenv("SONARR_API_KEY")
	if sonarrKey == "" {
		return nil, fmt.Errorf("SONARR_API_KEY is required")
	}

	serverPortString := os.Getenv("PORT")
	if serverPortString == "" {
		serverPortString = "8080"
	}
	serverPort, err := strconv.ParseInt(serverPortString, 0, 0)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse port: %w", err)
	}

	return &Config{
		RadarrHost:   *radarrHost,
		RadarrAPIKey: radarrKey,
		SonarrHost:   *sonarrHost,
		SonarrAPIKey: sonarrKey,
		ServerPort:   serverPort,
	}, nil
}

// loadDotEnv parses a file named ".env" in the current working directory.
// It follows the same rules as the popular godotenv package:
//   - ignore empty lines and lines that start with '#'
//   - split on the first '=', trim surrounding whitespace and optional quotes
func loadDotEnv() error {
	f, err := os.Open(".env")
	if err != nil {
		return err
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
