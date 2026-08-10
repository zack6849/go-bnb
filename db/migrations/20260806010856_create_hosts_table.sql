-- +goose Up
CREATE TABLE hosts
(
    id          uuid not null primary key unique default uuidv7(),
    user_id     uuid references users(id),
    name        text not null,
    description text not null,
    superhost   boolean not null default false,
    slug        text not null unique,
    profile_id  bigint not null unique,
    host_since  timestamp not null,
    location    text not null
);

CREATE INDEX hosts_user_id_idx ON hosts(user_id);

-- +goose Down
DROP TABLE hosts;
