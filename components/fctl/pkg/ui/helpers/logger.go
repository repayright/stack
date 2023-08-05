package helpers

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	filePath = "./debug.log"
)

type Logger struct {
	filePath string
	prefix   string

	file *os.File
}

func NewLogger(prefix string) *Logger {
	f, err := tea.LogToFile("./debug.log", "table")
	if err != nil {
		os.Exit(1)
	}
	return &Logger{
		filePath: filePath,
		prefix:   prefix,
		file:     f,
	}
}

func (l *Logger) Log(msg ...string) {
	l.file.WriteString(l.prefix + ": " + fmt.Sprintf("%s\n", strings.Join(msg, " ")))
}
