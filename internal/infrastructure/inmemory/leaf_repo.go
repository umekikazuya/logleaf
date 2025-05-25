package repository

import (
	"context"

	"github.com/umekikazuya/logleaf/internal/domain"
	"github.com/umekikazuya/logleaf/internal/interface/repository"
)

type InMemoryLeafRepository struct {
	leaves map[string]*domain.Leaf
}

func NewInMemoryLeafRepository() *InMemoryLeafRepository {
	return &InMemoryLeafRepository{
		leaves: make(map[string]*domain.Leaf),
	}
}

func NewInMemoryLeafRepositoryWithSample() *InMemoryLeafRepository {
	repo := NewInMemoryLeafRepository()
	repo.leaves["1"] = &domain.Leaf{
		ID:       "1",
		Title:    "サンプル記事1",
		URL:      "https://example.com/1",
		Platform: "Qiita",
		Tags:     []string{"Go", "サンプル"},
		Read:     false,
	}
	repo.leaves["2"] = &domain.Leaf{
		ID:       "2",
		Title:    "サンプル記事2",
		URL:      "https://example.com/2",
		Platform: "Zenn",
		Tags:     []string{"AWS", "DynamoDB"},
		Read:     true,
	}
	return repo
}

func (r *InMemoryLeafRepository) Get(ctx context.Context, id string) (*domain.Leaf, error) {
	leaf, ok := r.leaves[id]
	if !ok {
		return nil, nil
	}
	return leaf, nil
}

func (r *InMemoryLeafRepository) List(ctx context.Context, opts repository.ListOptions) ([]domain.Leaf, error) {
	result := make([]domain.Leaf, 0, len(r.leaves))
	for _, leaf := range r.leaves {
		result = append(result, *leaf)
	}
	return result, nil
}

func (r *InMemoryLeafRepository) Put(ctx context.Context, leaf *domain.Leaf) error {
	r.leaves[leaf.ID] = leaf
	return nil
}

func (r *InMemoryLeafRepository) Update(ctx context.Context, id string, update *domain.Leaf) error {
	if _, ok := r.leaves[id]; !ok {
		return nil
	}
	r.leaves[id] = update
	return nil
}

func (r *InMemoryLeafRepository) Delete(ctx context.Context, id string) error {
	delete(r.leaves, id)
	return nil
}
