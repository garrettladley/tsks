package main

import (
	"fmt"

	"github.com/garrettladley/tsks/internal/db"
	"github.com/spf13/cobra"
)

func migrateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "migrate",
		Short: "Apply pending database migrations",
		RunE: func(cmd *cobra.Command, args []string) error {
			sqlDB, _, err := db.Open("tsks.db")
			if err != nil {
				return err
			}
			defer sqlDB.Close()

			fmt.Println("Migrations applied successfully")
			return nil
		},
	}
}
