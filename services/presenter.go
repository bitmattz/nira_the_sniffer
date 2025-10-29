package services

import (
	"os"
	"os/exec"
	"runtime"

	"github.com/bitmattz/nira_the_sniffer/models"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ApplicationPresenter models.ApplicationPresenter

var stylePurple = lipgloss.NewStyle().Foreground(lipgloss.Color("12")).Bold(true)

const (
	menu int = iota
	scanPorts
	scanSpecificPort
	history
)

func initialModel() ApplicationPresenter {
	return ApplicationPresenter{
		Choices:  []string{"Scan ports from IP", "Scan specific port", "History"},
		Cursor:   0,
		Selected: make(map[int]struct{}),
		Page:     menu,
	}
}

func (m ApplicationPresenter) Init() tea.Cmd {
	return nil
}

func (m ApplicationPresenter) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:

		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.Cursor > 0 {
				m.Cursor--
			}

		case "down", "j":
			if m.Cursor < len(m.Choices)-1 {
				m.Cursor++
			}

		case "enter", " ":
			if m.Page == menu {
				switch m.Cursor {
				case 0:
					m.Page = scanPorts
				case 1:
					m.Page = scanSpecificPort
				case 2:
					m.Page = history
				}
			} else {
				m.Page = menu
			}
		}
	}

	return m, nil
}

func (m ApplicationPresenter) View() string {

	switch m.Page {
	case menu:
		return loadMenu(m)
	case scanPorts:
		return stylePurple.Render("Scan Ports from IP - [Functionality to be implemented]\nPress Enter to go back to menu.")
	case scanSpecificPort:
		return stylePurple.Render("Scan Specific Port - [Functionality to be implemented]\nPress Enter to go back to menu.")
	case history:
		return stylePurple.Render("History - [Functionality to be implemented]\nPress Enter to go back to menu.")
	}

	return ""
}

func loadMenu(m ApplicationPresenter) string {
	s := stylePurple.Render("Nira The Sniffer\n")
	s += "\nWhat do you want to do?\n\n"

	// Iterate over choices
	for i, choice := range m.Choices {
		cursor := " " // no cursor
		if m.Cursor == i {
			cursor = stylePurple.Render(">")
			choice = stylePurple.Render(choice)
		}

		// Render the row
		s += cursor + " " + choice + "\n"
	}

	s += "\nPress q to quit.\n"

	return s
}

func clear() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func StartApplicationPresenter() {
	clear()
	p := tea.NewProgram(initialModel())
	if err := p.Start(); err != nil {
		panic(err)
	}
}

// TODO:
// 1. Scan ports from IP
// Ask user for IP Address and returns open ports with details, and there is a option to export to a file
// 2. Scan specific port
// Ask user for IP Address and returns open ports with details, and there is a option to export to a file
// 3. History
// Show previous scans done with date and time, and option to export to a file
// 4. Export to file
// Export scan results to a file in formats like CSV, JSON, or TXT
// 5. Save history using local database, sqlite
