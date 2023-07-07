package modelutils

import (
	"strings"
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

func FillCharBeforeChar(str string, charToFill string, beforeChar string, max int) string {
	splits := strings.Split(str, beforeChar)
	if len(splits[0]) >= max {
		return str
	}

	toFill := max - len(splits[0])

	return splits[0] + strings.Repeat(charToFill, toFill) + beforeChar + strings.TrimPrefix(splits[1], " ")
}
