package usecase

import (
	"context"

	"github.com/umekikazuya/logleaf/internal/domain"
	"github.com/umekikazuya/logleaf/internal/interface/repository"
)

type LeafUsecase struct {
	Repo repository.LeafRepository
}

func NewLeafUsecase(repo repository.LeafRepository) *LeafUsecase {
	return &LeafUsecase{Repo: repo}
}

func (u *LeafUsecase) ListLeaves(ctx context.Context, opts repository.ListOptions) ([]domain.Leaf, error) {
	return u.Repo.List(ctx, opts)
}

func (u *LeafUsecase) GetLeaf(ctx context.Context, id string) (*domain.Leaf, error) {
	return u.Repo.Get(ctx, id)
}

func (u *LeafUsecase) AddLeaf(ctx context.Context, leaf *domain.Leaf) error {
	return u.Repo.Put(ctx, leaf)
}

func (u *LeafUsecase) UpdateLeaf(ctx context.Context, id string, update *domain.Leaf) error {
	return u.Repo.Update(ctx, id, update)
}

func (u *LeafUsecase) DeleteLeaf(ctx context.Context, id string) error {
	return u.Repo.Delete(ctx, id)
}
