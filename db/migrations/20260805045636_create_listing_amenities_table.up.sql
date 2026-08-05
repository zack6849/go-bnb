create table listing_amenities
(
    listing_id uuid references listings (id),
    amenity_id uuid references amenities (id)
);
