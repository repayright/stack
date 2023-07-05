package ui

import (
	"io"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/fctl/pkg/ui/modelutils"
)

// https://github.com/charmbracelet/bubbletea/blob/master/examples/pager/main.go#L92
type modelManager struct {
	vp      viewport.Model
	ready   bool
	content string
}

func (m modelManager) Init() tea.Cmd {
	return nil
}
func (m modelManager) GetListKeyMapHandler() *modelutils.KeyMapHandler {
	k := modelutils.NewKeyMapHandler()
	k.AddNewBinding(
		key.NewBinding(
			key.WithKeys("q", "esc", "ctrl+c"),
			key.WithHelp("q  ", "Quit the application"),
		),
	)
	k.AddNewBinding(
		key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "move up"),
		),
	)
	k.AddNewBinding(
		key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "move down"),
		),
	)
	k.AddNewBinding(
		key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("? ", "Toggle help"),
		),
	)
	k.AddNewBinding(
		key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter ", "show selected item"),
		),
	)
	k.AddNewBinding(
		key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter ", "show selected item"),
		),
	)
	k.AddNewBinding(
		key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter ", "show selected item"),
		),
	)
	k.AddNewBinding(
		key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter ", "show selected item"),
		),
	)
	k.AddNewBinding(
		key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter ", "show selected item"),
		),
	)
	k.AddNewBinding(
		key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter ", "show selected item"),
		),
	)

	k.AddNewBinding(
		key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter ", "show selected item"),
		),
	)

	k.AddNewBinding(
		key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter ", "show selected item"),
		),
	)

	return k
}
func NewViewPortManager(content string, out io.Writer, profile fctl.Profile) (*modelManager, error) {
	width := fctl.ViewWidth
	vp := viewport.New(width, fctl.ViewHeight)

	// This paramaeter is working well
	// It makes the terminal much smoother with a higher framerate
	// But it breaks bubbletea output
	// vp.HighPerformanceRendering = true
	vp.Style = fctl.WindowStyle

	renderer, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
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

	return &modelManager{
		vp:      vp,
		content: content,
	}, nil

}

func (m modelManager) View() string {
	return m.vp.View()
}

func (m modelManager) GetHelpViewHeight() int {
	h := NewHeader()
	return h.GetMaxPossibleHeight()
}

func (m modelManager) Update(msg tea.Msg) (Model, tea.Cmd) {

	var (
		cmd  tea.Cmd
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
		if !m.ready {
			w, h := fctl.DocStyle.GetFrameSize()

			m.vp.SetContent(m.content)
			viewport.Sync(m.vp)
			m.vp.Width = msg.Width - w
			m.vp.Height = msg.Height - h - m.GetHelpViewHeight()
			m.ready = true
		} else {
			// width, height, err := terminal.GetSize(0)
			// fmt.Println(width, height, err)
			w, h := fctl.DocStyle.GetFrameSize()

			m.vp.SetContent(m.content)
			viewport.Sync(m.vp)
			m.vp.YPosition = m.GetHelpViewHeight() + 2
			m.vp.Width = msg.Width - w
			m.vp.Height = msg.Height - h - m.GetHelpViewHeight()
		}
	}

	m.vp, cmd = m.vp.Update(msg)

	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}
