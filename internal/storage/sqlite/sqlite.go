package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
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

func (s *Storage) CreateTask(ctx context.Context, task *models.Task) error {
	const op = "storage.sqlite.CreateTask"

	_, err := s.db.ExecContext(ctx, InsertNewTask, task.ID, task.AuthorID, task.Title, task.Description, task.Deadline)
	if err != nil {
		var sqliteErr sqlite3.Error

		if errors.As(err, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintNotNull {
			return fmt.Errorf("%s: %w", op, my_err.ErrEmptyTitle)
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) GetTask(ctx context.Context, author uuid.UUID) ([]*models.Task, error) {
	const op = "storage.sqlite.GetTask"

	var tasks []*models.Task
	rows, err := s.db.QueryContext(ctx, SelectTasksByAuthor, author)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("%s: execute statement: %w", op, err)
	}
	defer rows.Close()

	for rows.Next() {
		task := &models.Task{AuthorID: author}
		if err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Status, &task.Deadline); err != nil {
			return nil, fmt.Errorf("%s: scan row: %w", op, err)
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (s *Storage) UpdateTask(ctx context.Context, newTask *models.Task) error {
	const op = "storage.sqlite.UpdateTask"

	_, err := s.db.ExecContext(ctx, UpdateTaskByID, newTask.Title, newTask.Description, newTask.Status, newTask.Deadline, newTask.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return my_err.ErrTaskNotFound
		}

		return fmt.Errorf("%s: execute statement: %w", op, err)
	}

	return nil
}

func (s *Storage) DeleteTask(ctx context.Context, taskID, author uuid.UUID) error {
	const op = "storage.sqlite.DeleteTask"

	result, err := s.db.ExecContext(ctx, DeleteTaskByID, taskID, author)
	if err != nil {
		return fmt.Errorf("%s: execute statement: %w", op, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: get rows affected: %w", op, err)
	}

	if rowsAffected == 0 {
		return my_err.ErrTaskNotFound
	}

	return nil
}
