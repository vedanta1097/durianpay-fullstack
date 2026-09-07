package usecase

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/durianpay/fullstack-boilerplate/internal/entity"
	"github.com/durianpay/fullstack-boilerplate/internal/module/auth/repository"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase interface {
	Login(ctx context.Context, email string, password string) (string, *entity.User, error)
	VerifyToken(token string) error
}

type Auth struct {
	repo      repository.UserRepository
	jwtSecret []byte
	ttl       time.Duration
}

func NewAuthUsecase(repo repository.UserRepository, jwtSecret []byte, ttl time.Duration) *Auth {
	return &Auth{repo: repo, jwtSecret: jwtSecret, ttl: ttl}
}

// Login verifies email + password and returns a JWT.
func (a *Auth) Login(ctx context.Context, email string, password string) (string, *entity.User, error) {
	if strings.TrimSpace(email) == "" || password == "" {
		return "", nil, entity.ErrorBadRequest("email and password are required")
	}

	user, err := a.repo.GetUserByEmail(ctx, email)
	if err != nil {
		var appErr *entity.AppError
		if errors.As(err, &appErr) && appErr.Code == entity.ErrorCodeNotFound {
			return "", nil, entity.ErrorUnauthorized("invalid credentials")
		}
		return "", nil, err
	}
	if user.ID == "" {
		return "", nil, entity.ErrorNotFound("user not found")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", nil, entity.WrapError(err, entity.ErrorCodeUnauthorized, "invalid credentials")
	}

	claims := jwt.MapClaims{
		"sub":  user.ID,
		"role": user.Role,
		"exp":  time.Now().Add(a.ttl).Unix(),
		"iat":  time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(a.jwtSecret)
	if err != nil {
		return "", nil, entity.WrapError(err, entity.ErrorCodeUnauthorized, "invalid credentials")
	}
	return signed, user, nil
}

func (a *Auth) VerifyToken(rawToken string) error {
	if rawToken == "" {
		return entity.ErrorUnauthorized("missing or invalid token")
	}

	token, err := jwt.Parse(rawToken, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, entity.ErrorUnauthorized("missing or invalid token")
		}
		return a.jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return entity.ErrorUnauthorized("missing or invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return entity.ErrorUnauthorized("missing or invalid token")
	}
	subject, subjectOK := claims["sub"].(string)
	role, roleOK := claims["role"].(string)
	if !subjectOK || subject == "" || !roleOK || !entity.IsSupportedRole(role) {
		return entity.ErrorUnauthorized("missing or invalid token")
	}

	return nil
}
