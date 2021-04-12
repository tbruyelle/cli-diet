package xstrings

import (
	"strings"
	"unicode"
)

// AllOrSomeFilter filters elems out from the list as they  present in filterList and
// returns the remaning ones.
// if filterList is empty, all elems from list returned.
func AllOrSomeFilter(list, filterList []string) []string {
	if len(filterList) == 0 {
		return list
	}

	var elems []string

	for _, elem := range list {
		if !SliceContains(filterList, elem) {
			elems = append(elems, elem)
		}
	}

	return elems
}

// SliceContains returns with true if s is a member of ss.
func SliceContains(ss []string, s string) bool {
	for _, e := range ss {
		if e == s {
			return true
		}
	}

	return false
}

// List returns a slice of strings captured after the value returned by do which is
// called n times.
func List(n int, do func(i int) string) []string {
	var list []string

	for i := 0; i < n; i++ {
		list = append(list, do(i))
	}

	return list
}

// FormatUsername formats a username to make it usable as a variable
func FormatUsername(s string) string {
	return NoDash(NoNumberPrefix(s))
}

// NoDash removes dash from the string
func NoDash(s string) string {
	return strings.ReplaceAll(s, "-", "")
}

// NoNumberPrefix adds a underscore at the beginning of the string if it stars with a number
// this is used for package of proto files template because the package name can't start with a string
func NoNumberPrefix(s string) string {
	// Check if it starts with a digit
	if unicode.IsDigit(rune(s[0])) {
		return "_" + s
	}
	return s
}
