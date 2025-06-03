package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/SlashLight/todo-list/internal/domain/models"
)

func NewToken(user *models.User, secret string, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = user.ID
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(duration).Unix()

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func SessionFromToken(token *jwt.Token) *models.Session {
	sess := &models.Session{
		UserID: token.Claims.(jwt.MapClaims)["uid"].(uuid.UUID),
		Email:  token.Claims.(jwt.MapClaims)["email"].(string),
	}

	return sess
}
