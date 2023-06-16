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
	ViewWidth  = 40
	ViewHeight = 20
)

type modelManager struct {
	vp   viewport.Model
	keys *ListKeyMap
	// delegateKeys *ui.DelegateKeyMap
}

func (m modelManager) Init() tea.Cmd {
	return nil
}

func NewModelManager(content string, out io.Writer, profile *fctl.Profile, stack *membershipclient.Stack, versions *shared.GetVersionsResponse) (*modelManager, error) {
	width := 78
	vp := viewport.New(width, 20)

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
	return m.helpView() + "\n" + m.vp.View()
}

func (m modelManager) helpView() string {
	return HelpStyle.Render("Formance CLI: \n • ↑/↓: Navigate \n • q: Quit")
}

func (m modelManager) GetViewKeyHelper() *viewport.Model {
	return &m.vp
}

// func headerView() string {
// 	return HeaderStyle.Render("fctl")
// }

func (m modelManager) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
		w, h := DocStyle.GetFrameSize()
		m.vp.Style.Width(msg.Width - w)
		m.vp.Style.Height(msg.Width - h)
		return m, nil
	default:
		return m, nil
	}

	// var cmd tea.Cmd
	// m.GetCurrentModel().list, cmd = m.GetCurrentModel().list.Update(msg)
	// return m, cmd
}
