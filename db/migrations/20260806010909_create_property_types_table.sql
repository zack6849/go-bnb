-- +goose Up
CREATE TABLE property_types
(
    id   uuid primary key default uuidv7(),
    name text not null unique
);

insert into property_types (name)
values ('Entire rental unit'),
       ('Private room in home'),
       ('Entire home'),
       ('Private room in rental unit')
ON CONFLICT DO NOTHING;

-- +goose Down
DROP TABLE property_types;
