package main

import (
	"context"
	"os"

	"github.com/charmbracelet/fang"
	"github.com/spf13/cobra"
)

func main() {
	cmd := &cobra.Command{
		Use:   "seed",
		Short: "Seed the DB with some sample tasks",
	}
	if err := fang.Execute(context.Background(), cmd); err != nil {
		os.Exit(1)
	}
}
