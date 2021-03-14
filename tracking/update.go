package tracking

type userKey struct {
	UserID  string
	GuildID string
}

type UserUpdate struct {
	UserID    string
	GuildID   string
	Connected bool
}

func (u UserUpdate) getKey() userKey {
	return userKey{
		UserID:  u.UserID,
		GuildID: u.GuildID,
	}
}
