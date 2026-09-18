-- +goose Up
create table geocoding_cache
(
    normalized_search_query text,
    location         geography(point, 4326) not null,
    place_id         text

);

-- +goose Down
drop table geocoding_cache;
