package domain

import (
	"github.com/google/uuid"
)

type DatabaseDrivenList struct {
	UUIDIdentifier `json:"id"`
	Name           string `json:"name"`
}

type UUIDIdentifier struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey;default:uuidv7()"`
}
