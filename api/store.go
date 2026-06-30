package main

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrDuplicateKey = errors.New("duplicate key")

type presentationStore interface {
	create(ctx context.Context, id, title string) error
}

type pgPresentationStore struct {
	pool *pgxpool.Pool
}

func newPresentationStore(pool *pgxpool.Pool) *pgPresentationStore {
	return &pgPresentationStore{pool: pool}
}

func (s *pgPresentationStore) create(ctx context.Context, id, title string) error {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	_, err := s.pool.Exec(ctx, `INSERT INTO presentations (id, title) VALUES ($1, $2)`, id, title)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrDuplicateKey
		}
	}
	return err
}
