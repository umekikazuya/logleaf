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
	Tags     []string
	SyncedAt string
}

func LeafInputDTOToDomain(dto *LeafInputDTO) *domain.Leaf {
	return &domain.Leaf{
		ID:       dto.ID,
		Title:    dto.Title,
		URL:      dto.URL,
		Platform: dto.Platform,
		Tags:     dto.Tags,
	}
}

func LeafDomainToOutputDTO(leaf *domain.Leaf) *LeafOutputDTO {
	return &LeafOutputDTO{
		ID:       leaf.ID,
		Title:    leaf.Title,
		URL:      leaf.URL,
		Platform: leaf.Platform,
		Tags:     leaf.Tags,
		SyncedAt: leaf.SyncedAt.Format(time.RFC3339),
	}
}
