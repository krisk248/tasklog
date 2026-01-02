package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/krisk248/nexus/internal/domain"
)

// StorageSchema represents the data structure saved to disk
type StorageSchema struct {
	Version  string                           `json:"version"`
	Tasks    domain.TaskTree                  `json:"tasks"`
	Timeline domain.Timeline                  `json:"timeline"`
	Settings Settings                         `json:"settings"`
}

// Settings holds user preferences
type Settings struct {
	Theme          string `json:"theme"`
	DateFormat     string `json:"dateFormat"`
	TimeFormat     string `json:"timeFormat"`
	SkippedVersion string `json:"skippedVersion,omitempty"`
}

// Storage handles data persistence
type Storage struct {
	DataPath string
}

// NewStorage creates a new storage instance
func NewStorage() (*Storage, error) {
	dataPath, err := getDataPath()
	if err != nil {
		return nil, err
	}

	// Ensure directory exists
	dir := filepath.Dir(dataPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	return &Storage{DataPath: dataPath}, nil
}

// getDataPath returns the platform-specific data file path
func getDataPath() (string, error) {
	var configDir string

	switch runtime.GOOS {
	case "darwin":
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		configDir = filepath.Join(home, "Library", "Application Support", "nexus")
	case "windows":
		configDir = filepath.Join(os.Getenv("APPDATA"), "nexus")
	default: // Linux and others
		// Check XDG_DATA_HOME first
		if xdg := os.Getenv("XDG_DATA_HOME"); xdg != "" {
			configDir = filepath.Join(xdg, "nexus")
		} else {
			home, err := os.UserHomeDir()
			if err != nil {
				return "", err
			}
			configDir = filepath.Join(home, ".local", "share", "nexus")
		}
	}

	return filepath.Join(configDir, "data.json"), nil
}

// Load reads the data file and returns the stored data
func (s *Storage) Load() (*StorageSchema, error) {
	data, err := os.ReadFile(s.DataPath)
	if err != nil {
		if os.IsNotExist(err) {
			// Return default schema
			return s.defaultSchema(), nil
		}
		return nil, err
	}

	var schema StorageSchema
	if err := json.Unmarshal(data, &schema); err != nil {
		return nil, err
	}

	// Hydrate dates (JSON stores as strings, need to parse)
	s.hydrateDates(&schema)

	return &schema, nil
}

// Save writes the data to disk
func (s *Storage) Save(schema *StorageSchema) error {
	schema.Version = "1.0.0"

	data, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.DataPath, data, 0644)
}

// Backup creates a backup of the current data file
func (s *Storage) Backup() error {
	if _, err := os.Stat(s.DataPath); os.IsNotExist(err) {
		return nil // Nothing to backup
	}

	timestamp := time.Now().Format("20060102-150405")
	backupPath := s.DataPath + ".backup-" + timestamp

	data, err := os.ReadFile(s.DataPath)
	if err != nil {
		return err
	}

	return os.WriteFile(backupPath, data, 0644)
}

// defaultSchema returns an empty default schema
func (s *Storage) defaultSchema() *StorageSchema {
	return &StorageSchema{
		Version:  "1.0.0",
		Tasks:    make(domain.TaskTree),
		Timeline: make(domain.Timeline),
		Settings: Settings{
			Theme:      "ultraviolet",
			DateFormat: "January 2, 2006",
			TimeFormat: "12h",
		},
	}
}

// hydrateDates converts date strings back to proper Date objects if needed
func (s *Storage) hydrateDates(schema *StorageSchema) {
	// Tasks and Timeline already use time.Time which JSON handles correctly
	// with proper ISO format. This method is a placeholder for any
	// additional date hydration needed.
}

