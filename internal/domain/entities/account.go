package entities

import "time"

type Account struct {
	ID        string    `gorm:"type:uuid;primaryKey" json:"id"`
	UserID    string    `gorm:"type:uuid;not null;index" json:"user_id"`
	Balance   int64     `gorm:"not null;default:0;check:balance_non_negative,balance >= 0" json:"balance"`
	Currency  string    `gorm:"size:10;not null;default:THB;index" json:"currency"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
