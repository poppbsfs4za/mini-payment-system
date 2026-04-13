package entities

import "time"

type Transaction struct {
	ID            string    `gorm:"type:uuid;primaryKey" json:"id"`
	FromAccountID string    `gorm:"type:uuid;not null;index" json:"from_account_id"`
	ToAccountID   string    `gorm:"type:uuid;not null;index" json:"to_account_id"`
	Amount        int64     `gorm:"not null;check:amount_positive,amount > 0" json:"amount"`
	Currency      string    `gorm:"size:10;not null;index" json:"currency"`
	Status        string    `gorm:"size:30;not null;index" json:"status"`
	Reference     string    `gorm:"size:120;index" json:"reference"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
