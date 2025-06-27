package profile

import (
	"fmt"
	"regexp"
)

// Validator validates profile configurations
type Validator struct {
	idPattern *regexp.Regexp
}

// NewValidator creates a new profile validator
func NewValidator() *Validator {
	return &Validator{
		idPattern: regexp.MustCompile(`^[a-z0-9-]+$`),
	}
}

// Validate performs validation on a profile
func (v *Validator) Validate(p *Profile) error {
	if p == nil {
		return &ValidationError{Field: "profile", Message: "profile cannot be nil"}
	}

	// Validate ID
	if p.ID == "" {
		return &ValidationError{Field: "id", Message: "profile ID is required"}
	}
	if !v.idPattern.MatchString(p.ID) {
		return &ValidationError{Field: "id", Message: "invalid profile ID format (use lowercase letters, numbers, and hyphens)"}
	}
	if len(p.ID) > 64 {
		return &ValidationError{Field: "id", Message: "profile ID must be 64 characters or less"}
	}

	// Validate name
	if p.Name == "" {
		return &ValidationError{Field: "name", Message: "profile name is required"}
	}
	if len(p.Name) > 100 {
		return &ValidationError{Field: "name", Message: "profile name must be 100 characters or less"}
	}

	// Validate inheritance (check for circular dependencies)
	if len(p.Inherits) > 0 {
		seen := make(map[string]bool)
		seen[p.ID] = true
		for _, parentID := range p.Inherits {
			if seen[parentID] {
				return &ValidationError{Field: "inherits", Message: "circular inheritance detected"}
			}
			seen[parentID] = true
		}
	}

	// Validate auto-detect patterns
	if p.AutoDetect != nil {
		for i, pattern := range p.AutoDetect.Patterns {
			// Test if pattern is valid
			if pattern == "" {
				return &ValidationError{
					Field:   fmt.Sprintf("auto_detect.patterns[%d]", i),
					Message: "empty pattern not allowed",
				}
			}
		}
	}

	return nil
}
