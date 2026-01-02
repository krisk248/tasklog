package theme

import "github.com/charmbracelet/lipgloss"

// ColorScheme defines all colors used in the application
type ColorScheme struct {
	// Base colors
	Background lipgloss.Color
	Surface    lipgloss.Color
	Border     lipgloss.Color

	// Text colors
	TextPrimary   lipgloss.Color
	TextSecondary lipgloss.Color
	TextMuted     lipgloss.Color

	// Accent colors
	Primary   lipgloss.Color
	Secondary lipgloss.Color
	Accent    lipgloss.Color

	// Status colors
	Success lipgloss.Color
	Warning lipgloss.Color
	Error   lipgloss.Color

	// Calendar specific
	CalendarHeader       lipgloss.Color
	CalendarToday        lipgloss.Color
	CalendarSelected     lipgloss.Color
	CalendarDayWithTasks lipgloss.Color
	CalendarOtherMonth   lipgloss.Color

	// Task states
	TaskTodo      lipgloss.Color
	TaskCompleted lipgloss.Color
	TaskDelegated lipgloss.Color
	TaskDelayed   lipgloss.Color
	TaskRunning   lipgloss.Color

	// Priority colors
	PriorityHigh lipgloss.Color
	PriorityMed  lipgloss.Color
	PriorityLow  lipgloss.Color

	// Timeline
	TimelineConnector lipgloss.Color
	TimelineTimestamp lipgloss.Color

	// UI elements
	FocusIndicator lipgloss.Color
	Separator      lipgloss.Color
	ModalOverlay   lipgloss.Color
}

// Theme represents a complete theme with name and colors
type Theme struct {
	Name   string
	Colors ColorScheme
}

// Available themes
var themes = map[string]Theme{
	"ultraviolet": Ultraviolet,
	"terminal":    Terminal,
	"minimal":     Minimal,
	"nord":        Nord,
}

// GetTheme returns a theme by name
func GetTheme(name string) Theme {
	if theme, ok := themes[name]; ok {
		return theme
	}
	return Ultraviolet // Default
}

// GetThemeNames returns all available theme names
func GetThemeNames() []string {
	names := make([]string, 0, len(themes))
	for name := range themes {
		names = append(names, name)
	}
	return names
}

// Styles holds all lipgloss styles derived from a theme
type Styles struct {
	// Base
	App    lipgloss.Style
	Pane   lipgloss.Style
	Header lipgloss.Style

	// Calendar
	CalendarHeader    lipgloss.Style
	CalendarDay       lipgloss.Style
	CalendarToday     lipgloss.Style
	CalendarSelected  lipgloss.Style
	CalendarOther     lipgloss.Style
	CalendarWithTasks lipgloss.Style

	// Tasks
	TaskItem      lipgloss.Style
	TaskSelected  lipgloss.Style
	TaskCompleted lipgloss.Style
	TaskDelegated lipgloss.Style
	TaskDelayed   lipgloss.Style
	TaskRunning   lipgloss.Style
	TaskPriority1 lipgloss.Style
	TaskPriority2 lipgloss.Style
	TaskPriority3 lipgloss.Style

	// Timeline
	TimelineItem      lipgloss.Style
	TimelineTimestamp lipgloss.Style
	TimelineConnector lipgloss.Style

	// UI
	Separator      lipgloss.Style
	SeparatorFocus lipgloss.Style
	KeyboardHint   lipgloss.Style
	Modal          lipgloss.Style
	ModalTitle     lipgloss.Style
	Dialog         lipgloss.Style
	Input          lipgloss.Style
	InputFocused   lipgloss.Style
}

// NewStyles creates all styles from a theme
func NewStyles(t Theme) Styles {
	c := t.Colors

	return Styles{
		// Base styles - no explicit background, let terminal handle it
		App: lipgloss.NewStyle().
			Foreground(c.TextPrimary),

		Pane: lipgloss.NewStyle().
			Foreground(c.TextPrimary).
			Padding(0, 1),

		Header: lipgloss.NewStyle().
			Foreground(c.Primary).
			Bold(true),

		// Calendar styles
		CalendarHeader: lipgloss.NewStyle().
			Foreground(c.CalendarHeader).
			Bold(true),

		CalendarDay: lipgloss.NewStyle().
			Foreground(c.TextPrimary),

		CalendarToday: lipgloss.NewStyle().
			Foreground(c.CalendarToday).
			Bold(true),

		CalendarSelected: lipgloss.NewStyle().
			Foreground(c.CalendarSelected).
			Bold(true),

		CalendarOther: lipgloss.NewStyle().
			Foreground(c.CalendarOtherMonth),

		CalendarWithTasks: lipgloss.NewStyle().
			Foreground(c.CalendarDayWithTasks),

		// Task styles
		TaskItem: lipgloss.NewStyle().
			Foreground(c.TaskTodo),

		TaskSelected: lipgloss.NewStyle().
			Foreground(c.FocusIndicator).
			Bold(true),

		TaskCompleted: lipgloss.NewStyle().
			Foreground(c.TaskCompleted).
			Strikethrough(true),

		TaskDelegated: lipgloss.NewStyle().
			Foreground(c.TaskDelegated),

		TaskDelayed: lipgloss.NewStyle().
			Foreground(c.TaskDelayed),

		TaskRunning: lipgloss.NewStyle().
			Foreground(c.TaskRunning).
			Bold(true),

		TaskPriority1: lipgloss.NewStyle().
			Foreground(c.PriorityHigh).
			Bold(true),

		TaskPriority2: lipgloss.NewStyle().
			Foreground(c.PriorityMed),

		TaskPriority3: lipgloss.NewStyle().
			Foreground(c.PriorityLow),

		// Timeline styles
		TimelineItem: lipgloss.NewStyle().
			Foreground(c.TextPrimary),

		TimelineTimestamp: lipgloss.NewStyle().
			Foreground(c.TimelineTimestamp),

		TimelineConnector: lipgloss.NewStyle().
			Foreground(c.TimelineConnector),

		// UI styles
		Separator: lipgloss.NewStyle().
			Foreground(c.Separator),

		SeparatorFocus: lipgloss.NewStyle().
			Foreground(c.FocusIndicator),

		KeyboardHint: lipgloss.NewStyle().
			Foreground(c.TextMuted),

		Modal: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(c.Border).
			Padding(1, 2),

		ModalTitle: lipgloss.NewStyle().
			Foreground(c.Primary).
			Bold(true).
			MarginBottom(1),

		Dialog: lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(c.Primary).
			Padding(1, 2),

		Input: lipgloss.NewStyle().
			Foreground(c.TextPrimary).
			Padding(0, 1),

		InputFocused: lipgloss.NewStyle().
			Foreground(c.TextPrimary).
			Border(lipgloss.NormalBorder()).
			BorderForeground(c.Primary).
			Padding(0, 1),
	}
}
