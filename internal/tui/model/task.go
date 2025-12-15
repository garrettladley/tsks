package model

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/garrettladley/tsks/internal/schemas"
	"github.com/garrettladley/tsks/internal/sqlc"
)

type task struct {
	t sqlc.Task
}

var _ list.Item = (*task)(nil)

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

type taskDelegate struct{}

var _ list.ItemDelegate = (*taskDelegate)(nil)

func newTaskDelegate() taskDelegate {
	return taskDelegate{}
}

func (d taskDelegate) Height() int {
	return 2
}

func (d taskDelegate) Spacing() int {
	return 0
}

func (d taskDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return nil
}

func (d taskDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	t, ok := item.(task)
	if !ok {
		return
	}

	var (
		icon        string
		statusStyle lipgloss.Style
		titleStyle  lipgloss.Style
		descStyle   lipgloss.Style
	)

	isSelected := index == m.Index()

	switch t.t.Status {
	case schemas.TaskStatusCompleted:
		icon = "✓"
		statusStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Bold(true)
		titleStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Strikethrough(true)
		descStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	case schemas.TaskStatusPending:
		icon = "○"
		statusStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("3"))
		titleStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("15"))
		descStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
	}

	if isSelected {
		statusStyle = statusStyle.Underline(true)
		titleStyle = titleStyle.Underline(true).Bold(true)
		descStyle = descStyle.Underline(true)
	}

	cursor := "  "
	if isSelected {
		cursor = "> "
	}

	title := fmt.Sprintf("%s%s %s",
		cursor,
		statusStyle.Render(icon),
		titleStyle.Render(t.t.Title),
	)

	var desc string
	if t.t.Description != nil && *t.t.Description != "" {
		desc = "\n  " + descStyle.Render(*t.t.Description)
	}

	_, _ = fmt.Fprintf(w, "%s%s\n", title, desc)
}
