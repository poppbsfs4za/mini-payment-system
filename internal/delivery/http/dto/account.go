package dto

type CreateAccountRequest struct {
	UserID         string `json:"user_id" binding:"required,uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	InitialBalance int64  `json:"initial_balance" binding:"gte=0" example:"100000"`
	Currency       string `json:"currency" example:"THB"`
}

type UpdateAccountRequest struct {
	Currency string `json:"currency,omitempty" example:"THB"`
}
