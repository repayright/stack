package utils

func Map[SRC any, DST any](srcs []SRC, mapper func(SRC) DST) []DST {
	ret := make([]DST, 0)
	for _, src := range srcs {
		ret = append(ret, mapper(src))
	}
	return ret
}

func Min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

func Max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}

func BoolToString(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

func BoolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func BoolToAbscisse(b bool) int {
	if b {
		return 1
	}
	return -1
}
