package main

import (
	"fmt"

	"github.com/garrettladley/tsks/internal/db"
	"github.com/spf13/cobra"
)

func truncateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "truncate",
		Short: "Truncate all tasks from the database",
		RunE: func(cmd *cobra.Command, args []string) error {
			sqlDB, querier, err := db.Open("tsks.db")
			if err != nil {
				return err
			}
			defer func() {
				_ = sqlDB.Close()
			}()

			ctx := cmd.Context()
			err = querier.TruncateTasks(ctx)
			if err != nil {
				return fmt.Errorf("failed to truncate tasks: %w", err)
			}

			fmt.Println("Successfully truncated all tasks")
			return nil
		},
	}
}
