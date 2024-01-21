package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"infotecs_trainee_task/internal/entity"
	"infotecs_trainee_task/internal/repo"
	"infotecs_trainee_task/internal/repo/repoerrors"
	"log/slog"
)

type WalletService struct {
	walletRepo      repo.Wallet
	transactionRepo repo.Transaction
}

func NewWalletService(walletRepo repo.Wallet, transactionRepo repo.Transaction) *WalletService {
	return &WalletService{
		walletRepo:      walletRepo,
		transactionRepo: transactionRepo,
	}
}

func (s *WalletService) CreateWallet(ctx context.Context) (uuid.UUID, error) {

	wallet := entity.Wallet{
		Balance: 100.0,
	}

	walletUUID, err := s.walletRepo.CreateWallet(ctx, wallet)
	if err != nil {
		if errors.Is(err, repoerrors.ErrAlreadyExist) {
			return uuid.Nil, ErrWalletAlreadyExists
		}
		slog.Error("WalletService.CreateWallet", err)
		return uuid.Nil, ErrCannotCreateWallet
	}

	return walletUUID, nil
}

func (s *WalletService) MakeTransaction(ctx context.Context, sender, receiver uuid.UUID, amount float64) error {

	senderWallet, err := s.walletRepo.GetWalletStateById(ctx, sender)
	if err != nil {
		if errors.Is(err, repoerrors.ErrNotFound) {
			return ErrCannotGetWallet
		}
		slog.Error("WalletService.MakeTransaction", err)
		return ErrCannotCreateTransaction
	}

	receiverWallet, err := s.walletRepo.GetWalletStateById(ctx, receiver)
	if err != nil {
		if errors.Is(err, repoerrors.ErrNotFound) {
			return ErrCannotGetWallet
		}
		slog.Error("WalletService.MakeTransaction", err)
		return ErrCannotCreateTransaction
	}

	if senderWallet.Balance-amount < 0 {
		return fmt.Errorf("not enought money to send: %f", senderWallet.Balance)
	}

	err = s.walletRepo.CashTransfer(ctx, senderWallet, receiverWallet, amount)
	if err != nil {
		if errors.Is(err, pgx.ErrTxCommitRollback) {
			return fmt.Errorf("transaction rollbacked: %v", err)
		}
	}

	transaction := entity.Transaction{
		Sender:   sender,
		Receiver: receiver,
		Amount:   amount,
	}

	err = s.transactionRepo.CreateTransaction(ctx, transaction)
	if err != nil {
		if errors.Is(err, repoerrors.ErrAlreadyExist) {
			return ErrTransactionAlreadyExists
		}
		slog.Error("WalletService.MakeTransaction", err)
		return ErrCannotCreateTransaction
	}

	return nil
}

func (s *WalletService) GetWalletState(ctx context.Context, walletUUID uuid.UUID) (entity.Wallet, error) {
	return s.walletRepo.GetWalletStateById(ctx, walletUUID)
}

func (s *WalletService) GetTransactionsHistory(ctx context.Context, walletUUID uuid.UUID) ([]entity.Transaction, error) {
	_, err := s.walletRepo.GetWalletStateById(ctx, walletUUID)
	if err != nil {
		if errors.Is(err, repoerrors.ErrNotFound) {
			return []entity.Transaction{}, ErrCannotGetWallet
		}
	}
	return s.transactionRepo.GetWalletHistory(ctx, walletUUID)
}
