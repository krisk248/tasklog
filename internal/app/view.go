package app

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/krisk248/tasklog/internal/domain"
)

// View renders the entire application
func (m Model) View() string {
	if m.Width == 0 || m.Height == 0 {
		return "Loading..."
	}

	// Handle overview mode
	if m.ShowOverview {
		return m.renderOverview()
	}

	// Handle dialogs
	if m.ActiveDialog != DialogNone {
		return m.renderWithDialog()
	}

	// Main three-pane layout
	return m.renderMainLayout()
}

// renderMainLayout renders the three-pane layout
func (m Model) renderMainLayout() string {
	s := m.Styles

	// Calculate widths
	totalWidth := m.Width - 2 // Account for padding
	calendarWidth := totalWidth * 20 / 100
	timelineWidth := totalWidth * 30 / 100
	taskWidth := totalWidth - calendarWidth - timelineWidth - 2 // 2 for separators

	// Ensure minimum widths
	if calendarWidth < 22 {
		calendarWidth = 22
	}
	if timelineWidth < 25 {
		timelineWidth = 25
	}
	taskWidth = totalWidth - calendarWidth - timelineWidth - 2

	contentHeight := m.Height - 3 // Account for hints bar

	// Render each pane
	calendarPane := m.renderCalendarPane(calendarWidth, contentHeight)
	taskPane := m.renderTaskPane(taskWidth, contentHeight)
	timelinePane := m.renderTimelinePane(timelineWidth, contentHeight)

	// Separators
	sep1Style := s.Separator
	sep2Style := s.Separator
	if m.ActivePane == PaneCalendar || m.ActivePane == PaneTasks {
		sep1Style = s.SeparatorFocus
	}
	if m.ActivePane == PaneTasks || m.ActivePane == PaneTimeline {
		sep2Style = s.SeparatorFocus
	}

	separator1 := sep1Style.Render(strings.Repeat("│\n", contentHeight))
	separator2 := sep2Style.Render(strings.Repeat("│\n", contentHeight))

	// Join panes horizontally
	mainContent := lipgloss.JoinHorizontal(
		lipgloss.Top,
		calendarPane,
		separator1,
		taskPane,
		separator2,
		timelinePane,
	)

	// Add keyboard hints at bottom
	hints := m.renderKeyboardHints()

	// Add exit confirmation if active
	if m.ExitConfirm {
		hints = s.Header.Render("Press Ctrl+C again or 'y' to exit, any other key to cancel")
	}

	// Build final view
	view := lipgloss.JoinVertical(
		lipgloss.Left,
		mainContent,
		hints,
	)

	// Apply background
	return s.App.Width(m.Width).Height(m.Height).Render(view)
}

// renderCalendarPane renders the calendar pane
func (m Model) renderCalendarPane(width, height int) string {
	s := m.Styles
	c := m.CurrentTheme.Colors

	var b strings.Builder

	// Header
	title := "CALENDAR"
	if m.ActivePane == PaneCalendar {
		title = s.Header.Render(title)
	} else {
		title = lipgloss.NewStyle().Foreground(c.TextMuted).Render(title)
	}
	b.WriteString(title + "\n\n")

	// Month/Year header
	monthYear := fmt.Sprintf("%s %d", m.ViewingMonth.MonthName(), m.ViewingMonth.Year)
	b.WriteString(s.CalendarHeader.Render(monthYear) + "\n\n")

	// Weekday headers
	weekdays := []string{"Su", "Mo", "Tu", "We", "Th", "Fr", "Sa"}
	for _, wd := range weekdays {
		b.WriteString(s.CalendarDay.Render(wd))
	}
	b.WriteString("\n")

	// Calendar grid
	grid := domain.GenerateCalendarGrid(m.ViewingMonth)
	for _, week := range grid.Weeks {
		for _, day := range week {
			dayStr := fmt.Sprintf("%2d", day.Day)

			// Determine style
			var style lipgloss.Style
			isCurrentMonth := grid.IsCurrentMonth(day)
			isSelected := day.Equals(m.SelectedDate)
			isToday := day.IsToday()
			hasTasks := len(m.Tasks.GetTasksForDate(day.String())) > 0

			switch {
			case isSelected:
				style = s.CalendarSelected
			case isToday:
				style = s.CalendarToday
			case !isCurrentMonth:
				style = s.CalendarOther
			case hasTasks:
				style = s.CalendarWithTasks
			default:
				style = s.CalendarDay
			}

			b.WriteString(style.Render(dayStr))
		}
		b.WriteString("\n")
	}

	// Pad to fill height
	content := b.String()
	lines := strings.Split(content, "\n")
	for len(lines) < height {
		lines = append(lines, "")
	}

	return s.Pane.Width(width).Height(height).Render(strings.Join(lines[:height], "\n"))
}

