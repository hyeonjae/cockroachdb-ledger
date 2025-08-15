package db

import (
	"database/sql"
	"fmt"

	"mini-ledger/internal/config"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Database struct {
	*sqlx.DB
}

func New(cfg *config.Config) (*Database, error) {
	return NewDatabase(cfg.DatabaseURL)
}

func NewDatabase(databaseURL string) (*Database, error) {
	db, err := sqlx.Connect("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	database := &Database{DB: db}
	if err := database.runMigrations(); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return database, nil
}

func (db *Database) runMigrations() error {
	migrations := []string{
		`CREATE DATABASE IF NOT EXISTS mini_ledger`,
		`CREATE TABLE IF NOT EXISTS accounts (
		    id SERIAL PRIMARY KEY,
		    account_number STRING NOT NULL UNIQUE,
		    balance DECIMAL(15,2) NOT NULL,
		    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS holdings (
		    id SERIAL PRIMARY KEY,
		    account_id INT NOT NULL REFERENCES accounts(id),
		    stock_code STRING NOT NULL,
		    quantity INT NOT NULL,
		    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		    UNIQUE(account_id, stock_code)
		)`,
		`CREATE TABLE IF NOT EXISTS orders (
		    id SERIAL PRIMARY KEY,
		    account_id INT NOT NULL REFERENCES accounts(id),
		    stock_code STRING NOT NULL,
		    type STRING NOT NULL,
		    direction STRING NOT NULL,
		    quantity INT NOT NULL,
		    price DECIMAL(15,2),
		    filled_quantity INT NOT NULL DEFAULT 0,
		    status STRING NOT NULL,
		    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`INSERT INTO accounts (id, account_number, balance) VALUES (1, 'AC001', 1000000) ON CONFLICT (id) DO NOTHING`,
		`INSERT INTO holdings (account_id, stock_code, quantity) VALUES (1, 'STOCK01', 100) ON CONFLICT (account_id, stock_code) DO NOTHING`,
	}

	for i, migration := range migrations {
		if _, err := db.Exec(migration); err != nil {
			return fmt.Errorf("failed to execute migration %d: %w", i+1, err)
		}
	}

	return nil
}

func (db *Database) BeginTx() (*sqlx.Tx, error) {
	return db.Beginx()
}

type Querier interface {
	Get(dest interface{}, query string, args ...interface{}) error
	Select(dest interface{}, query string, args ...interface{}) error
	Exec(query string, args ...interface{}) (sql.Result, error)
	NamedExec(query string, arg interface{}) (sql.Result, error)
}