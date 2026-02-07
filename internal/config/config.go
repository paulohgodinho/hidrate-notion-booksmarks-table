package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config holds the application configuration
type Config struct {
	NotionAPIKey  string
	BookmarksDBID string
	TagsDBID      string
}

// Load reads configuration from environment variables
// It first attempts to load from a .env file, then reads from environment
func Load() (*Config, error) {
	// Load .env file if it exists (ignore error if file doesn't exist)
	_ = godotenv.Load()

	cfg := &Config{
		NotionAPIKey:  os.Getenv("NOTION_API_KEY"),
		BookmarksDBID: os.Getenv("NOTION_BOOKMARKS_DB_ID"),
		TagsDBID:      os.Getenv("NOTION_TAGS_DB_ID"),
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Validate checks that all required configuration values are present
func (c *Config) Validate() error {
	if c.NotionAPIKey == "" {
		return fmt.Errorf("NOTION_API_KEY is required")
	}
	if c.BookmarksDBID == "" {
		return fmt.Errorf("NOTION_BOOKMARKS_DB_ID is required")
	}
	if c.TagsDBID == "" {
		return fmt.Errorf("NOTION_TAGS_DB_ID is required")
	}
	return nil
}
