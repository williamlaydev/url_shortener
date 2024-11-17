package handlers

import (
	"log"
	"net/http"
	"urlShortener/service"

	"github.com/jackc/pgx/v5/pgxpool"
)

type getShortenUrlHandler struct {
	dbConn *pgxpool.Pool
}

func NewGetShortenUrlHandler(conn *pgxpool.Pool) *getShortenUrlHandler {
	return &getShortenUrlHandler{
		dbConn: conn,
	}
}

func (h *getShortenUrlHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
	service := service.NewUrlRetrievalService(h.dbConn)
	url, err := service.RetrieveUrlFromShortened(shortUrl)

	if err != nil {
		log.Printf("Error retrieving original URL for %s: %v", shortUrl, err)
		http.Error(w, "Invalid shortened URL", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, url, http.StatusMovedPermanently)
}
