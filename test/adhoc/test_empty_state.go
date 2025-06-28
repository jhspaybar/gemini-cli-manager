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

	fmt.Println("EmptyState Component Visual Test")
	fmt.Println("================================")
	fmt.Println()

	// Test 1: Extensions empty state
	fmt.Println("1. Extensions Empty State:")
	fmt.Println(strings.Repeat("-", 80))
	
	extEmpty := components.NewEmptyState(80).
		SetIcon("📦").
		SetTitle("No extensions installed").
		SetAction("Press 'n' to install your first extension")
	
	fmt.Println(extEmpty.Render())
	fmt.Println()

	// Test 2: Profiles empty state
	fmt.Println("2. Profiles Empty State:")
	fmt.Println(strings.Repeat("-", 80))
	
	profEmpty := components.NewEmptyState(80).
		SetIcon("👤").
		SetTitle("No profiles configured").
		SetAction("Press 'n' to create your first profile")
	
	fmt.Println(profEmpty.Render())
	fmt.Println()

	// Test 3: Search no results
	fmt.Println("3. Search No Results:")
	fmt.Println(strings.Repeat("-", 80))
	
	searchEmpty := components.NewEmptyState(80).
		NoItemsFound().
		SetAction("Press '/' to modify your search")
	
	fmt.Println(searchEmpty.Render())
	fmt.Println()

	// Test 4: Generic no data
	fmt.Println("4. Generic No Data:")
	fmt.Println(strings.Repeat("-", 80))
	
	genericEmpty := components.NewEmptyState(80).
		NoData("tools available")
	
	fmt.Println(genericEmpty.Render())
	fmt.Println()

	// Test 5: Coming soon
	fmt.Println("5. Coming Soon Feature:")
	fmt.Println(strings.Repeat("-", 80))
	
	soonEmpty := components.NewEmptyState(80).
		ComingSoon()
	
	fmt.Println(soonEmpty.Render())
	fmt.Println()

	// Test 6: Custom styled empty state
	fmt.Println("6. Custom Styled Empty State:")
	fmt.Println(strings.Repeat("-", 80))
	
	customBox := lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		BorderForeground(theme.Primary()).
		Padding(3, 6).
		Align(lipgloss.Center)
	
	customIcon := lipgloss.NewStyle().
		Foreground(theme.Warning()).
		Bold(true)
	
	customTitle := lipgloss.NewStyle().
		Foreground(theme.Error()).
		Bold(true).
		Underline(true)
	
	customEmpty := components.NewEmptyState(80).
		SetIcon("⚠️").
		SetTitle("Attention Required").
		SetDescription("This requires your immediate action").
		SetStyles(customBox, customIcon, customTitle, 
			lipgloss.NewStyle(), lipgloss.NewStyle())
	
	fmt.Println(customEmpty.Render())
	fmt.Println()

	// Test 7: Different widths
	fmt.Println("7. Different Widths:")
	fmt.Println(strings.Repeat("-", 80))
	
	fmt.Println("Width 40:")
	narrow := components.NewEmptyState(40).
		SetIcon("📱").
		SetTitle("Narrow view").
		SetDescription("Responsive sizing")
	fmt.Println(narrow.Render())
	
	fmt.Println("\nWidth 100:")
	wide := components.NewEmptyState(100).
		SetIcon("🖥️").
		SetTitle("Wide view").
		SetDescription("Still constrained to max width")
	fmt.Println(wide.Render())
	
	fmt.Println()

	// Test 8: No centering
	fmt.Println("8. Left-aligned (not centered):")
	fmt.Println(strings.Repeat("-", 80))
	
	leftAligned := components.NewEmptyState(80).
		SetIcon("←").
		SetTitle("Left aligned").
		SetDescription("Not centered in container").
		SetCentered(false)
	
	fmt.Println(leftAligned.Render())
	fmt.Println()

	// Test 9: Minimal empty state
	fmt.Println("9. Minimal Empty State (title only):")
	fmt.Println(strings.Repeat("-", 80))
	
	minimal := components.NewEmptyState(80).
		SetTitle("Nothing here")
	
	fmt.Println(minimal.Render())
}