package bookmarks

import (
	"context"
	"fmt"

	"github.com/jomei/notionapi"
	"github.com/pgodinho/hidrate-notion-bookmarks/pkg/notion"
)

// Service handles CRUD operations for bookmarks
type Service struct {
	client *notion.Client
}

// NewService creates a new bookmarks service
func NewService(client *notion.Client) *Service {
	return &Service{
		client: client,
	}
}

// Create creates a new bookmark in Notion
func (s *Service) Create(ctx context.Context, bookmark *Bookmark) (*Bookmark, error) {
	if bookmark.Title == "" {
		return nil, fmt.Errorf("bookmark title is required")
	}

	props := ToNotionProperties(bookmark)

	req := &notionapi.PageCreateRequest{
		Parent: notionapi.Parent{
			Type:       notionapi.ParentTypeDatabaseID,
			DatabaseID: s.client.BookmarksDB(),
		},
		Properties: props,
	}

	page, err := s.client.API().Page.Create(ctx, req)
	if err != nil {
		return nil, notion.NewError("create bookmark", err, fmt.Sprintf("failed to create bookmark: %s", bookmark.Title))
	}

	return ToBookmark(page)
}

// Get retrieves a bookmark by its ID
func (s *Service) Get(ctx context.Context, id string) (*Bookmark, error) {
	if id == "" {
		return nil, fmt.Errorf("bookmark ID is required")
	}

	page, err := s.client.API().Page.Get(ctx, notionapi.PageID(id))
	if err != nil {
		return nil, notion.NewError("get bookmark", err, fmt.Sprintf("failed to get bookmark with ID: %s", id))
	}

	return ToBookmark(page)
}

// Update updates an existing bookmark
func (s *Service) Update(ctx context.Context, id string, bookmark *Bookmark) (*Bookmark, error) {
	if id == "" {
		return nil, fmt.Errorf("bookmark ID is required")
	}

	props := ToNotionProperties(bookmark)

	req := &notionapi.PageUpdateRequest{
		Properties: props,
	}

	page, err := s.client.API().Page.Update(ctx, notionapi.PageID(id), req)
	if err != nil {
		return nil, notion.NewError("update bookmark", err, fmt.Sprintf("failed to update bookmark with ID: %s", id))
	}

	return ToBookmark(page)
}

// Delete archives a bookmark (Notion doesn't support true deletion)
func (s *Service) Delete(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("bookmark ID is required")
	}

	req := &notionapi.PageUpdateRequest{
		Archived: true,
	}

	_, err := s.client.API().Page.Update(ctx, notionapi.PageID(id), req)
	if err != nil {
		return notion.NewError("delete bookmark", err, fmt.Sprintf("failed to delete bookmark with ID: %s", id))
	}

	return nil
}

// List retrieves all bookmarks with optional filtering
func (s *Service) List(ctx context.Context, filter *Filter) ([]*Bookmark, error) {
	query := &notionapi.DatabaseQueryRequest{}

	// Apply filters
	if filter != nil {
		var filters []notionapi.PropertyFilter

		if filter.TitleContains != "" {
			filters = append(filters, notionapi.PropertyFilter{
				Property: "",
				RichText: &notionapi.TextFilterCondition{
					Contains: filter.TitleContains,
				},
			})
		}

		if filter.URLContains != "" {
			filters = append(filters, notionapi.PropertyFilter{
				Property: "URL",
				RichText: &notionapi.TextFilterCondition{
					Contains: filter.URLContains,
				},
			})
		}

		if filter.SummaryContains != "" {
			filters = append(filters, notionapi.PropertyFilter{
				Property: "Summary",
				RichText: &notionapi.TextFilterCondition{
					Contains: filter.SummaryContains,
				},
			})
		}

		if filter.AuthorContains != "" {
			filters = append(filters, notionapi.PropertyFilter{
				Property: "Author",
				RichText: &notionapi.TextFilterCondition{
					Contains: filter.AuthorContains,
				},
			})
		}

		if filter.HasTag != "" {
			filters = append(filters, notionapi.PropertyFilter{
				Property: "Tags",
				Relation: &notionapi.RelationFilterCondition{
					Contains: filter.HasTag,
				},
			})
		}

		// If we have multiple filters, combine them with AND
		if len(filters) > 0 {
			if len(filters) == 1 {
				query.Filter = &filters[0]
			} else {
				// Convert PropertyFilters to Filters
				filterInterfaces := make([]notionapi.Filter, len(filters))
				for i := range filters {
					filterInterfaces[i] = &filters[i]
				}
				query.Filter = notionapi.AndCompoundFilter(filterInterfaces)
			}
		}

		if filter.Limit > 0 {
			query.PageSize = filter.Limit
		}
	}

	// Add default sorting by Date Added (descending)
	query.Sorts = []notionapi.SortObject{
		{
			Property:  "Date Added",
			Direction: notionapi.SortOrderDESC,
		},
	}

	resp, err := s.client.API().Database.Query(ctx, s.client.BookmarksDB(), query)
	if err != nil {
		return nil, notion.NewError("list bookmarks", err, "failed to list bookmarks")
	}

	bookmarks := make([]*Bookmark, 0, len(resp.Results))
	for _, page := range resp.Results {
		bookmark, err := ToBookmark(&page)
		if err != nil {
			continue // Skip invalid bookmarks
		}
		bookmarks = append(bookmarks, bookmark)
	}

	return bookmarks, nil
}

