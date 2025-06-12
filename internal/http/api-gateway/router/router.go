package router

import (
	"net/http"

	"github.com/SlashLight/todo-list/internal/http/api-gateway/middleware"
)

type API interface {
	HandleLogin(w http.ResponseWriter, r *http.Request)
	HandleRegister(w http.ResponseWriter, r *http.Request)

	HandleCreateTask(w http.ResponseWriter, r *http.Request)
	HandleGetTask(w http.ResponseWriter, r *http.Request)
}

func New(api API, secret string) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/auth/login", api.HandleLogin)
	mux.HandleFunc("/auth/register", api.HandleRegister)

	mux.Handle("/tasks/create", middleware.AuthMiddleware(http.HandlerFunc(api.HandleCreateTask), secret))
	mux.Handle("/tasks/get", middleware.AuthMiddleware(http.HandlerFunc(api.HandleGetTask), secret))

	return mux
}
