package entity

import (
	"time"

	"github.com/google/uuid"
)

type Entry struct {
	ID        int64     `json:"id"`
	AccountID uuid.UUID `json:"account_id"`
	Amount    int64     `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}
