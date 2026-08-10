-- +goose Up
create table listing_amenities
(
    listing_id uuid not null references listings (id),
    amenity_id uuid not null references amenities (id),
    constraint listing_amenities_pkey primary key (listing_id, amenity_id)
);

CREATE INDEX idx_listing_amenities_amenity ON listing_amenities(amenity_id);

-- +goose Down
DROP TABLE listing_amenities;
