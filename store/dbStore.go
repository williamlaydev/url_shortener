package store

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type dbStore struct {
	conn *pgxpool.Pool
}

func NewDbStore(conn *pgxpool.Pool) *dbStore {
	return &dbStore{conn: conn}
}

func (d *dbStore) CreateNewShortenedUrl(url string, shortenedUrl string) error {
	query := `INSERT INTO urls (link, shortened) VALUES ($1, $2)`

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	_, err := d.conn.Exec(ctx, query, url, shortenedUrl)

	if err != nil {
		return err
	}

	return nil
}

func (d *dbStore) RetrieveUrlFromShortened(shortenedUrl string) (string, error) {
	query := `SELECT link FROM urls WHERE shortened=@$1`

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	rows, err := d.conn.Query(ctx, query, shortenedUrl)

	if err != nil {
		log.Printf("err1: %v\n", err)
		return "", err
	}

	data, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[Url])

	if err != nil {
		log.Printf("err2: %v", err)
		return "", err
	}

	log.Printf("link found: %v", data.Link)

	return data.Link, nil
}
