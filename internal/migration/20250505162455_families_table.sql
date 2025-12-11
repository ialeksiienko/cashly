-- +goose Up
CREATE TABLE families (
    id SERIAL PRIMARY KEY,
    created_by INT NOT NULL,
    name VARCHAR(20) NOT NULL,
    created_at TIMESTAMPZ DEFAULT (NOW() AT TIME ZONE 'utc'::text),
        CONSTRAINT fk_created_by FOREIGN KEY (created_by)
            REFERENCES users(id),
        CONSTRAINT user_family_name_unq UNIQUE (created_by, name)
);

-- +goose Down
DROP TABLE families;
