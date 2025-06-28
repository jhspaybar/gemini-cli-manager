package components

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jhspaybar/gemini-cli-manager/internal/theme"
)

// SearchBar represents a search input component with consistent styling
type SearchBar struct {
	textInput   textinput.Model
	width       int
	active      bool
	placeholder string
	prompt      string
	
	// Styling
	borderStyle       lipgloss.Border
	borderColor       lipgloss.TerminalColor
	borderColorFocus  lipgloss.TerminalColor
	containerPadding  [2]int // vertical, horizontal
}

// NewSearchBar creates a new search bar component
func NewSearchBar(width int) *SearchBar {
	ti := textinput.New()
	ti.Placeholder = "Type to search..."
	ti.CharLimit = 100
	ti.Prompt = "🔍 "
	
	return &SearchBar{
		textInput:        ti,
		width:            width,
		placeholder:      "Type to search...",
		prompt:           "🔍 ",
		borderStyle:      lipgloss.RoundedBorder(),
		borderColor:      theme.Border(),
		borderColorFocus: theme.BorderFocus(),
		containerPadding: [2]int{0, 1},
	}
}

// Init initializes the search bar
func (sb *SearchBar) Init() tea.Cmd {
	return nil
}

// Update handles search bar updates
func (sb *SearchBar) Update(msg tea.Msg) (*SearchBar, tea.Cmd) {
	if !sb.active {
		return sb, nil
	}
	
	var cmd tea.Cmd
	sb.textInput, cmd = sb.textInput.Update(msg)
	return sb, cmd
}

// View renders just the text input (for backward compatibility)
func (sb *SearchBar) View() string {
	return sb.textInput.View()
}

// Render renders the complete search bar with border
func (sb *SearchBar) Render() string {
	// Calculate available width for text input
	// Account for border (2) and padding
	innerWidth := sb.width - 2 - (sb.containerPadding[1] * 2)
	if innerWidth < 10 {
		innerWidth = 10
	}
	sb.textInput.Width = innerWidth
	
	// Choose border color based on focus state
	borderColor := sb.borderColor
	if sb.active {
		borderColor = sb.borderColorFocus
	}
	
	// Create container style
	style := lipgloss.NewStyle().
		Border(sb.borderStyle).
		BorderForeground(borderColor).
		Padding(sb.containerPadding[0], sb.containerPadding[1]).
		Width(sb.width)
	
	return style.Render(sb.textInput.View())
}

// Focus activates the search bar
func (sb *SearchBar) Focus() tea.Cmd {
	sb.active = true
	return sb.textInput.Focus()
}

// Blur deactivates the search bar
func (sb *SearchBar) Blur() {
	sb.active = false
	sb.textInput.Blur()
}

// SetWidth updates the search bar width
func (sb *SearchBar) SetWidth(width int) *SearchBar {
	sb.width = width
	return sb
}

// SetPlaceholder sets the placeholder text
func (sb *SearchBar) SetPlaceholder(placeholder string) *SearchBar {
	sb.placeholder = placeholder
	sb.textInput.Placeholder = placeholder
	return sb
}

// SetPrompt sets the prompt (e.g., "🔍 ")
func (sb *SearchBar) SetPrompt(prompt string) *SearchBar {
	sb.prompt = prompt
	sb.textInput.Prompt = prompt
	return sb
}

// SetValue sets the search query
func (sb *SearchBar) SetValue(value string) *SearchBar {
	sb.textInput.SetValue(value)
	return sb
}

// SetCharLimit sets the character limit
func (sb *SearchBar) SetCharLimit(limit int) *SearchBar {
	sb.textInput.CharLimit = limit
	return sb
}

// SetBorderStyle allows customization of the border
func (sb *SearchBar) SetBorderStyle(style lipgloss.Border) *SearchBar {
	sb.borderStyle = style
	return sb
}

// SetBorderColors sets the border colors for normal and focused states
func (sb *SearchBar) SetBorderColors(normal, focused lipgloss.TerminalColor) *SearchBar {
	sb.borderColor = normal
	sb.borderColorFocus = focused
	return sb
}

// SetPadding sets the container padding (vertical, horizontal)
func (sb *SearchBar) SetPadding(vertical, horizontal int) *SearchBar {
	sb.containerPadding = [2]int{vertical, horizontal}
	return sb
}

// Value returns the current search query
func (sb *SearchBar) Value() string {
	return sb.textInput.Value()
}

// Clear clears the search query
func (sb *SearchBar) Clear() *SearchBar {
	sb.textInput.SetValue("")
	return sb
}

// IsActive returns whether the search bar is active
func (sb *SearchBar) IsActive() bool {
	return sb.active
}

// SetActive sets the active state without triggering focus
func (sb *SearchBar) SetActive(active bool) *SearchBar {
	sb.active = active
	if active {
		sb.textInput.Focus()
	} else {
		sb.textInput.Blur()
	}
	return sb
}

// TextInput returns the underlying text input model for advanced usage
func (sb *SearchBar) TextInput() *textinput.Model {
	return &sb.textInput
}