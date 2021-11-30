package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

type Word struct {
	Word      string
	Frequency int
}

const countOfResults = 10

func Top10(stringForAnalyse string) []string {
	var wordsFreequencyMap map[string]Word
	var wordsArray []string
	var wordsArraySorted []string
	var wordsFreequencyArray []Word
	var wordsFreequencyArrayIndx int
	var wordsFreequencyArrayLen int
	var trimmedWord string
	var value string
	var tempArray []string
	var tempFrequency int
	var wordFromArray Word
	wordsFreequencyMap = make(map[string]Word)
	wordsArray = strings.Fields(stringForAnalyse)
	for _, value = range wordsArray {
		trimmedWord = strings.Trim(value, " ")
		if trimmedWord != "" {
			if word, ok := wordsFreequencyMap[trimmedWord]; ok {
				word.Frequency++
				wordsFreequencyMap[trimmedWord] = word
			} else {
				wordsFreequencyMap[trimmedWord] = Word{trimmedWord, 1}
			}
		}
	}
	wordsFreequencyArray = make([]Word, len(wordsFreequencyMap))

	wordsFreequencyArrayIndx = 0
	for _, value := range wordsFreequencyMap {
		wordsFreequencyArray[wordsFreequencyArrayIndx] = value
		wordsFreequencyArrayIndx++
	}

	sort.SliceStable(wordsFreequencyArray, func(i, j int) bool {
		return wordsFreequencyArray[i].Frequency > wordsFreequencyArray[j].Frequency
	})

	wordsFreequencyArrayLen = len(wordsFreequencyArray)
	for i := 0; i < wordsFreequencyArrayLen; i++ {
		wordFromArray = wordsFreequencyArray[i]
		if i == 0 {
			tempFrequency = wordsFreequencyArray[i].Frequency
			tempArray = make([]string, 0)
			tempArray = append(tempArray, wordFromArray.Word)
			if i == wordsFreequencyArrayLen-1 { // конец массива
				sort.Strings(tempArray)
				wordsArraySorted = append(wordsArraySorted, tempArray...)
			}
			continue
		}
		if wordFromArray.Frequency == tempFrequency && i != wordsFreequencyArrayLen-1 { // еще не конец массива
			tempArray = append(tempArray, wordFromArray.Word)
			continue
		}

		if wordFromArray.Frequency != tempFrequency { // не совпадает частотность
			sort.Strings(tempArray)
			wordsArraySorted = append(wordsArraySorted, tempArray...)
			tempArray = make([]string, 0)
			tempArray = append(tempArray, wordFromArray.Word)
			tempFrequency = wordsFreequencyArray[i].Frequency
			if i > countOfResults {
				break // нет смысла дальше выполнять проход так как достаточно элементов в результрующем массиве для результата
			}
		}

		if i == wordsFreequencyArrayLen-1 { // конец массива
			sort.Strings(tempArray)
			wordsArraySorted = append(wordsArraySorted, tempArray...)
		}
	}

	if len(wordsArraySorted) >= countOfResults {
		wordsArraySorted = wordsArraySorted[0:countOfResults]
	}

	return wordsArraySorted
}
