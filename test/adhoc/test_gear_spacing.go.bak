package main

import (
	"fmt"
)

func main() {
	// Test the exact formatting
	icon := "⚙️"
	title := "Settings"
	
	content := fmt.Sprintf("%s %s", icon, title)
	
	fmt.Printf("Icon: '%s'\n", icon)
	fmt.Printf("Title: '%s'\n", title)
	fmt.Printf("Combined: '%s'\n", content)
	fmt.Printf("Length of combined: %d\n", len(content))
	
	// Check each byte
	fmt.Println("\nByte analysis:")
	for i, b := range []byte(content) {
		fmt.Printf("  [%d]: %02x (%c)\n", i, b, b)
	}
	
	// Check runes
	fmt.Println("\nRune analysis:")
	for i, r := range content {
		fmt.Printf("  [%d]: U+%04X (%c)\n", i, r, r)
	}
}