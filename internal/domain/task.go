package domain

import (
	"time"

	"github.com/google/uuid"
)

type TaskState string

const (
	TaskStateTodo      TaskState = "todo"
	TaskStateCompleted TaskState = "completed"
	TaskStateDelegated TaskState = "delegated"
	TaskStateDelayed   TaskState = "delayed"
)

type TaskPriority int

const (
	PriorityNone TaskPriority = 0
	PriorityHigh TaskPriority = 1 // P1 - Critical
	PriorityMed  TaskPriority = 2 // P2 - Important
	PriorityLow  TaskPriority = 3 // P3 - Normal
)

type Task struct {
	ID        string       `json:"id"`
	Title     string       `json:"title"`
	State     TaskState    `json:"state"`
	Priority  TaskPriority `json:"priority"`
	CreatedAt time.Time    `json:"createdAt"`
	UpdatedAt time.Time    `json:"updatedAt"`
	StartTime *time.Time   `json:"startTime,omitempty"`
	EndTime   *time.Time   `json:"endTime,omitempty"`
	Children  []*Task      `json:"children,omitempty"`
	ParentID  string       `json:"parentId,omitempty"`
	Date      string       `json:"date"` // YYYY-MM-DD format
	Expanded  bool         `json:"-"`    // UI state, not persisted
}

func NewTask(title, date string) *Task {
	now := time.Now()
	return &Task{
		ID:        uuid.New().String(),
		Title:     title,
		State:     TaskStateTodo,
		Priority:  PriorityNone,
		CreatedAt: now,
		UpdatedAt: now,
		Date:      date,
		Children:  make([]*Task, 0),
		Expanded:  true,
	}
}

func (t *Task) SetState(state TaskState) {
	t.State = state
	t.UpdatedAt = time.Now()
	if state == TaskStateCompleted {
		now := time.Now()
		t.EndTime = &now
	}
}

func (t *Task) SetPriority(priority TaskPriority) {
	t.Priority = priority
	t.UpdatedAt = time.Now()
}

func (t *Task) Start() {
	now := time.Now()
	t.StartTime = &now
	t.State = TaskStateTodo
	t.UpdatedAt = now
}

func (t *Task) Stop() {
	now := time.Now()
	t.EndTime = &now
	t.UpdatedAt = now
}

func (t *Task) AddChild(child *Task) {
	child.ParentID = t.ID
	child.Date = t.Date
	t.Children = append(t.Children, child)
	t.UpdatedAt = time.Now()
}

func (t *Task) RemoveChild(childID string) bool {
	for i, child := range t.Children {
		if child.ID == childID {
			t.Children = append(t.Children[:i], t.Children[i+1:]...)
			t.UpdatedAt = time.Now()
			return true
		}
	}
	return false
}

func (t *Task) IsRunning() bool {
	return t.StartTime != nil && t.EndTime == nil && t.State == TaskStateTodo
}

func (t *Task) ToggleExpanded() {
	t.Expanded = !t.Expanded
}

// TaskTree represents tasks organized by date
type TaskTree map[string][]*Task

func (tt TaskTree) GetTasksForDate(date string) []*Task {
	if tasks, ok := tt[date]; ok {
		return tasks
	}
	return make([]*Task, 0)
}

func (tt TaskTree) AddTask(task *Task) {
	tt[task.Date] = append(tt[task.Date], task)
}

func (tt TaskTree) RemoveTask(date, taskID string) bool {
	tasks := tt[date]
	for i, task := range tasks {
		if task.ID == taskID {
			tt[date] = append(tasks[:i], tasks[i+1:]...)
			return true
		}
		// Check children recursively
		if removeTaskFromChildren(task, taskID) {
			return true
		}
	}
	return false
}

func removeTaskFromChildren(parent *Task, taskID string) bool {
	for i, child := range parent.Children {
		if child.ID == taskID {
			parent.Children = append(parent.Children[:i], parent.Children[i+1:]...)
			return true
		}
		if removeTaskFromChildren(child, taskID) {
			return true
		}
	}
	return false
}

// FlattenedTask represents a task with its depth for rendering
type FlattenedTask struct {
	Task  *Task
	Depth int
}

// FlattenTasks flattens a task tree into a slice for rendering
func FlattenTasks(tasks []*Task, depth int, expandedOnly bool) []FlattenedTask {
	var result []FlattenedTask
	for _, task := range tasks {
		result = append(result, FlattenedTask{Task: task, Depth: depth})
		if len(task.Children) > 0 && (!expandedOnly || task.Expanded) {
			result = append(result, FlattenTasks(task.Children, depth+1, expandedOnly)...)
		}
	}
	return result
}

// GetTaskStats returns completion stats for a list of tasks
func GetTaskStats(tasks []*Task) (total, completed int) {
	for _, task := range tasks {
		total++
		if task.State == TaskStateCompleted {
			completed++
		}
		childTotal, childCompleted := GetTaskStats(task.Children)
		total += childTotal
		completed += childCompleted
	}
	return
}
