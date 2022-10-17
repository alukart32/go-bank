package entity

import (
	"time"

	"github.com/google/uuid"
)

type Account struct {
	ID        uuid.UUID `json:"id"`
	Owner     string    `json:"owner"`
	Balance   int64     `json:"balance"`
	Currency  Currency  `json:"currency"`
	CreatedAt time.Time `json:"created_at"`
}
