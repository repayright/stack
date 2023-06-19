package fctl

import (
	"errors"
	"os/exec"
	"runtime"
)

// Reduce is a generic reduce function
// It takes a slice of T and a function that takes a M and a T and returns a M
func Reduce[T, M any](s []T, f func(M M, T T, err error) (M, error), initValue M) (M, error) {
	var err error
	acc := initValue
	for _, v := range s {
		if err != nil {
			return acc, err
		}
		acc, err = f(acc, v, err)
	}
	return acc, nil
}

// If is a ternary operator
func If[T any](cond bool, vtrue, vfalse T) T {
	if cond {
		return vtrue
	}
	return vfalse
}

func Map[SRC any, DST any](srcs []SRC, mapper func(SRC) DST) []DST {
	ret := make([]DST, 0)
	for _, src := range srcs {
		ret = append(ret, mapper(src))
	}
	return ret
}

func MapKeys[K comparable, V any](m map[K]V) []K {
	ret := make([]K, 0)
	for k := range m {
		ret = append(ret, k)
	}
	return ret
}

func Prepend[V any](array []V, items ...V) []V {
	return append(items, array...)
}

func Open(url string) error {
	var (
		cmd  string
		args []string
	)

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	_, err := exec.LookPath(cmd)
	if err == nil {
		args = append(args, url)
		return exec.Command(cmd, args...).Start()
	}

	SetSharedAdditionnalData("browser_url", url)
	return errors.New("error_opening_browser")
}
