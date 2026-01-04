package app

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/krisk248/seyal/internal/domain"
	"github.com/krisk248/seyal/internal/storage"
	"github.com/krisk248/seyal/internal/theme"
)

// Version is the current version of seyal
var Version = "0.1.0"

// Model is the main application model
type Model struct {
	// Window size
	Width  int
	Height int

	// Current state
	ActivePane   Pane
	CurrentMode  Mode
	ActiveDialog Dialog

	// Data
	Tasks        domain.TaskTree
	Timeline     domain.Timeline
	SelectedDate domain.CalendarDate

	// Task pane state
	SelectedTaskIndex int
	FlattenedTasks    []domain.FlattenedTask
	TaskScrollOffset  int
	EditingTask       *domain.Task

	// Timeline pane state
	TimelineScrollOffset int

	// Calendar state (for month view navigation)
	ViewingMonth domain.CalendarDate

	// Input
	TextInput textinput.Model

	// Search & Filter
	SearchQuery  string
	FilterType   FilterType
	FilterValue  string
	IsSearching  bool
	IsFiltering  bool

	// Theme
	CurrentTheme theme.Theme
	Styles       theme.Styles

	// UI state
	ShowOverview    bool
	ShowHelp        bool
	ExitConfirm     bool
	ExitConfirmTime int64

	// Undo stack (simplified - stores full state)
	UndoStack []UndoState
	MaxUndo   int

	// Storage
	DataPath string
	IsDirty  bool
}

// UndoState stores state for undo
type UndoState struct {
	Tasks    domain.TaskTree
	Timeline domain.Timeline
}

// NewModel creates a new application model
func NewModel() Model {
	ti := textinput.New()
	ti.Placeholder = "Enter task..."
	ti.CharLimit = 256
	ti.Width = 50

	currentTheme := theme.GetTheme("ultraviolet")

	return Model{
		// Default size (will be updated on WindowSizeMsg)
		Width:  80,
		Height: 24,

		// Initial state
		ActivePane:   PaneTasks,
		CurrentMode:  ModeNormal,
		ActiveDialog: DialogNone,

		// Data
		Tasks:        make(domain.TaskTree),
		Timeline:     make(domain.Timeline),
		SelectedDate: domain.Today(),
		ViewingMonth: domain.Today(),

		// Task state
		SelectedTaskIndex: 0,
		FlattenedTasks:    make([]domain.FlattenedTask, 0),
		TaskScrollOffset:  0,

		// Input
		TextInput: ti,

		// Theme
		CurrentTheme: currentTheme,
		Styles:       theme.NewStyles(currentTheme),

		// Undo
		UndoStack: make([]UndoState, 0),
		MaxUndo:   50,

		// Not dirty initially
		IsDirty: false,
	}
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		tea.EnterAltScreen,
		m.loadData(),
	)
}

// loadData returns a command to load saved data
func (m Model) loadData() tea.Cmd {
	return func() tea.Msg {
		store, err := storage.NewStorage()
		if err != nil {
			// Return empty data on error
			return LoadedMsg{
				Tasks:    make(domain.TaskTree),
				Timeline: make(domain.Timeline),
				Theme:    "ultraviolet",
			}
		}

		schema, err := store.Load()
		if err != nil {
			return LoadedMsg{
				Tasks:    make(domain.TaskTree),
				Timeline: make(domain.Timeline),
				Theme:    "ultraviolet",
			}
		}

		return LoadedMsg{
			Tasks:    schema.Tasks,
			Timeline: schema.Timeline,
			Theme:    schema.Settings.Theme,
		}
	}
}

// saveData saves the current state to disk
func (m Model) saveData() tea.Cmd {
	return func() tea.Msg {
		store, err := storage.NewStorage()
		if err != nil {
			return SavedMsg{Success: false, Error: err}
		}

		schema := &storage.StorageSchema{
			Tasks:    m.Tasks,
			Timeline: m.Timeline,
			Settings: storage.Settings{
				Theme:      "ultraviolet",
				DateFormat: "January 2, 2006",
				TimeFormat: "12h",
			},
		}

		err = store.Save(schema)
		return SavedMsg{Success: err == nil, Error: err}
	}
}

// UpdateFlattenedTasks updates the flattened task list for rendering
func (m *Model) UpdateFlattenedTasks() {
	tasks := m.Tasks.GetTasksForDate(m.SelectedDate.String())
	m.FlattenedTasks = domain.FlattenTasks(tasks, 0, true)

	// Apply search filter
	if m.IsSearching && m.SearchQuery != "" {
		m.FlattenedTasks = m.filterTasksBySearch(m.FlattenedTasks)
	}

	// Apply state/priority filter
	if m.IsFiltering && m.FilterValue != "" {
		m.FlattenedTasks = m.filterTasksByType(m.FlattenedTasks)
	}

	// Ensure selected index is valid
	if m.SelectedTaskIndex >= len(m.FlattenedTasks) {
		m.SelectedTaskIndex = max(0, len(m.FlattenedTasks)-1)
	}
}

