package main

import (
	"flag"
	"fmt"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"
	"github.com/muesli/termenv"
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
	var host = flag.String("host", "0.0.0.0", "Host address for SSH server to listen")
	var port = flag.Int("port", 22, "Port for SSH server to listen")

	// Set up the Wish SSH server with the Bubble Tea TUI
	sshServer, err := wish.NewServer(
		wish.WithAddress(fmt.Sprintf("%s:%d", *host, *port)),
		wish.WithHostKeyPath(".ssh/id_ed25519"), // Ensure you generate an SSH host key or set this path to an existing one
		wish.WithMiddleware(
			bubbletea.MiddlewareWithColorProfile(func(s ssh.Session) (tea.Model, []tea.ProgramOption) {
				return initialModel(), []tea.ProgramOption{tea.WithAltScreen()}
			}, termenv.TrueColor),
			logging.Middleware(),
		),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create SSH server: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Starting SSH server on %s:%d\n", *host, *port)

	// Start the SSH server
	if err := sshServer.ListenAndServe(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to start SSH server: %v\n", err)
		os.Exit(1)
	}
}
