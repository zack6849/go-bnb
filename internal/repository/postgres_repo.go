package repository

import (
	"context"

	"github.com/google/uuid"
)

func Create[T any](ctx context.Context, record *T) error {
	return nil
}

func GetById[T any](ctx context.Context, id uuid.UUID) (*T, error) {
	return nil, nil
}

func FindManyById[T any](ctx context.Context, ids []uuid.UUID) {

}
