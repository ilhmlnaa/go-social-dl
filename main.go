package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"twitter-down/handlers"
	"twitter-down/middleware"

	twitterscraper "github.com/imperatrona/twitter-scraper"
)

func main() {
	_ = godotenv.Load()

	authToken := os.Getenv("TWITTER_AUTH_TOKEN")
	csrfToken := os.Getenv("TWITTER_CSRF_TOKEN")
	
	if authToken == "" || csrfToken == "" {
		panic("TWITTER_AUTH_TOKEN dan TWITTER_CSRF_TOKEN harus di-set di environment")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	scraper := twitterscraper.New()
	scraper.SetAuthToken(twitterscraper.AuthToken{
		Token:     authToken,
		CSRFToken: csrfToken,
	})

	if !scraper.IsLoggedIn() {
		panic("AuthToken tidak valid")
	}

	mux := http.NewServeMux()

	mux.Handle("/dl", middleware.CORS(handlers.TwitterDownloadHandler(scraper)))

	fmt.Printf("Server running at http://localhost:%s\n", port)
	err := http.ListenAndServe(":"+port, mux)
	if err != nil {
		panic(err)
	}
}
