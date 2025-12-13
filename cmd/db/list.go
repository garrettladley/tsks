package main

import (
	"fmt"

	"github.com/garrettladley/tsks/internal/db"
	"github.com/spf13/cobra"
)

func listCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all tasks",
		RunE: func(cmd *cobra.Command, args []string) error {
			sqlDB, querier, err := db.Open("tsks.db")
			if err != nil {
				return err
			}
			defer func() {
				_ = sqlDB.Close()
			}()

			ctx := cmd.Context()
			tasks, err := querier.ListTasks(ctx)
			if err != nil {
				return fmt.Errorf("failed to list tasks: %w", err)
			}

			if len(tasks) == 0 {
				fmt.Println("No tasks found")
				return nil
			}

			fmt.Println("\nTasks:")
			fmt.Println("─────────────────────────────────────────────────────────────")
			for _, task := range tasks {
				archivedAt := ""
				if task.ArchivedAt != nil {
					archivedAt = fmt.Sprintf(" [ARCHIVED: %s]", task.ArchivedAt.Format("2006-01-02"))
				}
				description := ""
				if task.Description != nil {
					description = *task.Description
				}
				fmt.Printf("ID:          %s\n", task.ID)
				fmt.Printf("Title:       %s%s\n", task.Title, archivedAt)
				fmt.Printf("Description: %s\n", description)
				fmt.Printf("Status:      %s\n", task.Status)
				fmt.Printf("Created:     %s\n", task.CreatedAt.Format("2006-01-02 15:04:05"))
				fmt.Println("─────────────────────────────────────────────────────────────")
			}

			fmt.Printf("\nTotal: %d task(s)\n", len(tasks))
			return nil
		},
	}
}
