-- +goose Up
ALTER TABLE users
ADD COLUMN phone_number TEXT UNIQUE NOT NULL;

-- +goose Down
ALTER TABLE users
DROP COLUMN phone_number;