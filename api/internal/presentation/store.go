package presentation

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"classdir/api/internal/shared/cfg"
)

type Store interface {
	Create(ctx context.Context, id, title string) error
}

var ErrDuplicateKey = errors.New("duplicate key")

type pgPresentationStore struct {
	pool *pgxpool.Pool
}

func NewStore(pool *pgxpool.Pool) Store {
	return &pgPresentationStore{pool: pool}
}

func (s *pgPresentationStore) Create(ctx context.Context, id, title string) error {
	ctx, cancel := context.WithTimeout(ctx, cfg.DbTimeout)
	defer cancel()
	_, err := s.pool.Exec(ctx, `INSERT INTO presentations (id, title) VALUES ($1, $2)`, id, title)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == cfg.PgErrUniqueViolation {
			return ErrDuplicateKey
		}
	}
	return err
}
