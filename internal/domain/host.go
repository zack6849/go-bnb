package domain

import (
	"time"
)

type Host struct {
	UUIDIdentifier
	Name        string
	Description string
	Superhost   bool
	Slug        string
	ProfileID   int64
	HostSince   time.Time
	Location    string
}
