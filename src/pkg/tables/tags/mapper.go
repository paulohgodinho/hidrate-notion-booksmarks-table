package tags

import (
	"time"

	"github.com/jomei/notionapi"
	"github.com/pgodinho/hidrate-notion-bookmarks/pkg/notion"
)

// ToTag converts a Notion page to a Tag
func ToTag(page *notionapi.Page) (*Tag, error) {
	if page == nil {
		return nil, nil
	}

	tag := &Tag{
		ID:        string(page.ID),
		CreatedAt: time.Time(page.CreatedTime),
		UpdatedAt: time.Time(page.LastEditedTime),
	}

	// Extract Name from title property
	if titleProp, ok := page.Properties[PropertyName].(*notionapi.TitleProperty); ok {
		tag.Name = notion.GetTitleText(titleProp.Title)
	}

	return tag, nil
}

// ToNotionProperties converts a Tag to Notion page properties
func ToNotionProperties(tag *Tag) notionapi.Properties {
	props := notionapi.Properties{}

	if tag.Name != "" {
		props[PropertyName] = notionapi.TitleProperty{
			Title: notion.StringToRichText(tag.Name),
		}
	}

	return props
}
