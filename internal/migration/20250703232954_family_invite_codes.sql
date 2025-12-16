-- +goose Up
CREATE TABLE family_invite_codes (
    id SERIAL PRIMARY KEY,
    family_id INT NOT NULL,
    code VARCHAR(6) NOT NULL UNIQUE,
    created_by INT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT (NOW() AT TIME ZONE 'utc'::text),
    expires_at TIMESTAMPTZ NOT NULL,
    CONSTRAINT fk_family_id FOREIGN KEY (family_id)
        REFERENCES families(id),
    CONSTRAINT fk_created_by FOREIGN KEY (created_by)
        REFERENCES users(id)
);

-- +goose Down
DROP TABLE family_invite_codes;