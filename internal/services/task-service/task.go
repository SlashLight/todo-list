package task_service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"

	"github.com/SlashLight/todo-list/internal/domain/models"
)

type TaskProvider interface {
	CreateTask(ctx context.Context, task *models.Task) error
	GetTask(ctx context.Context, author uuid.UUID) ([]*models.Task, error)
	UpdateTask(ctx context.Context, newTask *models.Task) error
	DeleteTask(ctx context.Context, taskID, author uuid.UUID) error
}

type Service struct {
	TaskProvider TaskProvider
	logger       *slog.Logger
}

func New(taskProvider TaskProvider, log *slog.Logger) *Service {
	return &Service{
		TaskProvider: taskProvider,
		logger:       log,
	}
}

func (ts *Service) CreateTask(ctx context.Context, title, description string, deadline time.Time) (string, error) {
	const op = "task.CreateTask"

	log := ts.logger.With(
		slog.String("op", op),
		slog.String("title", title),
	)

	log.Info("creating task")

	session, err := models.SessionFromContext(ctx)
	if err != nil {
		log.Error("failed to get session from context", err)
		return "", fmt.Errorf("%s: %w", op, err)
	}

	task := &models.Task{
		ID:          uuid.New(),
		AuthorID:    session.UserID,
		Title:       title,
		Description: description,
		Deadline:    deadline,
	}

	err = ts.TaskProvider.CreateTask(ctx, task)
	if err != nil {
		//TODO ...
		log.Error("failed to create task", err)
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return task.ID.String(), nil
}

func (ts *Service) GetTasks(ctx context.Context) ([]*models.Task, error) {
	const op = "task.GetTask"

	session, err := models.SessionFromContext(ctx)
	if err != nil {
		ts.logger.Error("failed to get session from context", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log := ts.logger.With(
		slog.String("op", op),
		slog.String("author_id", session.UserID.String()),
	)

	log.Info("getting task")

	tasks, err := ts.TaskProvider.GetTask(ctx, session.UserID)
	if err != nil {
		//TODO ...
		log.Error("failed to get task", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return tasks, nil
}

func (ts *Service) UpdateTask(ctx context.Context, newTask *models.Task) error {
	const op = "task.UpdateTask"

	log := ts.logger.With(
		slog.String("op", op),
		slog.String("task_id", newTask.ID.String()),
	)

	log.Info("updating task")

	err := ts.TaskProvider.UpdateTask(ctx, newTask)
	if err != nil {
		//TODO ...
		log.Error("failed to update task", err)
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (ts *Service) DeleteTask(ctx context.Context, taskID uuid.UUID) error {
	const op = "task.DeleteTask"

	log := ts.logger.With(
		slog.String("op", op),
		slog.String("task_id", taskID.String()),
	)

	log.Info("deleting task")

	session, err := models.SessionFromContext(ctx)
	if err != nil {
		log.Error("failed to get session from context", err)
		return fmt.Errorf("%s: %w", op, err)
	}

	err = ts.TaskProvider.DeleteTask(ctx, taskID, session.UserID)
	if err != nil {
		//TODO ...
		log.Error("failed to delete task", err)
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
