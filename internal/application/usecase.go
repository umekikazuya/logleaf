package application

import (
	"context"
	"errors"

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
	leaf, err := domain.NewLeaf("", dto.Note, dto.URL, dto.Platform, dto.Tags)
	if err != nil {
		return nil, err
	}
	return u.repo.Put(ctx, leaf)
}

func (u *LeafUsecase) UpdateLeaf(ctx context.Context, update *LeafInputDTO) error {
	// 既存Leaf取得
	leaf, err := u.repo.Get(ctx, update.ID)
	if err != nil {
		return err
	}
	if leaf == nil {
		return errors.New("更新対象のLeafが見つかりません")
	}
	// Noteの更新
	if err := leaf.UpdateNote(update.Note); err != nil {
		return err
	}
	if err := leaf.UpdateNote(update.Note); err != nil {
		return err
	}
	// Platformの更新
	if err := leaf.UpdatePlatform(update.Platform); err != nil {
		return err
	}
	// タグ変換
	tags := make([]domain.Tag, 0, len(update.Tags))
	for _, v := range update.Tags {
		t, err := domain.NewTag(v)
		if err != nil {
			return err
		}
		tags = append(tags, t)
	}
	if err := leaf.UpdateTags(tags); err != nil {
		return err
	}
	return u.repo.Update(ctx, leaf)
}

func (u *LeafUsecase) ReadLeaf(ctx context.Context, id string) error {
	leaf, err := u.repo.Get(ctx, id)
	if err != nil {
		return err
	}
	if leaf == nil {
		return errors.New("leaf not found")
	}
	leaf.MarkAsRead()
	return u.repo.Update(ctx, leaf)
}

func (u *LeafUsecase) DeleteLeaf(ctx context.Context, id string) error {
	return u.repo.Delete(ctx, id)
}
