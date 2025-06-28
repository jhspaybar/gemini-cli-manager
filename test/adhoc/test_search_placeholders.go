package main

import (
	"fmt"
	"strings"

	"github.com/jhspaybar/gemini-cli-manager/internal/theme"
	"github.com/jhspaybar/gemini-cli-manager/internal/ui/components"
)

func main() {
	// Initialize theme
	theme.SetTheme("github-dark")

	fmt.Println("SearchBar Context-Aware Placeholders Test")
	fmt.Println("========================================")
	fmt.Println()

	// Test 1: Extensions search
	fmt.Println("1. Extensions View Search Bar:")
	fmt.Println(strings.Repeat("-", 80))
	
	extSearch := components.NewSearchBar(60).
		SetPlaceholder("Search extensions by name or description...")
	fmt.Println(extSearch.Render())
	fmt.Println()

	// Test 2: Profiles search
	fmt.Println("2. Profiles View Search Bar:")
	fmt.Println(strings.Repeat("-", 80))
	
	profSearch := components.NewSearchBar(60).
		SetPlaceholder("Search profiles by name or tags...")
	fmt.Println(profSearch.Render())
	fmt.Println()

	// Test 3: Generic search (for comparison with old behavior)
	fmt.Println("3. Generic Search Bar (old behavior):")
	fmt.Println(strings.Repeat("-", 80))
	
	genericSearch := components.NewSearchBar(60).
		SetPlaceholder("Search extensions, profiles...")
	fmt.Println(genericSearch.Render())
	fmt.Println()

	// Test 4: Settings search (potential future use)
	fmt.Println("4. Settings View Search Bar (example):")
	fmt.Println(strings.Repeat("-", 80))
	
	settingsSearch := components.NewSearchBar(60).
		SetPlaceholder("Search settings...")
	fmt.Println(settingsSearch.Render())
	fmt.Println()

	// Show how placeholder changes don't affect existing value
	fmt.Println("5. Changing Placeholder with Existing Value:")
	fmt.Println(strings.Repeat("-", 80))
	
	dynamicSearch := components.NewSearchBar(60).
		SetValue("test query").
		SetActive(true)
	
	fmt.Println("With generic placeholder:")
	dynamicSearch.SetPlaceholder("Search everything...")
	fmt.Println(dynamicSearch.Render())
	
	fmt.Println("\nWith specific placeholder (value unchanged):")
	dynamicSearch.SetPlaceholder("Search extensions...")
	fmt.Println(dynamicSearch.Render())
}