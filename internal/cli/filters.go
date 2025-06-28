package cli

import (
	"strings"

	"github.com/jhspaybar/gemini-cli-manager/internal/extension"
	"github.com/jhspaybar/gemini-cli-manager/internal/profile"
)

// filterExtensions filters extensions based on the search query
func filterExtensions(extensions []*extension.Extension, query string) []*extension.Extension {
	if query == "" {
		return extensions
	}

	query = strings.ToLower(query)
	filtered := make([]*extension.Extension, 0)

	for _, ext := range extensions {
		// Search in name and description
		if strings.Contains(strings.ToLower(ext.Name), query) ||
			strings.Contains(strings.ToLower(ext.Description), query) {
			filtered = append(filtered, ext)
		}
	}

	return filtered
}

// filterProfiles filters profiles based on the search query
func filterProfiles(profiles []*profile.Profile, query string) []*profile.Profile {
	if query == "" {
		return profiles
	}

	query = strings.ToLower(query)
	filtered := make([]*profile.Profile, 0)

	for _, prof := range profiles {
		// Search in name, description, and tags
		if strings.Contains(strings.ToLower(prof.Name), query) ||
			strings.Contains(strings.ToLower(prof.Description), query) {
			filtered = append(filtered, prof)
			continue
		}

		// Search in tags
		for _, tag := range prof.Tags {
			if strings.Contains(strings.ToLower(tag), query) {
				filtered = append(filtered, prof)
				break
			}
		}
	}

	return filtered
}