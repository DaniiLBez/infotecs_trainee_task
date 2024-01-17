package service

import (
	"context"
	"github.com/google/uuid"
	"infotecs_trainee_task/internal/repo"
	"infotecs_trainee_task/pkg/hasher"
	"time"
)

type Auth interface {
	CreateUser(ctx context.Context, input struct {
		username string
		password string
	}) (uuid.UUID, error)

	GenerateToken(ctx context.Context, input struct {
		username string
		password string
	}) (string, error)

	ParseToken(string) (uuid.UUID, error)
}

type Services struct {
	AuthService Auth
}

type Dependencies struct {
	Repos    *repo.Repositories
	Hasher   hasher.PasswordHasher
	SignKey  string
	TokenTTL time.Duration
}
