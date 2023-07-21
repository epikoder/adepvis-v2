package str

import (
	"strings"
	"unicode"
)

type strType int

var (
	Alphabet strType = 1
	Number   strType = 2
)

func SnakeCase(s string) string {
	r := ""
	first := false
	for _, c := range s {
		if unicode.IsUpper(c) {
			if !first {
				first = !first
				r += string(unicode.ToLower(c))
				continue
			}
			r += "_" + string(unicode.ToLower(c))
			continue
		}
		r += string(unicode.ToLower(c))
	}
	return r
}

func FirstToLower(s string) string {
	r := ""
	first := false
	for _, c := range s {
		if !first {
			first = !first
			r += string(unicode.ToLower(c))
			continue
		}
		r += string(c)
	}
	return r
}

func Append(s, a string) string {
	if strings.HasSuffix(s, a) {
		return s
	}
	return s + a
}

func TrimRight(s, r string) string {
	if strings.HasSuffix(s, r) {
		return strings.Replace(s, r, "", (len(s) - len(r) + 1))
	}
	return s
}

func TrimLeft(s, l string) string {
	if strings.HasPrefix(s, l) {
		return strings.Replace(s, l, "", 1)
	}
	return s
}

func Only(s string, t strType) (i string) {
	for _, r := range s {
		switch t {
		case Alphabet:
			if unicode.IsLetter(r) {
				i += string(r)
			}
		case Number:
			if unicode.IsNumber(r) {
				i += string(r)
			}
		}

	}
	return i
}
