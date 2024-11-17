package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"urlShortener/service"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UrlHandler struct {
	DbConn *pgxpool.Pool
}

type postShortenUrlRequest struct {
	Url string `json:"url"`
}

func (h *UrlHandler) PostShortenUrl(w http.ResponseWriter, r *http.Request) {
	// Log which API has been made
	log.Print("POST shortenUrl has been made")
	var req postShortenUrlRequest

	// Decode the JSON body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate the body
	if len(req.Url) == 0 {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Pass the service
	s := service.NewUrlShortenerService(h.DbConn)

	newUrl, err := s.ShortenUrl(req.Url)

	if err != nil {
		log.Printf("error: %s", err)
		http.Error(w, "Error processing request", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(newUrl))
}

func (h *UrlHandler) GetShortenedUrl(w http.ResponseWriter, r *http.Request) {
	// Log which API has been made
	log.Print("GET shortenedUrl has been made")

	// Get the shortened URL from param
	shortUrl := r.URL.Path[len("/"):]

	// Validate the url
	if len(shortUrl) != 6 {
		http.Error(w, "Invalid short url length", http.StatusBadRequest)
		return
	}

	// Pass the service
	service := service.NewUrlRetrievalService(h.DbConn)
	url, err := service.RetrieveUrlFromShortened(shortUrl)

	if err != nil {
		log.Printf("Error retrieving original URL for %s: %v", shortUrl, err)
		http.Error(w, "Invalid shortened URL", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, url, http.StatusMovedPermanently)
	return
}
