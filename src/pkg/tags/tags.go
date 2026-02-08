package tags

import (
	"context"
	"fmt"

	"github.com/jomei/notionapi"
	"github.com/pgodinho/hidrate-notion-bookmarks/pkg/notion"
)

// Service handles CRUD operations for tags
type Service struct {
	client *notion.Client
}

// NewService creates a new tags service
func NewService(client *notion.Client) *Service {
	return &Service{
		client: client,
	}
}

// Create creates a new tag in Notion
func (s *Service) Create(ctx context.Context, tag *Tag) (*Tag, error) {
	if tag.Name == "" {
		return nil, fmt.Errorf("tag name is required")
	}

	props := ToNotionProperties(tag)

	req := &notionapi.PageCreateRequest{
		Parent: notionapi.Parent{
			Type:       notionapi.ParentTypeDatabaseID,
			DatabaseID: s.client.TagsDB(),
		},
		Properties: props,
	}

	page, err := s.client.API().Page.Create(ctx, req)
	if err != nil {
		return nil, notion.NewError("create tag", err, fmt.Sprintf("failed to create tag: %s", tag.Name))
	}

	return ToTag(page)
}

// Get retrieves a tag by its ID
func (s *Service) Get(ctx context.Context, id string) (*Tag, error) {
	if id == "" {
		return nil, fmt.Errorf("tag ID is required")
	}

	page, err := s.client.API().Page.Get(ctx, notionapi.PageID(id))
	if err != nil {
		return nil, notion.NewError("get tag", err, fmt.Sprintf("failed to get tag with ID: %s", id))
	}

	return ToTag(page)
}

// GetByName retrieves a tag by its name
func (s *Service) GetByName(ctx context.Context, name string) (*Tag, error) {
	if name == "" {
		return nil, fmt.Errorf("tag name is required")
	}

	// Query the database for a tag with the given name
	query := &notionapi.DatabaseQueryRequest{
		Filter: &notionapi.PropertyFilter{
			Property: "Name",
			RichText: &notionapi.TextFilterCondition{
				Equals: name,
			},
		},
		PageSize: 1,
	}

	resp, err := s.client.API().Database.Query(ctx, s.client.TagsDB(), query)
	if err != nil {
		return nil, notion.NewError("get tag by name", err, fmt.Sprintf("failed to query tag: %s", name))
	}

	if len(resp.Results) == 0 {
		return nil, nil // Tag not found
	}

	return ToTag(&resp.Results[0])
}

// Update updates an existing tag
func (s *Service) Update(ctx context.Context, id string, tag *Tag) (*Tag, error) {
	if id == "" {
		return nil, fmt.Errorf("tag ID is required")
	}

	props := ToNotionProperties(tag)

	req := &notionapi.PageUpdateRequest{
		Properties: props,
	}

	page, err := s.client.API().Page.Update(ctx, notionapi.PageID(id), req)
	if err != nil {
		return nil, notion.NewError("update tag", err, fmt.Sprintf("failed to update tag with ID: %s", id))
	}

	return ToTag(page)
}

// Delete archives a tag (Notion doesn't support true deletion)
func (s *Service) Delete(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("tag ID is required")
	}

	req := &notionapi.PageUpdateRequest{
		Archived: true,
	}

	_, err := s.client.API().Page.Update(ctx, notionapi.PageID(id), req)
	if err != nil {
		return notion.NewError("delete tag", err, fmt.Sprintf("failed to delete tag with ID: %s", id))
	}

	return nil
}

// List retrieves all tags with optional filtering
func (s *Service) List(ctx context.Context, filter *Filter) ([]*Tag, error) {
	query := &notionapi.DatabaseQueryRequest{}

	// Apply filters
	if filter != nil {
		if filter.NameContains != "" {
			query.Filter = &notionapi.PropertyFilter{
				Property: "Name",
				RichText: &notionapi.TextFilterCondition{
					Contains: filter.NameContains,
				},
			}
		}

		if filter.Limit > 0 {
			query.PageSize = filter.Limit
		}
	}

	// Add sorting by name
	query.Sorts = []notionapi.SortObject{
		{
			Property:  "Name",
			Direction: notionapi.SortOrderASC,
		},
	}

	resp, err := s.client.API().Database.Query(ctx, s.client.TagsDB(), query)
	if err != nil {
		return nil, notion.NewError("list tags", err, "failed to list tags")
	}

	tags := make([]*Tag, 0, len(resp.Results))
	for _, page := range resp.Results {
		tag, err := ToTag(&page)
		if err != nil {
			continue // Skip invalid tags
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

// FindOrCreate finds a tag by name or creates it if it doesn't exist
func (s *Service) FindOrCreate(ctx context.Context, name string) (*Tag, error) {
	// Try to find existing tag
	tag, err := s.GetByName(ctx, name)
	if err != nil {
		return nil, err
	}

	// If found, return it
	if tag != nil {
		return tag, nil
	}

	// Create new tag
	newTag := &Tag{
		Name: name,
	}

	return s.Create(ctx, newTag)
}
