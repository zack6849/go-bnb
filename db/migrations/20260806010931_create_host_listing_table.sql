-- +goose Up
create table host_listing
(
    listing_id uuid references listings (id),
    host_id uuid references hosts (id)
);

CREATE INDEX idx_host_listing_listing_id ON host_listing (listing_id);
CREATE INDEX idx_host_listing_host_id on host_listing (host_id);

-- +goose Down
drop table host_listing;