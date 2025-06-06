package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	twitterscraper "github.com/imperatrona/twitter-scraper"
	"github.com/joho/godotenv"
)

func main() {
	// Load env
	_ = godotenv.Load()

	authToken := os.Getenv("TWITTER_AUTH_TOKEN")
	csrfToken := os.Getenv("TWITTER_CSRF_TOKEN")
	if authToken == "" || csrfToken == "" {
		panic("TWITTER_AUTH_TOKEN dan TWITTER_CSRF_TOKEN harus di-set di environment")
	}

	scraper := twitterscraper.New()
	scraper.SetAuthToken(twitterscraper.AuthToken{
		Token:     authToken,
		CSRFToken: csrfToken,
	})

	if !scraper.IsLoggedIn() {
		panic("AuthToken tidak valid")
	}

	// Serve static img folder
	fs := http.FileServer(http.Dir("./img"))
	http.Handle("/img/", http.StripPrefix("/img/", fs))

	http.HandleFunc("/download", func(w http.ResponseWriter, r *http.Request) {
		urlTweet := r.URL.Query().Get("url")
		if urlTweet == "" {
			http.Error(w, "Parameter 'url' dibutuhkan", http.StatusBadRequest)
			return
		}

		tweetID, err := extractTweetID(urlTweet)
		if err != nil {
			http.Error(w, "URL tweet tidak valid", http.StatusBadRequest)
			return
		}

		tweet, err := scraper.GetTweet(tweetID)
		if err != nil {
			http.Error(w, "Gagal mengambil tweet: "+err.Error(), http.StatusInternalServerError)
			return
		}

		if len(tweet.Photos) == 0 {
			http.Error(w, "Tweet tidak mengandung gambar", http.StatusNotFound)
			return
		}

		type ImgResult struct {
			Filename string `json:"filename"`
			URL      string `json:"url"`
		}

		var results []ImgResult

		// Pastikan folder img ada
		os.MkdirAll("img", os.ModePerm)

		for i, photo := range tweet.Photos {
			// Ganti ukuran ke large
			imgURL := strings.Replace(photo.URL, "&name=small", "&name=large", 1)

			filename := fmt.Sprintf("tweet_%s_img_%d.jpg", tweetID, i+1)
			filePath := filepath.Join("img", filename)

			err := downloadFile(imgURL, filePath)
			if err != nil {
				http.Error(w, "Gagal download gambar: "+err.Error(), http.StatusInternalServerError)
				return
			}

			results = append(results, ImgResult{
				Filename: filename,
				URL:      "/img/" + filename,
			})
		}

		resp := map[string]interface{}{
			"status":  "success",
			"tweetID": tweetID,
			"images":  results,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	fmt.Println("Server running at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func extractTweetID(url string) (string, error) {
	// Contoh URL: https://x.com/pottsness/status/1927288679086583819
	// regex cari angka setelah /status/
	re := regexp.MustCompile(`status/(\d+)`)
	matches := re.FindStringSubmatch(url)
	if len(matches) < 2 {
		return "", fmt.Errorf("tidak ditemukan tweet id")
	}
	return matches[1], nil
}

func downloadFile(url, filepath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}
