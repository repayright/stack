package modelutils

import "golang.org/x/crypto/ssh/terminal"

func GetTerminalSize() (int, int, error) {
	return terminal.GetSize(0)
}

func GetTerminalLimits() (int, int, error) {
	return terminal.GetSize(0)
}
