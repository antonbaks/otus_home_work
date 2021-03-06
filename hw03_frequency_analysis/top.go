package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

const topNumber = 10

type wordsWithCount struct {
	Word  string
	Count int
}

var re = regexp.MustCompile(`[":,!.?]|\s-\s`)

func Top10(text string) []string {
	text = stringPreprocessing(text)
	words := strings.Fields(text)
	mapWordsWithCount := make(map[string]int)

	for _, word := range words {
		mapWordsWithCount[word]++
	}

	wordsWithCount := getStruct(mapWordsWithCount)

	sort.Slice(wordsWithCount, func(i, j int) bool {
		if wordsWithCount[i].Count == wordsWithCount[j].Count {
			return wordsWithCount[i].Word < wordsWithCount[j].Word
		}

		return wordsWithCount[i].Count > wordsWithCount[j].Count
	})

	return getTop(wordsWithCount)
}

func stringPreprocessing(text string) string {
	text = strings.ToLower(text)

	text = re.ReplaceAllString(text, " ")

	return text
}

func getStruct(mapWordsWithCount map[string]int) []wordsWithCount {
	structs := make([]wordsWithCount, 0, len(mapWordsWithCount))

	for word, count := range mapWordsWithCount {
		structs = append(structs, wordsWithCount{word, count})
	}

	return structs
}

func getTop(sliceWordsWithCount []wordsWithCount) []string {
	top := make([]string, 0, topNumber)

	for i, v := range sliceWordsWithCount {
		if i+1 > topNumber {
			break
		}

		top = append(top, v.Word)
	}

	return top
}
