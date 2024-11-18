package handlers

import (
	"context"
	"log"
	"net/http"
	"urlShortener/components"
	"urlShortener/service"

	"github.com/jackc/pgx/v5/pgxpool"
)

type postShortenUrlHandler struct {
	dbConn *pgxpool.Pool
}

func NewPostShortenUrlHandler(conn *pgxpool.Pool) *postShortenUrlHandler {
	return &postShortenUrlHandler{
		dbConn: conn,
	}
}

func (h *postShortenUrlHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Log which API has been made
	log.Print("POST shortenUrl has been made")

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	originalUrl := r.FormValue("url")

	if originalUrl == "" {
		http.Error(w, "No URL provided", http.StatusBadRequest)
		return
	}

	// Pass the service
	s := service.NewUrlShortenerService(h.dbConn)

	newUrl, err := s.ShortenUrl(originalUrl)

	if err != nil {
		log.Printf("error: %s", err)
		http.Error(w, "Error processing request", http.StatusInternalServerError)
		return
	}

	if err := components.ShortenedUrlDisplay(newUrl).Render(context.Background(), w); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("broke"))
	}

}
