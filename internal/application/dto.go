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
	tags := leaf.Tags()
	tagStrings := make([]string, len(tags))
	for i, tag := range tags {
		tagStrings[i] = tag.String()
	}

	return &LeafOutputDTO{
		ID:       leaf.ID().String(),
		Note:     leaf.Note(),
		URL:      leaf.URL().String(),
		Platform: leaf.Platform(),
		Read:     leaf.Read(),
		Tags:     tagStrings,
		SyncedAt: leaf.SyncedAt().Format(time.RFC3339),
	}
}
