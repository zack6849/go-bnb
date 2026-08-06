-- +goose Up
CREATE TABLE hosts
(
    id          uuid primary key unique default uuidv7(),
    user_id     uuid references users(id),
    name        text,
    description text,
    superhost   boolean                 default false,
    slug        text unique,
    profile_id  int unique,
    host_since  timestamp,
    location    text
);

CREATE INDEX hosts_user_id_idx ON hosts(user_id);

-- +goose Down
DROP TABLE hosts;
