package domain

import (
	"github.com/google/uuid"
)

type Listing struct {
	UUIDIdentifier
	PublicID       uuid.UUID `gorm:"type:uuid;default:uuidv4()"`
	PropertyTypeID uuid.UUID
	PropertyType   PropertyType
	RoomTypeID     uuid.UUID
	RoomType       RoomType
	AirbnbID       int
	Accommodates   int
	NumBaths       float32 //baths can be half
	NumBeds        int
	ListingURL     string
	PictureURL     string
	Tagline        string
	Description    string
	MinNights      int
	MaxNights      int
	Location       Location
	DistanceMeters float64   `gorm:"->;-:"`
	Hosts          []Host    `gorm:"many2many:host_listing"`
	Amenities      []Amenity `gorm:"many2many:listing_amenities"`
}
