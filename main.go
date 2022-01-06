package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type dataset struct {
	Names []string `yaml:"names"`
}

type updateMsg struct{}

const file = "data/names.yml"

// keyMap defines a set of keybindings. To work for help it must satisfy
// key.Map. It could also very easily be a map[string]key.Binding.
type keyMap struct {
	Status key.Binding
	Update key.Binding
	Help   key.Binding
	Quit   key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Status, k.Update}, // first column
		{k.Help, k.Quit},     // second column
	}
}

var keys = keyMap{
	Status: key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "show status"),
	),
	Update: key.NewBinding(
		key.WithKeys("u"),
		key.WithHelp("u", "update list"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

var (
	highlightStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF75B7")).Render
)

type model struct {
	keys     keyMap
	help     help.Model
	spinner  spinner.Model
	status   map[string]string
	names    []string
	view     rune
	updating bool
	quitting bool
}

func newModel() model {
	s := spinner.NewModel()
	s.Spinner = spinner.MiniDot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFB300"))
	return model{
		keys:    keys,
		help:    help.NewModel(),
		spinner: s,
		view:    's',
		status:  getStatus(),
		names:   loadNames(),
	}
}

func (m model) Init() tea.Cmd {
	return spinner.Tick
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// If we set a width on the help menu it can it can gracefully truncate
		// its view as needed.
		m.help.Width = msg.Width

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Status):
			m.view = 's'
			m.status = getStatus()
		case key.Matches(msg, m.keys.Update):
			m.view = 'u'
			m.updating = true
			return m, updateNames
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		case key.Matches(msg, m.keys.Quit):
			m.quitting = true
			return m, tea.Quit
		}

	case updateMsg:
		m.names = loadNames()
		m.updating = false
		return m, nil

	case spinner.TickMsg:
		if m.view == 'u' && !m.updating {
			return m, nil
		}
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m model) View() string {
	if m.quitting {
		return "Bye!\n"
	}

	var status string
	if m.view == 's' {
		status = fmt.Sprintf("%-10s: ", "Filename") + highlightStyle(m.status["path"])
		status += fmt.Sprintf("\n%-10s: ", "Updated") + highlightStyle(m.status["date"])
		status += fmt.Sprintf("\n%-10s: ", "Count") + highlightStyle(fmt.Sprintf("%d", len(m.names)))
	}
	if m.view == 'u' {
		if m.updating {
			status = "Runing update..." + m.spinner.View()
		} else {
			status = highlightStyle("Update done...")
		}
	}

	helpView := m.help.View(m.keys)
	height := 8 - strings.Count(status, "\n") - strings.Count(helpView, "\n")

	return "\n" + status + strings.Repeat("\n", height) + helpView
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	// search("Trnucha", loadNames())
	search("", []string{})

	if err := tea.NewProgram(newModel()).Start(); err != nil {
		fmt.Printf("Could not start program :(\n%v\n", err)
		os.Exit(1)
	}

}

func updateNames() tea.Msg {
	runUpdate()
	return updateMsg(struct{}{})
}
