package chat

import (
	"fmt"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"os"
	"strings"
)

var (
	userStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF69B4")). // Hot Pink
		Bold(true)

	activeInputStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#FF1493")). // Deep Pink
		Foreground(lipgloss.Color("#FFFFFF"))        // White text

	disabledInputStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#696969")). // Dim Gray
		Foreground(lipgloss.Color("#A9A9A9"))        // Dark Gray text

	statusStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#98FB98"))
)

type ChatInterface struct {
	messages []string
	viewport viewport.Model
	input    textinput.Model
	err      error
	quitting bool
	waiting  bool
	spinner  spinner.Model
	width    int
	height   int
	count	 int
	renderer glamour.TermRenderer
}

func (c *ChatInterface) sendMessage() tea.Msg {
	userMessage := c.input.Value()
	c.input.SetValue("")
	return userMessageMsg(userMessage)
}

type (
	userMessageMsg string
)

// Init implements tea.Model.
func (c *ChatInterface) Init() tea.Cmd {
	return textinput.Blink
}

// Update implements tea.Model.
func (c *ChatInterface) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	// print("i\n")
	switch msg := msg.(type) {
	// Key press event
	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c", "q", "/exit":
			c.quitting = true
			return c, tea.Quit

		case "enter":	
			// print(c.input.Value())
			if c.input.Value() != "" {
				userInput := c.input.Value()          // Capture input
				c.input.SetValue("")                  // Clear input field
				c.addMessage("You", userInput)        // Add to viewport immediately
				// c.waiting = true

				if c.count == 0 {
					cmds = append(cmds, func() tea.Msg {  // Send message event
						return userMessageMsg(userInput)
					})

					c.count++
				}
			}

		case "up", "down", "pgup", "pgdown":
			// print("hi")
			// print("input value ", c.input.Value())
			c.viewport, _ = c.viewport.Update(msg)
		}

	case tea.WindowSizeMsg:
		print("isUpdated")
		c.width, c.height = msg.Width, msg.Height
		c.viewport.Width = msg.Width
		c.viewport.Height = msg.Height - 3
		c.input.Width = msg.Width - 4
		c.updateViewportContent()
	
	case userMessageMsg:
		c.addMessage("You", string(msg))

	}

	var cmd tea.Cmd
	// print("reached here")
	c.input, cmd = c.input.Update(msg)
	// print(cmd)
	cmds = append(cmds, cmd)

	// fmt.Println("1 ", c.input.Value())
	//
	c.viewport, cmd = c.viewport.Update(msg)
	cmds = append(cmds, cmd)

	// fmt.Println("2 ", c.input.Value())

	return c, tea.Batch(cmds...)
}

func (c *ChatInterface) addMessage(sender, content string) {
	formattedMsg := userStyle.Render(sender+":") + " " + content
	c.messages = append(c.messages, formattedMsg)
	c.updateViewportContent()
}

func (c *ChatInterface) updateViewportContent() {
	content := strings.Join(c.messages, "\n\n")
	c.viewport.SetContent(content)
	c.viewport.GotoBottom()
}

func (c *ChatInterface) View() string {
	var status string
	var inputView string

	if c.waiting {
		status = fmt.Sprintf("%s AI is thinking...", c.spinner.View())
		inputView = disabledInputStyle.Render(c.input.View())
	} else {
		status = "Ready for your message"
		inputView = activeInputStyle.Render(c.input.View())
	}

	return lipgloss.JoinVertical(lipgloss.Left,
		c.viewport.View(),
		inputView,
		statusStyle.Render(status),
	)
}

func InitiateChatInterface() *ChatInterface {

	input := textinput.New()
	input.Placeholder = "Type to enter a message..."
	input.Focus()

	vp := viewport.New(80, 20)
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
		spinner:  s,
		renderer: *renderer,
		input: input,
		count: 0,
	}
}

func (c *ChatInterface) Run() (bool, error) {
	p := tea.NewProgram(c, tea.WithAltScreen())
	m, err := p.Run()
	if err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
	cI, _ := m.(*ChatInterface)
	return cI.quitting, cI.err
}
