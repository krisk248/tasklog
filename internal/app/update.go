package app

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/krisk248/tasklog/internal/domain"
)

// Update handles all messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyMsg(msg)

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		return m, nil

	case LoadedMsg:
		m.Tasks = msg.Tasks
		m.Timeline = msg.Timeline
		if msg.Theme != "" {
			m.SetTheme(msg.Theme)
		}
		m.UpdateFlattenedTasks()
		return m, nil

	case DateSelectedMsg:
		m.SelectedDate = msg.Date
		m.SelectedTaskIndex = 0
		m.TaskScrollOffset = 0
		m.UpdateFlattenedTasks()
		return m, nil

	case TaskAddedMsg:
		m.PushUndo()
		m.Tasks.AddTask(msg.Task)
		m.UpdateFlattenedTasks()
		m.IsDirty = true
		// Add timeline event
		event := domain.NewTimelineEvent(msg.Task.ID, msg.Task.Title, domain.EventCreated)
		m.Timeline.AddEvent(m.SelectedDate.String(), event)
		return m, nil

	case TaskStateChangedMsg:
		m.PushUndo()
		msg.Task.SetState(msg.NewState)
		m.UpdateFlattenedTasks()
		m.IsDirty = true
		// Add timeline event
		event := domain.NewStateChangeEvent(msg.Task.ID, msg.Task.Title, msg.PrevState, msg.NewState)
		m.Timeline.AddEvent(m.SelectedDate.String(), event)
		return m, nil

	case TaskPriorityChangedMsg:
		m.PushUndo()
		msg.Task.SetPriority(msg.Priority)
		m.UpdateFlattenedTasks()
		m.IsDirty = true
		return m, nil

	case TaskDeletedMsg:
		m.PushUndo()
		m.Tasks.RemoveTask(m.SelectedDate.String(), msg.TaskID)
		m.Timeline.RemoveEventsByTaskID(m.SelectedDate.String(), msg.TaskID)
		m.UpdateFlattenedTasks()
		m.IsDirty = true
		return m, nil

	case ToggleDialogMsg:
		if m.ActiveDialog == msg.Dialog {
			m.ActiveDialog = DialogNone
		} else {
			m.ActiveDialog = msg.Dialog
		}
		return m, nil

	case ThemeChangedMsg:
		m.SetTheme(msg.ThemeName)
		return m, nil

	case SearchMsg:
		m.SearchQuery = msg.Query
		m.UpdateFlattenedTasks()
		return m, nil

	case FilterMsg:
		m.FilterType = msg.FilterType
		m.FilterValue = msg.Value
		m.IsFiltering = msg.FilterType != FilterNone
		m.UpdateFlattenedTasks()
		return m, nil

	case UndoMsg:
		m.PopUndo()
		return m, nil

	case ErrorMsg:
		// TODO: Display error
		return m, nil
	}

	return m, tea.Batch(cmds...)
}

