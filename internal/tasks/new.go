package tasks

import (
	"time"

	"github.com/garrettladley/tsks/internal/ids"
	"github.com/garrettladley/tsks/internal/schemas"
	"github.com/garrettladley/tsks/internal/sqlc"
)

type Option func(*sqlc.Task)

func WithID(id string) Option {
	return func(t *sqlc.Task) {
		t.ID = id
	}
}

func WithDescription(description string) Option {
	return func(t *sqlc.Task) {
		t.Description = &description
	}
}

func WithArchivedAt(archivedAt time.Time) Option {
	return func(t *sqlc.Task) {
		t.ArchivedAt = &archivedAt
	}
}

func New(title string, status schemas.TaskStatus, options ...Option) *sqlc.Task {
	task := &sqlc.Task{
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
