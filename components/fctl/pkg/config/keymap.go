package config

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type KeyMapHandler struct {
	keyMapsAction map[*key.Binding]func() Controller
}

func NewKeyMapHandler() *KeyMapHandler {
	return &KeyMapHandler{
		keyMapsAction: make(map[*key.Binding]func() Controller),
	}
}

func (k *KeyMapHandler) GetKeyMapAction() map[*key.Binding]func() Controller {
	return k.keyMapsAction
}

func (k *KeyMapHandler) AddNewKeyBinding(key key.Binding, action func() Controller) *KeyMapHandler {
	k.keyMapsAction[&key] = action
	return k
}

func (k *KeyMapHandler) Reset() *KeyMapHandler {
	k.keyMapsAction = make(map[*key.Binding]func() Controller)
	return k
}

func (k *KeyMapHandler) Init() *tea.Cmd {
	return nil
}

func (k *KeyMapHandler) Update(msg tea.Msg) *tea.Cmd {
	return nil
}

func (k *KeyMapHandler) View() string {
	var s = ""
	for k := range k.keyMapsAction {
		h := k.Help()
		s += h.Key + ": " + h.Desc + "\n"
	}

	return s
}
