package pgdb

import (
	"context"
	"errors"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"infotecs_trainee_task/internal/entity"
	"infotecs_trainee_task/internal/repo/repoerrors"
	"infotecs_trainee_task/pkg/postgres"
)

type TransactionRepo struct {
	*postgres.Postgres
}

func NewTransactionRepo(pg *postgres.Postgres) *TransactionRepo {
	return &TransactionRepo{pg}
}

func (r *TransactionRepo) CreateTransaction(ctx context.Context, transaction entity.Transaction) error {
	sql, args, _ := r.Builder.
		Insert("transactions").
		Columns("sender", "receiver", "amount").
		Values(transaction.Sender, transaction.Receiver, transaction.Amount).
		ToSql()

	_, err := r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if ok := errors.As(err, &pgErr); ok {
			if pgErr.Code == "23505" {
				return repoerrors.ErrAlreadyExist
			}
		}
		return fmt.Errorf("TransactionRepo.CreateTransaction - r.Pool.Exec: %v", err)
	}

	return nil
}

func (r *TransactionRepo) GetWalletHistory(ctx context.Context, uuid uuid.UUID) ([]entity.Transaction, error) {
	sql, args, _ := r.Builder.
		Select("sender", "receiver", "created_at", "amount").
		From("transactions").
		Where(squirrel.Or{
			squirrel.Eq{"sender": uuid},
			squirrel.Eq{"receiver": uuid},
		}).
		ToSql()

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("TransactionRepo.GetWalletHistory - r.Pool.Query: %v", err)
	}
	defer rows.Close()

	var transactions []entity.Transaction
	for rows.Next() {
		var transaction entity.Transaction
		if err = rows.Scan(&transaction.Sender, &transaction.Receiver, &transaction.CreatedAt, &transaction.Amount); err != nil {
			return nil, fmt.Errorf("TransactionRepo.GetWalletHistory - rows.Scan: %v", err)
		}
		transactions = append(transactions, transaction)
	}

	return transactions, nil
}
