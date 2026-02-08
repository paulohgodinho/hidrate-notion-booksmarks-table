# Notion Bookmarks Manager

A Go application to interact with Notion's Bookmarks and Tags databases. This project provides a clean, idiomatic Go API for performing CRUD operations on bookmarks and tags stored in Notion.

## Features

- **Full CRUD Operations** for both Bookmarks and Tags
- **Tag Management** with find-or-create functionality
- **Advanced Filtering** and sorting capabilities
- **Type-safe Go structs** for Notion data
- **Relation handling** between Bookmarks and Tags
- **Clean package architecture** with separation of concerns
- **Automated Web Scraping** with metadata extraction
- **Docker support** with Docker Compose orchestration

## Project Structure

```
hidrate-notion-bookmarks/
├── src/                      # All Go source code
│   ├── main.go               # Main processor application
│   ├── config.go             # Configuration loading from .env
│   ├── pkg/
│   │   ├── notion/
│   │   │   ├── client.go     # Central Notion API client wrapper
│   │   │   ├── types.go      # Common utility functions and converters
│   │   │   ├── errors.go     # Custom error types
│   │   │   └── uploader.go   # Image uploader for Notion
│   │   ├── bookmarks/
│   │   │   ├── bookmarks.go  # Bookmark CRUD operations
│   │   │   ├── types.go      # Bookmark type definitions
│   │   │   └── mapper.go     # Notion API <-> Go struct mappings
│   │   ├── tags/
│   │   │   ├── tags.go       # Tag CRUD operations
│   │   │   ├── types.go      # Tag type definitions
│   │   │   └── mapper.go     # Notion API <-> Go struct mappings
│   │   └── scraper/
│   │       ├── client.go     # Webmeatscraper HTTP client
│   │       └── types.go      # Scraper request/response types
│   ├── go.mod
│   └── go.sum
├── docker/                   # Docker-related files
├── bin/                      # Build output (created by build.sh)
├── .env                      # Environment configuration (not committed)
├── .gitignore
├── docker-compose.yml        # Docker Compose configuration
├── Dockerfile                # Docker build configuration
├── build.sh                  # Build script
├── PLAN.md
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
- **Processed** (checkbox) - Whether the bookmark has been processed
- **Error** (rich_text) - Error message if processing failed

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

# Webmeatscraper service URL (optional, defaults to http://localhost:7878)
# Use http://webmeatscraper:7878 when running in Docker Compose
# Use http://localhost:7878 when running locally
WEBMEATSCRAPER_URL=http://localhost:7878

# Image Upload Configuration (optional)
# Enable uploading images to Notion storage (default: true)
UPLOAD_IMAGES_TO_NOTION=true

# Maximum time to wait for upload completion (default: 30s)
IMAGE_UPLOAD_TIMEOUT=30s

# How often to check upload status (default: 3s)
IMAGE_UPLOAD_POLL_INTERVAL=3s

# Use external URL if upload fails/times out (default: true)
FALLBACK_TO_EXTERNAL_URL=true

# Debug Configuration (optional)
# Print full JSON output from webmeatscraper (default: false)
DEBUG=false
```

### 6. Install Dependencies

```bash
go mod download
```

## Usage

### Bookmark Processor

The bookmark processor automatically fetches unprocessed bookmarks from Notion, scrapes their content using the webmeatscraper service, and updates the bookmark metadata.

#### Features
- Processes ALL unprocessed bookmarks in a single run
- Extracts metadata (author, image) from scraped content
- **Uploads images to Notion storage** for permanent hosting
- **Sets page covers** with uploaded images
- Only updates empty fields (non-destructive)
- Prints full JSON response for each bookmark
- Handles errors gracefully and logs them to Notion
- Shows progress and summary statistics

#### Running with Docker Compose (Recommended)

The easiest way to run the processor is with Docker Compose, which automatically starts both the scraper service and the processor:

```bash
# Build and run both services
docker-compose up --build

# View logs
docker-compose logs -f

# Stop services
docker-compose down
```

#### Running Locally

If you prefer to run locally:

```bash
# Terminal 1: Start the scraper service
docker run -p 7878:7878 --platform linux/amd64 ghcr.io/paulohgodinho/webmeatscraper:latest

# Terminal 2: Run the processor
cd src
go run main.go
```

Or build and run:

```bash
cd src
go build -o ../bin/bookmarks-processor .
../bin/bookmarks-processor
```

Or use the build script to build for all platforms:

```bash
./build.sh
# Creates: bin/hidrate-notion-bookmarks-linux-amd64
#          bin/hidrate-notion-bookmarks-darwin-arm64
#          bin/hidrate-notion-bookmarks-windows-amd64.exe
```

#### Configuration

The processor supports several configuration options via environment variables:

##### Basic Configuration

```env
# Use http://webmeatscraper:7878 when running in Docker Compose
# Use http://localhost:7878 when running locally
WEBMEATSCRAPER_URL=http://localhost:7878
```

##### Image Upload Configuration

The processor can automatically upload scraped images to Notion for permanent storage:

| Variable | Default | Description |
|----------|---------|-------------|
| `UPLOAD_IMAGES_TO_NOTION` | `true` | Enable/disable image uploads to Notion |
| `IMAGE_UPLOAD_TIMEOUT` | `30s` | Max time to wait for upload completion |
| `IMAGE_UPLOAD_POLL_INTERVAL` | `3s` | How often to check upload status |
| `FALLBACK_TO_EXTERNAL_URL` | `true` | Use external URL if upload fails |

