package main

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
)

var baseStyle = renderer.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

var helpStyle = renderer.NewStyle().Foreground(lipgloss.Color("#626262")).Render

var fileWindowStyle = renderer.NewStyle().
	BorderStyle(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("62")).
	Padding(1, 2).
	Width(60)

var centeredStyle = renderer.NewStyle().
	Margin(0, 2). // Add top/bottom and left/right margin
	Width(60).    // Set the width of the content
	Align(lipgloss.Center, lipgloss.Center)

// View renders the different states for the app
func (m model) View() string {
	// Create the centered container style with dynamic height and width
	centered := centeredStyle.
		Width(m.terminalWidth).
		Height(m.terminalHeight)

	switch m.state {
	case stateBooting:
		// Center the booting message
		return centered.Render(m.progress.View() + "\n\n" + helpStyle("Booting system..."))

	case stateListing:
		// Center the file table
		return centered.Render(baseStyle.Render(m.fileTable.View()) + "\n" + m.fileTable.HelpView())

	case stateViewing:

		fc := m.renderFileContent(m.currentFileContents)

		// Only show the file content, no table headers or rows
		return centered.Render(fileWindowStyle.Render(fmt.Sprintf("Viewing %s\n\n%s", m.viewingFile, fc)) +
			"\n\nPress 'esc' or 'q' to go back.")
	}
	return ""
}
