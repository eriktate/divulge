package pg

import (
	"context"
	"fmt"

	"github.com/eriktate/divulge"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

const insertUserQuery = `
INSERT INTO users (id, name, email)
VALUES (:id, :name, :email);
`

const updateUserQuery = `
UPDATE users
SET
	name = :name,
	email = :email
WHERE
	id = :id;
`

const removeUserQuery = `
UPDATE users
SET
	deleted_at = CURRENT_TIMESTAMP
WHERE
	id = $1;
`

const fetchUserQuery = `
SELECT *
FROM users
WHERE
	id = $1
	AND deleted_at IS NULL;
`

const listUsersQuery = `
SELECT *
FROM users
WHERE
	deleted_at IS NULL;
`

func (db DB) SaveUser(ctx context.Context, user divulge.User) (uuid.UUID, error) {
	// are we inserting?
	query := updateUserQuery
	if divulge.IsEmpty(user.ID) {
		user.ID = uuid.New()
		query = insertUserQuery
	}

	tx, err := db.db.BeginTxx(ctx, nil)
	if err != nil {
		return user.ID, fmt.Errorf("failed to create transaction: %w", err)
	}

	if _, err := sqlx.NamedExecContext(ctx, tx, query, &user); err != nil {
		tx.Rollback()
		return user.ID, fmt.Errorf("failed to execute query: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return user.ID, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return user.ID, nil
}

func (db DB) FetchUser(ctx context.Context, id uuid.UUID) (divulge.User, error) {
	var user divulge.User
	if err := db.db.GetContext(ctx, &user, fetchUserQuery, id); err != nil {
		return user, fmt.Errorf("failed to select: %w", err)
	}

	return user, nil
}

func (db DB) ListUsers(ctx context.Context) ([]divulge.User, error) {
	var users []divulge.User
	if err := db.db.SelectContext(ctx, &users, listUsersQuery); err != nil {
		return nil, fmt.Errorf("failed to select: %w", err)
	}

	return users, nil
}

func (db DB) RemoveUser(ctx context.Context, id uuid.UUID) error {
	tx, err := db.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}

	if _, err := tx.ExecContext(ctx, removeUserQuery, id); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to execute query: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
