package presentation

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	"classdir/api/internal/shared/cfg"
)

func (s *pgPresentationStore) CreateSlide(ctx context.Context, presID, slideID, content string) error {
	ctx, cancel := context.WithTimeout(ctx, cfg.DbTimeout)
	defer cancel()

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `INSERT INTO slides (id, presentation_id, content) VALUES ($1, $2, $3)`, slideID, presID, content)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `UPDATE presentations SET slide_order = slide_order || $2, updated_at = NOW() WHERE id = $1`, presID, []string{slideID})
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (s *pgPresentationStore) GetSlide(ctx context.Context, presID, slideID string) (*Slide, error) {
	ctx, cancel := context.WithTimeout(ctx, cfg.DbTimeout)
	defer cancel()

	var slide Slide
	err := s.pool.QueryRow(ctx, `SELECT id, content FROM slides WHERE id = $1 AND presentation_id = $2`, slideID, presID).Scan(&slide.ID, &slide.Content)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &slide, nil
}

func (s *pgPresentationStore) UpdateSlide(ctx context.Context, presID, slideID, content string) error {
	ctx, cancel := context.WithTimeout(ctx, cfg.DbTimeout)
	defer cancel()

	tag, err := s.pool.Exec(ctx, `UPDATE slides SET content = $1, updated_at = NOW() WHERE id = $2 AND presentation_id = $3`, content, slideID, presID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

func (s *pgPresentationStore) DeleteSlide(ctx context.Context, presID, slideID string) error {
	ctx, cancel := context.WithTimeout(ctx, cfg.DbTimeout)
	defer cancel()

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	tag, err := tx.Exec(ctx, `DELETE FROM slides WHERE id = $1 AND presentation_id = $2`, slideID, presID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}

	_, err = tx.Exec(ctx, `UPDATE presentations SET slide_order = array_remove(slide_order, $2), updated_at = NOW() WHERE id = $1`, presID, slideID)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
