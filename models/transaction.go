package models

import (
	"context"
)

type Transaction struct {
	ID      int64  `db:"id"`
	UserID  string `db:"user_id"`
	GuildID string `db:"guild_id"`
	Amount  int64  `db:"amount"`
	Note    string `db:"note"`
}

func (r *repository) CreateTransaction(ctx context.Context, txn *Transaction) error {
	res, err := r.db.NamedExecContext(ctx, `
		INSERT INTO transactions (user_id, guild_id, amount, note)
		VALUES(:user_id, :guild_id, :amount, :note)
	`, txn)
	if err != nil {
		return err
	}
	txn.ID, err = res.LastInsertId()
	return err
}
