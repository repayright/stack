package modelutils

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type KeyMapHandler[T any] struct {
	keyMapsAction map[*key.Binding]func() *Controller[T]
}

func NewKeyMapHandler[T any]() *KeyMapHandler[T] {
	return &KeyMapHandler[T]{
		keyMapsAction: make(map[*key.Binding]func() *Controller[T]),
	}
}

func (k *KeyMapHandler[T]) GetKeyMapAction() map[*key.Binding]func() *Controller[T] {
	return k.keyMapsAction
}

func (k *KeyMapHandler[T]) AddNewKeyBinding(key key.Binding, action func() *Controller[T]) *KeyMapHandler[T] {
	// k.keyMapsAction = append(k.keyMapsAction, keymap)
	k.keyMapsAction[&key] = action
	return k
}

func (k *KeyMapHandler[T]) Reset() *KeyMapHandler[T] {
	k.keyMapsAction = make(map[*key.Binding]func() *Controller[T])
	return k
}

func (k *KeyMapHandler[T]) Init() *tea.Cmd {
	return nil
}

func (k *KeyMapHandler[T]) Update(msg tea.Msg) *tea.Cmd {
	return nil
}

func (k *KeyMapHandler[T]) View() string {
	var s string = ""
	for k, _ := range k.keyMapsAction {

		h := k.Help()
		s += h.Key + ": " + h.Desc + "\n"
	}

	return s
}

func GetFlatMappingKeys[T any](handlers ...KeyMapHandler[T]) []string {
	var keys []string
	for _, handler := range handlers {
		for k, _ := range handler.GetKeyMapAction() {
			keys = append(keys, strings.Join(k.Keys(), ", "))
		}
	}

	return keys
}
