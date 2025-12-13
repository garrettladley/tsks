package model

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/garrettladley/tsks/internal/sqlc"
)

type task struct {
	t sqlc.Task
}

func newTask(t sqlc.Task) *task {
	return &task{
		t: t,
	}
}

func (t task) Title() string { return t.t.Title }

func (t task) FilterValue() string { return t.t.Title }

func (t task) Description() string {
	if t.t.Description == nil {
		return ""
	}
	return *t.t.Description
}

var _ list.Item = (*task)(nil)
