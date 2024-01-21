package service

import (
	"context"
	"github.com/google/uuid"
	"infotecs_trainee_task/internal/entity"
	"infotecs_trainee_task/internal/repo"
	"infotecs_trainee_task/pkg/hasher"
	"time"
)

type InputData struct {
	Username string
	Password string
}

type Auth interface {
	CreateUser(ctx context.Context, input InputData) (uuid.UUID, error)
	GenerateToken(ctx context.Context, input InputData) (string, error)
	ParseToken(string) (uuid.UUID, error)
}

type Wallet interface {
	CreateWallet(ctx context.Context) (uuid.UUID, error)
	MakeTransaction(ctx context.Context, sender, receiver uuid.UUID, amount float64) error
	GetWalletState(ctx context.Context, walletUUID uuid.UUID) (entity.Wallet, error)
	GetTransactionsHistory(ctx context.Context, walletUUID uuid.UUID) ([]entity.Transaction, error)
}

type Services struct {
	Auth   Auth
	Wallet Wallet
}

type Dependencies struct {
	Repos    *repo.Repositories
	Hasher   hasher.PasswordHasher
	SignKey  string
	TokenTTL time.Duration
}

func NewServices(deps *Dependencies) *Services {
	return &Services{
		Auth:   NewAuthService(deps.Repos.User, deps.Hasher, deps.SignKey, deps.TokenTTL),
		Wallet: NewWalletService(deps.Repos.Wallet, deps.Repos.Transaction),
	}
}
