package ui

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Version struct {
	content string
}

func NewVersion(version string) *Version {
	return &Version{
		content: version,
	}
}
func (v *Version) GetMaxPossibleHeight() int {
	return 1
}
func (v *Version) Init() tea.Cmd {
	return nil
}

func (v *Version) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return nil, nil
}

func (v *Version) View() string {
	t := "Version : " + v.content

	return t
}
