-- +goose Up
CREATE TABLE users (
    id INT PRIMARY KEY NOT NULL,
    username VARCHAR NOT NULL UNIQUE,
    firstname VARCHAR(50),
    joined_at TIMESTAMPTZ DEFAULT (NOW() AT TIME ZONE 'utc'::text)
);

-- +goose Down
DROP TABLE users;
