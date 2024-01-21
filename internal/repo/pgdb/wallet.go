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
	"infotecs_trainee_task/internal/repo/repoerrors"
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
				return uuid.Nil, repoerrors.ErrAlreadyExist
			}
		}
		return uuid.Nil, fmt.Errorf("WalletRepo.CreateWallet - r.Pool.QueryRow: %v", err)
	}

	return walletUUID, nil
}

func (r *WalletRepo) GetWalletStateById(ctx context.Context, uuid uuid.UUID) (entity.Wallet, error) {
	sql, args, _ := r.Builder.
		Select("uuid", "balance").
		From("wallets").
		Where(squirrel.Eq{"uuid": uuid}).
		ToSql()

	var wallet entity.Wallet
	err := r.Pool.QueryRow(ctx, sql, args...).Scan(&wallet.UUID, &wallet.Balance)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Wallet{}, repoerrors.ErrNotFound
		}
		return entity.Wallet{}, fmt.Errorf("WalletRepo.GetWalletStateById - r.Pool.QueryRow: %v", err)
	}

	return wallet, nil
}

func (r *WalletRepo) CashTransfer(ctx context.Context, sender, receiver entity.Wallet, amount float64) error {
	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("WalletRepo.CashTransfer - r.Pool.Begin: %v", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	sql, args, _ := r.Builder.
		Update("wallets").
		Set("balance", sender.Balance-amount).
		Where(squirrel.Eq{"uuid": sender.UUID}).
		ToSql()

	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("WalletRepo.CashTransfer - tx.Exec: %v", err)
	}

	sql, args, _ = r.Builder.
		Update("wallets").
		Set("balance", receiver.Balance+amount).
		Where(squirrel.Eq{"uuid": receiver.UUID}).
		ToSql()

	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("WalletRepo.CashTransfer - tx.Exec: %v", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("WalletRepo.CashTransfer - tx.Commit: %v\n", err)
	}

	return nil
}
