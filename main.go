package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type dataset struct {
	Names []string `yaml:"names"`
}

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

type model struct {
	keys       keyMap
	help       help.Model
	inputStyle lipgloss.Style
	status     map[string]string
	names      []string
	view       rune
	quitting   bool
}

func newModel() model {
	return model{
		keys:       keys,
		help:       help.NewModel(),
		inputStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("#FF75B7")),
		view:       's',
		status:     getStatus(),
		names:      loadNames(),
	}
}

func (m model) Init() tea.Cmd {
	return nil
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
		case key.Matches(msg, m.keys.Update):
			m.view = 'u'
			runUpdate()
			m.names = loadNames()
			m.status = getStatus()
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		case key.Matches(msg, m.keys.Quit):
			m.quitting = true
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	if m.quitting {
		return "Bye!\n"
	}

	var status string
	status = "hi!"
	if m.view == 's' {
		status = ""
		for k, v := range m.status {
			status += fmt.Sprintf("%-10s: ", k) + m.inputStyle.Render(v) + "\n"
		}
		status += fmt.Sprintf("%-10s: ", "Count") + m.inputStyle.Render(fmt.Sprintf("%d", len(m.names))) + "\n"
	}
	if m.view == 'u' {
		status = "Runing update..."
	}

	helpView := m.help.View(m.keys)
	height := 8 - strings.Count(status, "\n") - strings.Count(helpView, "\n")

	return "\n" + status + strings.Repeat("\n", height) + helpView
}

func main() {
	// search("Trnucha", loadNames())
	search("", []string{})
	if err := tea.NewProgram(newModel()).Start(); err != nil {
		fmt.Printf("Could not start program :(\n%v\n", err)
		os.Exit(1)
	}
}
