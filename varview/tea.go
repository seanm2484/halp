package varview

// A simple example demonstrating the use of multiple text input components
// from the Bubbles component library.

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/seanch0n/halp/cheats"
	"github.com/seanch0n/halp/globals"
)

var (
	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle         = focusedStyle.Copy()
	noStyle             = lipgloss.NewStyle()
	helpStyle           = blurredStyle.Copy()
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))

	focusedButton = focusedStyle.Copy().Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
)

type Model struct {
	focusIndex int
	inputs     []textinput.Model
	cursorMode textinput.CursorMode
	Cheat      cheats.CheatSearch
}

func (m Model) InitialModel(vars []string, cmd string) Model {
	m.inputs = make([]textinput.Model, len(vars))
	m.Cheat.Desc = cmd
	m.Cheat.Vars = vars

	var t textinput.Model
	for i, v := range vars {
		t = textinput.New()
		t.CursorStyle = cursorStyle
		t.CharLimit = 100

		t.Placeholder = v
		if i == 0 {
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		}
		m.inputs[i] = t
	}

	return m
}
func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

type Done struct {
	Err error
	Cmd string
}

// the return func() isn't working and I cannot figure out why. So...
// just set the globals in here and return a dummy return because who cares
func (m Model) buildCmd(varFields []textinput.Model) tea.Cmd {
	cmd := m.Cheat.Desc
	// walk over variables. Replace each one with user input
	varValues := make([]string, 0)
	for _, v := range varFields {
		varValues = append(varValues, v.Value())
	}
	for idx, v := range m.Cheat.Vars {

		cmd = strings.Replace(cmd, string("_"+v+"_"), varValues[idx], 1)
	}
	globals.Set(cmd)

	return func() tea.Msg {
		return Done{Cmd: "hi"}
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit

		// Change cursor mode
		case "ctrl+r":
			m.cursorMode++
			if m.cursorMode > textinput.CursorHide {
				m.cursorMode = textinput.CursorBlink
			}
			cmds := make([]tea.Cmd, len(m.inputs))
			for i := range m.inputs {
				cmds[i] = m.inputs[i].SetCursorMode(m.cursorMode)
			}
			return m, tea.Batch(cmds...)

		// Set focus to next input
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			// Did the user press enter while the submit button was focused?
			// If so, exit.
			if s == "enter" && m.focusIndex == len(m.inputs) {
				// they hit submit, so lets get all the variables, and build the
				// command
				m.buildCmd(m.inputs)
				return m, tea.Quit
			}

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > len(m.inputs) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs)
			}

			cmds := make([]tea.Cmd, len(m.inputs))
			for i := 0; i <= len(m.inputs)-1; i++ {
				if i == m.focusIndex {
					// Set focused state
					cmds[i] = m.inputs[i].Focus()
					m.inputs[i].PromptStyle = focusedStyle
					m.inputs[i].TextStyle = focusedStyle
					continue
				}
				// Remove focused state
				m.inputs[i].Blur()
				m.inputs[i].PromptStyle = noStyle
				m.inputs[i].TextStyle = noStyle
			}
			return m, tea.Batch(cmds...)
		}
	case Done:
		globals.Set(msg.Cmd)
		return m, tea.Quit
	default:
		if len(m.inputs) == 0 {
			m.buildCmd(m.inputs)
			return m, tea.Quit
		}

	}

	// Handle character input and blinking
	cmd := m.updateInputs(msg)

	return m, cmd
}

func (m *Model) updateInputs(msg tea.Msg) tea.Cmd {
	var cmds = make([]tea.Cmd, len(m.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m Model) View() string {
	var b strings.Builder

	b.WriteString(helpStyle.Render("Set the following variables:"))
	fmt.Fprintf(&b, "\n")

	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	button := &blurredButton
	if m.focusIndex == len(m.inputs) {
		button = &focusedButton
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", *button)

	b.WriteString(helpStyle.Render("cursor mode is "))
	b.WriteString(cursorModeHelpStyle.Render(m.cursorMode.String()))
	b.WriteString(helpStyle.Render(" (ctrl+r to change style)"))

	return b.String()
}

// func main() {
// 	if err := tea.NewProgram(initialModel()).Start(); err != nil {
// 		fmt.Printf("could not start program: %s\n", err)
// 		os.Exit(1)
// 	}
// }
