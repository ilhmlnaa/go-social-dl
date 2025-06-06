package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	twitterscraper "github.com/imperatrona/twitter-scraper"
)

func TwitterDownloadHandler(scraper *twitterscraper.Scraper) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		var urls []string
		for _, photo := range tweet.Photos {
			imgURL := strings.Replace(photo.URL, "&name=small", "&name=large", 1)
			urls = append(urls, imgURL)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "success",
			"urls":   urls,
		})
	}
}

func extractTweetID(url string) (string, error) {
	re := regexp.MustCompile(`status/(\d+)`)
	matches := re.FindStringSubmatch(url)
	if len(matches) < 2 {
		return "", fmt.Errorf("tidak ditemukan tweet id")
	}
	return matches[1], nil
}
