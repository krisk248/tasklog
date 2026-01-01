package theme

import "github.com/charmbracelet/lipgloss"

// Terminal theme - Classic green/amber on black
var Terminal = Theme{
	Name: "terminal",
	Colors: ColorScheme{
		// Base
		Background: lipgloss.Color("#000000"),
		Surface:    lipgloss.Color("#0a0a0a"),
		Border:     lipgloss.Color("#333333"),

		// Text
		TextPrimary:   lipgloss.Color("#00ff00"),
		TextSecondary: lipgloss.Color("#00cc00"),
		TextMuted:     lipgloss.Color("#006600"),

		// Accents
		Primary:   lipgloss.Color("#00ff00"),
		Secondary: lipgloss.Color("#ffb000"),
		Accent:    lipgloss.Color("#ff6600"),

		// Status
		Success: lipgloss.Color("#00ff00"),
		Warning: lipgloss.Color("#ffb000"),
		Error:   lipgloss.Color("#ff0000"),

		// Calendar
		CalendarHeader:       lipgloss.Color("#00ff00"),
		CalendarToday:        lipgloss.Color("#ffb000"),
		CalendarSelected:     lipgloss.Color("#00ff00"),
		CalendarDayWithTasks: lipgloss.Color("#00cc00"),
		CalendarOtherMonth:   lipgloss.Color("#333333"),

		// Tasks
		TaskTodo:      lipgloss.Color("#00ff00"),
		TaskCompleted: lipgloss.Color("#006600"),
		TaskDelegated: lipgloss.Color("#ffb000"),
		TaskDelayed:   lipgloss.Color("#ff6600"),
		TaskRunning:   lipgloss.Color("#00ffff"),

		// Priority
		PriorityHigh: lipgloss.Color("#ff0000"),
		PriorityMed:  lipgloss.Color("#ffb000"),
		PriorityLow:  lipgloss.Color("#006600"),

		// Timeline
		TimelineConnector: lipgloss.Color("#333333"),
		TimelineTimestamp: lipgloss.Color("#00cc00"),

		// UI
		FocusIndicator: lipgloss.Color("#00ff00"),
		Separator:      lipgloss.Color("#333333"),
		ModalOverlay:   lipgloss.Color("#000000"),
	},
}

// Minimal theme - Grayscale
var Minimal = Theme{
	Name: "minimal",
	Colors: ColorScheme{
		// Base
		Background: lipgloss.Color("#1a1a1a"),
		Surface:    lipgloss.Color("#262626"),
		Border:     lipgloss.Color("#404040"),

		// Text
		TextPrimary:   lipgloss.Color("#ffffff"),
		TextSecondary: lipgloss.Color("#a3a3a3"),
		TextMuted:     lipgloss.Color("#737373"),

		// Accents (minimal uses white/gray)
		Primary:   lipgloss.Color("#ffffff"),
		Secondary: lipgloss.Color("#d4d4d4"),
		Accent:    lipgloss.Color("#a3a3a3"),

		// Status
		Success: lipgloss.Color("#ffffff"),
		Warning: lipgloss.Color("#d4d4d4"),
		Error:   lipgloss.Color("#a3a3a3"),

		// Calendar
		CalendarHeader:       lipgloss.Color("#ffffff"),
		CalendarToday:        lipgloss.Color("#ffffff"),
		CalendarSelected:     lipgloss.Color("#ffffff"),
		CalendarDayWithTasks: lipgloss.Color("#d4d4d4"),
		CalendarOtherMonth:   lipgloss.Color("#525252"),

		// Tasks
		TaskTodo:      lipgloss.Color("#ffffff"),
		TaskCompleted: lipgloss.Color("#737373"),
		TaskDelegated: lipgloss.Color("#a3a3a3"),
		TaskDelayed:   lipgloss.Color("#d4d4d4"),
		TaskRunning:   lipgloss.Color("#ffffff"),

		// Priority
		PriorityHigh: lipgloss.Color("#ffffff"),
		PriorityMed:  lipgloss.Color("#d4d4d4"),
		PriorityLow:  lipgloss.Color("#a3a3a3"),

		// Timeline
		TimelineConnector: lipgloss.Color("#404040"),
		TimelineTimestamp: lipgloss.Color("#a3a3a3"),

		// UI
		FocusIndicator: lipgloss.Color("#ffffff"),
		Separator:      lipgloss.Color("#404040"),
		ModalOverlay:   lipgloss.Color("#1a1a1a"),
	},
}

// Nord theme - Arctic, north-bluish color palette
var Nord = Theme{
	Name: "nord",
	Colors: ColorScheme{
		// Base - Polar Night
		Background: lipgloss.Color("#2e3440"),
		Surface:    lipgloss.Color("#3b4252"),
		Border:     lipgloss.Color("#4c566a"),

		// Text - Snow Storm
		TextPrimary:   lipgloss.Color("#eceff4"),
		TextSecondary: lipgloss.Color("#e5e9f0"),
		TextMuted:     lipgloss.Color("#d8dee9"),

		// Accents - Frost
		Primary:   lipgloss.Color("#88c0d0"), // Nord8
		Secondary: lipgloss.Color("#81a1c1"), // Nord9
		Accent:    lipgloss.Color("#5e81ac"), // Nord10

		// Status - Aurora
		Success: lipgloss.Color("#a3be8c"), // Nord14
		Warning: lipgloss.Color("#ebcb8b"), // Nord13
		Error:   lipgloss.Color("#bf616a"), // Nord11

		// Calendar
		CalendarHeader:       lipgloss.Color("#88c0d0"),
		CalendarToday:        lipgloss.Color("#a3be8c"),
		CalendarSelected:     lipgloss.Color("#88c0d0"),
		CalendarDayWithTasks: lipgloss.Color("#81a1c1"),
		CalendarOtherMonth:   lipgloss.Color("#4c566a"),

		// Tasks
		TaskTodo:      lipgloss.Color("#eceff4"),
		TaskCompleted: lipgloss.Color("#a3be8c"),
		TaskDelegated: lipgloss.Color("#b48ead"), // Nord15
		TaskDelayed:   lipgloss.Color("#ebcb8b"),
		TaskRunning:   lipgloss.Color("#88c0d0"),

		// Priority
		PriorityHigh: lipgloss.Color("#bf616a"),
		PriorityMed:  lipgloss.Color("#ebcb8b"),
		PriorityLow:  lipgloss.Color("#4c566a"),

		// Timeline
		TimelineConnector: lipgloss.Color("#4c566a"),
		TimelineTimestamp: lipgloss.Color("#d8dee9"),

		// UI
		FocusIndicator: lipgloss.Color("#88c0d0"),
		Separator:      lipgloss.Color("#4c566a"),
		ModalOverlay:   lipgloss.Color("#2e3440"),
	},
}
