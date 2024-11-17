package service

import (
	"urlShortener/store"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UrlRetrievalService struct {
	conn *pgxpool.Pool
}

func NewUrlRetrievalService(c *pgxpool.Pool) *UrlRetrievalService {
	return &UrlRetrievalService{
		conn: c,
	}
}

func (s *UrlRetrievalService) RetrieveUrlFromShortened(shortenedUrl string) (string, error) {
	db := store.NewDbStore(s.conn)

	url, err := db.RetrieveUrlFromShortened(shortenedUrl)

	if err != nil {
		return "", err
	}

	return url, nil
}
