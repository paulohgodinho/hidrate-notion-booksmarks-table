package bookmarks

import (
	"time"

	"github.com/jomei/notionapi"
	"github.com/pgodinho/hidrate-notion-bookmarks/pkg/notion"
)

// ToBookmark converts a Notion page to a Bookmark
func ToBookmark(page *notionapi.Page) (*Bookmark, error) {
	if page == nil {
		return nil, nil
	}

	bookmark := &Bookmark{
		ID:        string(page.ID),
		CreatedAt: time.Time(page.CreatedTime),
		UpdatedAt: time.Time(page.LastEditedTime),
	}

	// Extract Title (unnamed field)
	if titleProp, ok := page.Properties[""].(*notionapi.TitleProperty); ok {
		bookmark.Title = notion.GetTitleText(titleProp.Title)
	}

	// Extract URL
	if urlProp, ok := page.Properties["URL"].(*notionapi.URLProperty); ok {
		bookmark.URL = urlProp.URL
	}

	// Extract Summary
	if summaryProp, ok := page.Properties["Summary"].(*notionapi.RichTextProperty); ok {
		bookmark.Summary = notion.RichTextToString(summaryProp.RichText)
	}

	// Extract Author
	if authorProp, ok := page.Properties["Author"].(*notionapi.RichTextProperty); ok {
		bookmark.Author = notion.RichTextToString(authorProp.RichText)
	}

	// Extract Image URL
	if imageProp, ok := page.Properties["image"].(*notionapi.URLProperty); ok {
		bookmark.ImageURL = imageProp.URL
	}

	// Extract Date Added
	if dateProp, ok := page.Properties["Date Added"].(*notionapi.DateProperty); ok {
		bookmark.DateAdded = notion.NotionDateToTime(dateProp.Date)
	}

	// Extract Tags relation
	if tagsProp, ok := page.Properties["Tags"].(*notionapi.RelationProperty); ok {
		bookmark.TagIDs = make([]string, len(tagsProp.Relation))
		for i, rel := range tagsProp.Relation {
			bookmark.TagIDs[i] = string(rel.ID)
		}
	}

	// Extract Processed checkbox
	if processedProp, ok := page.Properties["Processed"].(*notionapi.CheckboxProperty); ok {
		bookmark.Processed = processedProp.Checkbox
	}

	// Extract Error
	if errorProp, ok := page.Properties["Error"].(*notionapi.RichTextProperty); ok {
		bookmark.Error = notion.RichTextToString(errorProp.RichText)
	}

	return bookmark, nil
}

// ToNotionProperties converts a Bookmark to Notion page properties
func ToNotionProperties(bookmark *Bookmark) notionapi.Properties {
	props := notionapi.Properties{}

	if bookmark.Title != "" {
		props[""] = notionapi.TitleProperty{
			Title: notion.StringToRichText(bookmark.Title),
		}
	}

	if bookmark.URL != "" {
		props["URL"] = notionapi.URLProperty{
			URL: bookmark.URL,
		}
	}

	if bookmark.Summary != "" {
		props["Summary"] = notionapi.RichTextProperty{
			RichText: notion.StringToRichText(bookmark.Summary),
		}
	}

	if bookmark.Author != "" {
		props["Author"] = notionapi.RichTextProperty{
			RichText: notion.StringToRichText(bookmark.Author),
		}
	}

	if bookmark.ImageURL != "" {
		props["image"] = notionapi.URLProperty{
			URL: bookmark.ImageURL,
		}
	}

	if !bookmark.DateAdded.IsZero() {
		props["Date Added"] = notionapi.DateProperty{
			Date: notion.DateToNotionDate(bookmark.DateAdded),
		}
	}

	if len(bookmark.TagIDs) > 0 {
		relations := make([]notionapi.Relation, len(bookmark.TagIDs))
		for i, tagID := range bookmark.TagIDs {
			relations[i] = notionapi.Relation{
				ID: notionapi.PageID(tagID),
			}
		}
		props["Tags"] = notionapi.RelationProperty{
			Relation: relations,
		}
	}

	// Always set Processed checkbox (defaults to false if not set)
	props["Processed"] = notionapi.CheckboxProperty{
		Checkbox: bookmark.Processed,
	}

	// Set Error field (can be empty string)
	props["Error"] = notionapi.RichTextProperty{
		RichText: notion.StringToRichText(bookmark.Error),
	}

	return props
}