// renderTaskPane renders the task pane
func (m Model) renderTaskPane(width, height int) string {
	s := m.Styles
	c := m.CurrentTheme.Colors

	var b strings.Builder

	// Header
	title := "TASKS"
	if m.ActivePane == PaneTasks {
		title = s.Header.Render(title)
	} else {
		title = lipgloss.NewStyle().Foreground(c.TextMuted).Render(title)
	}
	b.WriteString(title + "\n")

	// Date and stats
	tasks := m.Tasks.GetTasksForDate(m.SelectedDate.String())
	total, completed := domain.GetTaskStats(tasks)
	percentage := 0
	if total > 0 {
		percentage = completed * 100 / total
	}

	dateStr := m.SelectedDate.Format("January 2, 2006")
	statsStr := fmt.Sprintf("(%d%%)", percentage)
	b.WriteString(fmt.Sprintf("%s %s\n", dateStr, lipgloss.NewStyle().Foreground(c.TextMuted).Render(statsStr)))
	b.WriteString(strings.Repeat("─", width-2) + "\n")

	// Input field if in input mode
	if m.CurrentMode == ModeInput {
		prompt := "> "
		if m.EditingTask != nil {
			prompt = "Edit: "
		}
		b.WriteString(prompt + m.TextInput.View() + "\n")
	}

	// Search indicator
	if m.IsSearching {
		b.WriteString(lipgloss.NewStyle().Foreground(c.Primary).Render(fmt.Sprintf("Search: %s", m.SearchQuery)) + "\n")
	}

	// Filter indicator
	if m.IsFiltering {
		filterStr := fmt.Sprintf("Filter: %s=%s", m.FilterType, m.FilterValue)
		b.WriteString(lipgloss.NewStyle().Foreground(c.Secondary).Render(filterStr) + "\n")
	}

	// Task list
	visibleRows := height - 6 // Account for headers
	startIdx := m.TaskScrollOffset
	endIdx := min(startIdx+visibleRows, len(m.FlattenedTasks))

	// Scroll indicator (top)
	if m.TaskScrollOffset > 0 {
		b.WriteString(lipgloss.NewStyle().Foreground(c.TextMuted).Render("  ↑ more above") + "\n")
		visibleRows--
	}

	for i := startIdx; i < endIdx; i++ {
		ft := m.FlattenedTasks[i]
		task := ft.Task
		isSelected := i == m.SelectedTaskIndex && m.ActivePane == PaneTasks

		// Indentation
		indent := strings.Repeat("  ", ft.Depth)

		// Selection indicator
		selector := "  "
		if isSelected {
			selector = "> "
		}

		// Checkbox
		checkbox := m.getTaskCheckbox(task)

		// Priority indicator
		priority := m.getPriorityIndicator(task)

		// Expand/collapse indicator
		expandIcon := ""
		if len(task.Children) > 0 {
			if task.Expanded {
				expandIcon = "▼ "
			} else {
				expandIcon = "► "
			}
		}

		// Task title with appropriate style
		titleStyle := m.getTaskStyle(task, isSelected)
		title := titleStyle.Render(task.Title)

		// Running indicator
		runningIndicator := ""
		if task.IsRunning() {
			runningIndicator = lipgloss.NewStyle().Foreground(c.TaskRunning).Render(" ●")
		}

		line := fmt.Sprintf("%s%s%s%s%s%s%s", selector, indent, checkbox, priority, expandIcon, title, runningIndicator)

		// Truncate if too long
		if len(line) > width-2 {
			line = line[:width-5] + "..."
		}

		b.WriteString(line + "\n")
	}

	// Scroll indicator (bottom)
	if endIdx < len(m.FlattenedTasks) {
		b.WriteString(lipgloss.NewStyle().Foreground(c.TextMuted).Render("  ↓ more below") + "\n")
	}

	// Empty state
	if len(m.FlattenedTasks) == 0 {
		emptyMsg := "No tasks for this day. Press 'a' to add one."
		b.WriteString(lipgloss.NewStyle().Foreground(c.TextMuted).Render(emptyMsg) + "\n")
	}

	// Pad to fill height
	content := b.String()
	lines := strings.Split(content, "\n")
	for len(lines) < height {
		lines = append(lines, "")
	}

	return s.Pane.Width(width).Height(height).Render(strings.Join(lines[:height], "\n"))
}

