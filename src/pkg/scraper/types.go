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
	Title       string  `json:"title,omitempty"`
	Description string  `json:"description,omitempty"`
	Author      string  `json:"author,omitempty"`
	SiteName    string  `json:"site_name,omitempty"`
	Language    string  `json:"language,omitempty"`
	Keywords    string  `json:"keywords,omitempty"`
	Copyright   string  `json:"copyright,omitempty"`
	Image       *string `json:"image,omitempty"`

	// OpenGraph metadata
	OGTitle       string  `json:"og_title,omitempty"`
	OGDescription string  `json:"og_description,omitempty"`
	OGImage       *string `json:"og_image,omitempty"`
	OGType        string  `json:"og_type,omitempty"`
	OGURL         string  `json:"og_url,omitempty"`
	OGSiteName    string  `json:"og_site_name,omitempty"`

	// Twitter Card metadata
	TwitterCard        string  `json:"twitter_card,omitempty"`
	TwitterSite        string  `json:"twitter_site,omitempty"`
	TwitterCreator     string  `json:"twitter_creator,omitempty"`
	TwitterTitle       string  `json:"twitter_title,omitempty"`
	TwitterDescription string  `json:"twitter_description,omitempty"`
	TwitterImage       *string `json:"twitter_image,omitempty"`

	// Date metadata
	PublishedDate *time.Time `json:"published_date,omitempty"`
	ModifiedDate  *time.Time `json:"modified_date,omitempty"`

	// Article metadata
	ArticleAuthor     string   `json:"article_author,omitempty"`
	ArticleSection    string   `json:"article_section,omitempty"`
	ArticleTags       []string `json:"article_tags,omitempty"`
	ArticlePublished  string   `json:"article_published,omitempty"`
	ArticleModified   string   `json:"article_modified,omitempty"`
	ArticleExpiryDate string   `json:"article_expiry_date,omitempty"`
	ArticlePublisher  string   `json:"article_publisher,omitempty"`

	// Platform-specific metadata

	// YouTube
	YouTubeChannelID   string `json:"youtube_channel_id,omitempty"`
	YouTubeChannelName string `json:"youtube_channel_name,omitempty"`
	YouTubeVideoID     string `json:"youtube_video_id,omitempty"`
	YouTubeDuration    string `json:"youtube_duration,omitempty"`
	YouTubeViews       string `json:"youtube_views,omitempty"`

	// Twitter/X
	TwitterUsername string `json:"twitter_username,omitempty"`
	TwitterTweetID  string `json:"twitter_tweet_id,omitempty"`

	// Amazon
	AmazonASIN   string `json:"amazon_asin,omitempty"`
	AmazonPrice  string `json:"amazon_price,omitempty"`
	AmazonRating string `json:"amazon_rating,omitempty"`

	// Reddit
	RedditSubreddit string `json:"reddit_subreddit,omitempty"`
	RedditPostID    string `json:"reddit_post_id,omitempty"`
	RedditUpvotes   int    `json:"reddit_upvotes,omitempty"`
	RedditComments  int    `json:"reddit_comments,omitempty"`
	RedditAuthor    string `json:"reddit_author,omitempty"`

	// GitHub
	GitHubRepo   string `json:"github_repo,omitempty"`
	GitHubStars  int    `json:"github_stars,omitempty"`
	GitHubForks  int    `json:"github_forks,omitempty"`
	GitHubIssues int    `json:"github_issues,omitempty"`
}

// HealthResponse represents the response from the health check endpoint
type HealthResponse struct {
	Status  string `json:"status"`
	Version string `json:"version,omitempty"`
	Uptime  string `json:"uptime,omitempty"`
}
