package main

import (
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"path/filepath"
)

// Update handles each keypress and updates the model's state
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch m.state {
	// Handle booting state (progress bar)
	case stateBooting:
		switch msg := msg.(type) {
		case tea.WindowSizeMsg:
			m.terminalWidth = msg.Width
			m.terminalHeight = msg.Height

			m.progress.Width = msg.Width - padding*2 - 4
			if m.progress.Width > maxWidth {
				m.progress.Width = maxWidth
			}
			return m, nil

		case tickMsg:
			if m.progress.Percent() >= 1.0 {
				// Switch to listing state once progress completes
				m.state = stateListing
				return m, nil
			}

			// Increase progress bar incrementally
			cmd := m.progress.IncrPercent(0.25)
			return m, tea.Batch(tickCmd(), cmd)

		case progress.FrameMsg:
			progressModel, cmd := m.progress.Update(msg)
			m.progress = progressModel.(progress.Model)
			return m, cmd
		}

	// Handle directory listing state
	case stateListing:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "esc":
				if m.fileTable.Focused() {
					m.fileTable.Blur()
				} else {
					m.fileTable.Focus()
				}
			case "q", "ctrl+c":
				return m, tea.Quit
			case "enter":
				selectedRow := m.fileTable.SelectedRow()
				selectedName := selectedRow[0]
				selectedType := selectedRow[1]

				// Track current directory before changing
				m.previousDir = m.currentDir

				if selectedName == ".. (Go Back)" {
					// Go back to the parent directory
					m.currentDir = filepath.Dir(m.currentDir)
				} else if selectedType == "Directory" {
					// Go into the selected directory
					m.currentDir = filepath.Join(m.currentDir, selectedName)
				} else {
					// Set the file name and retrieve the content
					m.viewingFile = selectedName
					fullFilePath := filepath.Join(m.currentDir, selectedName)

					if content, ok := m.fileContents[fullFilePath]; ok {
						m.currentFileContents = content
					} else {
						m.currentFileContents = "File content not found."
					}

					m.state = stateViewing
				}

				// Update the table with the new directory content
				m.fileTable.SetRows(generateRows(m.currentDir, m.directories))
			}
		}
		m.fileTable, cmd = m.fileTable.Update(msg)
		return m, cmd

	// Handle file viewing state
	case stateViewing:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			if msg.String() == "esc" || msg.String() == "q" {
				// Go back to the listing state
				m.state = stateListing
			}
		}
	}

	return m, cmd
}
