package ui

import (
	"io"

	"github.com/formancehq/fctl/pkg/config"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/formancehq/fctl/pkg/ui/theme"
)

// https://github.com/charmbracelet/bubbletea/blob/master/examples/pager/main.go#L92
type ModelManager struct {
	vp      viewport.Model
	ready   bool
	content string
}

func (m ModelManager) Init() tea.Cmd {
	return nil
}
func (m ModelManager) GetListKeyMapHandler() *config.KeyMapHandler {
	k := config.NewKeyMapHandler()

	return k
}
func NewViewPortManager(content string, out io.Writer) (*ModelManager, error) {
	width := theme.ViewWidth
	vp := viewport.New(width, theme.ViewHeight)

	// This parameter is working well
	// It makes the terminal much smoother with a higher framerate
	// But it breaks bubbletea output
	// vp.HighPerformanceRendering = true

	vp.Style = theme.WindowStyle

	renderer, err := glamour.NewTermRenderer(
		//glamour.WithAutoStyle(),
		glamour.WithWordWrap(width),
	)

	if err != nil {
		return nil, err
	}

	str, err := renderer.Render(content)
	if err != nil {
		return nil, err
	}
	vp.SetContent(str)

	return &ModelManager{
		vp:      vp,
		content: content,
	}, nil

}

func (m ModelManager) View() string {
	return m.vp.View()
}

func (m ModelManager) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var (
		//cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		default:
			var cmd tea.Cmd
			m.vp, cmd = m.vp.Update(msg)
			return m, cmd
		}
	case tea.WindowSizeMsg:
		m.vp.Width = msg.Width
		m.vp.Height = msg.Height + 1

		renderer, err := glamour.NewTermRenderer(
			glamour.WithAutoStyle(),
			glamour.WithWordWrap(msg.Width),
		)
		if err != nil {
			return nil, tea.Quit
		}

		str, err := renderer.Render(m.content)
		if err != nil {
			return nil, tea.Quit
		}
		m.vp.SetContent(str)

	}

	return m, tea.Batch(cmds...)
}
