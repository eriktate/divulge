package pg

import (
	"context"
	"fmt"

	"github.com/eriktate/divulge"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

const insertAccountQuery = `
INSERT INTO accounts
	(id, name, owner_id)
VALUES
	(:id, :name, :owner_id);
`

const updateAccountQuery = `
UPDATE accounts
SET
	name = :name,
	owner_id = :owner_id
WHERE id = :id;
`

const fetchAccountQuery = `
SELECT *
FROM accounts
WHERE
	id = $1;
`

const listAccountsQuery = `
SELECT *
FROM accounts;
`

const removeAccountQuery = `
UPDATE accounts
SET
	deleted_at = CURRENT_TIMESTAMP
WHERE
	id = $1
	AND deleted_at IS NULL;
`

func (db DB) SaveAccount(ctx context.Context, account divulge.Account) (uuid.UUID, error) {
	// are we inserting?
	query := updateAccountQuery
	if divulge.IsEmpty(account.ID) {
		account.ID = uuid.New()
		query = insertAccountQuery
	}

	tx, err := db.db.BeginTxx(ctx, nil)
	if err != nil {
		return account.ID, fmt.Errorf("failed to create transaction: %w", err)
	}

	if _, err := sqlx.NamedExecContext(ctx, tx, query, &account); err != nil {
		tx.Rollback()
		return account.ID, fmt.Errorf("failed to execute query: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return account.ID, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return account.ID, nil
}

func (db DB) FetchAccount(ctx context.Context, id uuid.UUID) (divulge.Account, error) {
	var account divulge.Account
	if err := db.db.GetContext(ctx, &account, fetchAccountQuery, id); err != nil {
		return account, fmt.Errorf("failed to select: %w", err)
	}

	return account, nil
}

func (db DB) ListAccounts(ctx context.Context) ([]divulge.Account, error) {
	var accounts []divulge.Account
	if err := db.db.SelectContext(ctx, &accounts, listAccountsQuery); err != nil {
		return nil, fmt.Errorf("failed to select: %w", err)
	}

	return accounts, nil
}

func (db DB) RemoveAccount(ctx context.Context, id uuid.UUID) error {
	tx, err := db.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}

	if _, err := tx.ExecContext(ctx, removeAccountQuery, id); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to execute query: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
