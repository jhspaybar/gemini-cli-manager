package theme

import (
	"github.com/charmbracelet/lipgloss"
	tint "github.com/lrstanley/bubbletint"
)

// Available theme names mapped to their IDs
var themeMapping = map[string]string{
	"Solarized Light":  "builtin_solarized_light",
	"Solarized Dark":   "builtin_solarized_dark",
	"Dracula":          "dracula",
	"GitHub":           "github",
	"Nord":             "nord",
	"Gruvbox Dark":     "gruvbox_dark",
	"Gruvbox Light":    "gruvbox_light",
	"Tomorrow Night":   "tomorrow_night",
	"Tomorrow":         "tomorrow",
	"Tokyo Night":      "tokyo_night",
	"One Dark":         "one_dark",
	"Ayu Dark":         "ayu_dark",
	"Catppuccin Mocha": "catppuccin_mocha",
	"Catppuccin Latte": "catppuccin_latte",
	"Monokai":          "monokai",
}

// themeNames is an ordered list of theme display names
var themeNames []string

// Initialize with default theme
func init() {
	// Initialize the default registry with all built-in themes
	tint.NewDefaultRegistry()

	// Build ordered theme names list
	themeNames = []string{
		"Solarized Light",
		"Solarized Dark",
		"Dracula",
		"GitHub",
		"Nord",
		"Gruvbox Dark",
		"Gruvbox Light",
		"Tomorrow Night",
		"Tomorrow",
		"Tokyo Night",
		"One Dark",
		"Ayu Dark",
		"Catppuccin Mocha",
		"Catppuccin Latte",
		"Monokai",
	}

	// Default to Solarized Light
	setThemeByName("Solarized Light")
}

// setThemeByName applies a theme based on its display name
func setThemeByName(name string) {
	if id, ok := themeMapping[name]; ok {
		tint.SetTintID(id)
	}
}

// SetTheme changes the active theme
func SetTheme(themeName string) error {
	setThemeByName(themeName)
	return nil
}

// SetThemeByIndex changes theme by index
func SetThemeByIndex(index int) error {
	if index < 0 || index >= len(themeNames) {
		return nil
	}
	setThemeByName(themeNames[index])
	return nil
}

// GetAvailableThemes returns all available theme names
func GetAvailableThemes() []string {
	return themeNames
}

// GetCurrentTheme returns the active theme name
func GetCurrentTheme() string {
	currentID := tint.ID()

	// Map tint ID back to theme name
	for name, id := range themeMapping {
		if id == currentID {
			return name
		}
	}

	// Default if not found
	return themeNames[0]
}

// GetCurrentThemeIndex returns the index of the current theme
func GetCurrentThemeIndex() int {
	current := GetCurrentTheme()
	for i, name := range themeNames {
		if name == current {
			return i
		}
	}
	return 0
}

// Core colors from theme
// Note: These return lipgloss.TerminalColor which is compatible with lipgloss.Color
func Fg() lipgloss.TerminalColor           { return tint.Fg() }
func Bg() lipgloss.TerminalColor           { return tint.Bg() }
func Black() lipgloss.TerminalColor        { return tint.Black() }
func Red() lipgloss.TerminalColor          { return tint.Red() }
func Green() lipgloss.TerminalColor        { return tint.Green() }
func Yellow() lipgloss.TerminalColor       { return tint.Yellow() }
func Blue() lipgloss.TerminalColor         { return tint.Blue() }
func Purple() lipgloss.TerminalColor       { return tint.Purple() }
func Cyan() lipgloss.TerminalColor         { return tint.Cyan() }
func White() lipgloss.TerminalColor        { return tint.White() }
func BrightBlack() lipgloss.TerminalColor  { return tint.BrightBlack() }
func BrightRed() lipgloss.TerminalColor    { return tint.BrightRed() }
func BrightGreen() lipgloss.TerminalColor  { return tint.BrightGreen() }
func BrightYellow() lipgloss.TerminalColor { return tint.BrightYellow() }
func BrightBlue() lipgloss.TerminalColor   { return tint.BrightBlue() }
func BrightPurple() lipgloss.TerminalColor { return tint.BrightPurple() }
func BrightCyan() lipgloss.TerminalColor   { return tint.BrightCyan() }
func BrightWhite() lipgloss.TerminalColor  { return tint.BrightWhite() }

// Semantic colors for our application
func Primary() lipgloss.TerminalColor   { return Cyan() }
func Secondary() lipgloss.TerminalColor { return Blue() }
func Success() lipgloss.TerminalColor   { return Green() }
func Warning() lipgloss.TerminalColor   { return Yellow() }
func Error() lipgloss.TerminalColor     { return Red() }
func Info() lipgloss.TerminalColor      { return BrightBlue() }

// UI element colors
func TextPrimary() lipgloss.TerminalColor   { return Fg() }
func TextSecondary() lipgloss.TerminalColor { return BrightBlack() }
func TextMuted() lipgloss.TerminalColor     { return Black() }
func Border() lipgloss.TerminalColor        { return BrightBlack() }
func BorderFocus() lipgloss.TerminalColor   { return Primary() }
func Selection() lipgloss.TerminalColor     { return tint.SelectionBg() }
func Background() lipgloss.TerminalColor    { return Bg() }
