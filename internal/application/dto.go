package application

import (
	"time"

	"github.com/umekikazuya/logleaf/internal/domain"
)

type LeafInputDTO struct {
	ID       string
	Note     string
	URL      string
	Platform string
	Tags     []string
}

type LeafOutputDTO struct {
	ID       string
	Note     string
	URL      string
	Platform string
	Read     bool
	Tags     []string
	SyncedAt string
}

func LeafDomainToOutputDTO(leaf *domain.Leaf) *LeafOutputDTO {
	return &LeafOutputDTO{
		ID:       leaf.ID().String(),
		Note:     leaf.Note(),
		URL:      leaf.URL().String(),
		Platform: leaf.Platform(),
		Tags:     leaf.Tags(),
		SyncedAt: leaf.SyncedAt().Format(time.RFC3339),
	}
}
