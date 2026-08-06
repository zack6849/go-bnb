-- +goose Up
CREATE TABLE listings
(
    id               uuid unique primary key         default uuidv7(),
    -- Hide the main id since it exposes creation time
    -- just create another public uuid column with a v4 ID that doesn't have timestamp data
    public_id        uuid unique                     default uuidv4(),
    property_type_id uuid references property_types (id),
    room_type_id     uuid references room_types (id),
    airbnb_id        bigint                 null     default null,
    num_beds         int                    not null default 0,
    num_baths        float                  not null default 0,
    accommodates     int                    not null default 0,
    listing_url      text                   not null,
    tagline          text                   not null,
    description      text                   not null,
    min_nights       int                    not null default 1,
    max_nights       int                    not null default 7,
    location         geography(point, 4326) not null
);

-- Create indexes for columns I imagine will be filtered often
CREATE INDEX idx_listing_property_type_id ON listings(property_type_id);
CREATE INDEX idx_listing_room_type_id ON listings(room_type_id);
CREATE INDEX idx_listing_min_nights ON listings(min_nights);
CREATE INDEX idx_listing_max_nights ON listings(max_nights);
CREATE INDEX idx_listing_baths ON listings(num_baths);
CREATE INDEX idx_listing_beds ON listings(num_beds);
CREATE INDEX idx_listing_accommodates ON listings(accommodates);

-- Create a spatial index on the geography column
CREATE INDEX idx_listing_location
    ON listings USING GIST (location);


-- +goose Down
DROP TABLE listings;
