package main

import (
	"flag"
	"fmt"
	"io/ioutil"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

const createDatabaseQuery = "CREATE DATABASE divulge;"

func main() {
	// enable migrating down
	var down bool
	flag.BoolVar(&down, "down", false, "migrate down")
	flag.Parse()

	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{})

	logger.WithField("down", down).Info("migrate option")
	if down {
		migrateDown(logger)
		return
	}

	migrateUp(logger)
}

func migrateUp(logger *logrus.Logger) {
	logger.Info("starting migration...")
	up, err := ioutil.ReadFile("./migrations/up.sql")
	if err != nil {
		logger.WithError(err).Fatal("failed to read migration file")
	}

	// attempt to connect to database
	db, err := connect()
	if err != nil {
		logger.WithError(err).Error("failed to connect to database, attempting to create")
		if err := createDatabase(); err != nil {
			logger.WithError(err).Fatal("failed to create database")
		}

		db, err = connect()
		if err != nil {
			logger.WithError(err).Fatal("failed to connect to database")
		}
	}

	// create schema
	if _, err := db.Exec(string(up)); err != nil {
		logger.WithError(err).Fatal("failed to create schema")
	}

	logger.Info("successfully applied database schema")
}

func migrateDown(logger *logrus.Logger) {
	logger.Info("beginning teardown")
	down, err := ioutil.ReadFile("./migrations/down.sql")
	if err != nil {
		logger.WithError(err).Fatal("failed to read migration file")
	}

	db, err := connect()
	if err != nil {
		logger.WithError(err).Fatal("failed to connect to database")
	}

	if _, err := db.Exec(string(down)); err != nil {
		logger.WithError(err).Fatal("failed to execute migration")
	}
}

func connect() (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", "user=postgres password=password dbname=divulge sslmode=disable")
	if err != nil {
		return nil, err
	}

	return db, nil
}

func createDatabase() error {
	db, err := sqlx.Connect("postgres", "user=postgres password=password sslmode=disable")
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	if _, err := db.Exec(createDatabaseQuery); err != nil {
		return err
	}

	return err
}
