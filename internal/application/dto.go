package application

import (
	"time"

	"github.com/umekikazuya/logleaf/internal/domain"
)

type LeafInputDTO struct {
	ID       string
	Title    string
	URL      string
	Platform string
	Tags     []string
}

type LeafOutputDTO struct {
	ID       string
	Title    string
	URL      string
	Platform string
	Read     bool
	Tags     []string
	SyncedAt string
}

func LeafInputDTOToDomain(dto *LeafInputDTO) *domain.Leaf {
	leaf := domain.NewLeaf(dto.ID, dto.Title, dto.URL, dto.Platform)
	if dto.Tags != nil {
		leaf.UpdateTags(dto.Tags)
	}
	return leaf
}

func LeafDomainToOutputDTO(leaf *domain.Leaf) *LeafOutputDTO {
	return &LeafOutputDTO{
		ID:       leaf.ID,
		Title:    leaf.Title,
		URL:      leaf.URL,
		Platform: leaf.Platform,
		Read:     leaf.Read,
		Tags:     leaf.Tags,
		SyncedAt: leaf.SyncedAt.Format(time.RFC3339),
	}
}
