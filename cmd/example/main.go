package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/pgodinho/hidrate-notion-bookmarks/internal/config"
	"github.com/pgodinho/hidrate-notion-bookmarks/pkg/bookmarks"
	"github.com/pgodinho/hidrate-notion-bookmarks/pkg/notion"
	"github.com/pgodinho/hidrate-notion-bookmarks/pkg/tags"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Create Notion client
	client := notion.NewClient(cfg.NotionAPIKey, cfg.BookmarksDBID, cfg.TagsDBID)

	// Create services
	tagService := tags.NewService(client)
	bookmarkService := bookmarks.NewService(client)

	ctx := context.Background()

	// Example 1: Create tags
	fmt.Println("=== Creating Tags ===")
	tag1, err := tagService.FindOrCreate(ctx, "Go Programming")
	if err != nil {
		log.Fatalf("Failed to create tag: %v", err)
	}
	fmt.Printf("Created/Found tag: %s (ID: %s)\n", tag1.Name, tag1.ID)

	tag2, err := tagService.FindOrCreate(ctx, "Web Development")
	if err != nil {
		log.Fatalf("Failed to create tag: %v", err)
	}
	fmt.Printf("Created/Found tag: %s (ID: %s)\n", tag2.Name, tag2.ID)

	tag3, err := tagService.FindOrCreate(ctx, "APIs")
	if err != nil {
		log.Fatalf("Failed to create tag: %v", err)
	}
	fmt.Printf("Created/Found tag: %s (ID: %s)\n\n", tag3.Name, tag3.ID)

	// Example 2: List all tags
	fmt.Println("=== Listing All Tags ===")
	allTags, err := tagService.List(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to list tags: %v", err)
	}
	for _, tag := range allTags {
		fmt.Printf("- %s (ID: %s)\n", tag.Name, tag.ID)
	}
	fmt.Println()

	// Example 3: Create a bookmark
	fmt.Println("=== Creating Bookmark ===")
	bookmark := &bookmarks.Bookmark{
		Title:     "Go Official Documentation",
		URL:       "https://go.dev/doc/",
		Summary:   "Official Go programming language documentation with tutorials, guides, and reference materials.",
		Author:    "Go Team",
		DateAdded: time.Now(),
		TagIDs:    []string{tag1.ID, tag3.ID}, // Go Programming and APIs tags
		ImageURL:  "https://go.dev/images/gophers/ladder.svg",
	}

	createdBookmark, err := bookmarkService.Create(ctx, bookmark)
	if err != nil {
		log.Fatalf("Failed to create bookmark: %v", err)
	}
	fmt.Printf("Created bookmark: %s (ID: %s)\n", createdBookmark.Title, createdBookmark.ID)
	fmt.Printf("  URL: %s\n", createdBookmark.URL)
	fmt.Printf("  Tags: %d\n\n", len(createdBookmark.TagIDs))

	// Example 4: Get a bookmark by ID
	fmt.Println("=== Getting Bookmark by ID ===")
	fetchedBookmark, err := bookmarkService.Get(ctx, createdBookmark.ID)
	if err != nil {
		log.Fatalf("Failed to get bookmark: %v", err)
	}
	fmt.Printf("Fetched: %s\n", fetchedBookmark.Title)
	fmt.Printf("  Summary: %s\n\n", fetchedBookmark.Summary)

	// Example 5: Add more tags to the bookmark
	fmt.Println("=== Adding Tags to Bookmark ===")
	updatedBookmark, err := bookmarkService.AddTags(ctx, createdBookmark.ID, []string{tag2.ID})
	if err != nil {
		log.Fatalf("Failed to add tags: %v", err)
	}
	fmt.Printf("Updated bookmark now has %d tags\n\n", len(updatedBookmark.TagIDs))

	// Example 6: Update a bookmark
	fmt.Println("=== Updating Bookmark ===")
	fetchedBookmark.Summary = "The best resource for learning Go programming language!"
	updatedBookmark2, err := bookmarkService.Update(ctx, fetchedBookmark.ID, fetchedBookmark)
	if err != nil {
		log.Fatalf("Failed to update bookmark: %v", err)
	}
	fmt.Printf("Updated summary: %s\n\n", updatedBookmark2.Summary)

	// Example 7: List all bookmarks
	fmt.Println("=== Listing All Bookmarks ===")
	allBookmarks, err := bookmarkService.List(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to list bookmarks: %v", err)
	}
	fmt.Printf("Found %d bookmarks:\n", len(allBookmarks))
	for _, bm := range allBookmarks {
		fmt.Printf("- %s (%s)\n", bm.Title, bm.URL)
		fmt.Printf("  Tags: %d, Added: %s\n", len(bm.TagIDs), bm.DateAdded.Format("2006-01-02"))
	}
	fmt.Println()

	// Example 8: Filter bookmarks by tag
	fmt.Println("=== Filtering Bookmarks by Tag ===")
	filteredBookmarks, err := bookmarkService.List(ctx, &bookmarks.Filter{
		HasTag: tag1.ID, // Filter by "Go Programming" tag
		Limit:  10,
	})
	if err != nil {
		log.Fatalf("Failed to filter bookmarks: %v", err)
	}
	fmt.Printf("Found %d bookmarks with 'Go Programming' tag:\n", len(filteredBookmarks))
	for _, bm := range filteredBookmarks {
		fmt.Printf("- %s\n", bm.Title)
	}
	fmt.Println()

	// Example 9: Search bookmarks by title
	fmt.Println("=== Searching Bookmarks by Title ===")
	searchResults, err := bookmarkService.List(ctx, &bookmarks.Filter{
		TitleContains: "Go",
		Limit:         5,
	})
	if err != nil {
		log.Fatalf("Failed to search bookmarks: %v", err)
	}
	fmt.Printf("Found %d bookmarks with 'Go' in title:\n", len(searchResults))
	for _, bm := range searchResults {
		fmt.Printf("- %s\n", bm.Title)
	}
	fmt.Println()

	// Example 10: Query with sorting
	fmt.Println("=== Querying Bookmarks with Sorting ===")
	queryResults, err := bookmarkService.Query(ctx, &bookmarks.QueryOptions{
		SortBy:    bookmarks.SortByDateAdded,
		SortOrder: bookmarks.SortOrderDesc,
		Limit:     5,
	})
	if err != nil {
		log.Fatalf("Failed to query bookmarks: %v", err)
	}
	fmt.Printf("Latest %d bookmarks:\n", len(queryResults))
	for _, bm := range queryResults {
		fmt.Printf("- %s (Added: %s)\n", bm.Title, bm.DateAdded.Format("2006-01-02"))
	}
	fmt.Println()

	// Example 11: Remove tags from bookmark
	fmt.Println("=== Removing Tags from Bookmark ===")
	updatedBookmark3, err := bookmarkService.RemoveTags(ctx, createdBookmark.ID, []string{tag2.ID})
	if err != nil {
		log.Fatalf("Failed to remove tags: %v", err)
	}
	fmt.Printf("Bookmark now has %d tags (removed Web Development)\n\n", len(updatedBookmark3.TagIDs))

	// Example 12: Get tag by name
	fmt.Println("=== Getting Tag by Name ===")
	foundTag, err := tagService.GetByName(ctx, "Go Programming")
	if err != nil {
		log.Fatalf("Failed to get tag by name: %v", err)
	}
	if foundTag != nil {
		fmt.Printf("Found tag: %s (ID: %s)\n\n", foundTag.Name, foundTag.ID)
	}

	fmt.Println("=== All examples completed successfully! ===")
	fmt.Println("\nNote: The bookmark and tags created in this example are now in your Notion database.")
	fmt.Println("You can delete them manually from Notion if you wish.")
}
