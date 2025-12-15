package model

import (
	"context"
	"fmt"
	"slices"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/garrettladley/tsks/internal/sqlc"
	"github.com/garrettladley/tsks/internal/tui/pagination"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type Model struct {
	querier   sqlc.Querier
	list      list.Model
	window    *pagination.TaskListWindow
	err       error
	loading   bool
	lastIndex int
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
		window:  pagination.NewTaskListWindow(pagination.DefaultPageSize),
		loading: true,
	}
}

func (m Model) Init() tea.Cmd {
	return loadPageCmd(m.querier, pagination.TaskListCursor{}, pagination.DirectionForward, 0)
}

type pageLoadedMsg struct {
	page *pagination.TaskListPage
	err  error
}

func loadPageCmd(querier sqlc.Querier, cursor pagination.TaskListCursor, dir pagination.Direction, pageIndex int) tea.Cmd {
	return func() tea.Msg {
		var (
			ctx   = context.Background()
			tasks []sqlc.Task
			err   error
		)

		limit := int64(pagination.DefaultPageSize + 1)

		switch dir {
		case pagination.DirectionForward:
			tasks, err = querier.ListTasksPageForward(ctx, cursor.ForwardParams(limit))
		case pagination.DirectionBackward:
			tasks, err = querier.ListTasksPageBackward(ctx, cursor.BackwardParams(limit))
			slices.Reverse(tasks)
		}

		if err != nil {
			return pageLoadedMsg{err: err}
		}

		hasMore := len(tasks) > pagination.DefaultPageSize
		if hasMore {
			tasks = tasks[:pagination.DefaultPageSize]
		}

		page := &pagination.TaskListPage{
			Tasks:     tasks,
			PageIndex: pageIndex,
		}

		if len(tasks) > 0 {
			page.StartCursor = pagination.TaskListCursorFromTask(tasks[0])
			page.EndCursor = pagination.TaskListCursorFromTask(tasks[len(tasks)-1])
		}

		if dir == pagination.DirectionForward {
			page.HasNext = hasMore
			page.HasPrev = !cursor.IsZero()
		} else {
			page.HasPrev = hasMore
			page.HasNext = true
		}

		return pageLoadedMsg{page: page}
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)

	case pageLoadedMsg:
		m.loading = false
		if msg.err != nil {
			m.err = fmt.Errorf("failed to load tasks: %w", msg.err)
		} else if msg.page != nil {
			m.window.SetPage(msg.page)
			m.refreshListItems()
		}

	case tea.KeyMsg:
		if cmd := m.checkAndLoadMore(); cmd != nil {
			m.loading = true
			cmds = append(cmds, cmd)
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)

	newIndex := m.list.Index()
	if newIndex != m.lastIndex {
		m.lastIndex = newIndex
		pageIdx, _ := m.window.GlobalIndexToLocal(newIndex)
		m.window.SetCurrentPage(pageIdx)
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) checkAndLoadMore() tea.Cmd {
	if m.loading {
		return nil
	}

	currentIdx := m.list.Index()

	if m.window.ShouldLoadNext(currentIdx) {
		cursor, pageIdx, ok := m.window.GetNextPageCursor()
		if ok {
			return loadPageCmd(m.querier, cursor, pagination.DirectionForward, pageIdx)
		}
	}

	if m.window.ShouldLoadPrev(currentIdx) {
		cursor, pageIdx, ok := m.window.GetPrevPageCursor()
		if ok {
			return loadPageCmd(m.querier, cursor, pagination.DirectionBackward, pageIdx)
		}
	}

	return nil
}

func (m *Model) refreshListItems() {
	tasks := m.window.GetAllItems()
	items := make([]list.Item, len(tasks))
	for i, t := range tasks {
		items[i] = *newTask(t)
	}

	currentIdx := m.list.Index()
	m.list.SetItems(items)
	if currentIdx < len(items) {
		m.list.Select(currentIdx)
	}
}

func (m Model) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v\n", m.err)
	}

	if m.loading && m.window.ItemCount() == 0 {
		return "Loading tasks..."
	}

	view := m.list.View()
	if m.loading {
		view += "\n Loading more..."
	}
	return view
}
