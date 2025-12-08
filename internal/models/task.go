package models

import (
	"time"

	"github.com/garrettladley/tsks/internal/schemas"
	"github.com/garrettladley/tsks/internal/sqlc"
)

type Task struct {
	ID          string             `json:"id"`
	Version     time.Time          `json:"version"`
	Title       string             `json:"title"`
	Description *string            `json:"description"`
	Status      schemas.TaskStatus `json:"status"`
	CreatedAt   time.Time          `json:"created_at"`
	ArchivedAt  *time.Time         `json:"archived_at"`
}

func FromSQLCTask(t sqlc.Task) Task {
	return Task{
		ID:          t.ID,
		Version:     t.Version,
		Title:       t.Title,
		Description: t.Description,
		Status:      t.Status,
		CreatedAt:   t.CreatedAt,
		ArchivedAt:  t.ArchivedAt,
	}
}

func FromSQLCTaskPtr(t sqlc.Task) *Task {
	task := FromSQLCTask(t)
	return &task
}

func FromSQLCTasks(tasks []sqlc.Task) []Task {
	result := make([]Task, len(tasks))
	for i, t := range tasks {
		result[i] = FromSQLCTask(t)
	}
	return result
}

func (t Task) ToSQLCTask() sqlc.Task {
	return sqlc.Task{
		ID:          t.ID,
		Version:     t.Version,
		Title:       t.Title,
		Description: t.Description,
		Status:      t.Status,
		CreatedAt:   t.CreatedAt,
		ArchivedAt:  t.ArchivedAt,
	}
}
