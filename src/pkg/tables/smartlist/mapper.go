package smartlist

import (
	"time"

	"github.com/jomei/notionapi"
	"github.com/pgodinho/hidrate-notion-bookmarks/pkg/notion"
)

// ToSmartListItem converts a Notion page to a SmartListItem
func ToSmartListItem(page *notionapi.Page) (*SmartListItem, error) {
	if page == nil {
		return nil, nil
	}

	item := &SmartListItem{
		ID:        string(page.ID),
		CreatedAt: time.Time(page.CreatedTime),
		UpdatedAt: time.Time(page.LastEditedTime),
	}

	// Extract Name from title property
	if titleProp, ok := page.Properties[PropertyName].(*notionapi.TitleProperty); ok {
		item.Name = notion.GetTitleText(titleProp.Title)
	}

	return item, nil
}

// ToNotionProperties converts a SmartListItem to Notion page properties
func ToNotionProperties(item *SmartListItem) notionapi.Properties {
	props := notionapi.Properties{}

	if item.Name != "" {
		props[PropertyName] = notionapi.TitleProperty{
			Title: notion.StringToRichText(item.Name),
		}
	}

	return props
}