// GetExportPath returns the platform-specific export folder (Documents folder)
func GetExportPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	var exportDir string
	switch runtime.GOOS {
	case "darwin":
		exportDir = filepath.Join(home, "Documents", "nexus-exports")
	case "windows":
		// Use Documents folder on Windows
		exportDir = filepath.Join(home, "Documents", "nexus-exports")
	default: // Linux and others
		exportDir = filepath.Join(home, "Documents", "nexus-exports")
	}

	// Ensure directory exists
	if err := os.MkdirAll(exportDir, 0755); err != nil {
		return "", err
	}

	return exportDir, nil
}

// ExportToFile exports tasks to a file in the export folder
func (s *Storage) ExportToFile(tasks domain.TaskTree, format ExportFormat, scope string, filename string) (string, error) {
	content, err := s.Export(tasks, format, scope)
	if err != nil {
		return "", err
	}

	exportDir, err := GetExportPath()
	if err != nil {
		return "", err
	}

	filePath := filepath.Join(exportDir, filename)
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return "", err
	}

	return filePath, nil
}

// Export exports tasks to various formats
type ExportFormat int

const (
	FormatMarkdown ExportFormat = iota
	FormatJSON
	FormatPlainText
)

// Export exports tasks to the specified format (returns content as string)
func (s *Storage) Export(tasks domain.TaskTree, format ExportFormat, scope string) (string, error) {
	switch format {
	case FormatMarkdown:
		return s.exportMarkdown(tasks, scope)
	case FormatJSON:
		return s.exportJSON(tasks, scope)
	case FormatPlainText:
		return s.exportPlainText(tasks, scope)
	default:
		return s.exportMarkdown(tasks, scope)
	}
}

func (s *Storage) exportMarkdown(tasks domain.TaskTree, scope string) (string, error) {
	var result string

	for date, taskList := range tasks {
		if scope != "all" && scope != date {
			continue
		}

		result += "## " + date + "\n\n"
		for _, task := range taskList {
			result += s.taskToMarkdown(task, 0)
		}
		result += "\n"
	}

	return result, nil
}

func (s *Storage) taskToMarkdown(task *domain.Task, depth int) string {
	indent := ""
	for i := 0; i < depth; i++ {
		indent += "  "
	}

	checkbox := "[ ]"
	if task.State == domain.TaskStateCompleted {
		checkbox = "[x]"
	}

	priority := ""
	switch task.Priority {
	case domain.PriorityHigh:
		priority = " **P1**"
	case domain.PriorityMed:
		priority = " *P2*"
	case domain.PriorityLow:
		priority = " P3"
	}

	result := indent + "- " + checkbox + " " + task.Title + priority + "\n"

	for _, child := range task.Children {
		result += s.taskToMarkdown(child, depth+1)
	}

	return result
}

func (s *Storage) exportJSON(tasks domain.TaskTree, scope string) (string, error) {
	var filtered domain.TaskTree
	if scope == "all" {
		filtered = tasks
	} else {
		filtered = make(domain.TaskTree)
		if taskList, ok := tasks[scope]; ok {
			filtered[scope] = taskList
		}
	}

	data, err := json.MarshalIndent(filtered, "", "  ")
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (s *Storage) exportPlainText(tasks domain.TaskTree, scope string) (string, error) {
	var result string

	for date, taskList := range tasks {
		if scope != "all" && scope != date {
			continue
		}

		result += date + "\n"
		result += "─────────────────────\n"
		for _, task := range taskList {
			result += s.taskToPlainText(task, 0)
		}
		result += "\n"
	}

	return result, nil
}

func (s *Storage) taskToPlainText(task *domain.Task, depth int) string {
	indent := ""
	for i := 0; i < depth; i++ {
		indent += "  "
	}

	status := "○"
	if task.State == domain.TaskStateCompleted {
		status = "●"
	} else if task.State == domain.TaskStateDelegated {
		status = "→"
	} else if task.State == domain.TaskStateDelayed {
		status = "‖"
	}

	result := indent + status + " " + task.Title + "\n"

	for _, child := range task.Children {
		result += s.taskToPlainText(child, depth+1)
	}

	return result
}
