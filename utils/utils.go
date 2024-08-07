package utils

import (
	"math"
	"unicode"
)

func NumToHexString(num uint) string {
	if num == 0 {
		return "0"
	}
	hex := ""
	for num > 0 {
		remainder := uint8(num % 16)
		c := '0' + remainder
		if remainder >= 10 {
			remainder -= 10
			c = 'A' + remainder
		}
		num /= 16
		hex = string(c) + hex
	}

	return hex
}

func HexStringToInt(hexStr string) int {
	num := 0
	idx := len(hexStr) - 1
	for idx >= 0 {
		c := hexStr[idx]
		exponent := float64(len(hexStr) - idx - 1)
		if unicode.IsDigit(rune(c)) {
			num += int(c-'0') * int(math.Pow(16, exponent))
		} else {
			num += int((c-'A')+10) * int(math.Pow(16, exponent))
		}
		idx--
	}
	return num
}
