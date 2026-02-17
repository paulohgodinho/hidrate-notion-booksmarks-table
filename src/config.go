package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds the application configuration
type Config struct {
	NotionAPIKey      string
	BookmarksDBID     string
	TagsDBID          string
	ManualListDBID    string
	SmartListDBID     string
	WebmeatscraperURL string

	// Image upload configuration
	UploadImagesToNotion    bool
	ImageUploadTimeout      time.Duration
	ImageUploadPollInterval time.Duration
	FallbackToExternalURL   bool

	// Debug configuration
	Debug bool
}

// Load reads configuration from environment variables
// It first attempts to load from a .env file, then reads from environment
func Load() (*Config, error) {
	// Load .env file if it exists (ignore error if file doesn't exist)
	_ = godotenv.Load()

	cfg := &Config{
		NotionAPIKey:      os.Getenv("NOTION_API_KEY"),
		BookmarksDBID:     os.Getenv("NOTION_BOOKMARKS_DB_ID"),
		TagsDBID:          os.Getenv("NOTION_TAGS_DB_ID"),
		ManualListDBID:    os.Getenv("NOTION_MANUALLIST_DB_ID"),
		SmartListDBID:     os.Getenv("NOTION_SMARTLIST_DB_ID"),
		WebmeatscraperURL: os.Getenv("WEBMEATSCRAPER_URL"),

		// Parse image upload settings with defaults
		UploadImagesToNotion:    parseBoolWithDefault(os.Getenv("UPLOAD_IMAGES_TO_NOTION"), true),
		ImageUploadTimeout:      parseDurationWithDefault(os.Getenv("IMAGE_UPLOAD_TIMEOUT"), 30*time.Second),
		ImageUploadPollInterval: parseDurationWithDefault(os.Getenv("IMAGE_UPLOAD_POLL_INTERVAL"), 3*time.Second),
		FallbackToExternalURL:   parseBoolWithDefault(os.Getenv("FALLBACK_TO_EXTERNAL_URL"), true),

		// Parse debug settings with defaults
		Debug: parseBoolWithDefault(os.Getenv("DEBUG"), false),
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
	if c.ManualListDBID == "" {
		return fmt.Errorf("NOTION_MANUALLIST_DB_ID is required")
	}
	if c.SmartListDBID == "" {
		return fmt.Errorf("NOTION_SMARTLIST_DB_ID is required")
	}
	// WebmeatscraperURL is optional - will default to localhost:7878 if not set
	return nil
}

// parseBoolWithDefault parses a boolean string with a default value
func parseBoolWithDefault(value string, defaultVal bool) bool {
	if value == "" {
		return defaultVal
	}
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return defaultVal
	}
	return parsed
}

// parseDurationWithDefault parses a duration string with a default value
func parseDurationWithDefault(value string, defaultVal time.Duration) time.Duration {
	if value == "" {
		return defaultVal
	}
	parsed, err := time.ParseDuration(value)
	if err != nil {
		return defaultVal
	}
	return parsed
}
