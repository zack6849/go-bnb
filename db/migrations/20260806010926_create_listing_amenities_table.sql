-- +goose Up
create table listing_amenities
(
    listing_id uuid references listings (id),
    amenity_id uuid references amenities (id)
);

CREATE INDEX listing_amenities_listing ON listing_amenities(listing_id);
CREATE INDEX listing_amenities_amenity ON listing_amenities(amenity_id);

-- +goose Down
DROP TABLE listing_amenities;
