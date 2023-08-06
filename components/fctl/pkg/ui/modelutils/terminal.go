package modelutils

import "golang.org/x/term"

func GetTerminalSize() (int, int, error) {
	return term.GetSize(0)
}

func GetTerminalLimits() (int, int, error) {
	return term.GetSize(0)
}