**How it works:**
1. Extracts image URL from scraped content (OG image, Twitter Card, etc.)
2. Uploads image to Notion storage using the Indirect Import method
3. Sets the uploaded image as the page cover (always updates)
4. Also stores the original external URL in the `image` database property

**Benefits:**
- Images are permanently stored in Notion
- Prevents broken images if external URLs become unavailable
- Page covers display immediately in Notion UI
- Falls back to external URLs gracefully if upload fails

##### Debug Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `DEBUG` | `false` | Print full JSON output from webmeatscraper |

**When enabled:**
- Prints the complete scraped content JSON for each bookmark
- Useful for debugging and understanding what data is being extracted
- Automatically enabled in Docker Compose for troubleshooting
- Set to `false` in production for cleaner logs

#### How It Works

1. Connects to the scraper service and performs a health check
2. Fetches ALL unprocessed bookmarks (where Processed = false)
3. Iterates through each bookmark:
   - Scrapes the bookmark's URL using webmeatscraper
   - Prints the full JSON response to stdout
   - Updates the bookmark with scraped metadata:
     - Sets Author if empty
     - Sets Image URL if empty (tries multiple sources: OG image, Twitter image, etc.)
   - **Uploads image to Notion storage** (if enabled)
   - **Sets page cover** with uploaded image
   - Marks the bookmark as processed
   - On error, sets the Error field and continues to next bookmark
4. Displays summary statistics (total, successful, failed)

#### Example Output

**With DEBUG=false (default, clean output):**
```
=== Notion Bookmark Processor ===

✓ Image upload to Notion: ENABLED
  Debug mode: disabled

Checking scraper service at http://localhost:7878...
✓ Scraper service is healthy (status: ok)

Fetching unprocessed bookmarks...
Found 3 unprocessed bookmark(s)

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Processing bookmark 1 of 3
Title: Example Article
ID: abc123...
URL: https://example.com/article

Scraping content...

Updating bookmark with scraped metadata...
  ✓ Set author: John Doe
  📤 Uploading image to Notion... ✅ Uploaded (ID: abc-123-def)
  ✓ Set image property: https://example.com/image.jpg
✓ Bookmark metadata updated
  🖼️  Setting page cover... ✅ Cover set
✓ Bookmark marked as processed

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
=== Processing Complete ===
Total: 3 bookmarks
✓ Successfully processed: 3
✗ Failed: 0
```

**With DEBUG=true (full JSON output):**
```
=== Notion Bookmark Processor ===

✓ Image upload to Notion: ENABLED
✓ Debug mode: ENABLED (full JSON output)

Checking scraper service at http://localhost:7878...
✓ Scraper service is healthy (status: ok)

Fetching unprocessed bookmarks...
Found 3 unprocessed bookmark(s)

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Processing bookmark 1 of 3
Title: Example Article
ID: abc123...
URL: https://example.com/article

Scraping content...
=== SCRAPED CONTENT (Full JSON) ===
{
  "content": "This is the extracted content...",
  "image": "https://example.com/image.jpg",
  "metadata": {
    "title": "Example Article",
    "author": "John Doe",
    "og_image": "https://example.com/og-image.jpg",
    ...
  }
}
=== END OF SCRAPED CONTENT ===

Updating bookmark with scraped metadata...
  ✓ Set author: John Doe
  📤 Uploading image to Notion... ✅ Uploaded (ID: abc-123-def)
  ✓ Set image property: https://example.com/image.jpg
✓ Bookmark metadata updated
  🖼️  Setting page cover... ✅ Cover set
✓ Bookmark marked as processed

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
=== Processing Complete ===
Total: 3 bookmarks
✓ Successfully processed: 3
✗ Failed: 0
```

### Run the Application

```bash
cd src
go run main.go
```

Or build and run:

```bash
cd src
go build -o ../bin/bookmarks-processor .
../bin/bookmarks-processor
```

### Using in Your Code

#### Initialize Client

```go
package main

import (
    "context"
    "log"
    
    "github.com/pgodinho/hidrate-notion-bookmarks/pkg/bookmarks"
    "github.com/pgodinho/hidrate-notion-bookmarks/pkg/notion"
    "github.com/pgodinho/hidrate-notion-bookmarks/pkg/tags"
)

func main() {
    // Load configuration
    cfg, err := Load()
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
    Processed: false,
    Error:     "",
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
    ErrorContains:   "timeout",
    HasTag:          tagID,
    Processed:       &processedBool, // nil = no filter, true/false = filter
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

// Mark bookmark as processed
processed, err := bookmarkService.MarkAsProcessed(ctx, bookmarkID)

// Set error on bookmark
withError, err := bookmarkService.SetError(ctx, bookmarkID, "Failed to fetch: timeout")

// Get unprocessed bookmarks
unprocessed, err := bookmarkService.GetUnprocessed(ctx, 10)

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

#### `MarkAsProcessed(ctx, bookmarkID) (*Bookmark, error)`
Marks a bookmark as processed and clears any error message.

#### `SetError(ctx, bookmarkID, errorMsg) (*Bookmark, error)`
Sets an error message on a bookmark and marks it as not processed.

#### `GetUnprocessed(ctx, limit) ([]*Bookmark, error)`
Retrieves all unprocessed bookmarks.

#### `GetWithErrors(ctx, limit) ([]*Bookmark, error)`
Retrieves all bookmarks that have error messages.

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

- AI-powered tag generation from scraped content
- AI-powered summary generation
- Batch processing for multiple bookmarks
- Retry logic for failed scrapes
- Web UI for monitoring and manual triggering
- Metrics and analytics dashboard
- Scheduled processing with cron jobs
- Content storage in Notion pages
- Export functionality to various formats
