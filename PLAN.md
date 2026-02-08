# Implementation Plan: Notion Bookmark Processor with Webmeatscraper

## 📋 Overview

This document outlines the plan to integrate the webmeatscraper service with our Notion bookmark manager. The goal is to automatically fetch and process unprocessed bookmarks by scraping their content.

---

## 🎯 Requirements

### User Requirements
1. **A**: Use `DoesNotEqual` approach for checkbox filtering ✅
2. **A**: Print full JSON dump with all metadata ✅
3. **Process ALL**: Process all unprocessed bookmarks in a single run ✅ (Updated from "one at a time")
4. **Error Handling**: Mark as processed and add error to Error field on failure ✅
5. **Docker Image**: Use `ghcr.io/paulohgodinho/webmeatscraper:latest --platform linux/amd64` ✅
6. **Summary Field**: Do NOT modify the summary field ✅

### Functional Requirements
- Fetch ALL unprocessed bookmarks (Processed = false) ✅
- Iterate through each bookmark ✅
- Send bookmark URL to webmeatscraper service ✅
- Print full JSON response from scraper for each bookmark ✅
- Update bookmark metadata (Author, ImageURL) if empty ✅
- Mark bookmark as processed (success) or set error (failure) ✅
- Never modify Summary field ✅
- Show progress and summary statistics ✅

---

## 📁 Project Structure Changes

```
hidrate-notion-bookmarks/
├── cmd/
│   ├── example/
│   │   └── main.go                    # Existing example
│   └── processor/                     # ✨ NEW
│       └── main.go                    # Bookmark processor
├── internal/
│   └── config/
│       └── config.go                  # 🔧 Updated: Add WEBMEATSCRAPER_URL
├── pkg/
│   ├── notion/
│   │   ├── client.go
│   │   ├── types.go
│   │   └── errors.go
│   ├── bookmarks/
│   │   ├── bookmarks.go               # 🔧 Updated: Fix GetUnprocessed()
│   │   ├── types.go
│   │   └── mapper.go
│   ├── tags/
│   │   ├── tags.go
│   │   ├── types.go
│   │   └── mapper.go
│   └── scraper/                       # ✨ NEW
│       ├── client.go                  # Scraper HTTP client
│       └── types.go                   # Scraper request/response types
├── .env                               # 🔧 Updated: Add WEBMEATSCRAPER_URL
├── .gitignore
├── Dockerfile                         # ✨ NEW
├── docker-compose.yml                 # ✨ NEW
├── go.mod
├── go.sum
├── PLAN.md                            # ✨ THIS FILE
└── README.md                          # 🔧 Updated: Add processor docs
```

**Legend:**
- ✨ NEW: New file to create
- 🔧 Updated: Existing file to modify

---

## 🔨 Implementation Steps

### Step 1: Fix Checkbox Filter Bug ✅ CRITICAL
**File:** `pkg/bookmarks/bookmarks.go`

**Problem:** The current implementation uses `Equals: false` which gets omitted due to JSON's `omitempty` tag.

**Solution:** Change `GetUnprocessed()` method to use `DoesNotEqual: true`:

```go
func (s *Service) GetUnprocessed(ctx context.Context, limit int) ([]*Bookmark, error) {
    query := &notionapi.DatabaseQueryRequest{
        Filter: &notionapi.PropertyFilter{
            Property: "Processed",
            Checkbox: &notionapi.CheckboxFilterCondition{
                DoesNotEqual: true,  // Find all where Processed != true
            },
        },
    }
    
    if limit > 0 {
        query.PageSize = limit
    }
    
    resp, err := s.client.API().Database.Query(ctx, s.client.BookmarksDB(), query)
    if err != nil {
        return nil, notion.NewError("get unprocessed bookmarks", err, "failed to query unprocessed bookmarks")
    }

    bookmarks := make([]*Bookmark, 0, len(resp.Results))
    for _, page := range resp.Results {
        bookmark, err := ToBookmark(&page)
        if err != nil {
            continue
        }
        bookmarks = append(bookmarks, bookmark)
    }

    return bookmarks, nil
}
```

**Why this works:** `DoesNotEqual: true` finds all records where Processed is NOT true (includes false and null).

---

### Step 2: Create Scraper Package ✨ NEW

#### File: `pkg/scraper/types.go`

Defines the data structures for interacting with the webmeatscraper API.

**Key Types:**
- `ScrapedContent`: Main response from scraper
- `Metadata`: Contains all extracted metadata (title, author, dates, platform-specific fields)
- `ScrapeRequest`: Request body for scraper API

**Features:**
- Full support for all metadata fields from webmeatscraper
- Platform-specific fields (YouTube, Twitter, Amazon, Reddit)
- Optional image field (can be null)

