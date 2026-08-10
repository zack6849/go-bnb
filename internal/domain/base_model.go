package domain

import "github.com/google/uuid"

type DatabaseDrivenList struct {
	UUIDIdentifier
	Name string
}

type UUIDIdentifier struct {
	ID uuid.UUID `gorm:"primaryKey;default:uuidv7()"`
}
