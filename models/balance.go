package models

import (
	"context"
)

type Balance struct {
	UserID  string `db:"user_id"`
	GuildID string `db:"guild_id"`
	Balance int64  `db:"balance"`
}

func (r *repository) GetUserBalance(ctx context.Context, guildID, userID string) (*Balance, error) {
	balance := &Balance{}
	err := r.db.QueryRowxContext(ctx, `
		SELECT user_id, guild_id, balance FROM balances WHERE guild_id=? AND user_id=?
	`, guildID, userID).StructScan(balance)

	return balance, err
}
