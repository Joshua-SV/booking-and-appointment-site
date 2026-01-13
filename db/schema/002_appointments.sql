-- +goose Up
CREATE TABLE appointments(
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    appointment_time TIMESTAMPTZ NOT NULL UNIQUE,
    status TEXT NOT NULL CHECK (status IN ('scheduled', 'cancelled', 'completed')),
    notes TEXT NOT NULL
);

-- +goose Down
DROP TABLE appointments;