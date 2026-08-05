CREATE TABLE listings
(
    id          uuid unique primary key         default uuidv7(),
    airbnb_id   bigint                 null     default null,
    num_beds    int                    not null default 0,
    num_baths   float                  not null default 0,
    listing_url text                   not null,
    tagline     text                   not null,
    description text                   not null,
    min_nights  int                    not null default 1,
    max_nights  int                    not null default 7,
    location    geography(point, 4326) not null
);

CREATE INDEX idx_listing_id ON listings (id);
-- Create a spatial index on the geography column
CREATE INDEX idx_listing_location
    ON listings USING GIST (location);

CREATE INDEX idx_min_nights on listings (min_nights);
CREATE INDEX idx_max_nights on listings (max_nights);