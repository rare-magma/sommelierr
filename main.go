package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strings"
)

type Movie struct {
	ID        int     `json:"id"`
	Title     string  `json:"title"`
	Year      int     `json:"year"`
	Overview  string  `json:"overview"`
	Images    []Image `json:"images"`
}

type Image struct {
	CoverType string `json:"coverType"`
	URL       string `json:"url"`
	RemoteURL string `json:"remoteUrl"`
}

type DisplayMovie struct {
	Title    string
	Year     int
	Overview string
	Poster   string
}

func loadEnv() error {
	file, err := os.Open(".env")
	if err != nil {
		// If no .env file, proceed with environment variables if set
		return nil // Make it optional, but log or handle as needed
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		os.Setenv(key, value)
	}
	return scanner.Err()
}

func getMovies() ([]Movie, error) {
	host := os.Getenv("RADARR_HOST")
	if host == "" {
		return nil, fmt.Errorf("RADARR_HOST not set")
	}
	apikey := os.Getenv("RADARR_API_KEY")
	if apikey == "" {
		return nil, fmt.Errorf("RADARR_API_KEY not set")
	}

	url := fmt.Sprintf("%s/api/v3/movie?apikey=%s", host, apikey)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	var movies []Movie
	err = json.NewDecoder(resp.Body).Decode(&movies)
	if err != nil {
		return nil, err
	}

	return movies, nil
}

func getRandomMovie() (DisplayMovie, error) {
	movies, err := getMovies()
	if err != nil {
		return DisplayMovie{}, err
	}
	if len(movies) == 0 {
		return DisplayMovie{}, fmt.Errorf("no movies found in Radarr")
	}

	idx := rand.Intn(len(movies))
	m := movies[idx]

	// Find poster URL
	host := os.Getenv("RADARR_HOST")
	if host == "" {
		return DisplayMovie{}, fmt.Errorf("RADARR_HOST not set")
	}
	var poster string
	for _, img := range m.Images {
		if img.CoverType == "poster" {
			if img.URL != "" {
				poster = host + img.URL
			} else if img.RemoteURL != "" {
				poster = img.RemoteURL
			}
			break
		}
	}

	return DisplayMovie{
		Title:    m.Title,
		Year:     m.Year,
		Overview: m.Overview,
		Poster:   poster,
	}, nil
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	movie, err := getRandomMovie()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Random Movie from Radarr</title>
    <style>
        body { font-family: Arial, sans-serif; text-align: center; }
        img { max-width: 300px; height: auto; }
    </style>
</head>
<body>
    <h1>Random Movie</h1>
    <img height="450" width="300" src="%s" alt="%s poster">
    <h2>%s (%d)</h2>
    <p>%s</p>
    <button onclick="location.reload();">Get Another Random Movie</button>
</body>
</html>
`, movie.Poster, movie.Title, movie.Title, movie.Year, movie.Overview)
}

func main() {
	_ = loadEnv() // Load .env if present; ignore error if not found

	fmt.Println("Starting server on :8080")
	http.HandleFunc("/", homeHandler)
	http.ListenAndServe(":8080", nil)
}