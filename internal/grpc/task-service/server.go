package task_service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	todov1 "github.com/SlashLight/todo-list/api/gen/go/todo"
	"github.com/SlashLight/todo-list/internal/domain/models"
	"github.com/SlashLight/todo-list/pkg/my_err"
)

type Service interface {
	CreateTask(ctx context.Context, authorID uuid.UUID, title, description string, deadline time.Time) (string, error)
	GetTasks(ctx context.Context, authorID uuid.UUID) ([]*models.Task, error)
	UpdateTask(ctx context.Context, newTask *models.Task) error
	DeleteTask(ctx context.Context, taskID, authorID uuid.UUID) error
}

type serverAPI struct {
	todov1.UnimplementedTodoServer
	service Service
}

func RegisterServerAPI(gRPC *grpc.Server, service Service) {
	todov1.RegisterTodoServer(gRPC, &serverAPI{service: service})
}

const timeLayout = time.RFC1123

func (s *serverAPI) CreateTask(ctx context.Context, req *todov1.NewTaskRequest) (*todov1.NewTaskResponse, error) {
	if req.GetTitle() == "" {
		return nil, status.Error(codes.InvalidArgument, "title is empty")
	}

	var deadline time.Time
	if req.GetDeadline() != "" {
		var err error
		deadline, err = time.Parse(timeLayout, req.GetDeadline())

		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "deadline has invalid format")
		}
	}

	authorID, err := validateUID(req.GetAuthorId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("invalid author ID: %s", err))
	}

	taskID, err := s.service.CreateTask(ctx, authorID, req.GetTitle(), req.GetDescription(), deadline)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &todov1.NewTaskResponse{TaskId: taskID}, nil
}

func (s *serverAPI) GetTasks(ctx context.Context, req *todov1.TaskRequest) (*todov1.TaskResponse, error) {
	authorID, err := validateUID(req.GetAuthorId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("invalid author ID: %s", err))
	}

	tasks, err := s.service.GetTasks(ctx, authorID)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	protoTasks := make([]*todov1.Task, len(tasks))
	for idx, task := range tasks {
		protoTasks[idx] = &todov1.Task{
			Id:          task.ID.String(),
			AuthorId:    task.AuthorID.String(),
			Title:       task.Title,
			Description: task.Description,
			Status:      task.Status,
			Deadline:    task.Deadline.Format(timeLayout),
		}
	}

	return &todov1.TaskResponse{Tasks: protoTasks}, nil
}

func (s *serverAPI) DeleteTask(ctx context.Context, req *todov1.DeleteRequest) (*todov1.EmptyResponse, error) {
	id, err := validateUID(req.GetTaskId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("invalid task ID: %s", err))
	}

	authorID, err := validateUID(req.GetAuthorId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("invalid author ID: %s", err))
	}

	err = s.service.DeleteTask(ctx, id, authorID)
	if err != nil {
		//TODO ...
		return nil, status.Error(codes.Internal, "internal error")
	}

	return nil, nil
}

func (s *serverAPI) UpdateTask(ctx context.Context, req *todov1.UpdateRequest) (*todov1.EmptyResponse, error) {
	newTask, err := validateNewTask(req)
	if err != nil {
		return nil, err
	}

	err = s.service.UpdateTask(ctx, newTask)
	if err != nil {
		//TODO ..
		return nil, status.Error(codes.Internal, "internal error")
	}

	return nil, nil
}

func validateNewTask(req *todov1.UpdateRequest) (*models.Task, error) {
	var newTask *models.Task

	id, err := validateUID(req.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("invalid task ID: %s", err))
	}

	authorID, err := validateUID(req.GetAuthorId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("invalid author ID: %s", err))
	}

	newTask.ID = id
	newTask.AuthorID = authorID
	if req.GetNewDeadline() != "" {
		newDeadline, err := time.Parse(timeLayout, req.GetNewDeadline())
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "wrong deadline format")
		}

		newTask.Deadline = newDeadline
	}

	newTask.Description = req.GetNewDescription()
	newTask.Title = req.GetNewTitle()
	newTask.Status = req.GetNewStatus()

	return newTask, nil
}

func validateUID(UIDString string) (uuid.UUID, error) {
	if UIDString == "" {
		return uuid.Nil, my_err.ErrEmptyField
	}

	authorID, err := uuid.Parse(UIDString)
	if err != nil {
		return uuid.Nil, my_err.ErrParseUUID
	}

	return authorID, nil
}
