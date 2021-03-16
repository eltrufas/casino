-- migrate:up
CREATE VIEW balances AS
SELECT user_id, guild_id, sum(amount) AS balance
FROM transactions
GROUP BY user_id, guild_id;


-- migrate:down
DROP VIEW balances;
