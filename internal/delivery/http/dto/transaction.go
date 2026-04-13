package dto

type CreateTransactionRequest struct {
	FromAccountID string `json:"from_account_id" binding:"required,uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	ToAccountID   string `json:"to_account_id" binding:"required,uuid" example:"6ba7b810-9dad-11d1-80b4-00c04fd430c8"`
	Amount        int64  `json:"amount" binding:"required,gt=0" example:"5000"`
	Reference     string `json:"reference,omitempty" example:"invoice-1001"`
}
