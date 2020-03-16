package pg

import (
	"context"
	"fmt"

	"github.com/eriktate/divulge"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

const insertPostQuery = `
INSERT INTO posts
	(id, author_id, account_id, title, summary)
VALUES
	(:id, :author_id, :account_id, :title, :summary);
`

const updatePostQuery = `
UPDATE posts
SET
	title = :title,
	summary = :summary
WHERE
	id = :id
	AND deleted_at IS NULL;
`

const publishPostQuery = `
UPDATE posts
SET
	published_at = CURRENT_TIMESTAMP
WHERE
	id = $1;
`

const redactPostQuery = `
UPDATE posts
SET
	published_at = NULL
WHERE
	id = $1;
`

const fetchPostQuery = `
SELECT *
FROM posts
WHERE
	id = $1;
`

const listPostsByAccountQuery = `
SELECT *
FROM posts
WHERE
	account_id = $1;
`

const removePostQuery = `
DELETE FROM posts
WHERE id = $1;
`

func (db DB) SavePost(ctx context.Context, post divulge.Post) (uuid.UUID, error) {
	query := updatePostQuery
	if divulge.IsEmpty(post.ID) {
		post.ID = uuid.New()
		query = insertPostQuery
	}

	tx, err := db.db.BeginTxx(ctx, nil)
	if err != nil {
		return post.ID, fmt.Errorf("failed to create transaction: %w", err)
	}

	if _, err := sqlx.NamedExecContext(ctx, tx, query, &post); err != nil {
		tx.Rollback()
		return post.ID, fmt.Errorf("failed to execute query: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return post.ID, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return post.ID, nil
}

func (db DB) PublishPost(ctx context.Context, id uuid.UUID) error {
	tx, err := db.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}

	if _, err := tx.ExecContext(ctx, publishPostQuery, id); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to execute query: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (db DB) RedactPost(ctx context.Context, id uuid.UUID) error {
	tx, err := db.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}

	if _, err := tx.ExecContext(ctx, redactPostQuery, id); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to execute query: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (db DB) FetchPost(ctx context.Context, id uuid.UUID) (divulge.Post, error) {
	var post divulge.Post
	if err := db.db.GetContext(ctx, &post, fetchUserQuery, id); err != nil {
		return post, fmt.Errorf("failed to select: %w", err)
	}

	return post, nil
}

func (db DB) ListPostsByAccount(ctx context.Context, accountID uuid.UUID) ([]divulge.Post, error) {
	var posts []divulge.Post
	if err := db.db.SelectContext(ctx, &posts, listPostsByAccountQuery, accountID); err != nil {
		return posts, fmt.Errorf("failed to select: %w", err)
	}

	return posts, nil
}

func (db DB) RemovePost(ctx context.Context, id uuid.UUID) error {
	if _, err := db.db.ExecContext(ctx, removePostQuery, id); err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}

	return nil
}
