-- +goose Up
CREATE TABLE users
(
    id            uuid primary key unique default uuidv7(),
    email         text not null unique,
    password_hash text
);

-- +goose Down
DROP TABLE users;
