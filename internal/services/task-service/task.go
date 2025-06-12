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

func (ts *Service) CreateTask(ctx context.Context, authorID uuid.UUID, title, description string, deadline time.Time) (string, error) {
	const op = "task.CreateTask"

	log := ts.logger.With(
		slog.String("op", op),
		slog.String("title", title),
	)

	log.Info("creating task")

	task := &models.Task{
		ID:          uuid.New(),
		AuthorID:    authorID,
		Title:       title,
		Description: description,
		Deadline:    deadline,
	}

	err := ts.TaskProvider.CreateTask(ctx, task)
	if err != nil {
		//TODO ...
		log.Error("failed to create task", err)
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return task.ID.String(), nil
}

func (ts *Service) GetTasks(ctx context.Context, authorID uuid.UUID) ([]*models.Task, error) {
	const op = "task.GetTask"

	log := ts.logger.With(
		slog.String("op", op),
		slog.String("author_id", authorID.String()),
	)

	log.Info("getting task")

	tasks, err := ts.TaskProvider.GetTask(ctx, authorID)
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

func (ts *Service) DeleteTask(ctx context.Context, taskID, authorID uuid.UUID) error {
	const op = "task.DeleteTask"

	log := ts.logger.With(
		slog.String("op", op),
		slog.String("task_id", taskID.String()),
	)

	log.Info("deleting task")

	if err := ts.TaskProvider.DeleteTask(ctx, taskID, authorID); err != nil {
		//TODO ...
		log.Error("failed to delete task", err)
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
