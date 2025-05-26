package auth_service

import (
	"context"

	"github.com/google/uuid"
)

type UseCase interface {
	Register(ctx context.Context, email, password string) (uuid.UUID, error)
}

type Handler struct {
	uc UseCase
}

type (h *Handler) Register(ctx context.context, req *pb.Re)
