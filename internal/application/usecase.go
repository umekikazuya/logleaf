package application

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/umekikazuya/logleaf/internal/domain"
)

// LeafUsecase provides application-level operations for managing Leaf entities.
// It interacts with the LeafRepository to perform CRUD operations and other business logic.
type LeafUsecase struct {
	repo domain.LeafRepository
}

func NewLeafUsecase(repo domain.LeafRepository) *LeafUsecase {
	return &LeafUsecase{repo: repo}
}

func (u *LeafUsecase) ListLeaves(ctx context.Context, opts domain.ListOptions) ([]domain.Leaf, error) {
	return u.repo.List(ctx, opts)
}

func (u *LeafUsecase) GetLeaf(ctx context.Context, id string) (*domain.Leaf, error) {
	if id == "" {
		return nil, errors.New("leaf ID cannot be empty")
	}
	return u.repo.Get(ctx, id)
}

func (u *LeafUsecase) AddLeaf(ctx context.Context, dto *LeafInputDTO) (*domain.Leaf, error) {
	leaf := domain.Leaf{
		ID:       uuid.New().String(),
		Title:    dto.Title,
		URL:      dto.URL,
		Platform: dto.Platform,
		Tags:     dto.Tags,
	}
	return u.repo.Put(ctx, &leaf)
}

func (u *LeafUsecase) UpdateLeaf(ctx context.Context, update *LeafInputDTO) error {
	leaf := LeafInputDTOToDomain(update)
	return u.repo.Update(ctx, leaf)
}

func (u *LeafUsecase) DeleteLeaf(ctx context.Context, id string) error {
	return u.repo.Delete(ctx, id)
}
