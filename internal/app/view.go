package app

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/krisk248/seyal/internal/domain"
)

// View renders the entire application
func (m Model) View() string {
	if m.Width == 0 || m.Height == 0 {
		return "Loading..."
	}

	// Handle help screen (full screen)
	if m.ShowHelp {
		return m.renderHelpScreen()
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
	if calendarWidth < 30 {
		calendarWidth = 30
	}
	if timelineWidth < 28 {
		timelineWidth = 28
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

	// Apply styling to entire screen
	return lipgloss.NewStyle().
		Width(m.Width).
		Height(m.Height).
		Render(view)
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

	// Weekday headers (4 chars each to match day format)
	weekdays := []string{" Su ", " Mo ", " Tu ", " We ", " Th ", " Fr ", " Sa "}
	for _, wd := range weekdays {
		b.WriteString(lipgloss.NewStyle().Foreground(c.TextMuted).Render(wd))
	}
	b.WriteString("\n")

	// Calendar grid
	grid := domain.GenerateCalendarGrid(m.ViewingMonth)
	for _, week := range grid.Weeks {
		for _, day := range week {
			isCurrentMonth := grid.IsCurrentMonth(day)
			isSelected := day.Equals(m.SelectedDate)
			isToday := day.IsToday()
			hasTasks := len(m.Tasks.GetTasksForDate(day.String())) > 0

			// Format day with brackets for selected
			var dayStr string
			if isSelected {
				dayStr = fmt.Sprintf("[%2d]", day.Day)
			} else {
				dayStr = fmt.Sprintf(" %2d ", day.Day)
			}

			// Determine style
			var style lipgloss.Style
			switch {
			case isSelected:
				style = lipgloss.NewStyle().Foreground(c.CalendarSelected)
			case isToday:
				style = lipgloss.NewStyle().Foreground(c.CalendarToday).Bold(true)
			case !isCurrentMonth:
				style = lipgloss.NewStyle().Foreground(c.CalendarOtherMonth)
			case hasTasks:
				style = lipgloss.NewStyle().Foreground(c.CalendarDayWithTasks)
			default:
				style = lipgloss.NewStyle().Foreground(c.TextPrimary)
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
	if len(lines) > height {
		lines = lines[:height]
	}

	return lipgloss.NewStyle().
		Width(width).
		Height(height).
		Render(strings.Join(lines, "\n"))
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

		// Checkbox (plain text for length calculation)
		checkboxText := m.getTaskCheckboxText(task)

		// Priority indicator (plain text)
		priorityText := m.getPriorityText(task)

		// Expand/collapse indicator
		expandIcon := ""
		if len(task.Children) > 0 {
			if task.Expanded {
				expandIcon = "▼ "
			} else {
				expandIcon = "► "
			}
		}

		// Running indicator
		runningText := ""
		if task.IsRunning() {
			runningText = " ●"
		}

		// Pushed count indicator text (for width calculation)
		pushedText := ""
		if task.PushedCount > 0 {
			pushedText = fmt.Sprintf(" [↷%d]", task.PushedCount)
		}

		// Calculate available width for title (include all suffixes)
		prefixLen := len(selector) + len(indent) + len(checkboxText) + len(priorityText) + len(expandIcon) + len(runningText) + len(pushedText)
		availableWidth := width - prefixLen - 4 // margin

		// Truncate title if needed (on plain text, before styling)
		titleText := task.Title
		if len(titleText) > availableWidth && availableWidth > 3 {
			titleText = titleText[:availableWidth-3] + "..."
		}

		// Now apply styles
		checkbox := m.getTaskCheckbox(task)
		priority := m.getPriorityIndicator(task)
		titleStyle := m.getTaskStyle(task, isSelected)
		title := titleStyle.Render(titleText)

		runningIndicator := ""
		if task.IsRunning() {
			runningIndicator = lipgloss.NewStyle().Foreground(c.TaskRunning).Render(" ●")
		}

		// Pushed count indicator (styled)
		pushedIndicator := ""
		if task.PushedCount > 0 {
			pushedIndicator = lipgloss.NewStyle().Foreground(c.Warning).Render(pushedText)
		}

		line := fmt.Sprintf("%s%s%s%s%s%s%s%s", selector, indent, checkbox, priority, expandIcon, title, runningIndicator, pushedIndicator)
		b.WriteString(line + "\n")
	}

	// Scroll indicator (bottom)
	if endIdx < len(m.FlattenedTasks) {
		b.WriteString(lipgloss.NewStyle().Foreground(c.TextMuted).Render("  ↓ more below") + "\n")
	}

	// Empty state
	if len(m.FlattenedTasks) == 0 {
		emptyMsg := "No tasks. Press 'a' to add one."
		b.WriteString(lipgloss.NewStyle().Foreground(c.TextMuted).Render(emptyMsg) + "\n")
	}

	// Pad to fill height
	content := b.String()
	lines := strings.Split(content, "\n")
	for len(lines) < height {
		lines = append(lines, "")
	}
	if len(lines) > height {
		lines = lines[:height]
	}

	return lipgloss.NewStyle().
		Width(width).
		Height(height).
		Render(strings.Join(lines, "\n"))
}

// renderTimelinePane renders the timeline pane
func (m Model) renderTimelinePane(width, height int) string {
	s := m.Styles
	c := m.CurrentTheme.Colors

	var contentLines []string

	// Header
	title := "TIMELINE"
	if m.ActivePane == PaneTimeline {
		title = s.Header.Render(title)
	} else {
		title = lipgloss.NewStyle().Foreground(c.TextMuted).Render(title)
	}
	contentLines = append(contentLines, title)
	contentLines = append(contentLines, strings.Repeat("─", width-4))

	// Events
	events := m.Timeline.GetEventsForDate(m.SelectedDate.String())
	availableLines := height - 4

	// Each event takes ~3 lines (desc, timestamp, connector)
	visibleEvents := max(1, availableLines/3)

	startIdx := m.TimelineScrollOffset
	endIdx := min(startIdx+visibleEvents, len(events))

	// Scroll indicator (top)
	if m.TimelineScrollOffset > 0 {
		contentLines = append(contentLines, lipgloss.NewStyle().Foreground(c.TextMuted).Render("↑ more above"))
		visibleEvents--
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
		if len(desc) > width-18 {
			desc = desc[:width-21] + "..."
		}

		contentLines = append(contentLines, fmt.Sprintf(" %s %s", icon, desc))
		contentLines = append(contentLines, fmt.Sprintf("   %s", timestampStyled))

		// Connector line
		if !isLast && i < endIdx-1 {
			contentLines = append(contentLines, s.TimelineConnector.Render(" │"))
		}
	}

	// Scroll indicator (bottom)
	if endIdx < len(events) {
		contentLines = append(contentLines, lipgloss.NewStyle().Foreground(c.TextMuted).Render("↓ more below"))
	}

	// Empty state
	if len(events) == 0 {
		contentLines = append(contentLines, lipgloss.NewStyle().Foreground(c.TextMuted).Render("No activities yet."))
		contentLines = append(contentLines, lipgloss.NewStyle().Foreground(c.TextMuted).Render("Press 's' to start a task."))
	}

	// Build scrollbar if needed
	scrollbar := m.renderTimelineScrollbar(height-2, len(events), visibleEvents)

	// Pad content lines to fill height
	for len(contentLines) < height-2 {
		contentLines = append(contentLines, "")
	}

	// Combine content with scrollbar
	var result strings.Builder
	contentWidth := width - 3 // Leave space for scrollbar
	for i, line := range contentLines {
		if i >= height-2 {
			break
		}
		// Pad line to content width
		lineWidth := lipgloss.Width(line)
		if lineWidth < contentWidth {
			line = line + strings.Repeat(" ", contentWidth-lineWidth)
		}
		// Add scrollbar character
		scrollChar := " "
		if i < len(scrollbar) {
			scrollChar = scrollbar[i]
		}
		result.WriteString(line + scrollChar + "\n")
	}

	return lipgloss.NewStyle().
		Width(width).
		Height(height).
		Render(result.String())
}

// renderTimelineScrollbar renders a scrollbar for the timeline
func (m Model) renderTimelineScrollbar(height, totalEvents, visibleEvents int) []string {
	c := m.CurrentTheme.Colors

	scrollbar := make([]string, height)
	scrollStyle := lipgloss.NewStyle().Foreground(c.Border)
	thumbStyle := lipgloss.NewStyle().Foreground(c.Primary)

	// If no scrolling needed, return empty
	if totalEvents <= visibleEvents || totalEvents == 0 {
		for i := range scrollbar {
			scrollbar[i] = " "
		}
		return scrollbar
	}

	// Calculate thumb size and position
	thumbSize := max(1, height*visibleEvents/totalEvents)
	maxScroll := totalEvents - visibleEvents
	thumbPos := 0
	if maxScroll > 0 {
		thumbPos = m.TimelineScrollOffset * (height - thumbSize) / maxScroll
	}

	// Build scrollbar
	for i := 0; i < height; i++ {
		if i >= thumbPos && i < thumbPos+thumbSize {
			scrollbar[i] = thumbStyle.Render("█")
		} else {
			scrollbar[i] = scrollStyle.Render("░")
		}
	}

	return scrollbar
}

// renderKeyboardHints renders the keyboard hints bar
func (m Model) renderKeyboardHints() string {
	c := m.CurrentTheme.Colors

	// Style for keys (purple/violet) and descriptions (white)
	keyStyle := lipgloss.NewStyle().Foreground(c.Primary)
	descStyle := lipgloss.NewStyle().Foreground(c.TextPrimary)
	sepStyle := lipgloss.NewStyle().Foreground(c.TextMuted)

	var hintPairs [][]string
	switch m.ActivePane {
	case PaneCalendar:
		hintPairs = [][]string{
			{"h/l", "day"}, {"j/k", "week"}, {"n/p", "month"}, {"T", "today"}, {"Tab", "next"},
		}
	case PaneTasks:
		hintPairs = [][]string{
			{"j/k", "nav"}, {"a", "add"}, {"e", "edit"}, {"d", "del"}, {"v", "details"},
			{"Space", "done"}, {"D", "delegate"}, {"x", "delay"}, {"s", "start"},
			{"n", "next day"}, {"/", "search"}, {"1/2/3", "priority"},
		}
	case PaneTimeline:
		hintPairs = [][]string{
			{"j/k", "scroll"}, {"Shift+C", "clear"}, {"Tab", "next"},
		}
	}

	// Add global hints
	globalHints := [][]string{{"?", "help"}, {":", "overview"}, {"L", "logs"}}
	hintPairs = append(hintPairs, globalHints...)

	// Build hint string with styled keys and descriptions
	var parts []string
	for _, pair := range hintPairs {
		key := keyStyle.Render(pair[0])
		desc := descStyle.Render(pair[1])
		parts = append(parts, key+":"+desc)
	}

	sep := sepStyle.Render(" │ ")
	return strings.Join(parts, sep)
}

// renderWithDialog renders the main layout with a dialog overlay
func (m Model) renderWithDialog() string {
	// Render dialog only - simpler approach that works reliably
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
	case DialogTaskDetails:
		dialog = m.renderTaskDetailsDialog()
	}

	// Center dialog on screen
	return lipgloss.Place(
		m.Width,
		m.Height,
		lipgloss.Center,
		lipgloss.Center,
		dialog,
	)
}

// renderHelpDialog renders the help dialog
func (m Model) renderHelpDialog() string {
	s := m.Styles
	c := m.CurrentTheme.Colors

	var b strings.Builder

	// App header
	appName := lipgloss.NewStyle().Foreground(c.Primary).Bold(true).Render("seyal")
	appDesc := lipgloss.NewStyle().Foreground(c.TextMuted).Render(" - A beautiful terminal task manager")
	b.WriteString(appName + appDesc + "\n")
	b.WriteString(lipgloss.NewStyle().Foreground(c.TextMuted).Render("by krisk248 • github.com/krisk248/seyal") + "\n")
	b.WriteString(strings.Repeat("─", 45) + "\n\n")

	sections := []struct {
		title string
		keys  [][]string
	}{
		{
			title: "Global",
			keys: [][]string{
				{"Ctrl+C", "Exit (press twice)"},
				{"Ctrl+U", "Undo"},
				{"Ctrl+E", "Export"},
				{"?", "This help"},
				{":", "Month overview"},
				{"L", "Jump to logs"},
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
				{"v", "View full details"},
				{"Space", "Toggle complete"},
				{"D", "Delegate task"},
				{"x", "Toggle delayed"},
				{"s", "Start/stop timer"},
				{"n", "Push to next day"},
				{"1/2/3", "Set priority P1/P2/P3"},
				{"0", "Clear priority"},
				{"Enter", "Expand/collapse"},
			},
		},
		{
			title: "Timeline",
			keys: [][]string{
				{"j/k", "Scroll"},
				{"Shift+C", "Clear timeline"},
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

// renderTaskDetailsDialog renders the task details dialog
func (m Model) renderTaskDetailsDialog() string {
	s := m.Styles
	c := m.CurrentTheme.Colors

	task := m.GetSelectedTask()
	if task == nil {
		return s.Dialog.Render("No task selected")
	}

	var b strings.Builder

	// Title
	b.WriteString(s.ModalTitle.Render("Task Details") + "\n\n")

	// Full title (wrapped if needed)
	titleLabel := lipgloss.NewStyle().Foreground(c.Primary).Bold(true).Render("Title:")
	b.WriteString(titleLabel + "\n")
	b.WriteString(lipgloss.NewStyle().Foreground(c.TextPrimary).Render(task.Title) + "\n\n")

	// State
	stateLabel := lipgloss.NewStyle().Foreground(c.Primary).Bold(true).Render("State:")
	stateValue := string(task.State)
	stateStyle := lipgloss.NewStyle().Foreground(c.TextPrimary)
	switch task.State {
	case domain.TaskStateCompleted:
		stateStyle = lipgloss.NewStyle().Foreground(c.TaskCompleted)
	case domain.TaskStateDelegated:
		stateStyle = lipgloss.NewStyle().Foreground(c.TaskDelegated)
	case domain.TaskStateDelayed:
		stateStyle = lipgloss.NewStyle().Foreground(c.TaskDelayed)
	}
	b.WriteString(stateLabel + " " + stateStyle.Render(stateValue) + "\n")

	// Priority
	priorityLabel := lipgloss.NewStyle().Foreground(c.Primary).Bold(true).Render("Priority:")
	priorityValue := "None"
	priorityStyle := lipgloss.NewStyle().Foreground(c.TextMuted)
	switch task.Priority {
	case domain.PriorityHigh:
		priorityValue = "P1 (Critical)"
		priorityStyle = lipgloss.NewStyle().Foreground(c.Error)
	case domain.PriorityMed:
		priorityValue = "P2 (Important)"
		priorityStyle = lipgloss.NewStyle().Foreground(c.Warning)
	case domain.PriorityLow:
		priorityValue = "P3 (Normal)"
		priorityStyle = lipgloss.NewStyle().Foreground(c.TextPrimary)
	}
	b.WriteString(priorityLabel + " " + priorityStyle.Render(priorityValue) + "\n")

	// Created date
	createdLabel := lipgloss.NewStyle().Foreground(c.Primary).Bold(true).Render("Created:")
	createdValue := task.CreatedAt.Format("Jan 2, 2006 3:04 PM")
	b.WriteString(createdLabel + " " + lipgloss.NewStyle().Foreground(c.TextMuted).Render(createdValue) + "\n")

	// Pushed count (if any)
	if task.PushedCount > 0 {
		pushedLabel := lipgloss.NewStyle().Foreground(c.Primary).Bold(true).Render("Pushed:")
		pushedValue := fmt.Sprintf("%d time(s)", task.PushedCount)
		b.WriteString(pushedLabel + " " + lipgloss.NewStyle().Foreground(c.Warning).Render(pushedValue) + "\n")
	}

	// Running status
	if task.IsRunning() {
		runningLabel := lipgloss.NewStyle().Foreground(c.Primary).Bold(true).Render("Status:")
		b.WriteString(runningLabel + " " + lipgloss.NewStyle().Foreground(c.TaskRunning).Render("Running") + "\n")
	}

	b.WriteString("\n")
	b.WriteString(lipgloss.NewStyle().Foreground(c.TextMuted).Render("Press any key to close"))

	return s.Dialog.Render(b.String())
}

// renderHelpScreen renders the full-screen help
func (m Model) renderHelpScreen() string {
	c := m.CurrentTheme.Colors

	// Styles for the help screen
	borderStyle := lipgloss.NewStyle().
		Foreground(c.Border)
	headerStyle := lipgloss.NewStyle().
		Foreground(c.Primary).
		Bold(true)
	keyStyle := lipgloss.NewStyle().
		Foreground(c.Secondary).
		Width(12)
	descStyle := lipgloss.NewStyle().
		Foreground(c.TextPrimary)
	sectionStyle := lipgloss.NewStyle().
		Foreground(c.Primary).
		Bold(true)
	mutedStyle := lipgloss.NewStyle().
		Foreground(c.TextMuted)

	var b strings.Builder

	// Calculate content width (leave some margin)
	contentWidth := m.Width - 4
	if contentWidth > 80 {
		contentWidth = 80
	}

	// Top border
	b.WriteString(borderStyle.Render("╔" + strings.Repeat("═", contentWidth-2) + "╗") + "\n")

	// Header
	title := "SEYAL"
	titlePadding := (contentWidth - 2 - len(title)) / 2
	b.WriteString(borderStyle.Render("║") + strings.Repeat(" ", titlePadding) + headerStyle.Render(title) + strings.Repeat(" ", contentWidth-2-titlePadding-len(title)) + borderStyle.Render("║") + "\n")

	subtitle := "A Terminal Task Manager"
	subtitlePadding := (contentWidth - 2 - len(subtitle)) / 2
	b.WriteString(borderStyle.Render("║") + strings.Repeat(" ", subtitlePadding) + mutedStyle.Render(subtitle) + strings.Repeat(" ", contentWidth-2-subtitlePadding-len(subtitle)) + borderStyle.Render("║") + "\n")

	// Separator
	b.WriteString(borderStyle.Render("╠" + strings.Repeat("═", contentWidth-2) + "╣") + "\n")

	// Empty line
	b.WriteString(borderStyle.Render("║") + strings.Repeat(" ", contentWidth-2) + borderStyle.Render("║") + "\n")

	// Keyboard shortcuts in two columns
	leftCol := []struct {
		title string
		keys  [][]string
	}{
		{
			title: "GLOBAL",
			keys: [][]string{
				{"Ctrl+C", "Exit (press twice)"},
				{"Ctrl+U", "Undo"},
				{"Ctrl+E", "Export"},
				{"?", "Toggle help"},
				{":", "Month overview"},
				{"L", "Jump to logs"},
				{"/", "Search tasks"},
				{"Esc", "Clear search"},
				{"1/2/3", "Switch panes"},
				{"Tab", "Next pane"},
			},
		},
		{
			title: "TASKS",
			keys: [][]string{
				{"j/k", "Navigate"},
				{"a", "Add task"},
				{"e", "Edit task"},
				{"d", "Delete task"},
				{"Space", "Toggle complete"},
				{"D", "Delegate task"},
				{"x", "Toggle delayed"},
				{"s", "Start/stop timer"},
				{"n", "Push to next day"},
				{"1/2/3", "Set priority"},
				{"0", "Clear priority"},
				{"Enter", "Expand/collapse"},
			},
		},
	}

	rightCol := []struct {
		title string
		keys  [][]string
	}{
		{
			title: "CALENDAR",
			keys: [][]string{
				{"h/l", "Previous/next day"},
				{"j/k", "Previous/next week"},
				{"n/p", "Next/prev month"},
				{"T", "Jump to today"},
			},
		},
		{
			title: "TIMELINE",
			keys: [][]string{
				{"j/k", "Scroll"},
				{"Shift+C", "Clear timeline"},
			},
		},
	}

	// Render left column content
	var leftLines []string
	for _, section := range leftCol {
		leftLines = append(leftLines, sectionStyle.Render(section.title))
		leftLines = append(leftLines, mutedStyle.Render(strings.Repeat("─", 28)))
		for _, kv := range section.keys {
			leftLines = append(leftLines, keyStyle.Render(kv[0])+descStyle.Render(kv[1]))
		}
		leftLines = append(leftLines, "")
	}

	// Render right column content
	var rightLines []string
	for _, section := range rightCol {
		rightLines = append(rightLines, sectionStyle.Render(section.title))
		rightLines = append(rightLines, mutedStyle.Render(strings.Repeat("─", 28)))
		for _, kv := range section.keys {
			rightLines = append(rightLines, keyStyle.Render(kv[0])+descStyle.Render(kv[1]))
		}
		rightLines = append(rightLines, "")
	}

	// Combine columns
	maxLines := len(leftLines)
	if len(rightLines) > maxLines {
		maxLines = len(rightLines)
	}

	colWidth := (contentWidth - 6) / 2
	for i := 0; i < maxLines; i++ {
		left := ""
		right := ""
		if i < len(leftLines) {
			left = leftLines[i]
		}
		if i < len(rightLines) {
			right = rightLines[i]
		}

		// Pad left column
		leftLen := lipgloss.Width(left)
		if leftLen < colWidth {
			left = left + strings.Repeat(" ", colWidth-leftLen)
		}

		// Pad right column
		rightLen := lipgloss.Width(right)
		if rightLen < colWidth {
			right = right + strings.Repeat(" ", colWidth-rightLen)
		}

		b.WriteString(borderStyle.Render("║") + "  " + left + "  " + right + borderStyle.Render("║") + "\n")
	}

	// Empty line before footer
	b.WriteString(borderStyle.Render("║") + strings.Repeat(" ", contentWidth-2) + borderStyle.Render("║") + "\n")

	// Footer separator
	b.WriteString(borderStyle.Render("╠" + strings.Repeat("═", contentWidth-2) + "╣") + "\n")

	// Footer line 1: Description
	footerLine1 := "seyal - Track tasks, manage time, stay focused."
	footerPad1 := (contentWidth - 2 - len(footerLine1)) / 2
	b.WriteString(borderStyle.Render("║") + strings.Repeat(" ", footerPad1) + mutedStyle.Render(footerLine1) + strings.Repeat(" ", contentWidth-2-footerPad1-len(footerLine1)) + borderStyle.Render("║") + "\n")

	// Footer line 2: Author, GitHub, Version
	footerLine2 := fmt.Sprintf("by krisk248 • github.com/krisk248/seyal • v%s", Version)
	footerPad2 := (contentWidth - 2 - len(footerLine2)) / 2
	b.WriteString(borderStyle.Render("║") + strings.Repeat(" ", footerPad2) + mutedStyle.Render(footerLine2) + strings.Repeat(" ", contentWidth-2-footerPad2-len(footerLine2)) + borderStyle.Render("║") + "\n")

	// Bottom border
	b.WriteString(borderStyle.Render("╚" + strings.Repeat("═", contentWidth-2) + "╝") + "\n")

	// Close hint
	closeHint := "Press ? to close"
	b.WriteString(strings.Repeat(" ", (contentWidth-len(closeHint))/2) + mutedStyle.Render(closeHint))

	// Center the entire content
	return lipgloss.Place(
		m.Width,
		m.Height,
		lipgloss.Center,
		lipgloss.Center,
		b.String(),
	)
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

// getTaskCheckboxText returns plain text checkbox (for length calculation)
func (m Model) getTaskCheckboxText(task *domain.Task) string {
	switch task.State {
	case domain.TaskStateCompleted:
		return "[✓] "
	case domain.TaskStateDelegated:
		return "[→] "
	case domain.TaskStateDelayed:
		return "[‖] "
	default:
		return "[ ] "
	}
}

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

// getPriorityText returns plain text priority (for length calculation)
func (m Model) getPriorityText(task *domain.Task) string {
	switch task.Priority {
	case domain.PriorityHigh:
		return "P1 "
	case domain.PriorityMed:
		return "P2 "
	case domain.PriorityLow:
		return "P3 "
	default:
		return ""
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
	case domain.EventPushed:
		return c.Warning
	default:
		return c.TextPrimary
	}
}
