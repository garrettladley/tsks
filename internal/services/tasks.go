package services

import (
	"context"

	"github.com/garrettladley/tsks/internal/models"
	"github.com/garrettladley/tsks/internal/schemas"
	"github.com/garrettladley/tsks/internal/sqlc"
)

type TaskService struct {
	queries *sqlc.Queries
}

func NewTaskService(queries *sqlc.Queries) *TaskService {
	return &TaskService{queries: queries}
}

func (s *TaskService) CreateTask(ctx context.Context, id, title string, description *string, status schemas.TaskStatus) (*models.Task, error) {
	sqlcTask, err := s.queries.CreateTask(ctx, sqlc.CreateTaskParams{
		ID:          id,
		Title:       title,
		Description: description,
		Status:      status,
	})
	if err != nil {
		return nil, err
	}
	return models.FromSQLCTaskPtr(sqlcTask), nil
}

func (s *TaskService) GetTask(ctx context.Context, id string) (*models.Task, error) {
	sqlcTask, err := s.queries.GetTask(ctx, id)
	if err != nil {
		return nil, err
	}
	return models.FromSQLCTaskPtr(sqlcTask), nil
}

func (s *TaskService) ListTasks(ctx context.Context) ([]models.Task, error) {
	sqlcTasks, err := s.queries.ListTasks(ctx)
	if err != nil {
		return nil, err
	}
	return models.FromSQLCTasks(sqlcTasks), nil
}

func (s *TaskService) UpdateTask(ctx context.Context, id, title string, description *string, status schemas.TaskStatus) (*models.Task, error) {
	sqlcTask, err := s.queries.UpdateTask(ctx, sqlc.UpdateTaskParams{
		ID:          id,
		Title:       title,
		Description: description,
		Status:      status,
	})
	if err != nil {
		return nil, err
	}
	return models.FromSQLCTaskPtr(sqlcTask), nil
}

func (s *TaskService) DeleteTask(ctx context.Context, id string) error {
	return s.queries.DeleteTask(ctx, id)
}
