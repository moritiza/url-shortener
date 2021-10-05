package helper

import (
	"errors"
	"math"
	"strings"
)

const (
	base         uint64 = 62
	characterSet        = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
)

// ToBase62 get givven uint64 digit and convert it into base62 string with our custom 62 characters
func ToBase62(n uint64) string {
	encoded := ""
	for n > 0 {
		r := n % base
		n /= base
		encoded = string(characterSet[r]) + encoded

	}

	return encoded
}

// ToBase10 get givven string and convert it into bas10
func ToBase10(encoded string) (uint64, error) {
	var val uint64

	for index, char := range encoded {
		pow := len(encoded) - (index + 1)
		pos := strings.IndexRune(characterSet, char)
		if pos == -1 {
			return 0, errors.New("invalid character: " + string(char))
		}

		val += uint64(pos) * uint64(math.Pow(float64(base), float64(pow)))
	}

	return val, nil
}
