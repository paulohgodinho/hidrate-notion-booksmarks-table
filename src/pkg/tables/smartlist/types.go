package smartlist

import (
	"time"
)

// SmartListItem represents an item in the Notion Smart List database
type SmartListItem struct {
	ID        string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Property names for type-safe access to Notion properties
const (
	PropertyName = "name"
)

// Filter defines filtering options for listing smart list items
type Filter struct {
	NameContains string
	Limit        int
}
