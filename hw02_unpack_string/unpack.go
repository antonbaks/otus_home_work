package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

type symbolsAfter struct {
	oneAfter rune
	twoAfter rune
}

func Unpack(forUnpackText string) (string, error) {
	textRunes := []rune(forUnpackText)

	if hasErrors(textRunes) {
		return forUnpackText, ErrInvalidString
	}

	return unpackProcessing(textRunes), nil
}

func unpackProcessing(textRunes []rune) string {
	stringBuilder := strings.Builder{}
	i := 0

	for i < len(textRunes) {
		symbolsAfter := getSymbolsAfter(textRunes, i)

		if isSlash(textRunes[i]) && isSlash(symbolsAfter.oneAfter) && unicode.IsDigit(symbolsAfter.twoAfter) {
			digit, _ := strconv.Atoi(string(symbolsAfter.twoAfter))

			stringBuilder.WriteString(strings.Repeat(string(symbolsAfter.oneAfter), digit))

			i += 3

			continue
		}

		if isSlash(textRunes[i]) && unicode.IsDigit(symbolsAfter.oneAfter) && unicode.IsDigit(symbolsAfter.twoAfter) {
			digit, _ := strconv.Atoi(string(symbolsAfter.twoAfter))

			stringBuilder.WriteString(strings.Repeat(string(symbolsAfter.oneAfter), digit))

			i += 3

			continue
		}

		if isSlash(textRunes[i]) && isSlash(symbolsAfter.oneAfter) {
			stringBuilder.WriteRune(textRunes[i])

			i += 2

			continue
		}

		if isSlash(textRunes[i]) && unicode.IsDigit(symbolsAfter.oneAfter) {
			stringBuilder.WriteRune(symbolsAfter.oneAfter)

			i += 2

			continue
		}

		if unicode.IsDigit(symbolsAfter.oneAfter) {
			digit, _ := strconv.Atoi(string(symbolsAfter.oneAfter))

			stringBuilder.WriteString(strings.Repeat(string(textRunes[i]), digit))

			i += 2

			continue
		}

		stringBuilder.WriteRune(textRunes[i])

		i++
	}

	return stringBuilder.String()
}

func hasErrors(textRunes []rune) bool {
	i := 0
	hasError := false

	for i < len(textRunes) {
		if i == 0 && unicode.IsDigit(textRunes[i]) {
			hasError = true

			break
		}

		symbolsAfter := getSymbolsAfter(textRunes, i)

		if isSlash(textRunes[i]) && unicode.IsDigit(symbolsAfter.oneAfter) {
			i += 2

			continue
		}

		if isSlash(textRunes[i]) && !isSlash(symbolsAfter.oneAfter) && !unicode.IsDigit(symbolsAfter.oneAfter) {
			hasError = true

			break
		}

		if unicode.IsDigit(textRunes[i]) && unicode.IsDigit(symbolsAfter.oneAfter) {
			hasError = true

			break
		}

		i++
	}

	return hasError
}

func getSymbolsAfter(textRunes []rune, numIteration int) symbolsAfter {
	symbolsAfter := symbolsAfter{}

	if (numIteration + 1) < len(textRunes) {
		symbolsAfter.oneAfter = textRunes[numIteration+1]
	}

	if (numIteration + 2) < len(textRunes) {
		symbolsAfter.twoAfter = textRunes[numIteration+2]
	}

	return symbolsAfter
}

func isSlash(symbolCode rune) bool {
	var slash rune = 92

	return slash == symbolCode
}
