package tasks

import (
	"time"

	"github.com/garrettladley/tsks/internal/ids"
	"github.com/garrettladley/tsks/internal/models"
	"github.com/garrettladley/tsks/internal/schemas"
)

type Option func(*models.Task)

func WithID(id string) Option {
	return func(t *models.Task) {
		t.ID = id
	}
}

func WithDescription(description string) Option {
	return func(t *models.Task) {
		t.Description = &description
	}
}

func WithArchivedAt(archivedAt time.Time) Option {
	return func(t *models.Task) {
		t.ArchivedAt = &archivedAt
	}
}

func New(title string, status schemas.TaskStatus, options ...Option) *models.Task {
	task := &models.Task{
		ID:     newTaskID(),
		Title:  title,
		Status: status,
	}
	for _, option := range options {
		option(task)
	}
	return task
}

func newTaskID() string {
	const prefix = "task"
	return ids.New(prefix)
}