// renderTimelinePane renders the timeline pane
func (m Model) renderTimelinePane(width, height int) string {
	s := m.Styles
	c := m.CurrentTheme.Colors

	var b strings.Builder

	// Header
	title := "TIMELINE"
	if m.ActivePane == PaneTimeline {
		title = s.Header.Render(title)
	} else {
		title = lipgloss.NewStyle().Foreground(c.TextMuted).Render(title)
	}
	b.WriteString(title + "\n")
	b.WriteString(strings.Repeat("─", width-2) + "\n")

	// Events
	events := m.Timeline.GetEventsForDate(m.SelectedDate.String())
	visibleRows := height - 4

	startIdx := m.TimelineScrollOffset
	endIdx := min(startIdx+visibleRows, len(events))

	// Scroll indicator (top)
	if m.TimelineScrollOffset > 0 {
		b.WriteString(lipgloss.NewStyle().Foreground(c.TextMuted).Render("↑ more above") + "\n")
		visibleRows--
	}

	for i := startIdx; i < endIdx; i++ {
		event := events[i]
		isLast := i == len(events)-1

		// Event icon with color
		iconColor := m.getEventColor(event)
		icon := lipgloss.NewStyle().Foreground(iconColor).Render(event.GetEventIcon())

		// Timestamp
		timestamp := event.Timestamp.Format("3:04 PM")
		timestampStyled := s.TimelineTimestamp.Render(timestamp)

		// Event description
		desc := fmt.Sprintf("%s %s", event.GetEventDescription(), event.TaskTitle)
		if len(desc) > width-15 {
			desc = desc[:width-18] + "..."
		}

		b.WriteString(fmt.Sprintf(" %s %s\n", icon, desc))
		b.WriteString(fmt.Sprintf("   %s\n", timestampStyled))

		// Connector line
		if !isLast && i < endIdx-1 {
			b.WriteString(s.TimelineConnector.Render(" │") + "\n")
		}
	}

	// Scroll indicator (bottom)
	if endIdx < len(events) {
		b.WriteString(lipgloss.NewStyle().Foreground(c.TextMuted).Render("↓ more below") + "\n")
	}

	// Empty state
	if len(events) == 0 {
		emptyMsg := "No activity yet today."
		b.WriteString(lipgloss.NewStyle().Foreground(c.TextMuted).Render(emptyMsg) + "\n")
	}

	// Pad to fill height
	content := b.String()
	lines := strings.Split(content, "\n")
	for len(lines) < height {
		lines = append(lines, "")
	}

	return s.Pane.Width(width).Height(height).Render(strings.Join(lines[:height], "\n"))
}

// renderKeyboardHints renders the keyboard hints bar
func (m Model) renderKeyboardHints() string {
	s := m.Styles

	var hints []string
	switch m.ActivePane {
	case PaneCalendar:
		hints = []string{"h/l:day", "j/k:week", "n/p:month", "T:today", "Tab:next"}
	case PaneTasks:
		hints = []string{"j/k:nav", "a:add", "e:edit", "d:del", "Space:done", "/:search", "1/2/3:priority"}
	case PaneTimeline:
		hints = []string{"j/k:scroll", "C:clear", "Tab:next"}
	}

	globalHints := []string{"?:help", "Ctrl+T:theme", "::overview"}
	hints = append(hints, globalHints...)

	hintStr := strings.Join(hints, " │ ")
	return s.KeyboardHint.Render(hintStr)
}