func (m *Model) filterTasksBySearch(tasks []domain.FlattenedTask) []domain.FlattenedTask {
	var filtered []domain.FlattenedTask
	query := m.SearchQuery
	for _, ft := range tasks {
		if containsIgnoreCase(ft.Task.Title, query) {
			filtered = append(filtered, ft)
		}
	}
	return filtered
}

func (m *Model) filterTasksByType(tasks []domain.FlattenedTask) []domain.FlattenedTask {
	var filtered []domain.FlattenedTask
	for _, ft := range tasks {
		switch m.FilterType {
		case FilterState:
			if string(ft.Task.State) == m.FilterValue {
				filtered = append(filtered, ft)
			}
		case FilterPriority:
			// Priority filter: "1", "2", "3"
			var targetPriority domain.TaskPriority
			switch m.FilterValue {
			case "1":
				targetPriority = domain.PriorityHigh
			case "2":
				targetPriority = domain.PriorityMed
			case "3":
				targetPriority = domain.PriorityLow
			}
			if ft.Task.Priority == targetPriority {
				filtered = append(filtered, ft)
			}
		}
	}
	return filtered
}

// GetSelectedTask returns the currently selected task
func (m *Model) GetSelectedTask() *domain.Task {
	if m.SelectedTaskIndex >= 0 && m.SelectedTaskIndex < len(m.FlattenedTasks) {
		return m.FlattenedTasks[m.SelectedTaskIndex].Task
	}
	return nil
}

// PushUndo saves current state to undo stack
func (m *Model) PushUndo() {
	// Deep copy tasks and timeline
	state := UndoState{
		Tasks:    deepCopyTaskTree(m.Tasks),
		Timeline: deepCopyTimeline(m.Timeline),
	}
	m.UndoStack = append(m.UndoStack, state)
	if len(m.UndoStack) > m.MaxUndo {
		m.UndoStack = m.UndoStack[1:]
	}
}

// PopUndo restores the previous state
func (m *Model) PopUndo() bool {
	if len(m.UndoStack) == 0 {
		return false
	}
	state := m.UndoStack[len(m.UndoStack)-1]
	m.UndoStack = m.UndoStack[:len(m.UndoStack)-1]
	m.Tasks = state.Tasks
	m.Timeline = state.Timeline
	m.UpdateFlattenedTasks()
	m.IsDirty = true
	return true
}

// SetTheme changes the current theme
func (m *Model) SetTheme(themeName string) {
	m.CurrentTheme = theme.GetTheme(themeName)
	m.Styles = theme.NewStyles(m.CurrentTheme)
	m.IsDirty = true
}

// Helper functions

func containsIgnoreCase(s, substr string) bool {
	// Simple case-insensitive contains
	return len(s) >= len(substr) && (substr == "" || findIgnoreCase(s, substr))
}

func findIgnoreCase(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		match := true
		for j := 0; j < len(substr); j++ {
			sc := s[i+j]
			pc := substr[j]
			// Convert to lowercase for comparison
			if sc >= 'A' && sc <= 'Z' {
				sc += 32
			}
			if pc >= 'A' && pc <= 'Z' {
				pc += 32
			}
			if sc != pc {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}

func deepCopyTaskTree(tt domain.TaskTree) domain.TaskTree {
	copy := make(domain.TaskTree)
	for date, tasks := range tt {
		copy[date] = deepCopyTasks(tasks)
	}
	return copy
}

func deepCopyTasks(tasks []*domain.Task) []*domain.Task {
	result := make([]*domain.Task, len(tasks))
	for i, t := range tasks {
		result[i] = deepCopyTask(t)
	}
	return result
}

func deepCopyTask(t *domain.Task) *domain.Task {
	if t == nil {
		return nil
	}
	copy := *t
	copy.Children = deepCopyTasks(t.Children)
	return &copy
}

func deepCopyTimeline(tl domain.Timeline) domain.Timeline {
	copy := make(domain.Timeline)
	for date, events := range tl {
		eventsCopy := make([]*domain.TimelineEvent, len(events))
		for i, e := range events {
			eCopy := *e
			eventsCopy[i] = &eCopy
		}
		copy[date] = eventsCopy
	}
	return copy
}
