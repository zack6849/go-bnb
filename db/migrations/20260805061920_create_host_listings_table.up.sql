CREATE TABLE host_listings
(
    host_id    uuid references hosts (id),
    listing_id uuid references listings (id)
)