// renderWithDialog renders the main layout with a dialog overlay
func (m Model) renderWithDialog() string {
	// Render main layout (dimmed)
	main := m.renderMainLayout()

	// Render dialog
	var dialog string
	switch m.ActiveDialog {
	case DialogHelp:
		dialog = m.renderHelpDialog()
	case DialogTheme:
		dialog = m.renderThemeDialog()
	case DialogExport:
		dialog = m.renderExportDialog()
	case DialogClearTimeline:
		dialog = m.renderClearTimelineDialog()
	}

	// Center dialog over main content
	return m.overlayDialog(main, dialog)
}

// overlayDialog centers a dialog over the main content
func (m Model) overlayDialog(main, dialog string) string {
	mainLines := strings.Split(main, "\n")
	dialogLines := strings.Split(dialog, "\n")

	dialogWidth := 0
	for _, line := range dialogLines {
		if len(line) > dialogWidth {
			dialogWidth = len(line)
		}
	}

	// Calculate position
	startY := (m.Height - len(dialogLines)) / 2
	startX := (m.Width - dialogWidth) / 2

	// Overlay dialog
	for i, dialogLine := range dialogLines {
		lineY := startY + i
		if lineY >= 0 && lineY < len(mainLines) {
			// Insert dialog line into main line
			mainLine := mainLines[lineY]
			// Ensure main line is long enough
			for len(mainLine) < m.Width {
				mainLine += " "
			}
			// Replace portion with dialog
			if startX >= 0 && startX+len(dialogLine) <= len(mainLine) {
				mainLines[lineY] = mainLine[:startX] + dialogLine + mainLine[startX+len(dialogLine):]
			}
		}
	}

	return strings.Join(mainLines, "\n")
}

// renderHelpDialog renders the help dialog
func (m Model) renderHelpDialog() string {
	s := m.Styles
	c := m.CurrentTheme.Colors

	var b strings.Builder
	b.WriteString(s.ModalTitle.Render("Keyboard Shortcuts") + "\n\n")

	sections := []struct {
		title string
		keys  [][]string
	}{
		{
			title: "Global",
			keys: [][]string{
				{"Ctrl+C", "Exit (press twice)"},
				{"Ctrl+U", "Undo"},
				{"Ctrl+T", "Theme selector"},
				{"Ctrl+E", "Export"},
				{"?", "This help"},
				{":", "Month overview"},
				{"/", "Search tasks"},
				{"Esc", "Clear search/filter"},
				{"1/2/3", "Switch panes"},
				{"Tab", "Next pane"},
			},
		},
		{
			title: "Calendar",
			keys: [][]string{
				{"h/l", "Previous/next day"},
				{"j/k", "Previous/next week"},
				{"n/p", "Next/previous month"},
				{"T", "Jump to today"},
			},
		},
		{
			title: "Tasks",
			keys: [][]string{
				{"j/k", "Navigate up/down"},
				{"a", "Add task"},
				{"e", "Edit task"},
				{"d", "Delete task"},
				{"Space", "Toggle complete"},
				{"D", "Delegate task"},
				{"x", "Toggle delayed"},
				{"s", "Start/stop timer"},
				{"1/2/3", "Set priority P1/P2/P3"},
				{"0", "Clear priority"},
				{"Enter", "Expand/collapse"},
			},
		},
		{
			title: "Timeline",
			keys: [][]string{
				{"j/k", "Scroll"},
				{"C", "Clear timeline"},
			},
		},
	}

	for _, section := range sections {
		b.WriteString(lipgloss.NewStyle().Foreground(c.Primary).Bold(true).Render(section.title) + "\n")
		for _, kv := range section.keys {
			key := lipgloss.NewStyle().Foreground(c.Secondary).Width(12).Render(kv[0])
			desc := lipgloss.NewStyle().Foreground(c.TextPrimary).Render(kv[1])
			b.WriteString(fmt.Sprintf("  %s %s\n", key, desc))
		}
		b.WriteString("\n")
	}

	b.WriteString(lipgloss.NewStyle().Foreground(c.TextMuted).Render("Press any key to close"))

	return s.Modal.Render(b.String())
}

