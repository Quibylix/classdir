package presentation

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"classdir/api/internal/shared/cfg"
)

type Store interface {
	Create(ctx context.Context, id, title string) error
	GetByID(ctx context.Context, id string) (*Presentation, error)
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

func (s *pgPresentationStore) GetByID(ctx context.Context, id string) (*Presentation, error) {
	ctx, cancel := context.WithTimeout(ctx, cfg.DbTimeout)
	defer cancel()

	var pres Presentation
	err := s.pool.QueryRow(ctx, `SELECT id, title FROM presentations WHERE id = $1`, id).Scan(&pres.ID, &pres.Title)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	rows, err := s.pool.Query(ctx, `SELECT id, slide_number, content, metadata FROM slides WHERE presentation_id = $1 ORDER BY slide_number`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var slide Slide
		var metadata json.RawMessage
		if err := rows.Scan(&slide.ID, &slide.SlideNumber, &slide.Content, &metadata); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(metadata, &slide.Metadata); err != nil {
			return nil, err
		}
		pres.Slides = append(pres.Slides, slide)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	if pres.Slides == nil {
		pres.Slides = []Slide{}
	}

	return &pres, nil
}
