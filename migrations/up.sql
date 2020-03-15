CREATE TABLE IF NOT EXISTS users(
	id UUID PRIMARY KEY,
	name VARCHAR(512) NOT NULL,
	email VARCHAR(320) NOT NULL UNIQUE,
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	deleted_at TIMESTAMP DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS accounts(
	id UUID PRIMARY KEY,
	name VARCHAR(512) NOT NULL,
	owner_id UUID NOT NULL REFERENCES users(id),
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	deleted_at TIMESTAMP DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS user_accounts(
	user_id UUID NOT NULL REFERENCES users(id),
	account_id UUID NOT NULL REFERENCES accounts(id)
);

CREATE TABLE IF NOT EXISTS posts(
	id UUID PRIMARY KEY,
	author_id UUID NOT NULL REFERENCES users(id),
	account_id UUID NOT NULL REFERENCES accounts(id),
	title VARCHAR(256) NOT NULL,
	summary VARCHAR(512),
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	published_at TIMESTAMP DEFAULT NULL
);