// handleKeyMsg handles keyboard input
func (m Model) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Handle exit confirmation
	if m.ExitConfirm {
		if msg.String() == "ctrl+c" || msg.String() == "y" {
			return m, tea.Quit
		}
		m.ExitConfirm = false
		return m, nil
	}

	// Handle dialogs first
	if m.ActiveDialog != DialogNone {
		return m.handleDialogKeys(msg)
	}

	// Handle input mode
	if m.CurrentMode == ModeInput {
		return m.handleInputMode(msg)
	}

	// Handle search mode
	if m.CurrentMode == ModeSearch {
		return m.handleSearchMode(msg)
	}

	// Global keys
	switch msg.String() {
	case "ctrl+c":
		m.ExitConfirm = true
		m.ExitConfirmTime = time.Now().Unix()
		return m, nil

	case "ctrl+u":
		return m, func() tea.Msg { return UndoMsg{} }

	case "ctrl+t":
		return m, func() tea.Msg { return ToggleDialogMsg{Dialog: DialogTheme} }

	case "?":
		return m, func() tea.Msg { return ToggleDialogMsg{Dialog: DialogHelp} }

	case ":":
		m.ShowOverview = !m.ShowOverview
		return m, nil

	case "ctrl+e":
		return m, func() tea.Msg { return ToggleDialogMsg{Dialog: DialogExport} }

	case "/":
		m.CurrentMode = ModeSearch
		m.TextInput.SetValue("")
		m.TextInput.Focus()
		return m, nil

	case "1":
		if m.ActivePane == PaneTasks && m.GetSelectedTask() != nil {
			task := m.GetSelectedTask()
			return m, func() tea.Msg {
				return TaskPriorityChangedMsg{Task: task, Priority: domain.PriorityHigh}
			}
		}
		m.ActivePane = PaneCalendar
		return m, nil

	case "2":
		if m.ActivePane == PaneTasks && m.GetSelectedTask() != nil {
			task := m.GetSelectedTask()
			return m, func() tea.Msg {
				return TaskPriorityChangedMsg{Task: task, Priority: domain.PriorityMed}
			}
		}
		m.ActivePane = PaneTasks
		return m, nil

	case "3":
		if m.ActivePane == PaneTasks && m.GetSelectedTask() != nil {
			task := m.GetSelectedTask()
			return m, func() tea.Msg {
				return TaskPriorityChangedMsg{Task: task, Priority: domain.PriorityLow}
			}
		}
		m.ActivePane = PaneTimeline
		return m, nil

	case "tab":
		m.ActivePane = (m.ActivePane + 1) % 3
		return m, nil

	case "shift+tab":
		m.ActivePane = (m.ActivePane + 2) % 3
		return m, nil

	case "esc":
		// Clear search/filter
		m.IsSearching = false
		m.IsFiltering = false
		m.SearchQuery = ""
		m.FilterType = FilterNone
		m.FilterValue = ""
		m.UpdateFlattenedTasks()
		return m, nil
	}

	// Pane-specific keys
	switch m.ActivePane {
	case PaneCalendar:
		return m.handleCalendarKeys(msg)
	case PaneTasks:
		return m.handleTaskKeys(msg)
	case PaneTimeline:
		return m.handleTimelineKeys(msg)
	}

	return m, nil
}

// handleCalendarKeys handles calendar pane keyboard input
func (m Model) handleCalendarKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "h", "left":
		m.SelectedDate = m.SelectedDate.AddDays(-1)
		m.ViewingMonth = m.SelectedDate
		m.UpdateFlattenedTasks()
	case "l", "right":
		m.SelectedDate = m.SelectedDate.AddDays(1)
		m.ViewingMonth = m.SelectedDate
		m.UpdateFlattenedTasks()
	case "j", "down":
		m.SelectedDate = m.SelectedDate.AddDays(7)
		m.ViewingMonth = m.SelectedDate
		m.UpdateFlattenedTasks()
	case "k", "up":
		m.SelectedDate = m.SelectedDate.AddDays(-7)
		m.ViewingMonth = m.SelectedDate
		m.UpdateFlattenedTasks()
	case "n":
		m.SelectedDate = m.SelectedDate.AddMonths(1)
		m.ViewingMonth = m.SelectedDate
		m.UpdateFlattenedTasks()
	case "p":
		m.SelectedDate = m.SelectedDate.AddMonths(-1)
		m.ViewingMonth = m.SelectedDate
		m.UpdateFlattenedTasks()
	case "t", "T":
		m.SelectedDate = domain.Today()
		m.ViewingMonth = m.SelectedDate
		m.UpdateFlattenedTasks()
	}
	return m, nil
}

