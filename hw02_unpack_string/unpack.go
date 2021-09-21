package hw02unpackstring

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

const constBackSlash rune = '\\'

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	var resultSB strings.Builder
	var err error
	runes := []rune(str)
	for i := 0; i < len(runes); {
		currentRune := runes[i]
		if i == 0 {
			if unicode.IsDigit(currentRune) { // if first letter is digit, it is error
				return "", ErrInvalidString
			}
			resultSB.WriteRune(currentRune)
			i++
			continue
		}

		i, err = makeResultString(i, runes, currentRune, &resultSB)
		if err != nil {
			return "", err
		}
	}
	fmt.Println(resultSB.String())
	return resultSB.String(), nil
}

func makeResultString(i int, runes []rune, currentRune rune, resultSB *strings.Builder) (int, error) {
	var err error
	if unicode.IsDigit(currentRune) {
		digit, err := strconv.Atoi(string(currentRune))
		if err != nil {
			return 0, ErrInvalidString
		}
		prevRune := runes[i-1]
		if returnErrorIfDigit(prevRune) != nil {
			fmt.Println(prevRune)
			return 0, ErrInvalidString
		}

		if digit == 0 {
			stringWithoutSufix := strings.TrimSuffix(resultSB.String(), string(prevRune))
			resultSB.Reset()
			resultSB.WriteString(stringWithoutSufix)
			i++
		} else {
			resultSB.WriteString(strings.Repeat(string(prevRune), digit-1))
			i++
		}
	} else {
		if checkOnBackSlash(currentRune) {
			if i < len(runes)-1 && unicode.IsDigit(runes[i+1]) { // case \5 => 5
				if i < len(runes)-2 && unicode.IsDigit(runes[i+2]) { // case \45 => 44444
					digit2, err := strconv.Atoi(string(runes[i+2]))
					if err != nil {
						return 0, ErrInvalidString
					}
					resultSB.WriteString(strings.Repeat(string(runes[i+1]), digit2))
					i += 3
				} else {
					resultSB.WriteRune(runes[i+1])
					i += 2
				}
			} else {
				resultSB.WriteRune(currentRune)
				i++
			}
		} else {
			resultSB.WriteRune(currentRune)
			i++
		}
	}
	return i, err
}

func returnErrorIfDigit(char rune) error {
	if unicode.IsDigit(char) {
		return ErrInvalidString
	}
	return nil
}

func checkOnBackSlash(currentRune rune) bool {
	return currentRune == constBackSlash
}
