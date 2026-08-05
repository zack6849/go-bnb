create table users
(
    id            uuid primary key unique,
    email         text unique,
    password_hash text
)