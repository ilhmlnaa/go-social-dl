package handlers

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
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

		if len(tweet.Photos) == 1 {
			imgURL := strings.Replace(tweet.Photos[0].URL, "&name=small", "&name=large", 1)
			err = streamImage(w, imgURL)
			if err != nil {
				http.Error(w, "Gagal mengirim gambar: "+err.Error(), http.StatusInternalServerError)
			}
			return
		}

		err = streamZipImages(w, tweetID, tweet.Photos)
		if err != nil {
			http.Error(w, "Gagal mengirim zip gambar: "+err.Error(), http.StatusInternalServerError)
			return
		}
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

func streamImage(w http.ResponseWriter, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	w.Header().Set("Content-Disposition", "inline")

	_, err = io.Copy(w, resp.Body)
	return err
}

func streamZipImages(w http.ResponseWriter, tweetID string, photos []twitterscraper.Photo) error {
	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	for i, photo := range photos {
		imgURL := strings.Replace(photo.URL, "&name=small", "&name=large", 1)

		resp, err := http.Get(imgURL)
		if err != nil {
			zipWriter.Close()
			return err
		}

		f, err := zipWriter.Create(fmt.Sprintf("tweet_%s_img_%d.jpg", tweetID, i+1))
		if err != nil {
			resp.Body.Close()
			zipWriter.Close()
			return err
		}

		_, err = io.Copy(f, resp.Body)
		resp.Body.Close()
		if err != nil {
			zipWriter.Close()
			return err
		}
	}

	zipWriter.Close()

	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"tweet_%s_images.zip\"", tweetID))
	w.Write(buf.Bytes())
	return nil
}
