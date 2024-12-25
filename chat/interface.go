package chat

import (
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

type ChatInterface struct {
	messages []string
	viewport viewport.Model
	input textinput.Model
	err error
	quitting bool
	waiting bool 
	spinner spinner.Model
	width int
	height int
	renderer glamour.TermRenderer
}

func InitiateChatInterface() *ChatInterface {

	input := textinput.New();
	input.Placeholder = "Type to enter a message..."

	vp := viewport.New(80, 20);
	s := spinner.New()

	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFA500"))

	renderer, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(80),
	)

	return &ChatInterface{
		messages: []string{},
		viewport: vp,
		spinner: s,
		renderer: *renderer,
	}
}