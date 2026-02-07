package notion

import (
	"fmt"
)

// Error types for common Notion API errors
var (
	ErrNotFound     = fmt.Errorf("resource not found")
	ErrUnauthorized = fmt.Errorf("unauthorized: invalid API key or insufficient permissions")
	ErrRateLimited  = fmt.Errorf("rate limited: too many requests")
	ErrInvalidInput = fmt.Errorf("invalid input")
	ErrAPIError     = fmt.Errorf("notion API error")
)

// NotionError wraps a Notion API error with additional context
type NotionError struct {
	Operation string
	Err       error
	Message   string
}

func (e *NotionError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("%s failed: %s - %v", e.Operation, e.Message, e.Err)
	}
	return fmt.Sprintf("%s failed: %v", e.Operation, e.Err)
}

func (e *NotionError) Unwrap() error {
	return e.Err
}

// NewError creates a new NotionError
func NewError(operation string, err error, message string) error {
	return &NotionError{
		Operation: operation,
		Err:       err,
		Message:   message,
	}
}
