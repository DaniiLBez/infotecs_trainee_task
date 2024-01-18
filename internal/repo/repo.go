package repo

import (
	"context"
	"github.com/google/uuid"
	"infotecs_trainee_task/internal/entity"
	"infotecs_trainee_task/internal/repo/pgdb"
	"infotecs_trainee_task/pkg/postgres"
)

type User interface {
	CreateUser(ctx context.Context, user entity.User) (uuid.UUID, error)
	GetUserById(ctx context.Context, uuid uuid.UUID) (entity.User, error)
	GetUserByUsername(ctx context.Context, username string) (entity.User, error)
	GetUserByUsernameAndPassword(ctx context.Context, username, password string) (entity.User, error)
}

type Wallet interface {
	CreateWallet(ctx context.Context, wallet entity.Wallet) (uuid.UUID, error)
	ChangeBalance(ctx context.Context, amount float64) error
	GetWalletStateById(ctx context.Context, uuid uuid.UUID) (entity.Wallet, error)
}

type Transaction interface {
	CreateTransaction(ctx context.Context, transaction entity.Transaction) error
	GetWalletHistory(ctx context.Context, uuid uuid.UUID) ([]entity.Transaction, error)
}

type Repositories struct {
	User        User
	Wallet      Wallet
	Transaction Transaction
}

func NewRepositories(pg *postgres.Postgres) *Repositories {
	return &Repositories{
		User:        pgdb.NewUserRepo(pg),
		Wallet:      nil,
		Transaction: nil,
	}
}
