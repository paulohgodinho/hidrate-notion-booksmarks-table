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
