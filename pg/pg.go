package pg

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// A DB implements divulge services with a postgres backend.
type DB struct {
	db *sqlx.DB
}

// New creates a new pg.DB capable of working with divulge data.
func New(host, user, password string) (DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=divulge sslmode=disable", host, user, password)
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return DB{}, fmt.Errorf("failed to connect to database: %w", err)
	}

	return DB{
		db: db,
	}, nil
}