#### File: `pkg/scraper/client.go`

HTTP client for communicating with the webmeatscraper service.

**Key Methods:**
- `NewClient(baseURL string)`: Initialize client
- `Scrape(ctx, url)`: Send URL and get scraped content
- `Health(ctx)`: Check if scraper service is healthy

**Features:**
- 30-second timeout for scraping requests
- Context support for cancellation
- Proper error handling and status code checking
- JSON encoding/decoding

---

### Step 3: Update Configuration 🔧

#### File: `internal/config/config.go`

Add new field for scraper URL:

```go
type Config struct {
    NotionAPIKey      string
    BookmarksDBID     string
    TagsDBID          string
    WebmeatscraperURL string  // NEW
}
```

**Note:** WebmeatscraperURL is optional and will default to `http://localhost:7878` if not set.

#### File: `.env`

Add new environment variable:

```env
# Webmeatscraper service URL
# Use http://webmeatscraper:7878 when running in Docker
# Use http://localhost:7878 when running locally
WEBMEATSCRAPER_URL=http://localhost:7878
```

---

### Step 4: Create Dockerfile ✨ NEW

Multi-stage Docker build for the Go application:

**Stage 1 (Builder):**
- Use golang:1.21-alpine as base
- Copy go.mod/go.sum and download dependencies
- Copy source code
- Build processor binary with CGO_ENABLED=0

**Stage 2 (Runtime):**
- Use alpine:latest for minimal image size
- Install ca-certificates for HTTPS support
- Copy binary from builder stage
- Copy .env file
- Set processor as CMD

**Benefits:**
- Small final image size
- Secure (minimal attack surface)
- Fast builds with layer caching

---

### Step 5: Create Docker Compose Configuration ✨ NEW

#### File: `docker-compose.yml`

Defines two services:

**Service 1: webmeatscraper**
- Image: `ghcr.io/paulohgodinho/webmeatscraper:latest`
- Platform: `linux/amd64` (as specified)
- Port: 7878
- Health check: Polls `/health` endpoint every 10s
- Network: notion-network

**Service 2: notion-processor**
- Built from local Dockerfile
- Depends on webmeatscraper (waits for health check)
- Environment: Loaded from .env + WEBMEATSCRAPER_URL
- Network: notion-network
- Volume: Mounts .env as read-only

**Network:**
- Bridge network for service communication
- Services can reach each other by service name
- Isolated from host network

**Key Features:**
- Health check ensures scraper is ready before processor starts
- Internal DNS allows processor to call `http://webmeatscraper:7878`
- Automatic restart on failure (can be added if needed)

---

### Step 6: Create Processor Application ✨ NEW

#### File: `cmd/processor/main.go`

Main application that processes bookmarks.

**Flow:**
1. Load configuration from .env
2. Initialize Notion client, bookmark service, scraper client
3. Check scraper service health
4. Fetch ONE unprocessed bookmark (Processed = false)
5. If no bookmarks found, exit successfully
6. Scrape the bookmark's URL
7. Print full JSON response to stdout
8. On success:
   - Update Author (if empty)
   - Update ImageURL (if empty)
   - Set Processed = true
   - Clear Error field
9. On failure:
   - Set Error field with error message
   - Set Processed = false
10. Print summary and exit

**Key Features:**
- Detailed logging at each step
- Pretty-printed JSON output
- Non-destructive updates (only sets empty fields)
- Proper error handling with context
- Returns appropriate exit codes

**Example Output:**
```
=== Notion Bookmark Processor ===

Checking scraper service at http://webmeatscraper:7878...
✓ Scraper service is healthy

Fetching unprocessed bookmarks...
Found unprocessed bookmark: Example Article
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
    ...
  }
}
=== END OF SCRAPED CONTENT ===

Updating bookmark with scraped metadata...
  ✓ Set author: John Doe
  ✓ Set image: https://example.com/image.jpg
✓ Bookmark marked as processed

=== Processing Complete ===
```

---

### Step 7: Update README Documentation 🔧

Add new section after the example section covering:

1. **Bookmark Processor Overview**
   - What it does
   - How it works

2. **Running with Docker Compose**
   - Build and run command
   - View logs
   - Stop services

3. **Running Locally**
   - Start scraper manually
   - Run processor

4. **How It Works** (detailed flow)

5. **Configuration**
   - Environment variables
   - Docker vs local settings

6. **Example Output**
   - Sample console output
   - Success and error cases

---

## 🧪 Testing Plan

### Test 1: Checkbox Filter Fix
```bash
go run cmd/example/main.go
```
**Expected:** "Getting Unprocessed Bookmarks" section should succeed.

