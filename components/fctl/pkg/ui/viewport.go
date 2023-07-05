package ui

import (
	"io"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/formancehq/fctl/membershipclient"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
)

var (
	ViewWidth  = 300 //Number of characters
	ViewHeight = 120 // Number of lines
)

// https://github.com/charmbracelet/bubbletea/blob/master/examples/pager/main.go#L92
type modelManager struct {
	vp    viewport.Model
	keys  *ListKeyMap
	ready bool
	// delegateKeys *ui.DelegateKeyMap
}

func (m modelManager) Init() tea.Cmd {
	return nil
}

func NewViewPortManager(content string, out io.Writer, profile *fctl.Profile, stack *membershipclient.Stack, versions *shared.GetVersionsResponse) (*modelManager, error) {
	width := ViewWidth
	vp := viewport.New(width, ViewHeight)

	// This paramaeter is working well
	// It makes the terminal much smoother with a higher framerate
	// But it breaks bubbletea output
	// vp.HighPerformanceRendering = true
	vp.Style = WindowStyle

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
		vp: vp,
	}, nil

}

func (m modelManager) View() string {

	if !m.ready {
		return m.helpView() + "\n" + "Initializing..."
	}

	return m.helpView() + "\n" + m.vp.View()
}

func (m modelManager) helpView() string {
	return HelpStyle.Render("Formance CLI: \n • ↑/↓: Navigate \n • q: Quit")
}

// TODO: Need to be calculated depending of the Header content or BE Fixed
func (m modelManager) GetHelpViewHeight() int {
	return 5
}

// func headerView() string {
// 	return HeaderStyle.Render("fctl")
// }

func (m modelManager) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

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
			w, h := DocStyle.GetFrameSize()
			m.vp.Width = msg.Width - w
			m.vp.Height = msg.Height - h - m.GetHelpViewHeight()

			// m.vp.Style.Width().Height(msg.Height - h - m.GetHelpViewHeight())
			m.vp.YPosition = m.GetHelpViewHeight() + 1
			m.ready = true
		} else {
			m.vp.Style = m.vp.Style.Width(msg.Width).Height(msg.Height - m.GetHelpViewHeight())
		}
	}

	m.vp, cmd = m.vp.Update(msg)

	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}
