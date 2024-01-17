package repo

import (
	"context"
	"github.com/google/uuid"
	"infotecs_trainee_task/internal/entity"
)

type User interface {
	CreateUser(ctx context.Context, user entity.User) (uuid.UUID, error)
	GetUserById(ctx context.Context, uuid uuid.UUID) (entity.User, error)
	GetUserByUsername(ctx context.Context, username string) (entity.User, error)
	GetUserByUsernameAndPassword(ctx context.Context, username, password string) (entity.User, error)
}

type Repositories struct {
	User User
}
