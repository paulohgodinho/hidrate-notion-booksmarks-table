package notion

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"
)

// ImageUploader handles uploading images to Notion storage
type ImageUploader struct {
	notionToken  string
	httpClient   *http.Client
	timeout      time.Duration
	pollInterval time.Duration
}

// NewImageUploader creates a new image uploader with the specified configuration
func NewImageUploader(token string, timeout time.Duration, pollInterval time.Duration) *ImageUploader {
	return &ImageUploader{
		notionToken:  token,
		httpClient:   &http.Client{Timeout: 10 * time.Second},
		timeout:      timeout,
		pollInterval: pollInterval,
	}
}

// UploadImageFromURL uploads an image from an external URL to Notion storage
// Uses the Indirect Import method (mode: "external_url")
// Returns fileUploadID on success, error on failure
func (u *ImageUploader) UploadImageFromURL(ctx context.Context, imageURL string) (string, error) {
	// Extract filename from URL
	filename := extractFilenameFromURL(imageURL)

	// Step 1: Create FileUpload object
	fileUpload, err := u.createFileUpload(ctx, imageURL, filename)
	if err != nil {
		return "", fmt.Errorf("failed to create file upload: %w", err)
	}

	// Step 2: Poll for completion (with timeout)
	fileUploadID, err := u.pollForCompletion(ctx, fileUpload.ID)
	if err != nil {
		return "", fmt.Errorf("upload failed or timed out: %w", err)
	}

	return fileUploadID, nil
}

// createFileUpload creates a new file upload request in Notion
func (u *ImageUploader) createFileUpload(ctx context.Context, externalURL, filename string) (*FileUploadObject, error) {
	reqBody := CreateFileUploadRequest{
		Mode:        FileUploadModeExternalURL,
		ExternalURL: externalURL,
		Filename:    filename,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.notion.com/v1/file_uploads",
		bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+u.notionToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Notion-Version", "2022-06-28")

	resp, err := u.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
	}

	var fileUpload FileUploadObject
	if err := json.NewDecoder(resp.Body).Decode(&fileUpload); err != nil {
		return nil, err
	}

	return &fileUpload, nil
}

// retrieveFileUpload retrieves the current state of a file upload
func (u *ImageUploader) retrieveFileUpload(ctx context.Context, fileUploadID string) (*FileUploadObject, error) {
	req, err := http.NewRequestWithContext(ctx, "GET",
		fmt.Sprintf("https://api.notion.com/v1/file_uploads/%s", fileUploadID), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+u.notionToken)
	req.Header.Set("Notion-Version", "2022-06-28")

	resp, err := u.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
	}

	var fileUpload FileUploadObject
	if err := json.NewDecoder(resp.Body).Decode(&fileUpload); err != nil {
		return nil, err
	}

	return &fileUpload, nil
}

// pollForCompletion polls the file upload status until it completes or times out
func (u *ImageUploader) pollForCompletion(ctx context.Context, fileUploadID string) (string, error) {
	deadline := time.Now().Add(u.timeout)

	for {
		// Check if we've exceeded timeout
		if time.Now().After(deadline) {
			return "", fmt.Errorf("upload timed out after %v", u.timeout)
		}

		// Check context cancellation
		if ctx.Err() != nil {
			return "", ctx.Err()
		}

		// Retrieve FileUpload status
		fileUpload, err := u.retrieveFileUpload(ctx, fileUploadID)
		if err != nil {
			return "", err
		}

		switch fileUpload.Status {
		case FileUploadStatusUploaded:
			return fileUpload.ID, nil
		case FileUploadStatusFailed:
			return "", fmt.Errorf("file upload failed")
		case FileUploadStatusPending:
			// Continue polling
			time.Sleep(u.pollInterval)
		default:
			return "", fmt.Errorf("unknown status: %s", fileUpload.Status)
		}
	}
}

// extractFilenameFromURL extracts a filename from a URL
func extractFilenameFromURL(imageURL string) string {
	parsedURL, err := url.Parse(imageURL)
	if err != nil {
		return "image.png" // Default fallback
	}

	filename := path.Base(parsedURL.Path)

	// If no extension or invalid filename, add .png as default
	if filename == "" || filename == "." || filename == "/" || !strings.Contains(filename, ".") {
		filename = "image.png"
	}

	// Ensure filename is not too long (max 900 bytes recommended by Notion)
	if len(filename) > 100 {
		ext := path.Ext(filename)
		filename = filename[:100-len(ext)] + ext
	}

	return filename
}

// SetPageCover sets the cover of a Notion page using a FileUpload ID
func SetPageCover(ctx context.Context, notionToken string, pageID string, fileUploadID string) error {
	// Build raw JSON since library doesn't support file_upload type yet
	updatePayload := map[string]interface{}{
		"cover": map[string]interface{}{
			"type": "file_upload",
			"file_upload": map[string]string{
				"id": fileUploadID,
			},
		},
	}

	bodyBytes, err := json.Marshal(updatePayload)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "PATCH",
		fmt.Sprintf("https://api.notion.com/v1/pages/%s", pageID),
		bytes.NewReader(bodyBytes))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+notionToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Notion-Version", "2022-06-28")

	httpClient := &http.Client{Timeout: 10 * time.Second}
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
	}

	return nil
}
