package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/pgodinho/hidrate-notion-bookmarks/pkg/bookmarks"
	"github.com/pgodinho/hidrate-notion-bookmarks/pkg/notion"
	"github.com/pgodinho/hidrate-notion-bookmarks/pkg/scraper"
)

func main() {
	fmt.Println("=== Notion Bookmark Processor ===")
	fmt.Println()

	// Load configuration
	cfg, err := Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize clients
	notionClient := notion.NewClient(cfg.NotionAPIKey, cfg.BookmarksDBID, cfg.TagsDBID)
	bookmarkService := bookmarks.NewService(notionClient)

	// Default to localhost if not set in config
	scraperURL := cfg.WebmeatscraperURL
	if scraperURL == "" {
		scraperURL = "http://localhost:7878"
	}
	scraperClient := scraper.NewClient(scraperURL)

	// Initialize image uploader if enabled
	var imageUploader *notion.ImageUploader
	if cfg.UploadImagesToNotion {
		imageUploader = notion.NewImageUploader(
			cfg.NotionAPIKey,
			cfg.ImageUploadTimeout,
			cfg.ImageUploadPollInterval,
		)
		fmt.Println("✓ Image upload to Notion: ENABLED")
	} else {
		fmt.Println("  Image upload to Notion: disabled")
	}

	// Show debug mode status
	if cfg.Debug {
		fmt.Println("✓ Debug mode: ENABLED (full JSON output)")
	} else {
		fmt.Println("  Debug mode: disabled")
	}
	fmt.Println()

	ctx := context.Background()

	// Check scraper service health
	fmt.Printf("Checking scraper service at %s...\n", scraperClient.BaseURL())
	health, err := scraperClient.Health(ctx)
	if err != nil {
		log.Fatalf("Scraper service is not available: %v", err)
	}
	fmt.Printf("✓ Scraper service is healthy (status: %s)\n", health.Status)
	fmt.Println()

	// Fetch ALL unprocessed bookmarks
	fmt.Println("Fetching unprocessed bookmarks...")
	unprocessed, err := bookmarkService.GetUnprocessed(ctx, 0) // 0 = get all
	if err != nil {
		log.Fatalf("Failed to fetch unprocessed bookmarks: %v", err)
	}

	if len(unprocessed) == 0 {
		fmt.Println("No unprocessed bookmarks found.")
		fmt.Println()
		fmt.Println("=== Processing Complete ===")
		os.Exit(0)
	}

	fmt.Printf("Found %d unprocessed bookmark(s)\n", len(unprocessed))
	fmt.Println()

	// Process each bookmark
	successCount := 0
	errorCount := 0

	for i, bookmark := range unprocessed {
		fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
		fmt.Printf("Processing bookmark %d of %d\n", i+1, len(unprocessed))
		fmt.Printf("Title: %s\n", bookmark.Title)
		fmt.Printf("ID: %s\n", bookmark.ID)
		fmt.Printf("URL: %s\n", bookmark.URL)
		fmt.Println()

		// Scrape the bookmark
		fmt.Println("Scraping content...")
		result, err := scraperClient.Scrape(ctx, bookmark.URL)
		if err != nil {
			// On error: Set error field and mark as not processed
			errorMsg := fmt.Sprintf("Failed to scrape URL: %v", err)
			fmt.Printf("✗ %s\n", errorMsg)
			fmt.Println()
			fmt.Println("Updating bookmark with error...")

			_, updateErr := bookmarkService.SetError(ctx, bookmark.ID, errorMsg)
			if updateErr != nil {
				log.Printf("Failed to update bookmark with error: %v", updateErr)
				errorCount++
				fmt.Println()
				continue
			}

			fmt.Println("✓ Bookmark marked with error")
			fmt.Println()
			errorCount++
			continue
		}

		// Print full raw JSON response (only if DEBUG is enabled)
		if cfg.Debug {
			fmt.Println("=== SCRAPED CONTENT (Full JSON) ===")
			// Pretty print the raw JSON
			var prettyJSON bytes.Buffer
			if err := json.Indent(&prettyJSON, []byte(result.RawJSON), "", "  "); err == nil {
				fmt.Println(prettyJSON.String())
			} else {
				// Fallback to raw JSON if indentation fails
				fmt.Println(result.RawJSON)
			}
			fmt.Println("=== END OF SCRAPED CONTENT ===")
			fmt.Println()
		}

		// Use the parsed content
		content := result.Content

		// Update bookmark with scraped metadata
		fmt.Println("Updating bookmark with scraped metadata...")
		updated := false

		// Update Author if empty and available
		if bookmark.Author == "" && content.Metadata != nil && content.Metadata.Author != "" {
			bookmark.Author = content.Metadata.Author
			fmt.Printf("  ✓ Set author: %s\n", content.Metadata.Author)
			updated = true
		}

		// Extract image URL from scraped content (try multiple sources)
		var imageURL string

		// Priority order: content.Image → OGImage → TwitterImage → metadata.Image
		if content.Image != nil && *content.Image != "" {
			imageURL = *content.Image
		} else if content.Metadata != nil {
			// Try OpenGraph image
			if content.Metadata.OGImage != nil && *content.Metadata.OGImage != "" {
				imageURL = *content.Metadata.OGImage
			} else if content.Metadata.TwitterImage != nil && *content.Metadata.TwitterImage != "" {
				// Try Twitter Card image
				imageURL = *content.Metadata.TwitterImage
			} else if content.Metadata.Image != nil && *content.Metadata.Image != "" {
				// Try generic metadata image
				imageURL = *content.Metadata.Image
			}
		}

		// Upload image to Notion and set as page cover (if enabled)
		var coverFileUploadID string
		if imageURL != "" && cfg.UploadImagesToNotion && imageUploader != nil {
			fmt.Printf("  📤 Uploading image to Notion...")

			fileUploadID, err := imageUploader.UploadImageFromURL(ctx, imageURL)
			if err != nil {
				if cfg.FallbackToExternalURL {
					fmt.Printf(" ⚠️  Upload failed (%v), using external URL\n", err)
					// Keep imageURL as-is for database property
				} else {
					fmt.Printf(" ❌ Upload failed: %v\n", err)
					imageURL = "" // Don't set any image
				}
			} else {
				fmt.Printf(" ✅ Uploaded (ID: %s)\n", fileUploadID)
				coverFileUploadID = fileUploadID
			}
		}

		// Update ImageURL property if empty and available
		if bookmark.ImageURL == "" && imageURL != "" {
			bookmark.ImageURL = imageURL
			fmt.Printf("  ✓ Set image property: %s\n", imageURL)
			updated = true
		}

		// Mark as processed and clear error
		bookmark.Processed = true
		bookmark.Error = ""

		// Update the bookmark in Notion
		_, err = bookmarkService.Update(ctx, bookmark.ID, bookmark)
		if err != nil {
			log.Printf("Failed to update bookmark: %v", err)
			errorCount++
			fmt.Println()
			continue
		}

		if updated {
			fmt.Println("✓ Bookmark metadata updated")
		} else {
			fmt.Println("  (No metadata property updates needed)")
		}

		// Set page cover if we have a FileUpload ID (always update cover)
		if coverFileUploadID != "" {
			fmt.Printf("  🖼️  Setting page cover...")
			err = notion.SetPageCover(ctx, cfg.NotionAPIKey, string(bookmark.ID), coverFileUploadID)
			if err != nil {
				fmt.Printf(" ⚠️  Failed to set cover: %v\n", err)
			} else {
				fmt.Printf(" ✅ Cover set\n")
			}
		}

		fmt.Println("✓ Bookmark marked as processed")
		fmt.Println()
		successCount++
	}

	// Print summary
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	fmt.Println("=== Processing Complete ===")
	fmt.Printf("Total: %d bookmarks\n", len(unprocessed))
	fmt.Printf("✓ Successfully processed: %d\n", successCount)
	fmt.Printf("✗ Failed: %d\n", errorCount)
}
