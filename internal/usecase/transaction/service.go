package transaction

import (
	"context"
	"strings"

	"mini-payment-system/internal/domain/entities"
	"mini-payment-system/internal/domain/errs"
	"mini-payment-system/internal/domain/repositories"

	"github.com/google/uuid"
)

const transactionStatusCompleted = "completed"

type Service struct {
	accountRepo     repositories.AccountRepository
	transactionRepo repositories.TransactionRepository
	txManager       repositories.TxManager
}

func NewService(
	accountRepo repositories.AccountRepository,
	transactionRepo repositories.TransactionRepository,
	txManager repositories.TxManager,
) *Service {
	return &Service{
		accountRepo:     accountRepo,
		transactionRepo: transactionRepo,
		txManager:       txManager,
	}
}

func (s *Service) CreateTransfer(ctx context.Context, fromAccountID, toAccountID string, amount int64, reference string) (*entities.Transaction, error) {
	fromAccountID = strings.TrimSpace(fromAccountID)
	toAccountID = strings.TrimSpace(toAccountID)
	reference = strings.TrimSpace(reference)

	if fromAccountID == "" || toAccountID == "" {
		return nil, errs.InvalidInput("from_account_id and to_account_id are required")
	}
	if fromAccountID == toAccountID {
		return nil, errs.InvalidInput("source and destination accounts must be different")
	}
	if amount <= 0 {
		return nil, errs.InvalidInput("amount must be greater than zero")
	}

	var transaction *entities.Transaction

	err := s.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		firstLockID, secondLockID := orderedIDs(fromAccountID, toAccountID)

		firstAccount, err := s.accountRepo.GetByIDForUpdate(txCtx, firstLockID)
		if err != nil {
			return errs.NotFound("account not found")
		}
		secondAccount, err := s.accountRepo.GetByIDForUpdate(txCtx, secondLockID)
		if err != nil {
			return errs.NotFound("account not found")
		}

		fromAccount, toAccount := mapLockedAccounts(firstAccount, secondAccount, fromAccountID, toAccountID)

		if fromAccount.Currency != toAccount.Currency {
			return errs.Conflict("cross-currency transfer is not supported")
		}
		if fromAccount.Balance < amount {
			return errs.ErrInsufficientBalance
		}

		fromAccount.Balance -= amount
		toAccount.Balance += amount

		if err := s.accountRepo.Update(txCtx, fromAccount); err != nil {
			return err
		}
		if err := s.accountRepo.Update(txCtx, toAccount); err != nil {
			return err
		}

		transaction = &entities.Transaction{
			ID:            uuid.NewString(),
			FromAccountID: fromAccount.ID,
			ToAccountID:   toAccount.ID,
			Amount:        amount,
			Currency:      fromAccount.Currency,
			Status:        transactionStatusCompleted,
			Reference:     reference,
		}

		if err := s.transactionRepo.Create(txCtx, transaction); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

func (s *Service) GetByID(ctx context.Context, id string) (*entities.Transaction, error) {
	return s.transactionRepo.GetByID(ctx, id)
}

func (s *Service) List(ctx context.Context) ([]entities.Transaction, error) {
	return s.transactionRepo.List(ctx)
}

func orderedIDs(a, b string) (string, string) {
	if a < b {
		return a, b
	}
	return b, a
}

func mapLockedAccounts(first, second *entities.Account, fromID, toID string) (*entities.Account, *entities.Account) {
	if first.ID == fromID {
		return first, second
	}
	return second, first
}
