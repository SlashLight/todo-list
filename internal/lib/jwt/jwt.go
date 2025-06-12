package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/SlashLight/todo-list/internal/domain/models"
)

func NewToken(user *models.User, secret string, duration time.Duration) (string, error) {
	secret = "my_very_secret_key"

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

func ParseToken(tokenString, secret string) (*models.Session, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		strID := claims["uid"].(string)
		userID, err := uuid.Parse(strID)
		if err != nil {
			return nil, fmt.Errorf("failed to parse user ID: %w", err)
		}
		session := &models.Session{
			UserID: userID,
			Email:  claims["email"].(string),
		}
		return session, nil
	}

	return nil, fmt.Errorf("invalid token claims or token is not valid")
}
