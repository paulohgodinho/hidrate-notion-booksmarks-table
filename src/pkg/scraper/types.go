package scraper

import "time"

// ScrapeRequest represents a request to scrape a URL
type ScrapeRequest struct {
	URL string `json:"url"`
}

// ScrapedContent represents the full response from the webmeatscraper service
type ScrapedContent struct {
	Content  string    `json:"content"`
	Image    *string   `json:"image,omitempty"` // Optional: can be null
	Metadata *Metadata `json:"metadata,omitempty"`
}

// Metadata contains all extracted metadata from the scraped content
type Metadata struct {
	// Common metadata fields
	Author          string     `json:"author,omitempty"`
	Date            *time.Time `json:"date,omitempty"`
	DateModified    *time.Time `json:"dateModified,omitempty"`
	DatePublished   *time.Time `json:"datePublished,omitempty"`
	Description     string     `json:"description,omitempty"`
	Image           *string    `json:"image,omitempty"`
	Lang            *string    `json:"lang,omitempty"`
	Logo            *string    `json:"logo,omitempty"`
	Publisher       string     `json:"publisher,omitempty"`
	RedditAuthor    string     `json:"redditAuthor,omitempty"`
	RedditSubreddit *string    `json:"redditSubreddit,omitempty"`
	RedditUpvotes   *int       `json:"redditUpvotes,omitempty"`
	Subreddit       *string    `json:"subreddit,omitempty"`
	Title           string     `json:"title,omitempty"`
	URL             string     `json:"url,omitempty"`
}

// HealthResponse represents the response from the health check endpoint
type HealthResponse struct {
	Status  string `json:"status"`
	Version string `json:"version,omitempty"`
	Uptime  string `json:"uptime,omitempty"`
}
