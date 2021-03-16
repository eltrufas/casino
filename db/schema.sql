CREATE TABLE IF NOT EXISTS "schema_migrations" (version varchar(255) primary key);
CREATE TABLE transactions (
  id INTEGER PRIMARY KEY,
  user_id VARCHAR(256) NOT NULL,
  guild_id VARCHAR(256) NOT NULL,
  amount INTEGER NOT NULL,
  note VARCHAR(512)
);
CREATE INDEX transactions_user_id_guild_id
ON transactions(user_id, guild_id);
CREATE VIEW balances AS
SELECT user_id, guild_id, sum(amount) AS balance
FROM transactions
GROUP BY user_id, guild_id
/* balances(user_id,guild_id,balance) */;
-- Dbmate schema migrations
INSERT INTO "schema_migrations" (version) VALUES
  ('20210314074918'),
  ('20210314104217');
