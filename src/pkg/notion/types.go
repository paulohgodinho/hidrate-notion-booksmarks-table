package notion

import (
	"time"

	"github.com/jomei/notionapi"
)

// RichTextToString converts Notion rich text to a plain string
func RichTextToString(richText []notionapi.RichText) string {
	if len(richText) == 0 {
		return ""
	}

	result := ""
	for _, rt := range richText {
		result += rt.PlainText
	}
	return result
}

// StringToRichText converts a plain string to Notion rich text
func StringToRichText(text string) []notionapi.RichText {
	if text == "" {
		return []notionapi.RichText{}
	}

	return []notionapi.RichText{
		{
			Type: notionapi.ObjectTypeText,
			Text: &notionapi.Text{
				Content: text,
			},
			PlainText: text,
		},
	}
}

// DateToNotionDate converts a time.Time to a Notion date
func DateToNotionDate(t time.Time) *notionapi.DateObject {
	if t.IsZero() {
		return nil
	}

	date := notionapi.Date(t)
	return &notionapi.DateObject{
		Start: &date,
	}
}

// NotionDateToTime converts a Notion date to time.Time
func NotionDateToTime(date *notionapi.DateObject) time.Time {
	if date == nil || date.Start == nil {
		return time.Time{}
	}

	// notionapi.Date is a time.Time wrapper
	return time.Time(*date.Start)
}

// GetTitleText extracts the plain text from a title property
func GetTitleText(title []notionapi.RichText) string {
	return RichTextToString(title)
}

// GetURLString extracts the URL string from a URL property
func GetURLString(url *string) string {
	if url == nil {
		return ""
	}
	return *url
}

// FileUpload API types (not yet supported in notionapi library v1.13.3)
// These types extend the library to support the File Upload API

// FileUploadMode represents the mode of file upload
type FileUploadMode string

const (
	// FileUploadModeExternalURL imports a file from an external URL
	FileUploadModeExternalURL FileUploadMode = "external_url"
)

// FileUploadStatus represents the status of a file upload
type FileUploadStatus string

const (
	// FileUploadStatusPending indicates the upload is in progress
	FileUploadStatusPending FileUploadStatus = "pending"
	// FileUploadStatusUploaded indicates the upload completed successfully
	FileUploadStatusUploaded FileUploadStatus = "uploaded"
	// FileUploadStatusFailed indicates the upload failed
	FileUploadStatusFailed FileUploadStatus = "failed"
)

// CreateFileUploadRequest is the request body for creating a file upload
type CreateFileUploadRequest struct {
	Mode        FileUploadMode `json:"mode"`
	ExternalURL string         `json:"external_url"`
	Filename    string         `json:"filename"`
}

// FileUploadObject represents a file upload response from the Notion API
type FileUploadObject struct {
	Object         string           `json:"object"` // "file_upload"
	ID             string           `json:"id"`
	CreatedTime    string           `json:"created_time"`
	LastEditedTime string           `json:"last_edited_time"`
	ExpiryTime     string           `json:"expiry_time,omitempty"`
	UploadURL      string           `json:"upload_url,omitempty"`
	Archived       bool             `json:"archived"`
	Status         FileUploadStatus `json:"status"`
	Filename       *string          `json:"filename"`
	ContentType    *string          `json:"content_type"`
	ContentLength  *int             `json:"content_length"`
	RequestID      string           `json:"request_id"`
}