// handleTaskKeys handles task pane keyboard input
func (m Model) handleTaskKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "j", "down":
		if m.SelectedTaskIndex < len(m.FlattenedTasks)-1 {
			m.SelectedTaskIndex++
			m.ensureTaskVisible()
		}
	case "k", "up":
		if m.SelectedTaskIndex > 0 {
			m.SelectedTaskIndex--
			m.ensureTaskVisible()
		}
	case "a":
		// Add new task
		m.CurrentMode = ModeInput
		m.TextInput.SetValue("")
		m.TextInput.Focus()
		m.EditingTask = nil
	case "e":
		// Edit selected task
		if task := m.GetSelectedTask(); task != nil {
			m.CurrentMode = ModeInput
			m.TextInput.SetValue(task.Title)
			m.TextInput.Focus()
			m.EditingTask = task
		}
	case "d":
		// Delete selected task
		if task := m.GetSelectedTask(); task != nil {
			return m, func() tea.Msg { return TaskDeletedMsg{TaskID: task.ID} }
		}
	case " ":
		// Toggle completion
		if task := m.GetSelectedTask(); task != nil {
			prevState := task.State
			var newState domain.TaskState
			if task.State == domain.TaskStateCompleted {
				newState = domain.TaskStateTodo
			} else {
				newState = domain.TaskStateCompleted
			}
			return m, func() tea.Msg {
				return TaskStateChangedMsg{Task: task, PrevState: prevState, NewState: newState}
			}
		}
	case "D":
		// Delegate task
		if task := m.GetSelectedTask(); task != nil {
			prevState := task.State
			return m, func() tea.Msg {
				return TaskStateChangedMsg{Task: task, PrevState: prevState, NewState: domain.TaskStateDelegated}
			}
		}
	case "x":
		// Delay task
		if task := m.GetSelectedTask(); task != nil {
			prevState := task.State
			var newState domain.TaskState
			if task.State == domain.TaskStateDelayed {
				newState = domain.TaskStateTodo
			} else {
				newState = domain.TaskStateDelayed
			}
			return m, func() tea.Msg {
				return TaskStateChangedMsg{Task: task, PrevState: prevState, NewState: newState}
			}
		}
	case "s":
		// Start/stop task
		if task := m.GetSelectedTask(); task != nil {
			m.PushUndo()
			if task.IsRunning() {
				task.Stop()
			} else {
				task.Start()
				// Add started event
				event := domain.NewTimelineEvent(task.ID, task.Title, domain.EventStarted)
				m.Timeline.AddEvent(m.SelectedDate.String(), event)
			}
			m.IsDirty = true
		}
	case "enter", "right":
		// Expand/collapse
		if task := m.GetSelectedTask(); task != nil && len(task.Children) > 0 {
			task.ToggleExpanded()
			m.UpdateFlattenedTasks()
		}
	case "left":
		// Collapse
		if task := m.GetSelectedTask(); task != nil {
			task.Expanded = false
			m.UpdateFlattenedTasks()
		}
	case "0":
		// Clear priority
		if task := m.GetSelectedTask(); task != nil {
			return m, func() tea.Msg {
				return TaskPriorityChangedMsg{Task: task, Priority: domain.PriorityNone}
			}
		}
	case "f":
		// Enter filter mode - next key determines filter type
		m.CurrentMode = ModeFilter
	}
	return m, nil
}

// handleTimelineKeys handles timeline pane keyboard input
func (m Model) handleTimelineKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	events := m.Timeline.GetEventsForDate(m.SelectedDate.String())
	maxScroll := max(0, len(events)-m.visibleTimelineRows())

	switch msg.String() {
	case "j", "down":
		if m.TimelineScrollOffset < maxScroll {
			m.TimelineScrollOffset++
		}
	case "k", "up":
		if m.TimelineScrollOffset > 0 {
			m.TimelineScrollOffset--
		}
	case "ctrl+d":
		// Page down
		m.TimelineScrollOffset = min(m.TimelineScrollOffset+10, maxScroll)
	case "C":
		// Clear timeline
		return m, func() tea.Msg { return ToggleDialogMsg{Dialog: DialogClearTimeline} }
	}
	return m, nil
}

