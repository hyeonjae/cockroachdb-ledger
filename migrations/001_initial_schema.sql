-- CockroachDB compatible schema
CREATE DATABASE IF NOT EXISTS mini_ledger;

-- accounts 테이블
CREATE TABLE IF NOT EXISTS accounts (
    id SERIAL PRIMARY KEY,
    account_number STRING NOT NULL UNIQUE,
    balance DECIMAL(15,2) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- holdings 테이블 (보유 주식)
CREATE TABLE IF NOT EXISTS holdings (
    id SERIAL PRIMARY KEY,
    account_id INT NOT NULL REFERENCES accounts(id),
    stock_code STRING NOT NULL,
    quantity INT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(account_id, stock_code)
);

-- orders 테이블 (주문)
CREATE TABLE IF NOT EXISTS orders (
    id SERIAL PRIMARY KEY,
    account_id INT NOT NULL REFERENCES accounts(id),
    stock_code STRING NOT NULL,
    type STRING NOT NULL,           -- 'LIMIT'
    direction STRING NOT NULL,      -- 'BUY' or 'SELL'
    quantity INT NOT NULL,
    price DECIMAL(15,2),
    filled_quantity INT NOT NULL DEFAULT 0,
    status STRING NOT NULL,         -- 'PENDING', 'PARTIAL', 'FILLED', 'CANCELED'
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);