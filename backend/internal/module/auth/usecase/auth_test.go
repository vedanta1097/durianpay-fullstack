package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/durianpay/fullstack-boilerplate/internal/entity"
)

type userRepositoryStub struct {
	user *entity.User
	err  error
}

func (s userRepositoryStub) GetUserByEmail(context.Context, string) (*entity.User, error) {
	return s.user, s.err
}

func TestAuthLoginAndVerifyToken(t *testing.T) {
	hash, err := passwordHash("password")
	if err != nil {
		t.Fatal(err)
	}
	auth := NewAuthUsecase(userRepositoryStub{user: &entity.User{ID: "user-1", Email: "cs@test.com", PasswordHash: hash, Role: entity.RoleCS}}, []byte("test-secret"), time.Hour)

	token, user, err := auth.Login(context.Background(), "cs@test.com", "password")
	if err != nil {
		t.Fatalf("Login() error = %v", err)
	}
	if user.Email != "cs@test.com" {
		t.Fatalf("Login() user email = %q", user.Email)
	}
	if err := auth.VerifyToken(token); err != nil {
		t.Fatalf("VerifyToken() error = %v", err)
	}
	if err := auth.VerifyToken(token + "tampered"); err == nil {
		t.Fatal("VerifyToken() accepted a tampered token")
	}
}

func TestAuthLoginRejectsUnknownUserAndWrongPassword(t *testing.T) {
	hash, err := passwordHash("password")
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name     string
		repo     userRepositoryStub
		password string
	}{
		{name: "unknown user", repo: userRepositoryStub{err: entity.ErrorNotFound("user not found")}, password: "password"},
		{name: "wrong password", repo: userRepositoryStub{user: &entity.User{ID: "user-1", PasswordHash: hash, Role: entity.RoleCS}}, password: "wrong"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			auth := NewAuthUsecase(tt.repo, []byte("test-secret"), time.Hour)
			_, _, err := auth.Login(context.Background(), "cs@test.com", tt.password)
			if err == nil {
				t.Fatal("Login() error = nil")
			}
			appErr, ok := err.(*entity.AppError)
			if !ok || appErr.Code != entity.ErrorCodeUnauthorized {
				t.Fatalf("Login() error = %#v, want unauthorized AppError", err)
			}
		})
	}
}
