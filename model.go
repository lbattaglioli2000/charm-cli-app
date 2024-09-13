package main

import (
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"path/filepath"
)

var renderer = lipgloss.DefaultRenderer()

// Represents the application state
type model struct {
	state               appState            // Tracks the current state (booting or listing directories)
	progress            progress.Model      // Loading indicator on app start
	fileContents        map[string]string   // File contents for dummy file system
	directories         map[string][]string // Fake file system
	currentDir          string              // Current directory
	previousDir         string              // The previous directory
	fileTable           table.Model         // The Table model
	viewingFile         string              // The name of the file being viewed
	currentFileContents string              // The contents of the current file being viewed
	terminalWidth       int                 // Terminal width
	terminalHeight      int                 // Terminal height
	renderer            *lipgloss.Renderer
}

// Initialize the app state with some fake directories and files
func initialModel() model {
	directories := map[string][]string{
		"/":          {"documents", "downloads", "notes.txt"},
		"/documents": {"report.txt", "plan.md"},
		"/downloads": {"image.jpg", "archive.zip"},
	}

	fileContents := map[string]string{
		"/notes.txt":             "Here are some notes for the project...",
		"/documents/report.txt":  "This is the project report...",
		"/documents/plan.md":     "# Project Plan\n\n## Tasks\n1. Task one\n2. Task two\n> pee `pee poo` poo\n\npoo *poo* pee **pee**\n- one\n - two\n",
		"/downloads/image.jpg":   "[Binary image data]",
		"/downloads/archive.zip": "[Binary zip data]",
	}

	columns := []table.Column{
		{Title: "Name", Width: 100},
		{Title: "Type", Width: 50},
	}
	rows := generateRows("/", directories)

	directoryStructureTable := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	// Get the table's default styling
	tableStyles := table.DefaultStyles()

	// Customize the defaults
	tableStyles.Header = tableStyles.Header.BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(true)

	tableStyles.Selected = tableStyles.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("201")).
		Bold(false)

	directoryStructureTable.SetStyles(tableStyles)

	return model{
		renderer:     renderer,
		state:        stateBooting, // Start in the booting state
		progress:     progress.New(progress.WithDefaultGradient()),
		fileContents: fileContents,
		directories:  directories,
		currentDir:   "/",
		previousDir:  "/",
		fileTable:    directoryStructureTable,
	}
}

func generateRows(currentDirectory string, directories map[string][]string) []table.Row {
	var rows []table.Row

	if currentDirectory != "/" {
		rows = append(rows, table.Row{".. (Go Back)", "Directory"})
	}

	// Get the files and directories for the current directory
	for _, file := range directories[currentDirectory] {
		fileType := "File"
		if _, isDir := directories[filepath.Join(currentDirectory, file)]; isDir {
			fileType = "Directory"
		}
		rows = append(rows, table.Row{file, fileType})
	}

	return rows
}

func (m *model) renderFileContent(content string) string {
	fileExtension := filepath.Ext(m.viewingFile)

	if fileExtension == ".md" {
		rendered, err := glamour.Render(content, "dark")
		if err == nil {
			return rendered
		}
	}

	return content
}
