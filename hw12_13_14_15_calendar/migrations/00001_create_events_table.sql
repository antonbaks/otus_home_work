-- +goose Up
CREATE TABLE events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title text,
    description text,
    start_at TIMESTAMP,
    end_at TIMESTAMP,
    user_id INT,
    remind_for TIMESTAMP
);

-- +goose Down
DROP TABLE events;