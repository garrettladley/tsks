package model

import (
	"context"
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/garrettladley/tsks/internal/sqlc"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type Model struct {
	querier sqlc.Querier
	list    list.Model
	err     error
	loading bool
}

func New(querier sqlc.Querier) *Model {
	delegate := newTaskDelegate()
	l := list.New([]list.Item{}, delegate, 0, 0)
	l.Title = "Tasks"
	l.SetShowTitle(true)

	l.Styles.Title = lipgloss.NewStyle().
		Foreground(lipgloss.Color("15")).
		Bold(true).
		Padding(0, 1)

	return &Model{
		querier: querier,
		list:    l,
		loading: true,
	}
}

func (m Model) Init() tea.Cmd {
	return loadTasksCmd(m.querier)
}

type tasksLoadedMsg struct {
	tasks []sqlc.Task
	err   error
}

func loadTasksCmd(querier sqlc.Querier) tea.Cmd {
	return func() tea.Msg {
		tasks, err := querier.ListTasks(context.Background())
		return tasksLoadedMsg{tasks: tasks, err: err}
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
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
