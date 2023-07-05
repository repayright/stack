package modelutils

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type KeyMapHandler struct {
	list []key.Binding
}

func NewKeyMapHandler() *KeyMapHandler {
	return &KeyMapHandler{
		list: []key.Binding{},
	}
}

func (k *KeyMapHandler) GetListKeyMap() []key.Binding {
	return k.list
}

func (k *KeyMapHandler) AddNewBinding(keymap key.Binding) *KeyMapHandler {
	k.list = append(k.list, keymap)
	return k
}

func (k *KeyMapHandler) Reset() *KeyMapHandler {
	k.list = []key.Binding{}
	return k
}

func (k *KeyMapHandler) Init() *tea.Cmd {
	return nil
}

func (k *KeyMapHandler) Update(msg tea.Msg) *tea.Cmd {
	return nil
}

func (k *KeyMapHandler) View() string {
	var s string = ""
	for _, keymap := range k.list {

		h := keymap.Help()
		s += h.Key + ": " + h.Desc + "\n"
	}

	return s
}

func GetFlatMappingKeys(handlers ...KeyMapHandler) []string {
	var keys []string
	for _, handler := range handlers {
		for _, keymap := range handler.GetListKeyMap() {
			keys = append(keys, strings.Join(keymap.Keys(), ", "))
		}
	}

	return keys
}
