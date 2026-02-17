package bookmarks

import (
	"time"
)

// Bookmark represents a bookmark in the Notion Bookmarks database
type Bookmark struct {
	ID            string
	Title         string
	URL           string
	Summary       string
	TagIDs        []string // Related tag IDs
	DateAdded     time.Time
	DateProcessed time.Time
	Author        string
	ImageURL      string
	DatePublished string   // Date published as string from rich_text
	ManualListIDs []string // Related manual list IDs
	SmartListIDs  []string // Related smart list IDs
	Processed     bool     // Whether the bookmark has been processed
	Error         string   // Error message if processing failed
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// Property names for type-safe access to Notion properties
const (
	PropertyPage          = "page"
	PropertyURL           = "url"
	PropertySummary       = "summary"
	PropertyAuthor        = "author"
	PropertyImage         = "image"
	PropertyDateAdded     = "date_added"
	PropertyDateProcessed = "date_processed"
	PropertyDatePublished = "date_published"
	PropertyTag           = "tag"
	PropertyManualLists   = "manual_lists"
	PropertySmartLists    = "smart_lists"
	PropertyProcessed     = "processed"
	PropertyError         = "error"
)

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
	SortByDateAdded SortBy = PropertyDateAdded
	SortByTitle     SortBy = PropertyPage
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
