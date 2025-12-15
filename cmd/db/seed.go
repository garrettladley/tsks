package main

import (
	"fmt"

	"github.com/garrettladley/tsks/internal/db"
	"github.com/garrettladley/tsks/internal/ids"
	"github.com/garrettladley/tsks/internal/schemas"
	"github.com/garrettladley/tsks/internal/sqlc"
	"github.com/spf13/cobra"
)

func seedCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "seed",
		Short: "Seed the DB with some sample tasks",
		RunE: func(cmd *cobra.Command, args []string) error {
			tasks := []struct {
				name        string
				status      schemas.TaskStatus
				description *string
			}{
				{name: "Water the couch", status: schemas.TaskStatusPending},
				{name: "Walk my goldfish", status: schemas.TaskStatusPending, description: strPtr("They were parmesan flavored")},
				{name: "Rewrite everything in Go", status: schemas.TaskStatusCompleted, description: strPtr("Mission accomplished!")},
				{name: "Learn bubbletea", status: schemas.TaskStatusCompleted},
			}

			sqlDB, querier, err := db.Open("tsks.db")
			if err != nil {
				return err
			}
			defer func() {
				_ = sqlDB.Close()
			}()

			ctx := cmd.Context()
			for _, task := range tasks {
				taskID := ids.New("task")
				_, err := querier.CreateTask(ctx, sqlc.CreateTaskParams{
					ID:          taskID,
					Title:       task.name,
					Description: task.description,
					Status:      task.status,
				})
				if err != nil {
					return fmt.Errorf("failed to create task '%s': %w", task.name, err)
				}
				fmt.Printf("Created task: %s (ID: %s)\n", task.name, taskID)
			}

			fmt.Printf("\nSuccessfully seeded %d tasks\n", len(tasks))
			return nil
		},
	}
}

func strPtr(s string) *string {
	return &s
}
