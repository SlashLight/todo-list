package grpc

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	grpclog "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	grpcretry "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"

	taskv1 "github.com/SlashLight/todo-list/api/gen/go/todo"
	"github.com/SlashLight/todo-list/internal/domain/models"
)

type Client struct {
	api taskv1.TodoClient
	log *slog.Logger
}

func New(addr string, log *slog.Logger, retries int, timeout time.Duration) (*Client, error) {
	const op = "task.grpc.New"

	retryOpts := []grpcretry.CallOption{
		grpcretry.WithCodes(codes.NotFound, codes.DeadlineExceeded, codes.Aborted),
		grpcretry.WithMax(uint(retries)),
		grpcretry.WithPerRetryTimeout(timeout),
	}

	logOpts := []grpclog.Option{
		grpclog.WithLogOnEvents(grpclog.PayloadSent, grpclog.PayloadReceived),
	}

	conn, err := grpc.NewClient(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			grpclog.UnaryClientInterceptor(InterceptorLogger(log), logOpts...),
			grpcretry.UnaryClientInterceptor(retryOpts...),
		))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Client{api: taskv1.NewTodoClient(conn)}, nil
}

func InterceptorLogger(l *slog.Logger) grpclog.Logger {
	return grpclog.LoggerFunc(func(ctx context.Context, lvl grpclog.Level, msg string, fields ...any) {
		l.Log(ctx, slog.Level(lvl), msg, fields...)
	})
}

func (c *Client) CreateTask(ctx context.Context, authorID uuid.UUID, title, description, deadline string) (string, error) {
	const op = "task.grpc.CreateTask"

	resp, err := c.api.CreateTask(ctx, &taskv1.NewTaskRequest{
		AuthorId:    authorID.String(),
		Title:       title,
		Description: description,
		Deadline:    deadline,
	})
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return resp.TaskId, nil
}

func (c *Client) GetTask(ctx context.Context, authorID uuid.UUID) ([]*models.Task, error) {
	const op = "task.grpc.GetTask"

	resp, err := c.api.GetTask(ctx, &taskv1.TaskRequest{
		AuthorId: authorID.String(),
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	tasks := make([]*models.Task, len(resp.Tasks))
	for i := range resp.Tasks {
		var deadline time.Time
		if resp.Tasks[i].Deadline != "" {
			deadline, err = time.Parse(time.RFC3339, resp.Tasks[i].Deadline)
			if err != nil {
				return nil, fmt.Errorf("%s: failed to parse deadline: %w", op, err)
			}
		}

		task := &models.Task{
			ID:          uuid.MustParse(resp.Tasks[i].Id),
			AuthorID:    uuid.MustParse(resp.Tasks[i].AuthorId),
			Title:       resp.Tasks[i].Title,
			Description: resp.Tasks[i].Description,
			Status:      resp.Tasks[i].Status,
			Deadline:    deadline,
		}
		tasks[i] = task
	}

	return tasks, nil
}

func (c *Client) UpdateTask(ctx context.Context, taskID, authorID uuid.UUID, title, description, status, deadline string) error {
	const op = "task.grpc.UpdateTask"

	_, err := c.api.UpdateTask(ctx, &taskv1.UpdateRequest{
		Id:             taskID.String(),
		AuthorId:       authorID.String(),
		NewTitle:       title,
		NewDescription: description,
		NewStatus:      status,
		NewDeadline:    deadline,
	})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (c *Client) DeleteTask(ctx context.Context, taskID, authorID uuid.UUID) error {
	const op = "task.grpc.DeleteTask"

	_, err := c.api.DeleteTask(ctx, &taskv1.DeleteRequest{
		TaskId:   taskID.String(),
		AuthorId: authorID.String(),
	})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
