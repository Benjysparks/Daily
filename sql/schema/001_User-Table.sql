-- +goose Up
CREATE TABLE users (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    email TEXT NOT NULL UNIQUE,
    pword TEXT NOT NULL,
    full_name TEXT NOT NULL
);

-- +goose Down
DROP TABLE users;
