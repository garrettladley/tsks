package pagination

import (
	"time"

	"github.com/garrettladley/tsks/internal/sqlc"
)

const sqliteDatetimeFormat = "2006-01-02 15:04:05"

type TaskListCursor struct {
	CreatedAt time.Time
	ID        string
}

func TaskListCursorFromTask(t sqlc.Task) TaskListCursor {
	return TaskListCursor{CreatedAt: t.CreatedAt, ID: t.ID}
}

func (c TaskListCursor) IsZero() bool {
	return c.CreatedAt.IsZero() && c.ID == ""
}

func (c TaskListCursor) ForwardParams(limit int64) sqlc.ListTasksPageForwardParams {
	if c.IsZero() {
		return sqlc.ListTasksPageForwardParams{
			CursorCreatedAt: nil,
			CursorID:        nil,
			PageLimit:       limit,
		}
	}
	ca := c.CreatedAt.Format(sqliteDatetimeFormat)
	return sqlc.ListTasksPageForwardParams{
		CursorCreatedAt: &ca,
		CursorID:        &c.ID,
		PageLimit:       limit,
	}
}

func (c TaskListCursor) BackwardParams(limit int64) sqlc.ListTasksPageBackwardParams {
	return sqlc.ListTasksPageBackwardParams{
		CursorCreatedAt: c.CreatedAt.Format(sqliteDatetimeFormat),
		CursorID:        c.ID,
		PageLimit:       limit,
	}
}

type Direction byte

const (
	DirectionForward Direction = iota
	DirectionBackward
)

type TaskListPage struct {
	Tasks       []sqlc.Task
	StartCursor TaskListCursor
	EndCursor   TaskListCursor
	HasPrev     bool
	HasNext     bool
	PageIndex   int
}

const (
	DefaultPageSize  = 50
	MaxPagesInMemory = 3 // prev + current + next
	LoadThreshold    = 10
)

type TaskListWindow struct {
	pages       map[int]*TaskListPage // pageIndex -> Page
	currentPage int
	pageSize    int

	firstLoadedPage int
	lastLoadedPage  int
}

func NewTaskListWindow(pageSize int) *TaskListWindow {
	if pageSize <= 0 {
		pageSize = DefaultPageSize
	}
	return &TaskListWindow{
		pages:           make(map[int]*TaskListPage),
		pageSize:        pageSize,
		firstLoadedPage: -1,
		lastLoadedPage:  -1,
	}
}

func (w *TaskListWindow) SetPage(page *TaskListPage) {
	w.pages[page.PageIndex] = page
	w.updateLoadedRange()
	w.purgeDistantPages()
}

func (w *TaskListWindow) GetAllItems() []sqlc.Task {
	if w.firstLoadedPage < 0 {
		return nil
	}

	var items []sqlc.Task
	for i := w.firstLoadedPage; i <= w.lastLoadedPage; i++ {
		if p, ok := w.pages[i]; ok {
			items = append(items, p.Tasks...)
		}
	}
	return items
}

func (w *TaskListWindow) ItemCount() int {
	var count int
	for _, p := range w.pages {
		count += len(p.Tasks)
	}
	return count
}

func (w *TaskListWindow) GlobalIndexToLocal(globalIdx int) (pageIndex, localIndex int) {
	var offset int
	for i := w.firstLoadedPage; i < 0; i++ {
		if p, ok := w.pages[i]; ok {
			offset += len(p.Tasks)
		}
	}

	adjustedIdx := globalIdx - offset
	if adjustedIdx < 0 {
		remaining := globalIdx
		for i := w.firstLoadedPage; i < 0; i++ {
			if p, ok := w.pages[i]; ok {
				if remaining < len(p.Tasks) {
					return i, remaining
				}
				remaining -= len(p.Tasks)
			}
		}
	}

	pageIndex = adjustedIdx / w.pageSize
	localIndex = adjustedIdx % w.pageSize
	return
}

func (w *TaskListWindow) ShouldLoadNext(currentIdx int) bool {
	if w.lastLoadedPage < 0 {
		return false
	}

	lastPage, ok := w.pages[w.lastLoadedPage]
	if !ok || !lastPage.HasNext {
		return false
	}

	totalItems := w.ItemCount()
	itemsFromEnd := totalItems - 1 - currentIdx

	return itemsFromEnd < LoadThreshold
}

func (w *TaskListWindow) ShouldLoadPrev(currentIdx int) bool {
	if w.firstLoadedPage < 0 {
		return false
	}

	firstPage, ok := w.pages[w.firstLoadedPage]
	if !ok || !firstPage.HasPrev {
		return false
	}

	return currentIdx < LoadThreshold
}

func (w *TaskListWindow) GetNextPageCursor() (TaskListCursor, int, bool) {
	if p, ok := w.pages[w.lastLoadedPage]; ok && p.HasNext {
		return p.EndCursor, w.lastLoadedPage + 1, true
	}
	return TaskListCursor{}, 0, false
}

func (w *TaskListWindow) GetPrevPageCursor() (TaskListCursor, int, bool) {
	if p, ok := w.pages[w.firstLoadedPage]; ok && p.HasPrev {
		return p.StartCursor, w.firstLoadedPage - 1, true
	}
	return TaskListCursor{}, 0, false
}

func (w *TaskListWindow) PageSize() int {
	return w.pageSize
}

func (w *TaskListWindow) updateLoadedRange() {
	var (
		first = -1
		last  = -1
	)
	for idx := range w.pages {
		if first == -1 || idx < first {
			first = idx
		}
		if last == -1 || idx > last {
			last = idx
		}
	}
	w.firstLoadedPage = first
	w.lastLoadedPage = last
}

func (w *TaskListWindow) purgeDistantPages() {
	if len(w.pages) <= MaxPagesInMemory {
		return
	}

	var (
		minKeep = w.currentPage - 1
		maxKeep = w.currentPage + 1
	)

	for idx := range w.pages {
		if idx < minKeep || idx > maxKeep {
			delete(w.pages, idx)
		}
	}
	w.updateLoadedRange()
}

func (w *TaskListWindow) SetCurrentPage(pageIdx int) {
	if w.currentPage != pageIdx {
		w.currentPage = pageIdx
		w.purgeDistantPages()
	}
}

func (w *TaskListWindow) CurrentPage() int {
	return w.currentPage
}

func (w *TaskListWindow) FirstLoadedPage() int {
	return w.firstLoadedPage
}

func (w *TaskListWindow) LastLoadedPage() int {
	return w.lastLoadedPage
}
