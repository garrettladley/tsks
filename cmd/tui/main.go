package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/garrettladley/tsks/internal/db"
	"github.com/garrettladley/tsks/internal/tui/model"
)

func main() {
	sqlDB, querier, err := db.Open("tsks.db")
	if err != nil {
		log.Fatalf("failed to init the db: %v", err)
	}
	defer func() {
		_ = sqlDB.Close()
	}()

	m := model.New(querier)
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		log.Fatalf("error while running the program")
	}
}
