package pgdb

import (
	"context"
	"errors"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"infotecs_trainee_task/internal/entity"
	"infotecs_trainee_task/internal/repo"
	"infotecs_trainee_task/pkg/postgres"
)

type WalletRepo struct {
	*postgres.Postgres
}

func NewWalletRepo(pg *postgres.Postgres) *WalletRepo {
	return &WalletRepo{pg}
}

func (r *WalletRepo) CreateWallet(ctx context.Context, wallet entity.Wallet) (uuid.UUID, error) {
	sql, args, _ := r.Builder.
		Insert("wallets").
		Columns("balance").
		Values(wallet.Balance).
		Suffix("RETURNING uuid").
		ToSql()

	var walletUUID uuid.UUID
	err := r.Pool.QueryRow(ctx, sql, args...).Scan(&walletUUID)

	if err != nil {
		var pgErr *pgconn.PgError
		if ok := errors.As(err, &pgErr); ok {
			if pgErr.Code == "23505" {
				return uuid.Nil, repo.ErrAlreadyExist
			}
		}
		return uuid.Nil, fmt.Errorf("WalletRepo.CreateWallet - r.Pool.QueryRow: %v", err)
	}

	return walletUUID, nil
}

func (r *WalletRepo) ChangeBalance(ctx context.Context, wallet entity.Wallet, amount float64) error {
	sql, args, _ := r.Builder.
		Update("wallets").
		Set("wallet", amount).
		Where(squirrel.Eq{"uuid": wallet.UUID}).
		ToSql()

	_, err := r.Pool.Exec(ctx, sql, args...)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return repo.ErrNotFound
		}
		return fmt.Errorf("WalletRepo.ChangeBalance - r.Pool.Exec: %v", err)
	}

	return nil
}

func (r *WalletRepo) GetWalletStateById(ctx context.Context, uuid uuid.UUID) (entity.Wallet, error) {
	sql, args, _ := r.Builder.
		Select("uuid", "balance").
		From("wallets").
		Where(squirrel.Eq{"uuid": uuid}).
		ToSql()

	var wallet entity.Wallet
	err := r.Pool.QueryRow(ctx, sql, args...).Scan(&wallet)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Wallet{}, repo.ErrNotFound
		}
		return entity.Wallet{}, fmt.Errorf("WalletRepo.GetWalletStateById - r.Pool.QueryRow: %v", err)
	}

	return wallet, nil
}
