package auth

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type Service struct {
	repository *Repository
	secret     string
}

func NewService(repository *Repository, secret string) *Service {
	return &Service{repository: repository, secret: secret}
}
func (s *Service) Login(ctx context.Context, username, password string) (string, error) {
	admin, err := s.repository.FindByUsername(ctx, strings.TrimSpace(username))
	if err != nil || !admin.Active || bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(password)) != nil {
		return "", ErrInvalidCredentials
	}
	claims := jwt.MapClaims{"sub": admin.ID.Hex(), "username": admin.Username, "role": "admin", "exp": time.Now().Add(8 * time.Hour).Unix(), "iat": time.Now().Unix()}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(s.secret))
}
