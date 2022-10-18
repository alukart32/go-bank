package entity

import (
	"time"

	"github.com/google/uuid"
)

type Transfer struct {
	ID            int64     `json:"id"`
	FromAccountID uuid.UUID `json:"from_account_id"`
	ToAccountID   uuid.UUID `json:"to_account_id"`
	Amount        int64     `json:"amount"`
	FromEntryID   int64     `json:"from_entry_id"`
	ToEntryID     int64     `json:"to_entry_id"`
	CreatedAt     time.Time `json:"created_at"`
}

type TransferRes struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}
