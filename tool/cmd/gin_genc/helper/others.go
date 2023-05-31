package helper

import (
	"errors"
	"fmt"
)

func isUpper(c byte) bool {
	return c >= 'A' && c <= 'Z'
}

func isLower(c byte) bool {
	return !isUpper(c)
}

func toUpper(c byte) byte {
	return c - 32
}

func toLower(c byte) byte {
	return c + 32
}

// Underscore converts "CamelCasedString" to "camel_cased_string".
func Underscore(s string) string {
	r := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if isUpper(c) {
			if i > 0 && i+1 < len(s) && (isLower(s[i-1]) || isLower(s[i+1])) {
				r = append(r, '_', toLower(c))
			} else {
				r = append(r, toLower(c))
			}
		} else {
			r = append(r, c)
		}
	}
	return string(r)
}

func ToUpper(s string) string {
	if isUpperString(s) {
		return s
	}

	b := make([]byte, len(s))
	for i := range b {
		c := s[i]
		if c >= 'a' && c <= 'z' {
			c -= 'a' - 'A'
		}
		b[i] = c
	}
	return string(b)
}

func isUpperString(s string) bool {
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'a' && c <= 'z' {
			return false
		}
	}
	return true
}

func ToExported(s string) string {
	if len(s) == 0 {
		return s
	}
	if c := s[0]; isLower(c) {
		b := []byte(s)
		b[0] = toUpper(c)
		return string(b)
	}
	return s
}

func IsInSlice(slice []string, s string) bool {
	for _, thisS := range slice {
		if thisS == s {
			return true
		}
	}
	return false
}

func SliceUnderscore(strs []string) []string {
	rslt := []string{}
	for _, s := range strs {
		rslt = append(rslt, Underscore(s))
	}
	return rslt
}

func SliceExclude(strs []string, exclude []string) []string {
	rslt := []string{}
	for _, s := range strs {
		if !IsInSlice(exclude, s) {
			rslt = append(rslt, s)
		}
	}
	return rslt
}

func IsSliceContain(strs_a []string, strs_b []string) bool {
	for _, s := range strs_b {
		if !IsInSlice(strs_a, s) {
			return false
		}
	}
	return true
}

func IsSliceContainE(strs_a []string, strs_b []string) error {
	for _, s := range strs_b {
		if !IsInSlice(strs_a, s) {
			return errors.New(fmt.Sprintf("Not Contain %s", s))
		}
	}
	return nil
}

func SliceIntersection(strs_a []string, strs_b []string) []string {
	rslt := []string{}
	for _, s := range strs_a {
		if IsInSlice(strs_b, s) {
			rslt = append(rslt, s)
		}
	}
	return rslt
}

func SliceRemove(strs_a []string, strs_b []string) []string {
	rslt := []string{}
	for _, s := range strs_a {
		if !IsInSlice(strs_b, s) {
			rslt = append(rslt, s)
		}
	}
	return rslt
}
