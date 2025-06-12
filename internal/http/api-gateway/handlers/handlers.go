package handlers

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/google/uuid"

	"github.com/SlashLight/todo-list/internal/domain/models"
)

type AuthAPI interface {
	Register(ctx context.Context, email, password string) (string, error)
	Login(ctx context.Context, email, password string) (string, error)
}

type TaskAPI interface {
	CreateTask(ctx context.Context, authorID uuid.UUID, title, description, deadline string) (string, error)
	GetTask(ctx context.Context, authorID uuid.UUID) ([]*models.Task, error)
	UpdateTask(ctx context.Context, taskID, authorID uuid.UUID, title, description, status, deadline string) error
	DeleteTask(ctx context.Context, taskID, authorID uuid.UUID) error
}

type APIGateway struct {
	Auth AuthAPI
	Task TaskAPI
	log  *slog.Logger
}

func New(authAPI AuthAPI, taskAPI TaskAPI, log *slog.Logger) *APIGateway {
	return &APIGateway{
		Auth: authAPI,
		Task: taskAPI,
		log:  log,
	}
}

func (api *APIGateway) HandleRegister(w http.ResponseWriter, r *http.Request) {
	const op = "APIGateway.HandleRegister"

	log := api.log.With(slog.String("op", op))

	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		//TODO: ...
		log.Error("failed to decode request body", slog.String("error", err.Error()))
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	id, err := api.Auth.Register(r.Context(), req.Email, req.Password)
	if err != nil {
		//TODO: ...
		log.Error("failed to register user", slog.String("error", err.Error()))
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}

	log.Info("User registered successfully", "userID", id)
	w.WriteHeader(http.StatusCreated)
}

func (api *APIGateway) HandleLogin(w http.ResponseWriter, r *http.Request) {
	const op = "APIGateway.HandleLogin"

	log := api.log.With(slog.String("op", op))

	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error("failed to decode request body", slog.String("error", err.Error()))
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	token, err := api.Auth.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		log.Error("failed to login user", slog.String("error", err.Error()))
		http.Error(w, "Failed to login user", http.StatusInternalServerError)
		return
	}

	log.Info("User logged in successfully")
	w.Header().Set("Set-Cookie", "token="+token+"; HttpOnly; Secure; SameSite=Strict; Path=/")
	w.WriteHeader(http.StatusOK)
}

func (api *APIGateway) HandleCreateTask(w http.ResponseWriter, r *http.Request) {
	const op = "APIGateway.HandleCreateTask"

	sess, err := models.SessionFromContext(r.Context())
	if err != nil {
		api.log.Error("failed to get session from context", slog.String("error", err.Error()))
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	log := api.log.With(
		slog.String("op", op),
		slog.String("userID", sess.UserID.String()))

	var req struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Status      string `json:"status"`
		Deadline    string `json:"deadline"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error("failed to decode request body", slog.String("error", err.Error()))
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	taskID, err := api.Task.CreateTask(r.Context(), sess.UserID, req.Title, req.Description, req.Deadline)
	if err != nil {
		log.Error("failed to create task", slog.String("error", err.Error()))
		http.Error(w, "Failed to create task", http.StatusInternalServerError)
		return
	}

	log.Info("Task created successfully", "taskID", taskID)
	w.WriteHeader(http.StatusCreated)
}

func (api *APIGateway) HandleGetTask(w http.ResponseWriter, r *http.Request) {
	const op = "APIGateway.HandleGetTask"

	sess, err := models.SessionFromContext(r.Context())
	if err != nil {
		api.log.Error("failed to get session from context", slog.String("error", err.Error()))
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	log := api.log.With(
		slog.String("op", op),
		slog.String("userID", sess.UserID.String()))

	tasks, err := api.Task.GetTask(r.Context(), sess.UserID)
	if err != nil {
		log.Error("failed to get tasks", slog.String("error", err.Error()))
		http.Error(w, "Failed to get tasks", http.StatusInternalServerError)
		return
	}

	if len(tasks) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if err := json.NewEncoder(w).Encode(tasks); err != nil {
		log.Error("failed to encode tasks", slog.String("error", err.Error()))
		http.Error(w, "Failed to encode tasks", http.StatusInternalServerError)
		return
	}
}

func (api *APIGateway) HandleUpdateTask(w http.ResponseWriter, r *http.Request) {
	return
}

func (api *APIGateway) HandleDeleteTask(w http.ResponseWriter, r *http.Request) {
	return
}
