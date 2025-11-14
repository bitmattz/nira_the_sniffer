package view

import (
	"log"
	"os"
	"os/exec"
	"runtime"

	"strconv"

	"github.com/bitmattz/nira_the_sniffer/models"
	portHandler "github.com/bitmattz/nira_the_sniffer/services/ports"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ApplicationPresenter models.ApplicationPresenter

var stylePurple = lipgloss.NewStyle().Foreground(lipgloss.Color("12")).Bold(true)
var styleTitle = lipgloss.NewStyle().Foreground(lipgloss.Color("#04B575")).Bold(true)
var styleSubTitle = lipgloss.NewStyle().Foreground(lipgloss.Color("#95ffd8ff"))

const (
	menu int = iota
	scanPorts
	scanSpecificPort
	history
)

func initialModel() ApplicationPresenter {
	ti := textinput.New()
	ti.Placeholder = "..."
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	tableResult := table.New(
		table.WithFocused(true),
		table.WithHeight(7),
	)

	return ApplicationPresenter{
		Choices:     []string{"Scan ports from IP"},
		Cursor:      0,
		Selected:    make(map[int]struct{}),
		Page:        menu,
		InputMode:   false,
		TextInput:   ti,
		TableMode:   false,
		TableResult: tableResult,
	}

}

func (m ApplicationPresenter) Init() tea.Cmd {
	return textinput.Blink
}

func (m ApplicationPresenter) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	if key, ok := msg.(tea.KeyMsg); ok {
		if (m.Page == scanPorts || m.Page == scanSpecificPort) && m.InputMode {
			m.TextInput, cmd = m.TextInput.Update(key)

			// If user pressed Enter while typing, handle the value and exit input mode
			if key.String() == "enter" {
				address := m.TextInput.Value()
				switch m.Page {
				case scanPorts:

					result := portHandler.ScanPorts(address)
					m.TextInput.Blur()
					if result == nil || len(result) == 0 {
						m.InputMode = false
						return m, nil

					} else {
						columns := []table.Column{
							{Title: "Port", Width: 10},
							{Title: "State", Width: 10},
							{Title: "Service", Width: 200},
						}
						rows := make([]table.Row, len(result))
						for i, scan := range result {
							rows[i] = table.Row{strconv.Itoa(scan.Port), scan.State, scan.Service}
						}

						t := table.New(
							table.WithColumns(columns),
							table.WithRows(rows),
							table.WithFocused(true),
							table.WithHeight(7),
						)

						m.TableResult = t
					}

				case scanSpecificPort:
					var port, err = strconv.Atoi(address)
					if err != nil {
						return m, nil
					}
					portHandler.ScanPort("TCP", "localhost", port)
				}
				m.InputMode = false
			}

			return m, cmd
		}
	}
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

		case "enter":
			if m.Page == menu {
				switch m.Cursor {
				case 0:
					m.Page = scanPorts
					m.InputMode = true
					m.TextInput.Focus()
				case 1:
					m.Page = scanSpecificPort
					m.InputMode = true
					m.TextInput.Focus()
				case 2:
					m.Page = history
				}
			}

		case "esc", "backspace":
			if m.Page != menu {
				m.Page = menu
			}
		}

	}
	return m, cmd
}

func (m ApplicationPresenter) View() string {

	switch m.Page {
	case menu:
		return loadMenuView(m)
	case scanPorts:
		return scanPortsByIPView(m)
	}

	return ""
}

func loadMenuView(m ApplicationPresenter) string {
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

func scanPortsByIPView(m ApplicationPresenter) string {
	s := stylePurple.Render("Nira The Sniffer\n")
	s += styleTitle.Render("\nScan Ports by IP")
	s += stylePurple.Render(" > ")
	s += styleSubTitle.Render("Enter a IP address.\n")

	s += "\n" + m.TextInput.View() + "\n"
	var tableStyle = lipgloss.NewStyle().
		BorderForeground(lipgloss.Color("240"))

	s += "\n" + tableStyle.Render(m.TableResult.View()) + "\n"

	s += "\nPress esc or backspace to go back to menu.\n"
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
		log.Fatal(err)
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
