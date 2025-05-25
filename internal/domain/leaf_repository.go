package domain

import (
	"context"
)

type ListOptions struct {
	Platforms []string
	Tags      []string
	ReadOnly  bool
	Limit     int
	Offset    int
	SortBy    string
	SortDesc  bool
}

type LeafRepository interface {
	Get(ctx context.Context, id string) (*Leaf, error)
	List(ctx context.Context, opts ListOptions) ([]Leaf, error)
	Put(ctx context.Context, leaf *Leaf) error
	Update(ctx context.Context, id string, update *Leaf) error
	Delete(ctx context.Context, id string) error
}
