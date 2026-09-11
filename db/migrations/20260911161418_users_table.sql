-- +goose Up
CREATE TABLE users (
    id VARCHAR(36) PRIMARY KEY NOT NULL,
    username VARCHAR(25) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL
);

-- +goose Down
DROP TABLE users;
