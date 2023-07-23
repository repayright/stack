package config

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/formancehq/stack/libs/go-libs/collectionutils"
)

type KeyMapHandler struct {
	keyMapsAction map[*key.Binding]func(tea.Model) Controller
}

func NewKeyMapHandler() *KeyMapHandler {
	return &KeyMapHandler{
		keyMapsAction: make(map[*key.Binding]func(tea.Model) Controller),
	}
}

func (k *KeyMapHandler) GetKeyMapAction() map[*key.Binding]func(tea.Model) Controller {
	return k.keyMapsAction
}

func (k *KeyMapHandler) GetAction(teaKey tea.Key) func(tea.Model) Controller {
	for key := range k.keyMapsAction {
		if !collectionutils.Contains(key.Keys(), teaKey.String()) {
			return nil
		}
		return k.keyMapsAction[key]
	}
	return nil
}

func (k *KeyMapHandler) GetListKeys() []*key.Binding {
	var keys []*key.Binding
	for k := range k.keyMapsAction {
		keys = append(keys, k)
	}
	return keys
}

func (k *KeyMapHandler) AddNewKeyBinding(key key.Binding, action func(tea.Model) Controller) *KeyMapHandler {
	k.keyMapsAction[&key] = action
	return k
}

func (k *KeyMapHandler) Reset() *KeyMapHandler {
	k.keyMapsAction = make(map[*key.Binding]func(tea.Model) Controller)
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
