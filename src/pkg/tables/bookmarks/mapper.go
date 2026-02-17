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

	// Extract Title (page field)
	if titleProp, ok := page.Properties[PropertyPage].(*notionapi.TitleProperty); ok {
		bookmark.Title = notion.GetTitleText(titleProp.Title)
	}

	// Extract URL
	if urlProp, ok := page.Properties[PropertyURL].(*notionapi.URLProperty); ok {
		bookmark.URL = urlProp.URL
	}

	// Extract Summary
	if summaryProp, ok := page.Properties[PropertySummary].(*notionapi.RichTextProperty); ok {
		bookmark.Summary = notion.RichTextToString(summaryProp.RichText)
	}

	// Extract Author
	if authorProp, ok := page.Properties[PropertyAuthor].(*notionapi.RichTextProperty); ok {
		bookmark.Author = notion.RichTextToString(authorProp.RichText)
	}

	// Extract Image URL
	if imageProp, ok := page.Properties[PropertyImage].(*notionapi.URLProperty); ok {
		bookmark.ImageURL = imageProp.URL
	}

	// Extract Date Added
	if dateProp, ok := page.Properties[PropertyDateAdded].(*notionapi.DateProperty); ok {
		bookmark.DateAdded = notion.NotionDateToTime(dateProp.Date)
	}

	// Extract Date Processed
	if dateProcProp, ok := page.Properties[PropertyDateProcessed].(*notionapi.DateProperty); ok {
		bookmark.DateProcessed = notion.NotionDateToTime(dateProcProp.Date)
	}

	// Extract Date Published
	if datePubProp, ok := page.Properties[PropertyDatePublished].(*notionapi.RichTextProperty); ok {
		bookmark.DatePublished = notion.RichTextToString(datePubProp.RichText)
	}

	// Extract Tags relation
	if tagsProp, ok := page.Properties[PropertyTag].(*notionapi.RelationProperty); ok {
		bookmark.TagIDs = make([]string, len(tagsProp.Relation))
		for i, rel := range tagsProp.Relation {
			bookmark.TagIDs[i] = string(rel.ID)
		}
	}

	// Extract Manual Lists relation
	if manualProp, ok := page.Properties[PropertyManualLists].(*notionapi.RelationProperty); ok {
		bookmark.ManualListIDs = make([]string, len(manualProp.Relation))
		for i, rel := range manualProp.Relation {
			bookmark.ManualListIDs[i] = string(rel.ID)
		}
	}

	// Extract Smart Lists relation
	if smartProp, ok := page.Properties[PropertySmartLists].(*notionapi.RelationProperty); ok {
		bookmark.SmartListIDs = make([]string, len(smartProp.Relation))
		for i, rel := range smartProp.Relation {
			bookmark.SmartListIDs[i] = string(rel.ID)
		}
	}

	// Extract Processed checkbox
	if processedProp, ok := page.Properties[PropertyProcessed].(*notionapi.CheckboxProperty); ok {
		bookmark.Processed = processedProp.Checkbox
	}

	// Extract Error
	if errorProp, ok := page.Properties[PropertyError].(*notionapi.RichTextProperty); ok {
		bookmark.Error = notion.RichTextToString(errorProp.RichText)
	}

	return bookmark, nil
}

// ToNotionProperties converts a Bookmark to Notion page properties
func ToNotionProperties(bookmark *Bookmark) notionapi.Properties {
	props := notionapi.Properties{}

	if bookmark.Title != "" {
		props[PropertyPage] = notionapi.TitleProperty{
			Title: notion.StringToRichText(bookmark.Title),
		}
	}

	if bookmark.URL != "" {
		props[PropertyURL] = notionapi.URLProperty{
			URL: bookmark.URL,
		}
	}

	if bookmark.Summary != "" {
		props[PropertySummary] = notionapi.RichTextProperty{
			RichText: notion.StringToRichText(bookmark.Summary),
		}
	}

	if bookmark.Author != "" {
		props[PropertyAuthor] = notionapi.RichTextProperty{
			RichText: notion.StringToRichText(bookmark.Author),
		}
	}

	if bookmark.ImageURL != "" {
		props[PropertyImage] = notionapi.URLProperty{
			URL: bookmark.ImageURL,
		}
	}

	if !bookmark.DateAdded.IsZero() {
		props[PropertyDateAdded] = notionapi.DateProperty{
			Date: notion.DateToNotionDate(bookmark.DateAdded),
		}
	}

	if !bookmark.DateProcessed.IsZero() {
		props[PropertyDateProcessed] = notionapi.DateProperty{
			Date: notion.DateToNotionDate(bookmark.DateProcessed),
		}
	}

	if bookmark.DatePublished != "" {
		props[PropertyDatePublished] = notionapi.RichTextProperty{
			RichText: notion.StringToRichText(bookmark.DatePublished),
		}
	}

	if len(bookmark.TagIDs) > 0 {
		relations := make([]notionapi.Relation, len(bookmark.TagIDs))
		for i, tagID := range bookmark.TagIDs {
			relations[i] = notionapi.Relation{
				ID: notionapi.PageID(tagID),
			}
		}
		props[PropertyTag] = notionapi.RelationProperty{
			Relation: relations,
		}
	}

	if len(bookmark.ManualListIDs) > 0 {
		relations := make([]notionapi.Relation, len(bookmark.ManualListIDs))
		for i, id := range bookmark.ManualListIDs {
			relations[i] = notionapi.Relation{
				ID: notionapi.PageID(id),
			}
		}
		props[PropertyManualLists] = notionapi.RelationProperty{
			Relation: relations,
		}
	}

	if len(bookmark.SmartListIDs) > 0 {
		relations := make([]notionapi.Relation, len(bookmark.SmartListIDs))
		for i, id := range bookmark.SmartListIDs {
			relations[i] = notionapi.Relation{
				ID: notionapi.PageID(id),
			}
		}
		props[PropertySmartLists] = notionapi.RelationProperty{
			Relation: relations,
		}
	}

	// Always set Processed checkbox (defaults to false if not set)
	props[PropertyProcessed] = notionapi.CheckboxProperty{
		Checkbox: bookmark.Processed,
	}

	// Set Error field (can be empty string)
	props[PropertyError] = notionapi.RichTextProperty{
		RichText: notion.StringToRichText(bookmark.Error),
	}

	return props
}
