package student

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"classdir/api/internal/shared/cfg"
)

const (
	studentPKeyConstraint       = "students_pkey"
	studentNameUniqueConstraint = "students_presentation_id_name_key"
)

type Store interface {
	Create(ctx context.Context, presentationID, id, name string) error
	GetByID(ctx context.Context, presentationID, id string) (*Student, error)
	Update(ctx context.Context, presentationID, id, name string) error
	Delete(ctx context.Context, presentationID, id string) error
	List(ctx context.Context, presentationID string) ([]*Student, error)
}

var ErrDuplicateID = errors.New("duplicate student id")
var ErrDuplicateName = errors.New("duplicate student name")
var ErrNotFound = errors.New("not found")

type pgStudentStore struct {
	pool *pgxpool.Pool
}

func NewStore(pool *pgxpool.Pool) Store {
	return &pgStudentStore{pool: pool}
}

func (s *pgStudentStore) Create(ctx context.Context, presentationID, id, name string) error {
	ctx, cancel := context.WithTimeout(ctx, cfg.DbTimeout)
	defer cancel()
	_, err := s.pool.Exec(ctx, `INSERT INTO students (id, presentation_id, name) VALUES ($1, $2, $3)`, id, presentationID, name)
	if err != nil {
		if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok {
			switch {
			case pgErr.Code == cfg.PgErrUniqueViolation && pgErr.ConstraintName == studentPKeyConstraint:
				return ErrDuplicateID
			case pgErr.Code == cfg.PgErrUniqueViolation && pgErr.ConstraintName == studentNameUniqueConstraint:
				return ErrDuplicateName
			case pgErr.Code == cfg.PgErrForeignKeyViolation:
				return ErrNotFound
			}
		}
	}
	return err
}

func (s *pgStudentStore) GetByID(ctx context.Context, presentationID, id string) (*Student, error) {
	ctx, cancel := context.WithTimeout(ctx, cfg.DbTimeout)
	defer cancel()

	var st Student
	err := s.pool.QueryRow(ctx, `SELECT id, name FROM students WHERE presentation_id = $1 AND id = $2`, presentationID, id).Scan(&st.ID, &st.Name)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &st, nil
}

func (s *pgStudentStore) List(ctx context.Context, presentationID string) ([]*Student, error) {
	ctx, cancel := context.WithTimeout(ctx, cfg.DbTimeout)
	defer cancel()

	rows, err := s.pool.Query(ctx, `SELECT id, name FROM students WHERE presentation_id = $1 ORDER BY created_at ASC`, presentationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []*Student
	for rows.Next() {
		var st Student
		if err := rows.Scan(&st.ID, &st.Name); err != nil {
			return nil, err
		}
		students = append(students, &st)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	if students == nil {
		students = []*Student{}
	}

	return students, nil
}

func (s *pgStudentStore) Update(ctx context.Context, presentationID, id, name string) error {
	ctx, cancel := context.WithTimeout(ctx, cfg.DbTimeout)
	defer cancel()

	tag, err := s.pool.Exec(ctx, `UPDATE students SET name = $1, updated_at = NOW() WHERE presentation_id = $2 AND id = $3`, name, presentationID, id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == cfg.PgErrUniqueViolation && pgErr.ConstraintName == studentNameUniqueConstraint {
			return ErrDuplicateName
		}
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *pgStudentStore) Delete(ctx context.Context, presentationID, id string) error {
	ctx, cancel := context.WithTimeout(ctx, cfg.DbTimeout)
	defer cancel()

	tag, err := s.pool.Exec(ctx, `DELETE FROM students WHERE presentation_id = $1 AND id = $2`, presentationID, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
