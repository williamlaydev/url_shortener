package service

import (
	"crypto/rand"
	"fmt"
	"log"
	"urlShortener/store"

	"github.com/jackc/pgx/v5/pgxpool"
)

type urlShortenerService struct {
	conn *pgxpool.Pool
}

func NewUrlShortenerService(c *pgxpool.Pool) *urlShortenerService {
	return &urlShortenerService{
		conn: c,
	}
}

func (s *urlShortenerService) ShortenUrl(url string) (string, error) {
	log.Printf("url to shorten: %s", url)

	// Shorten url
	db := store.NewDbStore(s.conn)

	shortUrl, err := generateRandomString()

	if err != nil {
		return "", err
	}

	if err := db.CreateNewShortenedUrl(url, shortUrl); err != nil {
		return "", err
	}

	return shortUrl, err
}

func generateRandomString() (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, 6)
	charsetLength := len(charset)

	randomBytes := make([]byte, 6)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	for i, b := range randomBytes {
		result[i] = charset[b%byte(charsetLength)]
	}

	return string(result), nil
}
