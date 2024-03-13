package utils

import (
	"bufio"
	"errors"
	"os"
	"regexp"
	"strings"
)

var (
	badWordsList []string
	badWordRegex *regexp.Regexp
)

func CompileBadWordsPattern() error {
	var pattern strings.Builder
	pattern.WriteString(`(`)
	for i, word := range badWordsList {
		if word == "" {
			continue
		}
		pattern.WriteString(regexp.QuoteMeta(word))
		if i < len(badWordsList)-1 {
			pattern.WriteString(`|`)
		}
	}
	pattern.WriteString(`)`)

	var err error
	badWordRegex, err = regexp.Compile(pattern.String())
	return err
}

func CheckForBadWords(input string) (bool, error) {
	if badWordRegex == nil {
		return false, errors.New("bad words pattern not compiled")
	}

	return badWordRegex.MatchString(input), nil
}

// func CheckForBadWords(input string) (bool, error) {
// 	// TODO: Normalize input for comparison

// 	// TODO: consider parallelizing
// 	for _, word := range badWordsList {
// 		if word == "" {
// 			continue
// 		}

// 		// Check if the bad word is a substring of the input
// 		if strings.Contains(input, word) {
// 			return true, nil
// 		}
// 	}
// 	return false, nil
// }

// LoadBadWords loads bad words from a file into memory
func LoadBadWords(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		word := scanner.Text()
		badWordsList = append(badWordsList, word)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	CompileBadWordsPattern()
	return nil
}
