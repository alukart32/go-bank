-- Account
-- name: CreateAccount :exec
INSERT INTO accounts (
    id,
    owner,
    balance,
    currency
) VALUES (
  $1, $2, $3, $4
);

-- name: GetAccount :one
SELECT * FROM accounts
WHERE id = $1;

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

-- name: UpdateAccount :exec
UPDATE accounts
SET balance = $2
WHERE id = $1;

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
  amount
) VALUES (
  $1, $2, $3
) RETURNING *;

-- name: GetTransfer :one
SELECT * FROM transfers
WHERE id = $1;

-- name: ListTransfersByFromAccount :many
SELECT T.id, T.from_account_id, T.to_account_id,
T.amount, T.created_at FROM transfers AS T
JOIN (
    SELECT id FROM transfers as jt
    WHERE jt.from_account_id = $1
    LIMIT $2
    OFFSET $3
  ) as P
  ON P.id = T.id;

-- name: ListTransfersByToAccount :many
SELECT T.id, T.from_account_id, T.to_account_id,
T.amount, T.created_at FROM transfers AS T
JOIN (
    SELECT id FROM transfers as jt
    WHERE jt.to_account_id = $1
    LIMIT $2
    OFFSET $3
  ) as P
  ON P.id = T.id;

-- name: UpdateTransfer :exec
UPDATE transfers
SET amount = $2
WHERE id = $1;

-- name: DeleteTransfer :exec
DELETE FROM transfers
WHERE id = $1;