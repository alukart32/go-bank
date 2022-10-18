-- Account
-- name: CreateAccount :one
INSERT INTO accounts (
    id,
    owner,
    balance,
    currency
) VALUES (
  $1, $2, $3, $4
) RETURNING *;

-- name: GetAccount :one
SELECT * FROM accounts
WHERE id = $1;

-- name: UpdateAccountOwner :one
UPDATE accounts
SET owner = $2
WHERE id = $1
RETURNING *;

-- name: AddAccountBalance :one
UPDATE accounts
SET balance = balance + sqlc.arg(amount)
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: ListAccounts :many
SELECT A.id, A.owner, A.balance, A.currency, A.created_at FROM accounts as A
JOIN (
    SELECT id FROM accounts
    LIMIT $1
    OFFSET $2
  ) as P
  ON P.id = A.id;

-- name: DeleteAccount :exec
DELETE FROM accounts
WHERE id = $1;


-- Entry
-- name: CreateEntry :one
INSERT INTO entries (
  account_id,
  amount
) VALUES (
  $1, $2
) RETURNING *;

-- name: GetEntry :one
SELECT * FROM entries
WHERE id = $1;

-- name: ListEntriesByAccount :many
SELECT E.id, E.account_id, E.amount, E.created_at FROM entries as E
JOIN (
    SELECT id FROM entries as je
    WHERE je.account_id = $1
    LIMIT $2
    OFFSET $3
  ) as P
  ON P.id = E.id;

-- name: UpdateEntry :exec
UPDATE entries
SET amount = $2
WHERE id = $1;

-- name: DeleteEntry :exec
DELETE FROM entries
WHERE id = $1;


-- Transfer
-- name: CreateTransfer :one
INSERT INTO transfers (
  from_account_id,
  to_account_id,
  from_entry_id,
  to_entry_id,
  amount
) VALUES (
  $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetTransfer :one
SELECT * FROM transfers
WHERE id = $1;

-- name: ListTransfersByFromAccount :many
SELECT T.id, T.from_account_id, T.to_account_id, T.amount,
T.from_entry_id, T.to_entry_id, T.created_at FROM transfers AS T
JOIN (
    SELECT id FROM transfers as jt
    WHERE jt.from_account_id = $1
    LIMIT $2
    OFFSET $3
  ) as P
  ON P.id = T.id;

-- name: ListTransfersByToAccount :many
SELECT T.id, T.from_account_id, T.to_account_id, T.amount,
T.from_entry_id, T.to_entry_id, T.created_at FROM transfers AS T
JOIN (
    SELECT id FROM transfers as jt
    WHERE jt.to_account_id = $1
    LIMIT $2
    OFFSET $3
  ) as P
  ON P.id = T.id;

-- name: ListTransfersByAccounts :many
SELECT T.id, T.from_account_id, T.to_account_id, T.amount,
T.from_entry_id, T.to_entry_id, T.created_at FROM transfers AS T
JOIN (
    SELECT id FROM transfers as jt
    WHERE jt.to_account_id = $1 AND jt.from_account_id = $2
    LIMIT $3
    OFFSET $4
  ) as P
  ON P.id = T.id;

-- name: DeleteTransfer :exec
DELETE FROM transfers
WHERE id = $1;