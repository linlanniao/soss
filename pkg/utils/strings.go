package utils

import (
	"math/rand/v2"
)

const (
	lowerStr = "abcdefghijklmnopqrstuvwxyz"
	upperStr = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numStr   = "1234567890"
)

func RandStr(length int, lower, upper, number bool) string {
	var s string
	if lower {
		s += lowerStr
	}
	if upper {
		s += upperStr
	}
	if number {
		s += numStr
	}
	b := make([]byte, length)
	for i := range b {
		b[i] = s[rand.IntN(len(s))]
	}
	return string(b)
}

func RandLowerUpperNumStr(length int) string {
	return RandStr(length, true, true, true)
}

func RandLowerStr(length int) string {
	return RandStr(length, true, false, false)
}

func RandUpperStr(length int) string {
	return RandStr(length, false, true, false)
}

func RandNumStr(length int) string {
	return RandStr(length, false, false, true)
}
