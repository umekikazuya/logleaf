package usecase

import (
	"context"
	"errors"

	"github.com/umekikazuya/logleaf/internal/domain"
	"github.com/umekikazuya/logleaf/internal/interface/repository"
)

// LeafUsecase provides application-level operations for managing Leaf entities.
// It interacts with the LeafRepository to perform CRUD operations and other business logic.
type LeafUsecase struct {
	repo repository.LeafRepository
}

func NewLeafUsecase(repo repository.LeafRepository) *LeafUsecase {
	return &LeafUsecase{repo: repo}
}

func (u *LeafUsecase) ListLeaves(ctx context.Context, opts repository.ListOptions) ([]domain.Leaf, error) {
	return u.repo.List(ctx, opts)
}

func (u *LeafUsecase) GetLeaf(ctx context.Context, id string) (*domain.Leaf, error) {
	if id == "" {
		return nil, errors.New("leaf ID cannot be empty")
	}
	return u.repo.Get(ctx, id)
}

func (u *LeafUsecase) AddLeaf(ctx context.Context, leaf *domain.Leaf) error {
	return u.repo.Put(ctx, leaf)
}

func (u *LeafUsecase) UpdateLeaf(ctx context.Context, id string, update *domain.Leaf) error {
	return u.repo.Update(ctx, id, update)
}

func (u *LeafUsecase) DeleteLeaf(ctx context.Context, id string) error {
	return u.repo.Delete(ctx, id)
}
