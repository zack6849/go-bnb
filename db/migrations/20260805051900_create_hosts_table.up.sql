CREATE TABLE hosts
(
    id          uuid primary key unique,
    name        text,
    description text,
    superhost   boolean default false,
    profile_url text,
    profile_id  int unique,
    host_since  timestamp,
    location    text
)