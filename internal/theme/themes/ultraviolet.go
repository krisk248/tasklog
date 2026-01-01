package themes

import "github.com/charmbracelet/lipgloss"

// Ultraviolet color palette - Deep space black with electric violet accents
var UltravioletColors = struct {
	// Base
	Background lipgloss.Color
	Surface    lipgloss.Color
	Border     lipgloss.Color

	// Text
	TextPrimary   lipgloss.Color
	TextSecondary lipgloss.Color
	TextMuted     lipgloss.Color

	// Accents
	Primary   lipgloss.Color // Electric violet
	Secondary lipgloss.Color // Soft lavender
	Accent    lipgloss.Color // Hot pink/magenta

	// Status
	Success lipgloss.Color // Cyan glow
	Warning lipgloss.Color // Amber
	Error   lipgloss.Color // Rose

	// Additional violet shades
	VioletDark   lipgloss.Color
	VioletBright lipgloss.Color
	Magenta      lipgloss.Color
	Cyan         lipgloss.Color
}{
	// Base - Deep space blacks with purple undertones
	Background: lipgloss.Color("#0d0d14"),
	Surface:    lipgloss.Color("#1a1625"),
	Border:     lipgloss.Color("#2d2640"),

	// Text - Light grays for readability
	TextPrimary:   lipgloss.Color("#e5e7eb"),
	TextSecondary: lipgloss.Color("#9ca3af"),
	TextMuted:     lipgloss.Color("#6b7280"),

	// Accents - Violet spectrum
	Primary:   lipgloss.Color("#a855f7"), // Electric violet
	Secondary: lipgloss.Color("#c084fc"), // Soft lavender
	Accent:    lipgloss.Color("#e879f9"), // Hot pink/magenta

	// Status colors
	Success: lipgloss.Color("#22d3ee"), // Cyan glow
	Warning: lipgloss.Color("#f59e0b"), // Amber
	Error:   lipgloss.Color("#f43f5e"), // Rose

	// Additional shades
	VioletDark:   lipgloss.Color("#7c3aed"),
	VioletBright: lipgloss.Color("#d946ef"),
	Magenta:      lipgloss.Color("#ec4899"),
	Cyan:         lipgloss.Color("#06b6d4"),
}