### Test 2: Build Scraper Package
```bash
go build ./pkg/scraper
```
**Expected:** No compilation errors.

### Test 3: Test Processor Locally
```bash
# Terminal 1: Start scraper
docker run -p 7878:7878 --platform linux/amd64 ghcr.io/paulohgodinho/webmeatscraper:latest

# Terminal 2: Run processor
go run cmd/processor/main.go
```
**Expected:** 
- Connects to scraper
- Processes one bookmark
- Prints JSON
- Updates Notion

### Test 4: Test with Docker Compose
```bash
docker-compose up --build
```
**Expected:** Both services start, processor runs successfully.

### Test 5: Test Error Handling
1. Create bookmark with invalid URL in Notion
2. Run processor
3. **Expected:** 
   - Error caught and logged
   - Bookmark's Error field populated
   - Processed = false

---

## 📊 Data Flow Diagram

```
┌─────────────────┐
│  Notion API     │
│  (Bookmarks DB) │
└────────┬────────┘
         │
         │ 1. Fetch unprocessed
         │    (Processed = false)
         ▼
┌─────────────────┐
│   Processor     │
│   (Go App)      │
└────────┬────────┘
         │
         │ 2. Send URL
         ▼
┌─────────────────┐
│ Webmeatscraper  │
│   (Node.js)     │
└────────┬────────┘
         │
         │ 3. Scrape & Extract
         ▼
     Internet
         │
         │ 4. Return JSON
         ▼
┌─────────────────┐
│   Processor     │
│ (Print & Store) │
└────────┬────────┘
         │
         │ 5. Update metadata
         │    Mark processed
         ▼
┌─────────────────┐
│  Notion API     │
│  (Updated)      │
└─────────────────┘
```

---

## ⚠️ Important Notes

### What Gets Updated
- ✅ Author (only if empty)
- ✅ ImageURL (only if empty)
- ✅ Processed flag (set to true on success, false on error)
- ✅ Error field (populated on failure, cleared on success)

### What Does NOT Get Updated
- ❌ Summary (never touched)
- ❌ Title (never touched)
- ❌ URL (never touched)
- ❌ Tags (never touched)
- ❌ Date Added (never touched)

### Processing Rules
- Processes exactly ONE bookmark per run
- Only updates empty fields (non-destructive)
- Failed scrapes don't block - they're marked with errors
- Scraped content is printed but not stored

### Docker Considerations
- Scraper runs on `linux/amd64` platform
- Services communicate via internal Docker network
- Health check ensures scraper is ready before processing
- Logs available via `docker-compose logs`

---

## 🔄 Future Enhancements

After basic implementation is working, consider:

1. **Batch Processing**: Add flag to process N bookmarks in one run
2. **Retry Logic**: Automatically retry failed bookmarks after X time
3. **Content Storage**: Add new field to store scraped content if needed
4. **Tag Generation**: Use AI (OpenAI/Claude) to generate tags from content
5. **Summary Generation**: Use AI to create summaries from scraped content
6. **Scheduling**: Add cron job or systemd timer for automatic processing
7. **Web UI**: Create simple web interface to:
   - Trigger processing manually
   - View processing status
   - See failed bookmarks
   - Retry specific bookmarks
8. **Metrics**: Track processing stats (success rate, average time, etc.)
9. **Notifications**: Send alerts when processing completes or fails

---

## 🚀 Implementation Order

1. ✅ Create PLAN.md (this file)
2. Fix checkbox filter in bookmarks.go
3. Create pkg/scraper/types.go
4. Create pkg/scraper/client.go
5. Update internal/config/config.go
6. Update .env
7. Create Dockerfile
8. Create docker-compose.yml
9. Create cmd/processor/main.go
10. Test locally
11. Test with Docker Compose
12. Update README.md
13. Final testing and validation

---

## ✅ Success Criteria

The implementation is complete when:

- [ ] Checkbox filter bug is fixed
- [ ] Scraper package compiles and works
- [ ] Processor can fetch unprocessed bookmarks
- [ ] Processor can scrape URLs via webmeatscraper
- [ ] Full JSON response is printed to console
- [ ] Bookmarks are updated correctly (Author, ImageURL)
- [ ] Processed flag is set correctly
- [ ] Errors are caught and stored
- [ ] Docker Compose setup works end-to-end
- [ ] Documentation is complete and accurate
- [ ] All tests pass

---

## 📝 Notes

- This plan follows the user's specific requirements exactly
- Summary field is never modified (as requested)
- One bookmark at a time (as requested)
- Full JSON output (as requested)
- Error handling stores errors in Notion (as requested)
- Uses specified Docker image with platform flag (as requested)

---

**Plan Version:** 1.0  
**Date:** 2026-02-07  
**Status:** Ready for Implementation
