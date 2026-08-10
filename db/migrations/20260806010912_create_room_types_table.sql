-- +goose Up
CREATE TABLE room_types
(
    id   uuid not null primary key unique default uuidv7(),
    name text not null unique
);

insert into room_types (name)
values
    ('Private room'),
    ('Entire Home / Apt'),
    ('Shared Room'),
    ('Hotel Room') ON CONFLICT DO NOTHING ;

-- +goose Down
DROP TABLE room_types;
