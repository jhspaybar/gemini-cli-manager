package cli

import (
	"fmt"
	"strings"
)

// ErrorType represents different categories of errors
type ErrorType int

const (
	ErrorTypeValidation ErrorType = iota
	ErrorTypeFileSystem
	ErrorTypeNetwork
	ErrorTypeConfiguration
	ErrorTypePermission
	ErrorTypeNotFound
	ErrorTypeConflict
	ErrorTypeSystem
	ErrorTypeInfo // For informational messages
)

// UIError represents an error with user-friendly formatting
type UIError struct {
	Type    ErrorType
	Message string
	Details string
	Hint    string
}

// Error implements the error interface
func (e UIError) Error() string {
	return e.Message
}

// Format formats the error for display in the UI
func (e UIError) Format() string {
	var parts []string
	
	// Main message
	parts = append(parts, errorStyle.Render("Error: "+e.Message))
	
	// Details if available
	if e.Details != "" {
		parts = append(parts, textMutedStyle.Render("  "+e.Details))
	}
	
	// Hint if available
	if e.Hint != "" {
		parts = append(parts, helpStyle.Render("  💡 "+e.Hint))
	}
	
	return strings.Join(parts, "\n")
}

// Common error constructors

// NewValidationError creates a validation error
func NewValidationError(message, hint string) UIError {
	return UIError{
		Type:    ErrorTypeValidation,
		Message: message,
		Hint:    hint,
	}
}

// NewFileSystemError creates a filesystem error
func NewFileSystemError(operation, path string, err error) UIError {
	return UIError{
		Type:    ErrorTypeFileSystem,
		Message: fmt.Sprintf("Failed to %s file", operation),
		Details: fmt.Sprintf("Path: %s", path),
		Hint:    "Check file permissions and ensure the path exists",
	}
}

// NewPermissionError creates a permission error
func NewPermissionError(resource string) UIError {
	return UIError{
		Type:    ErrorTypePermission,
		Message: fmt.Sprintf("Permission denied accessing %s", resource),
		Hint:    "Try running with appropriate permissions or check file ownership",
	}
}

// NewNotFoundError creates a not found error
func NewNotFoundError(resource, name string) UIError {
	return UIError{
		Type:    ErrorTypeNotFound,
		Message: fmt.Sprintf("%s not found: %s", resource, name),
		Hint:    fmt.Sprintf("Use the list command to see available %ss", strings.ToLower(resource)),
	}
}

// NewConflictError creates a conflict error
func NewConflictError(resource, name string) UIError {
	return UIError{
		Type:    ErrorTypeConflict,
		Message: fmt.Sprintf("%s already exists: %s", resource, name),
		Hint:    "Choose a different name or delete the existing one first",
	}
}

// NewNetworkError creates a network error
func NewNetworkError(operation string, err error) UIError {
	hint := "Check your internet connection"
	if strings.Contains(err.Error(), "timeout") {
		hint = "The operation timed out. Check your connection or try again"
	} else if strings.Contains(err.Error(), "certificate") {
		hint = "There may be a certificate issue. Check the URL or your system time"
	}
	
	return UIError{
		Type:    ErrorTypeNetwork,
		Message: fmt.Sprintf("Network error during %s", operation),
		Details: err.Error(),
		Hint:    hint,
	}
}

// WrapError wraps a generic error with context
func WrapError(err error, context string) UIError {
	if err == nil {
		return UIError{}
	}
	
	// Check if it's already a UIError
	if uiErr, ok := err.(UIError); ok {
		return uiErr
	}
	
	// Try to determine error type from error message
	errStr := err.Error()
	
	if strings.Contains(errStr, "permission denied") || strings.Contains(errStr, "access denied") {
		return NewPermissionError(context)
	}
	
	if strings.Contains(errStr, "no such file") || strings.Contains(errStr, "not found") {
		return NewNotFoundError("Resource", context)
	}
	
	if strings.Contains(errStr, "already exists") {
		return NewConflictError("Resource", context)
	}
	
	if strings.Contains(errStr, "timeout") || strings.Contains(errStr, "connection") {
		return NewNetworkError(context, err)
	}
	
	// Generic error
	return UIError{
		Type:    ErrorTypeSystem,
		Message: fmt.Sprintf("Error in %s", context),
		Details: err.Error(),
	}
}