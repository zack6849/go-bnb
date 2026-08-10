package domain

import "github.com/google/uuid"

type Listing struct {
	UUIDIdentifier `gorm:"default:uuidv7()"`
	PublicID       uuid.UUID `gorm:"default:uuidv4()"`
	PropertyTypeID uuid.UUID `gorm:"default;"`
	PropertyType   PropertyType
	RoomTypeId     uuid.UUID `gorm:"default;"`
	RoomType       RoomType
	AirbnbID       int
	Accommodates   int
	NumBaths       float32 //baths can be half
	NumBeds        int
	ListingUrl     string
	Tagline        string
	Description    string
	MinNights      int
	MaxNights      int
	Location       Location
	Hosts          []Host    `gorm:"many2many:host_listing"`
	Amenities      []Amenity `gorm:"many2many:listing_amenities"`
}
