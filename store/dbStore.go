package store

import (
	"context"
	"log"

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
	query := `INSERT INTO urls (link, shortened) VALUES (@url, @shortenedUrl)`
	args := pgx.NamedArgs{
		"url":          url,
		"shortenedUrl": shortenedUrl,
	}

	_, err := d.conn.Exec(context.Background(), query, args)

	if err != nil {
		return err
	}

	return nil
}

func (d *dbStore) RetrieveUrlFromShortened(shortenedUrl string) (string, error) {
	query := `SELECT link FROM urls WHERE shortened=@shortenedUrl`
	args := pgx.NamedArgs{
		"shortenedUrl": shortenedUrl,
	}

	rows, err := d.conn.Query(context.Background(), query, args)

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
