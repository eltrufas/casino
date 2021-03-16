package models

import (
	"context"
	"github.com/jmoiron/sqlx"
)

type Repository interface {
	CreateTransaction(ctx context.Context, txn *Transaction) error
	GetUserBalance(ctx context.Context, guildID, userID string) (*Balance, error)
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{
		db: db,
	}
}
