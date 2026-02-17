package manuallist

import (
	"context"
	"fmt"

	"github.com/jomei/notionapi"
	"github.com/pgodinho/hidrate-notion-bookmarks/pkg/notion"
)

// Service handles CRUD operations for manual list items
type Service struct {
	client *notion.Client
}

// NewService creates a new manual list service
func NewService(client *notion.Client) *Service {
	return &Service{
		client: client,
	}
}

// Create creates a new manual list item in Notion
func (s *Service) Create(ctx context.Context, item *ManualListItem) (*ManualListItem, error) {
	if item.Name == "" {
		return nil, fmt.Errorf("item name is required")
	}

	props := ToNotionProperties(item)

	req := &notionapi.PageCreateRequest{
		Parent: notionapi.Parent{
			Type:       notionapi.ParentTypeDatabaseID,
			DatabaseID: s.client.ManualListDB(),
		},
		Properties: props,
	}

	page, err := s.client.API().Page.Create(ctx, req)
	if err != nil {
		return nil, notion.NewError("create manual list item", err, fmt.Sprintf("failed to create item: %s", item.Name))
	}

	return ToManualListItem(page)
}

// Get retrieves a manual list item by its ID
func (s *Service) Get(ctx context.Context, id string) (*ManualListItem, error) {
	if id == "" {
		return nil, fmt.Errorf("item ID is required")
	}

	page, err := s.client.API().Page.Get(ctx, notionapi.PageID(id))
	if err != nil {
		return nil, notion.NewError("get manual list item", err, fmt.Sprintf("failed to get item with ID: %s", id))
	}

	return ToManualListItem(page)
}

// GetByName retrieves a manual list item by its name
func (s *Service) GetByName(ctx context.Context, name string) (*ManualListItem, error) {
	if name == "" {
		return nil, fmt.Errorf("item name is required")
	}

	// Query the database for an item with the given name
	query := &notionapi.DatabaseQueryRequest{
		Filter: &notionapi.PropertyFilter{
			Property: PropertyName,
			RichText: &notionapi.TextFilterCondition{
				Equals: name,
			},
		},
		PageSize: 1,
	}

	resp, err := s.client.API().Database.Query(ctx, s.client.ManualListDB(), query)
	if err != nil {
		return nil, notion.NewError("get manual list item by name", err, fmt.Sprintf("failed to query item: %s", name))
	}

	if len(resp.Results) == 0 {
		return nil, nil // Item not found
	}

	return ToManualListItem(&resp.Results[0])
}

// Update updates an existing manual list item
func (s *Service) Update(ctx context.Context, id string, item *ManualListItem) (*ManualListItem, error) {
	if id == "" {
		return nil, fmt.Errorf("item ID is required")
	}

	props := ToNotionProperties(item)

	req := &notionapi.PageUpdateRequest{
		Properties: props,
	}

	page, err := s.client.API().Page.Update(ctx, notionapi.PageID(id), req)
	if err != nil {
		return nil, notion.NewError("update manual list item", err, fmt.Sprintf("failed to update item with ID: %s", id))
	}

	return ToManualListItem(page)
}

// Delete archives a manual list item (Notion doesn't support true deletion)
func (s *Service) Delete(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("item ID is required")
	}

	req := &notionapi.PageUpdateRequest{
		Archived: true,
	}

	_, err := s.client.API().Page.Update(ctx, notionapi.PageID(id), req)
	if err != nil {
		return notion.NewError("delete manual list item", err, fmt.Sprintf("failed to delete item with ID: %s", id))
	}

	return nil
}

// List retrieves all manual list items with optional filtering
func (s *Service) List(ctx context.Context, filter *Filter) ([]*ManualListItem, error) {
	query := &notionapi.DatabaseQueryRequest{}

	// Apply filters
	if filter != nil {
		if filter.NameContains != "" {
			query.Filter = &notionapi.PropertyFilter{
				Property: PropertyName,
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
			Property:  PropertyName,
			Direction: notionapi.SortOrderASC,
		},
	}

	resp, err := s.client.API().Database.Query(ctx, s.client.ManualListDB(), query)
	if err != nil {
		return nil, notion.NewError("list manual list items", err, "failed to list items")
	}

	items := make([]*ManualListItem, 0, len(resp.Results))
	for _, page := range resp.Results {
		item, err := ToManualListItem(&page)
		if err != nil {
			continue // Skip invalid items
		}
		items = append(items, item)
	}

	return items, nil
}

// FindOrCreate finds a manual list item by name or creates it if it doesn't exist
func (s *Service) FindOrCreate(ctx context.Context, name string) (*ManualListItem, error) {
	// Try to find existing item
	item, err := s.GetByName(ctx, name)
	if err != nil {
		return nil, err
	}

	// If found, return it
	if item != nil {
		return item, nil
	}

	// Create new item
	newItem := &ManualListItem{
		Name: name,
	}

	return s.Create(ctx, newItem)
}
