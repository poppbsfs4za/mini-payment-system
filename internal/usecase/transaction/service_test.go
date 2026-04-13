package transaction

import (
	"context"
	"errors"
	"testing"

	"mini-payment-system/internal/domain/entities"
	"mini-payment-system/internal/domain/errs"
)

type mockAccountRepo struct {
	accounts      map[string]*entities.Account
	failOnUpdate  bool
	lockedAccount []string
}

func (m *mockAccountRepo) Create(ctx context.Context, account *entities.Account) error { return nil }

func (m *mockAccountRepo) GetByID(ctx context.Context, id string) (*entities.Account, error) {
	acc, ok := m.accounts[id]
	if !ok {
		return nil, errors.New("account not found")
	}
	copy := *acc
	return &copy, nil
}

func (m *mockAccountRepo) GetByIDForUpdate(ctx context.Context, id string) (*entities.Account, error) {
	m.lockedAccount = append(m.lockedAccount, id)
	return m.GetByID(ctx, id)
}

func (m *mockAccountRepo) List(ctx context.Context) ([]entities.Account, error) { return nil, nil }

func (m *mockAccountRepo) Update(ctx context.Context, account *entities.Account) error {
	if m.failOnUpdate {
		return errors.New("update failed")
	}
	copied := *account
	m.accounts[account.ID] = &copied
	return nil
}

func (m *mockAccountRepo) Delete(ctx context.Context, id string) error { return nil }

type mockTransactionRepo struct {
	items        []entities.Transaction
	failOnCreate bool
}

func (m *mockTransactionRepo) Create(ctx context.Context, transaction *entities.Transaction) error {
	if m.failOnCreate {
		return errors.New("create failed")
	}
	m.items = append(m.items, *transaction)
	return nil
}

func (m *mockTransactionRepo) GetByID(ctx context.Context, id string) (*entities.Transaction, error) {
	return nil, nil
}

func (m *mockTransactionRepo) List(ctx context.Context) ([]entities.Transaction, error) {
	return nil, nil
}

type mockTxManager struct{}

func (m *mockTxManager) WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

func TestCreateTransfer_Success(t *testing.T) {
	accountRepo := &mockAccountRepo{accounts: map[string]*entities.Account{
		"from": {ID: "from", Balance: 1000, Currency: "THB"},
		"to":   {ID: "to", Balance: 300, Currency: "THB"},
	}}
	transactionRepo := &mockTransactionRepo{}
	service := NewService(accountRepo, transactionRepo, &mockTxManager{})

	trx, err := service.CreateTransfer(context.Background(), "from", "to", 250, "test-ref")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if trx.Amount != 250 {
		t.Fatalf("expected amount 250, got %d", trx.Amount)
	}
	if trx.Currency != "THB" {
		t.Fatalf("expected currency THB, got %s", trx.Currency)
	}
	if trx.Status != transactionStatusCompleted {
		t.Fatalf("expected status %s, got %s", transactionStatusCompleted, trx.Status)
	}
	if got := accountRepo.accounts["from"].Balance; got != 750 {
		t.Fatalf("expected from balance 750, got %d", got)
	}
	if got := accountRepo.accounts["to"].Balance; got != 550 {
		t.Fatalf("expected to balance 550, got %d", got)
	}
	if len(transactionRepo.items) != 1 {
		t.Fatalf("expected 1 transaction recorded, got %d", len(transactionRepo.items))
	}
	if len(accountRepo.lockedAccount) != 2 {
		t.Fatalf("expected two locked accounts, got %d", len(accountRepo.lockedAccount))
	}
}

func TestCreateTransfer_InsufficientBalance(t *testing.T) {
	accountRepo := &mockAccountRepo{accounts: map[string]*entities.Account{
		"from": {ID: "from", Balance: 100, Currency: "THB"},
		"to":   {ID: "to", Balance: 300, Currency: "THB"},
	}}
	transactionRepo := &mockTransactionRepo{}
	service := NewService(accountRepo, transactionRepo, &mockTxManager{})

	_, err := service.CreateTransfer(context.Background(), "from", "to", 250, "test-ref")
	if !errors.Is(err, errs.ErrInsufficientBalance) {
		t.Fatalf("expected insufficient balance error, got %v", err)
	}
	if got := accountRepo.accounts["from"].Balance; got != 100 {
		t.Fatalf("expected from balance unchanged, got %d", got)
	}
	if len(transactionRepo.items) != 0 {
		t.Fatalf("expected no transactions recorded, got %d", len(transactionRepo.items))
	}
}

func TestCreateTransfer_SameAccount(t *testing.T) {
	service := NewService(&mockAccountRepo{accounts: map[string]*entities.Account{}}, &mockTransactionRepo{}, &mockTxManager{})

	_, err := service.CreateTransfer(context.Background(), "same", "same", 100, "ref")
	if !errors.Is(err, errs.ErrInvalidInput) {
		t.Fatalf("expected invalid input error, got %v", err)
	}
}

func TestCreateTransfer_CurrencyMismatch(t *testing.T) {
	accountRepo := &mockAccountRepo{accounts: map[string]*entities.Account{
		"from": {ID: "from", Balance: 1000, Currency: "THB"},
		"to":   {ID: "to", Balance: 300, Currency: "USD"},
	}}
	service := NewService(accountRepo, &mockTransactionRepo{}, &mockTxManager{})

	_, err := service.CreateTransfer(context.Background(), "from", "to", 100, "ref")
	if !errors.Is(err, errs.ErrConflict) {
		t.Fatalf("expected conflict error, got %v", err)
	}
}

func TestCreateTransfer_FailureOnTransactionCreate(t *testing.T) {
	accountRepo := &mockAccountRepo{accounts: map[string]*entities.Account{
		"from": {ID: "from", Balance: 1000, Currency: "THB"},
		"to":   {ID: "to", Balance: 300, Currency: "THB"},
	}}
	transactionRepo := &mockTransactionRepo{failOnCreate: true}
	service := NewService(accountRepo, transactionRepo, &mockTxManager{})

	_, err := service.CreateTransfer(context.Background(), "from", "to", 200, "ref")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if len(transactionRepo.items) != 0 {
		t.Fatalf("expected no transactions recorded, got %d", len(transactionRepo.items))
	}
}
