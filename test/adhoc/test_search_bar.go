package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/jhspaybar/gemini-cli-manager/internal/theme"
	"github.com/jhspaybar/gemini-cli-manager/internal/ui/components"
)

func main() {
	// Initialize theme
	theme.SetTheme("github-dark")

	fmt.Println("SearchBar Component Visual Test")
	fmt.Println("===============================")
	fmt.Println()

	// Test 1: Basic search bar
	fmt.Println("1. Basic Search Bar:")
	fmt.Println(strings.Repeat("-", 80))
	
	search := components.NewSearchBar(60)
	fmt.Println(search.Render())
	fmt.Println()

	// Test 2: Active/focused search bar
	fmt.Println("2. Active/Focused Search Bar:")
	fmt.Println(strings.Repeat("-", 80))
	
	activeSearch := components.NewSearchBar(60).
		SetActive(true)
	fmt.Println(activeSearch.Render())
	fmt.Println()

	// Test 3: Search bar with value
	fmt.Println("3. Search Bar with Value:")
	fmt.Println(strings.Repeat("-", 80))
	
	valueSearch := components.NewSearchBar(60).
		SetValue("extension name").
		SetActive(true)
	fmt.Println(valueSearch.Render())
	fmt.Println()

	// Test 4: Custom placeholder
	fmt.Println("4. Custom Placeholder:")
	fmt.Println(strings.Repeat("-", 80))
	
	customSearch := components.NewSearchBar(60).
		SetPlaceholder("Search extensions by name or description...")
	fmt.Println(customSearch.Render())
	fmt.Println()

	// Test 5: Different widths
	fmt.Println("5. Different Widths:")
	fmt.Println(strings.Repeat("-", 80))
	
	fmt.Println("Width 40:")
	narrow := components.NewSearchBar(40)
	fmt.Println(narrow.Render())
	
	fmt.Println("\nWidth 80:")
	wide := components.NewSearchBar(80)
	fmt.Println(wide.Render())
	
	fmt.Println("\nWidth 100:")
	veryWide := components.NewSearchBar(100)
	fmt.Println(veryWide.Render())
	fmt.Println()

	// Test 6: Custom styling
	fmt.Println("6. Custom Styling:")
	fmt.Println(strings.Repeat("-", 80))
	
	fmt.Println("Double border:")
	doubleBorder := components.NewSearchBar(60).
		SetBorderStyle(lipgloss.DoubleBorder())
	fmt.Println(doubleBorder.Render())
	
	fmt.Println("\nThick border:")
	thickBorder := components.NewSearchBar(60).
		SetBorderStyle(lipgloss.ThickBorder())
	fmt.Println(thickBorder.Render())
	
	fmt.Println("\nCustom colors:")
	customColors := components.NewSearchBar(60).
		SetBorderColors(theme.Success(), theme.Warning()).
		SetActive(true)
	fmt.Println(customColors.Render())
	fmt.Println()

	// Test 7: Different prompts
	fmt.Println("7. Different Prompts:")
	fmt.Println(strings.Repeat("-", 80))
	
	fmt.Println("Filter prompt:")
	filterSearch := components.NewSearchBar(60).
		SetPrompt("🔽 ").
		SetPlaceholder("Filter results...")
	fmt.Println(filterSearch.Render())
	
	fmt.Println("\nCommand prompt:")
	cmdSearch := components.NewSearchBar(60).
		SetPrompt("> ").
		SetPlaceholder("Enter command...")
	fmt.Println(cmdSearch.Render())
	
	fmt.Println("\nNo prompt:")
	noPromptSearch := components.NewSearchBar(60).
		SetPrompt("")
	fmt.Println(noPromptSearch.Render())
	fmt.Println()

	// Test 8: Character limit demonstration
	fmt.Println("8. Character Limit:")
	fmt.Println(strings.Repeat("-", 80))
	
	limitSearch := components.NewSearchBar(60).
		SetCharLimit(20).
		SetValue("This text will be truncated at 20 characters").
		SetActive(true)
	fmt.Println(limitSearch.Render())
	fmt.Println("(Limited to 20 characters)")
	fmt.Println()

	// Test 9: Integration with theme
	fmt.Println("9. Theme Integration:")
	fmt.Println(strings.Repeat("-", 80))
	
	fmt.Println("Current theme: github-dark")
	themeSearch := components.NewSearchBar(60)
	fmt.Println("Inactive:")
	fmt.Println(themeSearch.Render())
	fmt.Println("Active:")
	fmt.Println(themeSearch.SetActive(true).Render())
	
	// Switch theme
	theme.SetTheme("solarized-light")
	fmt.Println("\nAfter switching to solarized-light:")
	newThemeSearch := components.NewSearchBar(60)
	fmt.Println("Inactive:")
	fmt.Println(newThemeSearch.Render())
	fmt.Println("Active:")
	fmt.Println(newThemeSearch.SetActive(true).Render())
}