// Query retrieves bookmarks with advanced filtering and sorting
func (s *Service) Query(ctx context.Context, options *QueryOptions) ([]*Bookmark, error) {
	query := &notionapi.DatabaseQueryRequest{}

	if options != nil {
		// Apply filter if provided
		if options.Filter != nil {
			bookmarks, err := s.List(ctx, options.Filter)
			if err != nil {
				return nil, err
			}

			// If a limit was specified in options but not in filter, apply it now
			if options.Limit > 0 && len(bookmarks) > options.Limit {
				bookmarks = bookmarks[:options.Limit]
			}

			return bookmarks, nil
		}

		// Apply sorting
		if options.SortBy != "" {
			sortOrder := notionapi.SortOrderDESC
			if options.SortOrder == SortOrderAsc {
				sortOrder = notionapi.SortOrderASC
			}

			if options.SortBy == SortByCreatedAt {
				query.Sorts = []notionapi.SortObject{
					{
						Timestamp: notionapi.TimestampCreated,
						Direction: sortOrder,
					},
				}
			} else if options.SortBy == SortByUpdatedAt {
				query.Sorts = []notionapi.SortObject{
					{
						Timestamp: notionapi.TimestampLastEdited,
						Direction: sortOrder,
					},
				}
			} else {
				query.Sorts = []notionapi.SortObject{
					{
						Property:  string(options.SortBy),
						Direction: sortOrder,
					},
				}
			}
		}

		if options.Limit > 0 {
			query.PageSize = options.Limit
		}
	}

	resp, err := s.client.API().Database.Query(ctx, s.client.BookmarksDB(), query)
	if err != nil {
		return nil, notion.NewError("query bookmarks", err, "failed to query bookmarks")
	}

	bookmarks := make([]*Bookmark, 0, len(resp.Results))
	for _, page := range resp.Results {
		bookmark, err := ToBookmark(&page)
		if err != nil {
			continue // Skip invalid bookmarks
		}
		bookmarks = append(bookmarks, bookmark)
	}

	return bookmarks, nil
}

// AddTags adds tags to a bookmark
func (s *Service) AddTags(ctx context.Context, bookmarkID string, tagIDs []string) (*Bookmark, error) {
	// Get the current bookmark to preserve existing tags
	bookmark, err := s.Get(ctx, bookmarkID)
	if err != nil {
		return nil, err
	}

	// Merge existing and new tag IDs (avoiding duplicates)
	tagMap := make(map[string]bool)
	for _, tagID := range bookmark.TagIDs {
		tagMap[tagID] = true
	}
	for _, tagID := range tagIDs {
		tagMap[tagID] = true
	}

	// Convert back to slice
	mergedTagIDs := make([]string, 0, len(tagMap))
	for tagID := range tagMap {
		mergedTagIDs = append(mergedTagIDs, tagID)
	}

	bookmark.TagIDs = mergedTagIDs

	return s.Update(ctx, bookmarkID, bookmark)
}

// RemoveTags removes tags from a bookmark
func (s *Service) RemoveTags(ctx context.Context, bookmarkID string, tagIDs []string) (*Bookmark, error) {
	// Get the current bookmark
	bookmark, err := s.Get(ctx, bookmarkID)
	if err != nil {
		return nil, err
	}

	// Create a map of tag IDs to remove
	removeMap := make(map[string]bool)
	for _, tagID := range tagIDs {
		removeMap[tagID] = true
	}

	// Filter out the tags to remove
	filteredTagIDs := make([]string, 0)
	for _, tagID := range bookmark.TagIDs {
		if !removeMap[tagID] {
			filteredTagIDs = append(filteredTagIDs, tagID)
		}
	}

	bookmark.TagIDs = filteredTagIDs

	return s.Update(ctx, bookmarkID, bookmark)
}
