package manuallist

import (
	"time"

	"github.com/jomei/notionapi"
	"github.com/pgodinho/hidrate-notion-bookmarks/pkg/notion"
)

// ToManualListItem converts a Notion page to a ManualListItem
func ToManualListItem(page *notionapi.Page) (*ManualListItem, error) {
	if page == nil {
		return nil, nil
	}

	item := &ManualListItem{
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

// ToNotionProperties converts a ManualListItem to Notion page properties
func ToNotionProperties(item *ManualListItem) notionapi.Properties {
	props := notionapi.Properties{}

	if item.Name != "" {
		props[PropertyName] = notionapi.TitleProperty{
			Title: notion.StringToRichText(item.Name),
		}
	}

	return props
}
