package view

import (
	"log"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

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

type scanFinishedMsg []models.PortScan

func scanPortsCmd(address string) tea.Cmd {
	return func() tea.Msg {
		results := portHandler.ScanPorts(address)
		return scanFinishedMsg(results)
	}
}

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
		IsLoading:   false,
	}
}

func (m ApplicationPresenter) Init() tea.Cmd {
	return textinput.Blink
}

func (m ApplicationPresenter) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case scanFinishedMsg:
		m.IsLoading = false
		res := []models.PortScan(msg)
		maxBannerWidth := len("Banner")

		if len(res) == 0 {
			return m, nil
		}

		rows := make([]table.Row, len(res))
		for i, scan := range res {
			port := "Unknown"
			if scan.Port != 0 {
				port = strconv.Itoa(scan.Port)
			}

			state := "Unknown"
			if strings.TrimSpace(scan.State) != "" && scan.State != "unknown" {
				state = scan.State
			}

			banner := "Unknown"
			if strings.TrimSpace(scan.Banner) != "" && scan.Banner != "unknown" {
				banner = scan.Banner
			}

			pid := "Unknown"
			if strings.TrimSpace(scan.PID) != "" && scan.PID != "unknown" {
				pid = scan.PID
			}

			if len(banner) > maxBannerWidth {
				maxBannerWidth = len(banner)

			}
			if maxBannerWidth > 30 {
				maxBannerWidth = 30
			}
			if m.TextInput.Value() == "localhost" || m.TextInput.Value() == "127.0.0.1" {
				rows[i] = table.Row{port, state, banner, pid}
			} else {
				rows[i] = table.Row{port, state, banner}

			}
		}
		if m.TextInput.Value() == "localhost" || m.TextInput.Value() == "127.0.0.1" {
			columns := []table.Column{
				{Title: "Port", Width: 8},
				{Title: "State", Width: 8},
				{Title: "Banner", Width: maxBannerWidth},
				{Title: "PID", Width: 20},
			}

			m.TableResult = table.New(
				table.WithColumns(columns),
				table.WithRows(rows),
				table.WithFocused(true),
				table.WithHeight(7),
			)
		} else {
			columns := []table.Column{
				{Title: "Port", Width: 8},
				{Title: "State", Width: 8},
				{Title: "Banner", Width: maxBannerWidth},
			}

			m.TableResult = table.New(
				table.WithColumns(columns),
				table.WithRows(rows),
				table.WithFocused(true),
				table.WithHeight(7),
			)
		}

		return m, nil

	case tea.KeyMsg:

		if (m.Page == scanPorts || m.Page == scanSpecificPort) && m.InputMode {
			m.TextInput, cmd = m.TextInput.Update(msg)

			if msg.String() == "enter" {
				address := m.TextInput.Value()

				switch m.Page {
				case scanPorts:
					m.InputMode = true
					m.IsLoading = true
					return m, scanPortsCmd(address)
				}
			}
			return m, cmd
		}

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
				}
			}

		case "esc":
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

	for i, choice := range m.Choices {
		cursor := " "
		if m.Cursor == i {
			cursor = stylePurple.Render(">")
			choice = stylePurple.Render(choice)
		}
		s += cursor + " " + choice + "\n"
	}

	s += "\nPress q to quit.\n"
	return s
}

func scanPortsByIPView(m ApplicationPresenter) string {
	s := stylePurple.Render("Nira The Sniffer\n")
	s += styleTitle.Render("\nScan Ports by IP")
	s += stylePurple.Render(" > ")
	s += styleSubTitle.Render("Enter an IP address.\n")

	s += "\n" + m.TextInput.View() + "\n"

	if m.IsLoading {
		s += "\nScanning...\n"
	}

	var tableStyle = lipgloss.NewStyle().
		BorderForeground(lipgloss.Color("240"))

	s += "\n" + tableStyle.Render(m.TableResult.View()) + "\n"

	s += "\nPress esc to go back to menu.\n"
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
