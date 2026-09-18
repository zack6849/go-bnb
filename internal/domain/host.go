package domain

import (
	"time"
)

type Host struct {
	UUIDIdentifier
	Name        string
	Description string
	SuperHost   bool `gorm:"column:superhost"`
	Slug        string
	ProfileID   int64
	HostSince   time.Time
	Location    string
	PictureURL  string `gorm:"column:picture_url"`
}
