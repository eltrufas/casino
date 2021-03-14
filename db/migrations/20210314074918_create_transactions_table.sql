-- migrate:up
CREATE TABLE transactions (
  id INTEGER PRIMARY KEY,
  user_id VARCHAR(256) NOT NULL,
  guild_id VARCHAR(256) NOT NULL,
  amount INTEGER NOT NULL,
  note VARCHAR(512)
);

CREATE INDEX transactions_user_id_guild_id
ON transactions(user_id, guild_id);

-- migrate:down

DROP INDEX transactions_user_id_guild_id;

DROP TABLE transactions;
