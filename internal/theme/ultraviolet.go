package theme

import "github.com/charmbracelet/lipgloss"

// Ultraviolet theme - Deep space black with electric violet accents
var Ultraviolet = Theme{
	Name: "ultraviolet",
	Colors: ColorScheme{
		// Base colors - Deep space blacks with purple undertones
		Background: lipgloss.Color("#0d0d14"),
		Surface:    lipgloss.Color("#1a1625"),
		Border:     lipgloss.Color("#2d2640"),

		// Text colors - Light grays for readability
		TextPrimary:   lipgloss.Color("#e5e7eb"),
		TextSecondary: lipgloss.Color("#9ca3af"),
		TextMuted:     lipgloss.Color("#6b7280"),

		// Accent colors - Violet spectrum
		Primary:   lipgloss.Color("#a855f7"), // Electric violet
		Secondary: lipgloss.Color("#c084fc"), // Soft lavender
		Accent:    lipgloss.Color("#e879f9"), // Hot pink/magenta

		// Status colors
		Success: lipgloss.Color("#22d3ee"), // Cyan glow
		Warning: lipgloss.Color("#f59e0b"), // Amber
		Error:   lipgloss.Color("#f43f5e"), // Rose

		// Calendar specific
		CalendarHeader:       lipgloss.Color("#a855f7"), // Electric violet
		CalendarToday:        lipgloss.Color("#22d3ee"), // Cyan
		CalendarSelected:     lipgloss.Color("#a855f7"), // Violet
		CalendarDayWithTasks: lipgloss.Color("#c084fc"), // Lavender
		CalendarOtherMonth:   lipgloss.Color("#4b5563"), // Dim gray

		// Task states
		TaskTodo:      lipgloss.Color("#e5e7eb"), // Light gray
		TaskCompleted: lipgloss.Color("#22d3ee"), // Cyan (success)
		TaskDelegated: lipgloss.Color("#c084fc"), // Lavender
		TaskDelayed:   lipgloss.Color("#f59e0b"), // Amber
		TaskRunning:   lipgloss.Color("#e879f9"), // Magenta (active)

		// Priority colors
		PriorityHigh: lipgloss.Color("#f43f5e"), // Rose/Red
		PriorityMed:  lipgloss.Color("#f59e0b"), // Amber
		PriorityLow:  lipgloss.Color("#9ca3af"), // Muted gray

		// Timeline
		TimelineConnector: lipgloss.Color("#2d2640"), // Border color
		TimelineTimestamp: lipgloss.Color("#9ca3af"), // Secondary text

		// UI elements
		FocusIndicator: lipgloss.Color("#a855f7"), // Electric violet
		Separator:      lipgloss.Color("#2d2640"), // Border color
		ModalOverlay:   lipgloss.Color("#0d0d14"), // Background
	},
}
