package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/seanch0n/halp/cheats"
	"github.com/seanch0n/halp/globals"
	"github.com/seanch0n/halp/varview"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type item struct {
	title, desc string
}
type GotCheat struct {
	Err   error
	Cheat cheats.CheatSearch
}

func (m model) fetchCheats(query string) tea.Cmd {
	return func() tea.Msg {
		cheat, err := m.Cheat.FindSelectedCheat(query)
		if err != nil {
			return GotCheat{Err: err}
		}
		return GotCheat{Cheat: cheat}
	}
}

type model struct {
	list      list.Model
	varView   varview.Model
	textInput textinput.Model
	Cheat     cheats.CheatSearch
	err       error
	loading   bool
	typing    bool
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		if msg.String() == "enter" {
			query := m.list.SelectedItem().FilterValue()
			if len(query) > 2 {
				//m.typing = false
				//m.loading = true
				// we are going to quit and start the next view
				// that gets the users input
				//return m, tea.Quit
				return m, tea.Cmd(
					m.fetchCheats(m.list.SelectedItem().FilterValue()),
				)
			}
		} else {
			m.loading = false
			if len(m.varView.Cheat.TitleField) > 0 {
				return m, tea.Quit
			}
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)

	case GotCheat:
		if err := msg.Err; err != nil {
			m.err = err
			return m, nil
		}
		m.Cheat = msg.Cheat
		globals.TheCheat = msg.Cheat
		return m, tea.Quit
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	// if !m.loading {
	// we are displaying the normal view where you can search
	return docStyle.Render(m.list.View())
	// } else {
	// 	if len(m.cheat.Vars) > 0 {
	// 		// we have found a cheat we want to use, so now go into editor mode
	// 		m.varView = m.varView.InitialModel(m.cheat.Vars, m.cheat.Desc)
	// 		//tea.NewProgram(m.varView).Start()
	// 		return docStyle.Render(m.varView.View())
	// 	}
	// 	return docStyle.Render(m.list.View())
	// 	//return fmt.Sprintf("%s\n", m.cheat)
	// }
}

func main() {
	// populate the list that we can fuzzysearch on
	items := cheats.GetList("./cheatFiles")

	t := textinput.NewModel()
	m := model{list: list.New(items, list.NewDefaultDelegate(), 0, 0)}
	m.textInput = t
	m.list.Title = "Available Cheatsheets"

	p := tea.NewProgram(m, tea.WithAltScreen())

	if err := p.Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
	// now get the vars
	mV := varview.Model{}
	m.varView = mV.InitialModel(globals.TheCheat.Vars, globals.TheCheat.Desc)
	p = tea.NewProgram(m.varView, tea.WithAltScreen())
	if err := p.Start(); err != nil {
		fmt.Println("error with var getter: ", err)
	}
	globals.Pr()
}
