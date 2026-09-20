-- +goose Up
CREATE TABLE cities (
    id            uuid primary key unique default uuidv7(),
    name text,
    location         geography(point, 4326) not null,
    -- https://en.wikipedia.org/wiki/ISO_3166-2 representation of where the city is
    subdivision text,
    full_name text,
    timezone text,
    hash text unique
);

-- +goose Down
DROP TABLE cities;
