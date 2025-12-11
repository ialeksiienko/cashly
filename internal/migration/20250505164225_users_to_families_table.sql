-- +goose Up
CREATE TABLE users_to_families (
     id SERIAL PRIMARY KEY,
     user_id INT NOT NULL,
     family_id INT NOT NULL,
     CONSTRAINT fk_user_id FOREIGN KEY (user_id)
        REFERENCES users(id),
    CONSTRAINT fk_family_id FOREIGN KEY (family_id)
        REFERENCES families(id),
    CONSTRAINT user_family_unq UNIQUE (user_id, family_id)
);
-- +goose Down
DROP TABLE users_to_families;
