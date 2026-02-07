# Notion Bookmarks Manager

A Go application to interact with Notion's Bookmarks and Tags databases. This project provides a clean, idiomatic Go API for performing CRUD operations on bookmarks and tags stored in Notion.

## Features

- **Full CRUD Operations** for both Bookmarks and Tags
- **Tag Management** with find-or-create functionality
- **Advanced Filtering** and sorting capabilities
- **Type-safe Go structs** for Notion data
- **Relation handling** between Bookmarks and Tags
- **Clean package architecture** with separation of concerns

## Project Structure

```
hidrate-notion-bookmarks/
├── cmd/
│   └── example/
│       └── main.go           # Example usage demonstrating all features
├── internal/
│   └── config/
│       └── config.go         # Configuration loading from .env
├── pkg/
│   ├── notion/
│   │   ├── client.go         # Central Notion API client wrapper
│   │   ├── types.go          # Common utility functions and converters
│   │   └── errors.go         # Custom error types
│   ├── bookmarks/
│   │   ├── bookmarks.go      # Bookmark CRUD operations
│   │   ├── types.go          # Bookmark type definitions
│   │   └── mapper.go         # Notion API <-> Go struct mappings
│   └── tags/
│       ├── tags.go           # Tag CRUD operations
│       ├── types.go          # Tag type definitions
│       └── mapper.go         # Notion API <-> Go struct mappings
├── .env                      # Environment configuration (not committed)
├── .gitignore               
├── go.mod                    
└── README.md
```

## Database Schema

### Bookmarks Database
- **Title** (title) - The bookmark title
- **URL** (url) - The bookmark URL
- **Summary** (rich_text) - Summary of the bookmark content
- **Author** (rich_text) - Author of the content
- **Date Added** (date) - When the bookmark was added
- **image** (url) - Image URL
- **Tags** (relation) - Related tags from Tags database

### Tags Database
- **Name** (title) - The tag name

## Setup

### 1. Prerequisites

- Go 1.21 or later
- A Notion account with API access
- Two Notion databases: Bookmarks and Tags

### 2. Create Notion Integration

