package auth_service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/SlashLight/todo-list/internal/domain/models"
	"github.com/SlashLight/todo-list/internal/lib/jwt"
	"github.com/SlashLight/todo-list/pkg/my_err"
)

type UserSaver interface {
	Register(ctx context.Context, user *models.User) error
}

type UserProvider interface {
	GetByEmail(ctx context.Context, email string) (*models.User, error)
}

type Service struct {
	userSaver    UserSaver
	UserProvider UserProvider
	logger       *slog.Logger
	tokenTTL     time.Duration
	tokenSecret  string
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

func New(userSaver UserSaver, userGetter UserProvider, log *slog.Logger, tokenTTL time.Duration) *Service {
	return &Service{
		userSaver:    userSaver,
		UserProvider: userGetter,
		logger:       log,
		tokenTTL:     tokenTTL,
	}
}

func (uc *Service) Register(ctx context.Context, email, password string) (string, error) {
	const op = "auth.Register"

	log := uc.logger.With(
		slog.String("op", op),
		slog.String("email", email),
	)

	log.Info("registering user")

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate password hash", err)
		return "", fmt.Errorf("%s: %w", op, err)
	}

	user := &models.User{
		ID:           uuid.New(),
		Email:        email,
		PasswordHash: string(hashedPass),
	}

	err = uc.userSaver.Register(ctx, user)
	return user.ID.String(), err
}

func (s *Service) Login(ctx context.Context, email string, pass string) (string, error) {
	const op = "auth.Login"

	log := s.logger.With(
		slog.String("op", op),
		slog.String("email", email),
	)

	log.Info("trying to login user")

	user, err := s.UserProvider.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, my_err.ErrUserNotFound) {
			s.logger.Warn("user not found", err)

			return "", fmt.Errorf("%s, %w", op, ErrInvalidCredentials)
		}

		s.logger.Error("failed to get user", err)

		return "", fmt.Errorf("%s: %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(pass)); err != nil {
		s.logger.Info("invalid credentials", err)

		return "", fmt.Errorf("%s, %w", op, ErrInvalidCredentials)
	}

	log.Info("user logged in successfully")

	token, err := jwt.NewToken(user, s.tokenSecret, s.tokenTTL)
	if err != nil {
		s.logger.Error("failed to generate token", err)

		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, err
}
