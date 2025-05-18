package domain

import (
	"time"
)

// Leaf represents a stock article or bookmarked knowledge unit.
type Leaf struct {
	ID       string    `json:"id"`
	Title    string    `json:"title"`
	URL      string    `json:"url"`
	Platform string    `json:"platform"`
	Tags     []string  `json:"tags"`
	Read     bool      `json:"read"`
	SyncedAt time.Time `json:"synced_at"`
}

// NewLeaf creates a new Leaf instance.
func NewLeaf(
	id, title string, url string, platform string,
) *Leaf {
	return &Leaf{
		ID:       id,
		Title:    title,
		URL:      url,
		Platform: platform,
		Tags:     []string{},
		Read:     false,
		SyncedAt: time.Now().UTC(),
	}
}

// MarkAsRead marks this leaf as read.
func (l *Leaf) MarkAsRead() {
	l.Read = true
}

// UpdateTags replaces the tags.
func (l *Leaf) UpdateTags(tags []string) {
	newTags := make([]string, len(tags))
	copy(newTags, tags)
	l.Tags = newTags
}
