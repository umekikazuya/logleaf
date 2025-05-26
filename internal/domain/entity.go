package domain

import (
	"time"
)

// Entity represents a leaf in the system.
type Leaf struct {
	ID       string
	Title    string
	URL      string
	Platform string
	Tags     []string
	Read     bool
	SyncedAt time.Time
}

// NewLeaf creates a new Leaf instance.
func NewLeaf(
	id string, title string, url string, platform string,
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
