package auth_service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/SlashLight/todo-list/internal/models"
	"github.com/SlashLight/todo-list/pkg/my_err"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserRepo interface {
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Register(ctx context.Context, user *models.User) error
}

type AuthUsecase struct {
	userRepo UserRepo
	logger   *slog.Logger
}

func (uc *AuthUsecase) Register(ctx context.Context, email, password string) (uuid.UUID, error) {
	existingUser, _ := uc.userRepo.GetByEmail(ctx, email)
	if existingUser != nil {
		return uuid.Nil, my_err.ErrUserExists
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		uc.logger.Error("Registration error: ", err)
		return uuid.Nil, fmt.Errorf("registration: %w", err)
	}

	user := &models.User{
		ID:           uuid.New(),
		Email:        email,
		PasswordHash: string(hashedPass),
	}

	err = uc.userRepo.Register(ctx, user)
	return user.ID, err
}
