package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/mattn/go-sqlite3"

	"github.com/SlashLight/todo-list/internal/domain/models"
	"github.com/SlashLight/todo-list/pkg/my_err"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.Sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	const op = "storage.sqlite.GetByEmail"

	user := &models.User{}

	err := s.db.QueryRowContext(ctx, SelectUserByEmail, email).Scan(&user.ID, &user.Email, &user.PasswordHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, my_err.ErrUserNotFound
		}

		return nil, fmt.Errorf("%s: execute statement: %w", op, err)
	}

	return user, nil
}

func (s Storage) Register(ctx context.Context, user *models.User) error {
	const op = "storage.sqlite.Register"

	_, err := s.db.ExecContext(ctx, InsertNewUser, user.ID, user.Email, user.PasswordHash)
	if err != nil {
		var sqliteErr sqlite3.Error

		if errors.As(err, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return fmt.Errorf("%s: %w", op, my_err.ErrUserExists)
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
