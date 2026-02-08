package bookmarks

import (
	"time"
)

// Bookmark represents a bookmark in the Notion Bookmarks database
type Bookmark struct {
	ID        string
	Title     string
	URL       string
	Summary   string
	TagIDs    []string // Related tag IDs
	DateAdded time.Time
	Author    string
	ImageURL  string
	Processed bool   // Whether the bookmark has been processed
	Error     string // Error message if processing failed
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Filter defines filtering options for listing bookmarks
type Filter struct {
	TitleContains   string
	URLContains     string
	SummaryContains string
	AuthorContains  string
	ErrorContains   string   // Filter by error message
	HasTag          string   // Filter by tag ID
	TagIDs          []string // Filter by multiple tag IDs
	Processed       *bool    // Filter by processed status (nil = no filter)
	Limit           int
}

// SortOrder defines the sort order for queries
type SortOrder string

const (
	SortOrderAsc  SortOrder = "ascending"
	SortOrderDesc SortOrder = "descending"
)

// SortBy defines what field to sort by
type SortBy string

const (
	SortByDateAdded SortBy = "Date Added"
	SortByTitle     SortBy = ""
	SortByCreatedAt SortBy = "created_time"
	SortByUpdatedAt SortBy = "last_edited_time"
)

// QueryOptions defines options for querying bookmarks
type QueryOptions struct {
	Filter    *Filter
	SortBy    SortBy
	SortOrder SortOrder
	Limit     int
}
