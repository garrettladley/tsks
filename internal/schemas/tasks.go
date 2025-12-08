package schemas

import (
	"database/sql/driver"
	"fmt"
)

type TaskStatus string

const (
	TaskStatusPending   TaskStatus = "pending"
	TaskStatusCompleted TaskStatus = "completed"
)

func (s TaskStatus) IsValid() bool {
	switch s {
	case TaskStatusPending, TaskStatusCompleted:
		return true
	}
	return false
}

func (s TaskStatus) String() string {
	return string(s)
}

func ParseTaskStatus(s string) (TaskStatus, error) {
	status := TaskStatus(s)
	if !status.IsValid() {
		return "", fmt.Errorf("invalid task status: %s", s)
	}
	return status, nil
}

func (s *TaskStatus) Scan(value any) error {
	if value == nil {
		*s = ""
		return nil
	}

	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("TaskStatus must scan from string, got %T", value)
	}

	parsed, err := ParseTaskStatus(str)
	if err != nil {
		return err
	}

	*s = parsed
	return nil
}

func (s TaskStatus) Value() (driver.Value, error) {
	if !s.IsValid() {
		return nil, fmt.Errorf("invalid TaskStatus: %s", s)
	}
	return string(s), nil
}

var (
	_ driver.Valuer = (*TaskStatus)(nil)
	_ driver.Valuer = TaskStatus("")
)

var _ interface{ Scan(any) error } = (*TaskStatus)(nil)
