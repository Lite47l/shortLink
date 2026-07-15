package postgres

import (
	"context"
	"errors"
	"shortLink/internal/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresStorage struct {
	pool *pgxpool.Pool
}

func New(ctx context.Context, dsn string) (*PostgresStorage, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	return &PostgresStorage{pool: pool}, nil
}

func (p *PostgresStorage) Save(ctx context.Context, originalU, shortU string) (string, error) {
	_, err := p.pool.Exec(ctx,
		"INSERT INTO urls (short_url, original_url) VALUES ($1, $2)",
		shortU, originalU)
	if err != nil {
		return shortU, err
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == "23505" {
			if pgErr.ConstraintName == "urls_original_url_key" {

				var existingShortU string
				err := p.pool.QueryRow(ctx,
					"SELECT short_url FROM urls WHERE original_url = $1", originalU).Scan(&existingShortU)
				if err != nil {
					return "", err
				}
				return existingShortU, nil
			}

			if pgErr.ConstraintName == "urls_pkey" {
				return "", model.ErrShortUrlCollision
			}
		}
	}

	return "", err
}

func (p *PostgresStorage) Get(ctx context.Context, shortURL string) (string, error) {
	var originalU string

	err := p.pool.QueryRow(ctx,
		"SELECT original_url FROM urls WHERE short_url = $1", shortURL).Scan(&originalU)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", model.ErrNotFound
	}

	if err != nil {
		return "", err
	}

	return originalU, nil
}

func (p *PostgresStorage) Close() error {
	p.pool.Close()
	return nil
}