// renderThemeDialog renders the theme selection dialog
func (m Model) renderThemeDialog() string {
	s := m.Styles
	c := m.CurrentTheme.Colors

	var b strings.Builder
	b.WriteString(s.ModalTitle.Render("Select Theme") + "\n\n")

	themes := []struct {
		key  string
		name string
		desc string
	}{
		{"1", "ultraviolet", "Deep space with violet accents"},
		{"2", "terminal", "Classic green/amber CRT"},
		{"3", "minimal", "Clean grayscale"},
		{"4", "nord", "Arctic blue palette"},
	}

	for _, t := range themes {
		keyStyle := lipgloss.NewStyle().Foreground(c.Primary).Bold(true)
		nameStyle := lipgloss.NewStyle().Foreground(c.TextPrimary)
		descStyle := lipgloss.NewStyle().Foreground(c.TextMuted)

		current := ""
		if t.name == m.CurrentTheme.Name {
			current = lipgloss.NewStyle().Foreground(c.Success).Render(" ✓")
		}

		b.WriteString(fmt.Sprintf("  %s  %s%s\n", keyStyle.Render(t.key), nameStyle.Render(t.name), current))
		b.WriteString(fmt.Sprintf("      %s\n\n", descStyle.Render(t.desc)))
	}

	b.WriteString(lipgloss.NewStyle().Foreground(c.TextMuted).Render("Press number to select, Esc to close"))

	return s.Modal.Render(b.String())
}

// renderExportDialog renders the export dialog
func (m Model) renderExportDialog() string {
	s := m.Styles
	c := m.CurrentTheme.Colors

	var b strings.Builder
	b.WriteString(s.ModalTitle.Render("Export Tasks") + "\n\n")

	b.WriteString(lipgloss.NewStyle().Foreground(c.Primary).Render("Format:") + "\n")
	b.WriteString("  1. Markdown\n")
	b.WriteString("  2. JSON\n")
	b.WriteString("  3. Plain Text\n\n")

	b.WriteString(lipgloss.NewStyle().Foreground(c.Primary).Render("Scope:") + "\n")
	b.WriteString("  d. Current day\n")
	b.WriteString("  m. Current month\n")
	b.WriteString("  a. All tasks\n\n")

	b.WriteString(lipgloss.NewStyle().Foreground(c.TextMuted).Render("Press format+scope (e.g., '1d' for Markdown/Day), Esc to close"))

	return s.Modal.Render(b.String())
}

// renderClearTimelineDialog renders the clear timeline confirmation
func (m Model) renderClearTimelineDialog() string {
	s := m.Styles
	c := m.CurrentTheme.Colors

	var b strings.Builder
	b.WriteString(s.ModalTitle.Render("Clear Timeline?") + "\n\n")
	b.WriteString("This will remove all timeline events for today.\n")
	b.WriteString("This action cannot be undone.\n\n")
	b.WriteString(lipgloss.NewStyle().Foreground(c.Success).Render("Y") + " - Yes, clear\n")
	b.WriteString(lipgloss.NewStyle().Foreground(c.Error).Render("N") + " - No, cancel\n")

	return s.Dialog.Render(b.String())
}

