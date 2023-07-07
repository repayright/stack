package modelutils

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"golang.org/x/mod/semver"
)

func GetMaxCharPosXinCharList(charList []string, char string) int {
	max := 0
	for _, str := range charList {
		split := strings.Split(str, char)
		if len(split[0]) >= max {
			max = len(split[0])
		}
	}
	return max
}

func FillCharBeforeChar(str string, charToFill string, char string, max int) string {
	splits := strings.Split(str, char)
	// fmt.Println(splits[0])
	if len(splits[0]) >= max {
		return str
	}

	toFill := max - len(splits[0])
	// fmt.Println(toFill)

	return splits[0] + strings.Repeat(charToFill, toFill) + char + strings.TrimPrefix(splits[1], " ")
}

func GetLatestVersion() string {
	return "0.0.1"
}

func GetActualversion() string {
	return "develop"
}

// Use LipGloss
func GetGradientSemverColor(version string) string {
	if !semver.IsValid(version) {
		return lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Render(version)
	}

	color := ""

	// Major RED
	if semver.Compare(version, GetLatestVersion()) == 0 {
		color = "196"
	}

	// Minor orange
	if color == "" && semver.Compare(version, GetLatestVersion()) == 0 {
		color = "202"
	}

	// Tiny Yellow
	if color == "" && semver.Compare(version, GetLatestVersion()) == 0 {

		color = "226"
	}

	// No Change Green
	if color == "" {
		color = "46"
	}

	return lipgloss.NewStyle().Foreground(lipgloss.Color(color)).Render(version)
}
