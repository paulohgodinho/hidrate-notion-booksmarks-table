package notion

import (
	"github.com/jomei/notionapi"
)

// Client wraps the Notion API client with our configuration
type Client struct {
	api          *notionapi.Client
	bookmarksDB  notionapi.DatabaseID
	tagsDB       notionapi.DatabaseID
	manualListDB notionapi.DatabaseID
	smartListDB  notionapi.DatabaseID
}

// NewClient creates a new Notion client with the provided configuration
func NewClient(apiKey, bookmarksDBID, tagsDBID, manualListDBID, smartListDBID string) *Client {
	return &Client{
		api:          notionapi.NewClient(notionapi.Token(apiKey)),
		bookmarksDB:  notionapi.DatabaseID(bookmarksDBID),
		tagsDB:       notionapi.DatabaseID(tagsDBID),
		manualListDB: notionapi.DatabaseID(manualListDBID),
		smartListDB:  notionapi.DatabaseID(smartListDBID),
	}
}

// API returns the underlying Notion API client
func (c *Client) API() *notionapi.Client {
	return c.api
}

// BookmarksDB returns the Bookmarks database ID
func (c *Client) BookmarksDB() notionapi.DatabaseID {
	return c.bookmarksDB
}

// TagsDB returns the Tags database ID
func (c *Client) TagsDB() notionapi.DatabaseID {
	return c.tagsDB
}

// ManualListDB returns the Manual List database ID
func (c *Client) ManualListDB() notionapi.DatabaseID {
	return c.manualListDB
}

// SmartListDB returns the Smart List database ID
func (c *Client) SmartListDB() notionapi.DatabaseID {
	return c.smartListDB
}