// renderOverview renders the month overview screen
func (m Model) renderOverview() string {
	s := m.Styles
	c := m.CurrentTheme.Colors

	var b strings.Builder

	// Header
	title := fmt.Sprintf("Overview: %s %d", m.ViewingMonth.MonthName(), m.ViewingMonth.Year)
	b.WriteString(s.Header.Render(title) + "\n")
	b.WriteString(strings.Repeat("═", m.Width-2) + "\n\n")

	// Generate calendar grid
	grid := domain.GenerateCalendarGrid(m.ViewingMonth)

	// Render each day with its tasks
	colWidth := (m.Width - 4) / 7

	for _, week := range grid.Weeks {
		var dayHeaders []string
		var taskLines [][]string

		// Prepare data for this week
		for _, day := range week {
			isCurrentMonth := grid.IsCurrentMonth(day)
			if !isCurrentMonth {
				dayHeaders = append(dayHeaders, strings.Repeat(" ", colWidth))
				taskLines = append(taskLines, []string{})
				continue
			}

			// Day header
			dayStr := fmt.Sprintf("%d", day.Day)
			tasks := m.Tasks.GetTasksForDate(day.String())
			total, completed := domain.GetTaskStats(tasks)

			headerStyle := lipgloss.NewStyle().Foreground(c.TextPrimary)
			if day.Equals(m.SelectedDate) {
				headerStyle = headerStyle.Foreground(c.Primary).Bold(true)
			}
			if day.IsToday() {
				headerStyle = headerStyle.Foreground(c.CalendarToday)
			}

			stats := ""
			if total > 0 {
				stats = fmt.Sprintf(" (%d/%d)", completed, total)
			}
			dayHeaders = append(dayHeaders, headerStyle.Width(colWidth).Render(dayStr+stats))

			// Task previews (just titles, truncated)
			var lines []string
			for i, task := range tasks {
				if i >= 3 {
					lines = append(lines, lipgloss.NewStyle().Foreground(c.TextMuted).Render(fmt.Sprintf("+%d more", len(tasks)-3)))
					break
				}
				prefix := "○ "
				if task.State == domain.TaskStateCompleted {
					prefix = "● "
				}
				line := prefix + task.Title
				if len(line) > colWidth-1 {
					line = line[:colWidth-4] + "..."
				}
				lines = append(lines, lipgloss.NewStyle().Width(colWidth).Render(line))
			}
			taskLines = append(taskLines, lines)
		}

		// Render week
		b.WriteString(strings.Join(dayHeaders, " ") + "\n")

		// Find max task lines for this week
		maxLines := 0
		for _, lines := range taskLines {
			if len(lines) > maxLines {
				maxLines = len(lines)
			}
		}

		for lineIdx := 0; lineIdx < maxLines; lineIdx++ {
			var rowParts []string
			for _, lines := range taskLines {
				if lineIdx < len(lines) {
					rowParts = append(rowParts, lines[lineIdx])
				} else {
					rowParts = append(rowParts, strings.Repeat(" ", colWidth))
				}
			}
			b.WriteString(strings.Join(rowParts, " ") + "\n")
		}

		b.WriteString("\n")
	}

	// Footer
	b.WriteString(strings.Repeat("─", m.Width-2) + "\n")
	b.WriteString(lipgloss.NewStyle().Foreground(c.TextMuted).Render("Press ':' or Esc to return"))

	return s.App.Width(m.Width).Height(m.Height).Render(b.String())
}

// Helper methods for view

func (m Model) getTaskCheckbox(task *domain.Task) string {
	c := m.CurrentTheme.Colors
	switch task.State {
	case domain.TaskStateCompleted:
		return lipgloss.NewStyle().Foreground(c.TaskCompleted).Render("[✓] ")
	case domain.TaskStateDelegated:
		return lipgloss.NewStyle().Foreground(c.TaskDelegated).Render("[→] ")
	case domain.TaskStateDelayed:
		return lipgloss.NewStyle().Foreground(c.TaskDelayed).Render("[‖] ")
	default:
		return lipgloss.NewStyle().Foreground(c.TaskTodo).Render("[ ] ")
	}
}

func (m Model) getPriorityIndicator(task *domain.Task) string {
	c := m.CurrentTheme.Colors
	switch task.Priority {
	case domain.PriorityHigh:
		return lipgloss.NewStyle().Foreground(c.PriorityHigh).Bold(true).Render("P1 ")
	case domain.PriorityMed:
		return lipgloss.NewStyle().Foreground(c.PriorityMed).Render("P2 ")
	case domain.PriorityLow:
		return lipgloss.NewStyle().Foreground(c.PriorityLow).Render("P3 ")
	default:
		return ""
	}
}

func (m Model) getTaskStyle(task *domain.Task, isSelected bool) lipgloss.Style {
	s := m.Styles
	c := m.CurrentTheme.Colors

	if isSelected {
		return s.TaskSelected
	}

	switch task.State {
	case domain.TaskStateCompleted:
		return s.TaskCompleted
	case domain.TaskStateDelegated:
		return s.TaskDelegated
	case domain.TaskStateDelayed:
		return s.TaskDelayed
	default:
		if task.IsRunning() {
			return s.TaskRunning
		}
		return lipgloss.NewStyle().Foreground(c.TaskTodo)
	}
}

func (m Model) getEventColor(event *domain.TimelineEvent) lipgloss.Color {
	c := m.CurrentTheme.Colors
	switch event.Type {
	case domain.EventStarted:
		return c.TaskRunning
	case domain.EventCompleted:
		return c.TaskCompleted
	case domain.EventDelegated:
		return c.TaskDelegated
	case domain.EventDelayed:
		return c.TaskDelayed
	default:
		return c.TextPrimary
	}
}
