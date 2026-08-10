-- +goose Up
CREATE TABLE amenities
(
    id   uuid not null primary key unique default uuidv7(),
    name text unique
);

INSERT INTO amenities (name)
values ('Microwave'),
       ('Stove'),
       ('WiFi'),
       ('Washer / Dryer'),
       ('Pool'),
       ('Hot Tub'),
       ('Grill'),
       ('Outdoor Seating'),
       ('Self Check-In'),
       ('Full Kitchen') ON CONFLICT DO NOTHING ;

-- +goose Down
DROP TABLE amenities;