// handleDialogKeys handles keyboard input when a dialog is open
func (m Model) handleDialogKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "q":
		m.ActiveDialog = DialogNone
		return m, nil
	}

	switch m.ActiveDialog {
	case DialogHelp:
		// Just close on any key
		if msg.String() != "" {
			// Keep open, scroll if needed
		}
	case DialogTheme:
		return m.handleThemeDialogKeys(msg)
	case DialogExport:
		return m.handleExportDialogKeys(msg)
	case DialogClearTimeline:
		return m.handleClearTimelineDialogKeys(msg)
	}

	return m, nil
}

// handleThemeDialogKeys handles theme selection dialog
func (m Model) handleThemeDialogKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	themes := []string{"ultraviolet", "terminal", "minimal", "nord"}
	switch msg.String() {
	case "1":
		return m, func() tea.Msg { return ThemeChangedMsg{ThemeName: themes[0]} }
	case "2":
		return m, func() tea.Msg { return ThemeChangedMsg{ThemeName: themes[1]} }
	case "3":
		return m, func() tea.Msg { return ThemeChangedMsg{ThemeName: themes[2]} }
	case "4":
		return m, func() tea.Msg { return ThemeChangedMsg{ThemeName: themes[3]} }
	}
	return m, nil
}

// handleExportDialogKeys handles export dialog
func (m Model) handleExportDialogKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// TODO: Implement export dialog navigation
	return m, nil
}

// handleClearTimelineDialogKeys handles clear timeline confirmation
func (m Model) handleClearTimelineDialogKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "y", "Y":
		m.PushUndo()
		m.Timeline.ClearDate(m.SelectedDate.String())
		m.ActiveDialog = DialogNone
		m.IsDirty = true
	case "n", "N":
		m.ActiveDialog = DialogNone
	}
	return m, nil
}

// handleInputMode handles text input mode
func (m Model) handleInputMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.CurrentMode = ModeNormal
		m.TextInput.Blur()
		m.EditingTask = nil
		return m, nil
	case "enter":
		value := m.TextInput.Value()
		if value != "" {
			if m.EditingTask != nil {
				// Editing existing task
				m.PushUndo()
				m.EditingTask.Title = value
				m.EditingTask.UpdatedAt = time.Now()
				m.IsDirty = true
			} else {
				// Creating new task
				task := domain.NewTask(value, m.SelectedDate.String())
				m.CurrentMode = ModeNormal
				m.TextInput.Blur()
				return m, func() tea.Msg { return TaskAddedMsg{Task: task} }
			}
		}
		m.CurrentMode = ModeNormal
		m.TextInput.Blur()
		m.EditingTask = nil
		m.UpdateFlattenedTasks()
		return m, nil
	}

	// Update text input
	var cmd tea.Cmd
	m.TextInput, cmd = m.TextInput.Update(msg)
	return m, cmd
}

// handleSearchMode handles search mode
func (m Model) handleSearchMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.CurrentMode = ModeNormal
		m.IsSearching = false
		m.SearchQuery = ""
		m.TextInput.Blur()
		m.UpdateFlattenedTasks()
		return m, nil
	case "enter":
		m.CurrentMode = ModeNormal
		m.IsSearching = true
		m.SearchQuery = m.TextInput.Value()
		m.TextInput.Blur()
		m.UpdateFlattenedTasks()
		return m, nil
	}

	var cmd tea.Cmd
	m.TextInput, cmd = m.TextInput.Update(msg)
	// Live search as typing
	m.SearchQuery = m.TextInput.Value()
	m.IsSearching = m.SearchQuery != ""
	m.UpdateFlattenedTasks()
	return m, cmd
}

// Helper methods

func (m *Model) ensureTaskVisible() {
	visibleRows := m.visibleTaskRows()
	if m.SelectedTaskIndex < m.TaskScrollOffset {
		m.TaskScrollOffset = m.SelectedTaskIndex
	} else if m.SelectedTaskIndex >= m.TaskScrollOffset+visibleRows {
		m.TaskScrollOffset = m.SelectedTaskIndex - visibleRows + 1
	}
}

func (m *Model) visibleTaskRows() int {
	// Account for header, hints, padding
	return max(5, m.Height-10)
}

func (m *Model) visibleTimelineRows() int {
	return max(5, m.Height-10)
}