1. Go to [Notion Integrations](https://www.notion.so/my-integrations)
2. Click "New integration"
3. Give it a name and select the workspace
4. Copy the "Internal Integration Token"

### 3. Share Databases with Integration

1. Open your Bookmarks database in Notion
2. Click the "..." menu → "Add connections"
3. Select your integration
4. Repeat for the Tags database

### 4. Get Database IDs

Database IDs are found in the URL when viewing the database:
```
https://notion.so/workspace/<database_id>?v=...
```

### 5. Configure Environment

Create a `.env` file in the project root:

```env
NOTION_API_KEY=your_notion_integration_token
NOTION_BOOKMARKS_DB_ID=your_bookmarks_database_id
NOTION_TAGS_DB_ID=your_tags_database_id
```

### 6. Install Dependencies

```bash
go mod download
```

## Usage

### Run the Example

```bash
go run cmd/example/main.go
```

Or build and run:

```bash
go build -o bin/example ./cmd/example
./bin/example
```

### Using in Your Code

#### Initialize Client

```go
package main

import (
    "context"
    "log"
    
    "github.com/pgodinho/hidrate-notion-bookmarks/internal/config"
    "github.com/pgodinho/hidrate-notion-bookmarks/pkg/bookmarks"
    "github.com/pgodinho/hidrate-notion-bookmarks/pkg/notion"
    "github.com/pgodinho/hidrate-notion-bookmarks/pkg/tags"
)

func main() {
    // Load configuration
    cfg, err := config.Load()
    if err != nil {
        log.Fatal(err)
    }
    
    // Create client
    client := notion.NewClient(cfg.NotionAPIKey, cfg.BookmarksDBID, cfg.TagsDBID)
    
    // Create services
    tagService := tags.NewService(client)
    bookmarkService := bookmarks.NewService(client)
    
    ctx := context.Background()
    
    // Use the services...
}
```

#### Working with Tags

```go
// Create a tag
tag, err := tagService.Create(ctx, &tags.Tag{
    Name: "Go Programming",
})

// Find or create a tag
tag, err := tagService.FindOrCreate(ctx, "Web Development")

// Get a tag by ID
tag, err := tagService.Get(ctx, "tag-id")

// Get a tag by name
tag, err := tagService.GetByName(ctx, "APIs")

// List all tags
allTags, err := tagService.List(ctx, nil)

// List tags with filter
filteredTags, err := tagService.List(ctx, &tags.Filter{
    NameContains: "Programming",
    Limit:        10,
})

// Update a tag
tag.Name = "Advanced Go"
updated, err := tagService.Update(ctx, tag.ID, tag)

// Delete a tag (archives it in Notion)
err := tagService.Delete(ctx, tag.ID)
```

#### Working with Bookmarks

```go
// Create a bookmark
bookmark, err := bookmarkService.Create(ctx, &bookmarks.Bookmark{
    Title:     "Go Documentation",
    URL:       "https://go.dev/doc/",
    Summary:   "Official Go documentation",
    Author:    "Go Team",
    DateAdded: time.Now(),
    TagIDs:    []string{tag1.ID, tag2.ID},
    ImageURL:  "https://example.com/image.png",
})

// Get a bookmark by ID
bookmark, err := bookmarkService.Get(ctx, "bookmark-id")

// Update a bookmark
bookmark.Summary = "Updated summary"
updated, err := bookmarkService.Update(ctx, bookmark.ID, bookmark)

// Delete a bookmark (archives it in Notion)
err := bookmarkService.Delete(ctx, bookmark.ID)

// List all bookmarks
allBookmarks, err := bookmarkService.List(ctx, nil)

// Filter bookmarks
filtered, err := bookmarkService.List(ctx, &bookmarks.Filter{
    TitleContains:   "Go",
    URLContains:     "go.dev",
    SummaryContains: "documentation",
    AuthorContains:  "Team",
    HasTag:          tagID,
    Limit:           10,
})

// Query with sorting
results, err := bookmarkService.Query(ctx, &bookmarks.QueryOptions{
    SortBy:    bookmarks.SortByDateAdded,
    SortOrder: bookmarks.SortOrderDesc,
    Limit:     5,
})

// Add tags to a bookmark
updated, err := bookmarkService.AddTags(ctx, bookmarkID, []string{tagID1, tagID2})

// Remove tags from a bookmark
updated, err := bookmarkService.RemoveTags(ctx, bookmarkID, []string{tagID1})
```

## API Reference

### Tags Service

#### `Create(ctx, tag) (*Tag, error)`
Creates a new tag in Notion.

#### `Get(ctx, id) (*Tag, error)`
Retrieves a tag by its ID.

#### `GetByName(ctx, name) (*Tag, error)`
Retrieves a tag by its name. Returns nil if not found.

#### `Update(ctx, id, tag) (*Tag, error)`
Updates an existing tag.

#### `Delete(ctx, id) error`
Archives a tag (Notion doesn't support true deletion).

#### `List(ctx, filter) ([]*Tag, error)`
Lists all tags with optional filtering.

#### `FindOrCreate(ctx, name) (*Tag, error)`
Finds a tag by name or creates it if it doesn't exist.

### Bookmarks Service

#### `Create(ctx, bookmark) (*Bookmark, error)`
Creates a new bookmark in Notion.

#### `Get(ctx, id) (*Bookmark, error)`
Retrieves a bookmark by its ID.

#### `Update(ctx, id, bookmark) (*Bookmark, error)`
Updates an existing bookmark.

#### `Delete(ctx, id) error`
Archives a bookmark (Notion doesn't support true deletion).

#### `List(ctx, filter) ([]*Bookmark, error)`
Lists all bookmarks with optional filtering and sorting.

#### `Query(ctx, options) ([]*Bookmark, error)`
Queries bookmarks with advanced filtering and sorting options.

#### `AddTags(ctx, bookmarkID, tagIDs) (*Bookmark, error)`
Adds tags to a bookmark (preserves existing tags).

#### `RemoveTags(ctx, bookmarkID, tagIDs) (*Bookmark, error)`
Removes specific tags from a bookmark.

## Dependencies

- [github.com/jomei/notionapi](https://github.com/jomei/notionapi) - Official Notion SDK for Go
- [github.com/joho/godotenv](https://github.com/joho/godotenv) - Environment variable loading

## Error Handling

The package provides custom error types through `pkg/notion/errors.go`:

- `ErrNotFound` - Resource not found
- `ErrUnauthorized` - Invalid API key or insufficient permissions
- `ErrRateLimited` - Too many requests
- `ErrInvalidInput` - Invalid input data
- `ErrAPIError` - General Notion API error

Errors include context about the operation that failed:

```go
bookmark, err := bookmarkService.Get(ctx, "invalid-id")
if err != nil {
    // Error message includes operation context
    log.Printf("Error: %v", err)
    // Output: "get bookmark failed: failed to get bookmark with ID: invalid-id - ..."
}
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is open source and available under the MIT License.

## Future Enhancements

- Web scraping functionality to extract content from URLs
- Automatic metadata generation
- Tag suggestion based on content
- Batch operations for bulk imports
- Export functionality to various formats
- CLI tool for command-line operations
