package repository

import (
	"context"

	"github.com/umekikazuya/logleaf/internal/domain"
)

type LeafRepository interface {
	Get(ctx context.Context, id string) (*domain.Leaf, error)
	List(ctx context.Context) ([]domain.Leaf, error)
	Put(ctx context.Context, leaf *domain.Leaf) error
	Update(ctx context.Context, id string, update *domain.Leaf) error
	Delete(ctx context.Context, id string) error
}
