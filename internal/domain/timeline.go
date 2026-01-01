package domain

import (
	"time"

	"github.com/google/uuid"
)

type TimelineEventType string

const (
	EventCreated   TimelineEventType = "created"
	EventStarted   TimelineEventType = "started"
	EventCompleted TimelineEventType = "completed"
	EventDelegated TimelineEventType = "delegated"
	EventDelayed   TimelineEventType = "delayed"
	EventUpdated   TimelineEventType = "updated"
)

type TimelineEvent struct {
	ID            string            `json:"id"`
	TaskID        string            `json:"taskId"`
	TaskTitle     string            `json:"taskTitle"`
	Type          TimelineEventType `json:"type"`
	Timestamp     time.Time         `json:"timestamp"`
	PreviousState TaskState         `json:"previousState,omitempty"`
	NewState      TaskState         `json:"newState,omitempty"`
}

func NewTimelineEvent(taskID, taskTitle string, eventType TimelineEventType) *TimelineEvent {
	return &TimelineEvent{
		ID:        uuid.New().String(),
		TaskID:    taskID,
		TaskTitle: taskTitle,
		Type:      eventType,
		Timestamp: time.Now(),
	}
}

func NewStateChangeEvent(taskID, taskTitle string, prevState, newState TaskState) *TimelineEvent {
	var eventType TimelineEventType
	switch newState {
	case TaskStateCompleted:
		eventType = EventCompleted
	case TaskStateDelegated:
		eventType = EventDelegated
	case TaskStateDelayed:
		eventType = EventDelayed
	default:
		eventType = EventUpdated
	}

	return &TimelineEvent{
		ID:            uuid.New().String(),
		TaskID:        taskID,
		TaskTitle:     taskTitle,
		Type:          eventType,
		Timestamp:     time.Now(),
		PreviousState: prevState,
		NewState:      newState,
	}
}

// Timeline represents events organized by date
type Timeline map[string][]*TimelineEvent

func (t Timeline) GetEventsForDate(date string) []*TimelineEvent {
	if events, ok := t[date]; ok {
		return events
	}
	return make([]*TimelineEvent, 0)
}

func (t Timeline) AddEvent(date string, event *TimelineEvent) {
	t[date] = append(t[date], event)
}

func (t Timeline) RemoveEventsByTaskID(date, taskID string) {
	events := t[date]
	filtered := make([]*TimelineEvent, 0)
	for _, event := range events {
		if event.TaskID != taskID {
			filtered = append(filtered, event)
		}
	}
	t[date] = filtered
}

func (t Timeline) ClearDate(date string) {
	t[date] = make([]*TimelineEvent, 0)
}

// GetEventIcon returns an icon for the event type
func (e *TimelineEvent) GetEventIcon() string {
	switch e.Type {
	case EventStarted:
		return "○"
	case EventCompleted:
		return "●"
	case EventDelegated:
		return "→"
	case EventDelayed:
		return "‖"
	case EventCreated:
		return "+"
	default:
		return "•"
	}
}

// GetEventDescription returns a human-readable description
func (e *TimelineEvent) GetEventDescription() string {
	switch e.Type {
	case EventStarted:
		return "started"
	case EventCompleted:
		return "completed"
	case EventDelegated:
		return "delegated"
	case EventDelayed:
		return "delayed"
	case EventCreated:
		return "created"
	default:
		return "updated"
	}
}
