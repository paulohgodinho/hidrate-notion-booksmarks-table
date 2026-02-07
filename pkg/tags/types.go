package tags

import (
	"time"
)

// Tag represents a tag in the Notion Tags database
type Tag struct {
	ID        string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Filter defines filtering options for listing tags
type Filter struct {
	NameContains string
	Limit        int
}
