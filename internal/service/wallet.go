package service

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"infotecs_trainee_task/internal/entity"
	"infotecs_trainee_task/internal/repo"
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
		if errors.Is(err, repo.ErrAlreadyExist) {
			return uuid.Nil, ErrWalletAlreadyExists
		}
		slog.Error("WalletService.CreateWallet", err)
		return uuid.Nil, ErrCannotCreateWallet
	}

	return walletUUID, nil
}

func (s *WalletService) MakeTransaction(ctx context.Context, sender, receiver uuid.UUID, amount float64) error {
	transaction := entity.Transaction{
		Sender:   sender,
		Receiver: receiver,
		Amount:   amount,
	}

	err := s.transactionRepo.CreateTransaction(ctx, transaction)
	if err != nil {
		if errors.Is(err, repo.ErrAlreadyExist) {
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
	return s.transactionRepo.GetWalletHistory(ctx, walletUUID)
}
