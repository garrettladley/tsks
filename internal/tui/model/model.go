package model

import (
	"context"
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/garrettladley/tsks/internal/sqlc"
)

type Model struct {
	querier sqlc.Querier
	list    list.Model
	err     error
	loading bool
}

type tasksLoadedMsg struct {
	tasks []sqlc.Task
	err   error
}

func New(querier sqlc.Querier) *Model {
	l := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Tasks"
	return &Model{
		querier: querier,
		list:    l,
		loading: true,
	}
}

func loadTasksCmd(querier sqlc.Querier) tea.Cmd {
	return func() tea.Msg {
		tasks, err := querier.ListTasks(context.Background())
		return tasksLoadedMsg{tasks: tasks, err: err}
	}
}

func (m Model) Init() tea.Cmd {
	return loadTasksCmd(m.querier)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height)
	case tasksLoadedMsg:
		m.loading = false
		if msg.err != nil {
			m.err = fmt.Errorf("failed to list all tasks: %w", msg.err)
		} else {
			tasks := make([]list.Item, len(msg.tasks))
			for i, t := range msg.tasks {
				tasks[i] = *newTask(t)
			}
			m.list.SetItems(tasks)
		}
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v\n", m.err)
	}
	if m.loading {
		return "Loading tasks..."
	}
	return m.list.View()
}
