package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Tab represents a single tab
type Tab struct {
	Title  string
	Icon   string
	ID     string // Optional identifier
}

// TabBar manages the rendering of a horizontal tab bar
type TabBar struct {
	tabs              []Tab
	activeIndex       int
	width             int
	activeTabStyle    lipgloss.Style
	inactiveTabStyle  lipgloss.Style
	activeTabBorder   lipgloss.Border
	inactiveTabBorder lipgloss.Border
	borderColor       lipgloss.TerminalColor
}

// NewTabBar creates a new tab bar component
func NewTabBar(tabs []Tab, width int) *TabBar {
	return &TabBar{
		tabs:        tabs,
		activeIndex: 0,
		width:       width,
	}
}

// SetActiveIndex sets the active tab by index
func (tb *TabBar) SetActiveIndex(index int) {
	if index >= 0 && index < len(tb.tabs) {
		tb.activeIndex = index
	}
}

// SetActiveByID sets the active tab by ID
func (tb *TabBar) SetActiveByID(id string) {
	for i, tab := range tb.tabs {
		if tab.ID == id {
			tb.activeIndex = i
			break
		}
	}
}

// GetActiveIndex returns the current active tab index
func (tb *TabBar) GetActiveIndex() int {
	return tb.activeIndex
}

// GetActiveTab returns the current active tab
func (tb *TabBar) GetActiveTab() *Tab {
	if tb.activeIndex >= 0 && tb.activeIndex < len(tb.tabs) {
		return &tb.tabs[tb.activeIndex]
	}
	return nil
}

// SetStyles allows customization of tab styles
func (tb *TabBar) SetStyles(activeStyle, inactiveStyle lipgloss.Style, borderColor lipgloss.TerminalColor) {
	tb.activeTabStyle = activeStyle
	tb.inactiveTabStyle = inactiveStyle
	tb.borderColor = borderColor
	
	// Set up tab borders
	tb.inactiveTabBorder = tabBorderWithBottom("┴", "─", "┴")
	tb.activeTabBorder = tabBorderWithBottom("┘", " ", "└")
}

// SetWidth updates the width of the tab bar
func (tb *TabBar) SetWidth(width int) {
	tb.width = width
}

// Render renders the complete tab bar with gap filler
func (tb *TabBar) Render() string {
	// Render individual tabs without gaps
	var renderedTabs []string
	for i, tab := range tb.tabs {
		isFirst := i == 0
		isActive := i == tb.activeIndex
		
		var style lipgloss.Style
		if isActive {
			style = tb.activeTabStyle
		} else {
			style = tb.inactiveTabStyle
		}
		
		// Adjust border for first/last tabs
		border, _, _, _, _ := style.GetBorder()
		if isFirst && isActive {
			border.BottomLeft = "│"
		} else if isFirst && !isActive {
			border.BottomLeft = "├"
		}
		// Note: We don't modify the last tab's bottom right when active
		// It should keep the default "└" from activeTabBorder
		// We only keep the default "┴" for inactive tabs
		style = style.Border(border)
		
		content := fmt.Sprintf("%s %s", tab.Icon, tab.Title)
		renderedTabs = append(renderedTabs, style.Render(content))
	}
	
	// Join tabs horizontally (no gaps)
	tabRow := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	
	// Create a simple gap filler with just a bottom line
	tabRowWidth := lipgloss.Width(tabRow)
	gapWidth := max(0, tb.width-tabRowWidth)
	
	if gapWidth > 0 {
		// Simple approach: just draw the bottom line
		gapLine := lipgloss.NewStyle().
			Foreground(tb.borderColor).
			Render(strings.Repeat("─", gapWidth))
		
		// We need to position this at the bottom of the tab row
		// Create empty lines for the top part of the gap
		emptyLine := strings.Repeat(" ", gapWidth)
		gapContent := emptyLine + "\n" + emptyLine + "\n" + gapLine
		
		tabRow = lipgloss.JoinHorizontal(lipgloss.Top, tabRow, gapContent)
	}
	
	return tabRow
}

// RenderWithContent renders tabs with content area below
func (tb *TabBar) RenderWithContent(content string, contentHeight int) string {
	tabRow := tb.Render()
	
	// Create content box with no top border (connects to tabs)
	contentStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(tb.borderColor).
		UnsetBorderTop(). // This removes the top border completely
		Padding(1, 2).
		Width(tb.width).
		Height(contentHeight)
	
	contentBox := contentStyle.Render(content)
	
	// Combine tabs and content using string builder (no gap)
	var result strings.Builder
	result.WriteString(tabRow)
	result.WriteString("\n")
	result.WriteString(contentBox)
	
	return result.String()
}

// Helper function for tab borders
func tabBorderWithBottom(left, middle, right string) lipgloss.Border {
	border := lipgloss.RoundedBorder()
	border.BottomLeft = left
	border.Bottom = middle
	border.BottomRight = right
	return border
}

// Helper function
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}