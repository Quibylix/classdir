package presentation

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"classdir/api/internal/shared/cfg"
)

type Store interface {
	Create(ctx context.Context, id, title string) error
	GetByID(ctx context.Context, id string) (*Presentation, error)
	UpdateTitle(ctx context.Context, id, title string) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]*PresentationPreview, error)
	CreateSlide(ctx context.Context, presID, slideID, content string) error
	GetSlide(ctx context.Context, presID, slideID string) (*Slide, error)
	UpdateSlide(ctx context.Context, presID, slideID, content string) error
	DeleteSlide(ctx context.Context, presID, slideID string) error
}

var ErrDuplicateKey = errors.New("duplicate key")
var ErrNotFound = errors.New("not found")

type pgPresentationStore struct {
	pool *pgxpool.Pool
}

func NewStore(pool *pgxpool.Pool) Store {
	return &pgPresentationStore{pool: pool}
}

func (s *pgPresentationStore) Create(ctx context.Context, id, title string) error {
	ctx, cancel := context.WithTimeout(ctx, cfg.DbTimeout)
	defer cancel()
	_, err := s.pool.Exec(ctx, `INSERT INTO presentations (id, title, slide_order) VALUES ($1, $2, '{}')`, id, title)
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
	err := s.pool.QueryRow(ctx, `SELECT id, title, slide_order FROM presentations WHERE id = $1`, id).Scan(&pres.ID, &pres.Title, &pres.SlideOrder)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	if pres.SlideOrder == nil {
		pres.SlideOrder = []string{}
	}

	rows, err := s.pool.Query(ctx, `SELECT s.id, s.content FROM slides s JOIN presentations p ON p.id = s.presentation_id WHERE s.presentation_id = $1 ORDER BY array_position(p.slide_order, s.id)`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var slide Slide
		if err := rows.Scan(&slide.ID, &slide.Content); err != nil {
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

func (s *pgPresentationStore) List(ctx context.Context) ([]*PresentationPreview, error) {
	ctx, cancel := context.WithTimeout(ctx, cfg.DbTimeout)
	defer cancel()

	rows, err := s.pool.Query(ctx, `SELECT id, title FROM presentations ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var presentations []*PresentationPreview
	for rows.Next() {
		var pres PresentationPreview
		if err := rows.Scan(&pres.ID, &pres.Title); err != nil {
			return nil, err
		}
		presentations = append(presentations, &pres)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	if presentations == nil {
		presentations = []*PresentationPreview{}
	}

	return presentations, nil
}

func (s *pgPresentationStore) UpdateTitle(ctx context.Context, id, title string) error {
	ctx, cancel := context.WithTimeout(ctx, cfg.DbTimeout)
	defer cancel()
	tag, err := s.pool.Exec(ctx, `UPDATE presentations SET title = $1, updated_at = NOW() WHERE id = $2`, title, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *pgPresentationStore) Delete(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, cfg.DbTimeout)
	defer cancel()
	tag, err := s.pool.Exec(ctx, `DELETE FROM presentations WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
