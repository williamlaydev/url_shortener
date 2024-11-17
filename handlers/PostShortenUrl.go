package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"urlShortener/service"

	"github.com/jackc/pgx/v5/pgxpool"
)

type postShortenUrlHandler struct {
	dbConn *pgxpool.Pool
}

type postShortenUrlRequest struct {
	Url string `json:"url"`
}

func NewPostShortenUrlHandler(conn *pgxpool.Pool) *postShortenUrlHandler {
	return &postShortenUrlHandler{
		dbConn: conn,
	}
}

func (h *postShortenUrlHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
	s := service.NewUrlShortenerService(h.dbConn)

	newUrl, err := s.ShortenUrl(req.Url)

	if err != nil {
		log.Printf("error: %s", err)
		http.Error(w, "Error processing request", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(newUrl))
}
