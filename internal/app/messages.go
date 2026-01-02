package app

import (
	"github.com/krisk248/nexus/internal/domain"
)

// Message types for Bubbletea

// Pane represents which pane is active
type Pane int

const (
	PaneCalendar Pane = iota
	PaneTasks
	PaneTimeline
)

// Mode represents the current input mode
type Mode int

const (
	ModeNormal Mode = iota
	ModeInput
	ModeSearch
	ModeFilter
)

// Dialog represents which dialog is open
type Dialog int

const (
	DialogNone Dialog = iota
	DialogHelp
	DialogTheme
	DialogExport
	DialogClearTimeline
	DialogConfirmExit
)

// Messages

// SwitchPaneMsg switches to a different pane
type SwitchPaneMsg struct {
	Pane Pane
}

// DateSelectedMsg is sent when a date is selected in the calendar
type DateSelectedMsg struct {
	Date domain.CalendarDate
}

// TaskAddedMsg is sent when a new task is added
type TaskAddedMsg struct {
	Task *domain.Task
}

// TaskUpdatedMsg is sent when a task is updated
type TaskUpdatedMsg struct {
	Task *domain.Task
}

// TaskDeletedMsg is sent when a task is deleted
type TaskDeletedMsg struct {
	TaskID string
}

// TaskStateChangedMsg is sent when a task's state changes
type TaskStateChangedMsg struct {
	Task      *domain.Task
	PrevState domain.TaskState
	NewState  domain.TaskState
}

// TaskPriorityChangedMsg is sent when a task's priority changes
type TaskPriorityChangedMsg struct {
	Task     *domain.Task
	Priority domain.TaskPriority
}

// TimelineEventMsg is sent when a timeline event occurs
type TimelineEventMsg struct {
	Event *domain.TimelineEvent
}

// ToggleDialogMsg toggles a dialog
type ToggleDialogMsg struct {
	Dialog Dialog
}

// ThemeChangedMsg is sent when the theme changes
type ThemeChangedMsg struct {
	ThemeName string
}

// SearchMsg is sent when searching
type SearchMsg struct {
	Query string
}

// FilterMsg is sent when filtering
type FilterMsg struct {
	FilterType FilterType
	Value      string
}

// FilterType represents different filter types
type FilterType int

const (
	FilterNone FilterType = iota
	FilterState
	FilterPriority
)

// ExportMsg is sent when exporting
type ExportMsg struct {
	Format ExportFormat
	Scope  ExportScope
}

// ExportFormat represents export format
type ExportFormat int

const (
	ExportMarkdown ExportFormat = iota
	ExportJSON
	ExportPlainText
)

// ExportScope represents export scope
type ExportScope int

const (
	ExportCurrentDay ExportScope = iota
	ExportCurrentMonth
	ExportAll
)

// UndoMsg triggers an undo action
type UndoMsg struct{}

// SaveMsg triggers a save
type SaveMsg struct{}

// SavedMsg is sent after data is saved
type SavedMsg struct {
	Success bool
	Error   error
}

// LoadedMsg is sent when data is loaded
type LoadedMsg struct {
	Tasks    domain.TaskTree
	Timeline domain.Timeline
	Theme    string
}

// ErrorMsg represents an error
type ErrorMsg struct {
	Error error
}

// TickMsg is sent for periodic updates (debounce, etc.)
type TickMsg struct{}

// WindowSizeMsg is sent when window size changes (handled by Bubbletea)
// We use tea.WindowSizeMsg directly
