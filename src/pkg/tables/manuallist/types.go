package manuallist

import (
	"time"
)

// ManualListItem represents an item in the Notion Manual List database
type ManualListItem struct {
	ID        string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Property names for type-safe access to Notion properties
const (
	PropertyName = "name"
)

// Filter defines filtering options for listing manual list items
type Filter struct {
	NameContains string
	Limit        int
}
