package entity

import (
	"github.com/google/uuid"
)

type Wallet struct {
	UUID    uuid.UUID `db:"uuid" json:"id"`
	Balance float64   `db:"balance" json:"balance"`
}
