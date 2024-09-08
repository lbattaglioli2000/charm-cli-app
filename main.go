package main

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	padding  = 2
	maxWidth = 80
)

type tickMsg time.Time

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// Init is called once at the start of the app
func (m model) Init() tea.Cmd {
	return tickCmd()
}

func main() {
	// Create a new BubbleTea app from the initial model defined
	p := tea.NewProgram(
		initialModel(),
		tea.WithAltScreen(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
