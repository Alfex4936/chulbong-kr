package util

import (
	"bufio"
	"context"
	"errors"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/rrethy/ahocorasick"
)

var (
	badWordsList []string
	badWordRegex *regexp.Regexp
	matcher      *ahocorasick.Matcher
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
		CompileBadWordsPattern()
		return false, errors.New("bad words pattern not compiled")
	}

	return badWordRegex.MatchString(input), nil
}

// USE CheckForBadWords
func CheckForBadWordsWithGo(input string) (bool, error) {
	for _, word := range badWordsList {
		if word == "" {
			continue
		}

		// Check if the bad word is a substring of the input
		if strings.Contains(input, word) {
			return true, nil
		}
	}
	return false, nil
}

// USE CheckForBadWords
func CheckForBadWordsWithGoRoutine(input string) (bool, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Ensures context is canceled once we return.

	resultChan := make(chan bool)
	var wg sync.WaitGroup

	for _, word := range badWordsList {
		wg.Add(1)
		go func(w string) {
			defer wg.Done()
			select {
			case <-ctx.Done():
				return // Early exit on context cancellation.
			default:
				if w == "" {
					return // Skip empty words.
				}
				if strings.Contains(input, w) {
					resultChan <- true
					cancel() // Found a bad word, signal to cancel other goroutines.
				}
			}
		}(word)
	}

	// Close the resultChan once all goroutines have finished.
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Process results.
	for result := range resultChan {
		if result {
			return true, nil
		}
	}

	return false, nil
}

func ReplaceBadWords(input string) (string, error) {
	if badWordRegex == nil {
		return input, errors.New("bad words pattern not compiled")
	}

	// Use ReplaceAllStringFunc to replace bad words with asterisks
	return badWordRegex.ReplaceAllStringFunc(input, func(match string) string {
		// Replace each character of the bad word with an asterisk
		return strings.Repeat("*", len([]rune(match)))
	}), nil
}

func RemoveURLs(input string) string {
	// Compile the regular expression for matching URLs
	urlRegex := regexp.MustCompile(`\bhttps?:\/\/\S+\b`)
	// Replace URLs with an empty string
	return urlRegex.ReplaceAllString(input, "")
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

// LoadBadWords loads bad words from a file into memory with optimizations.
func LoadBadWords(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Estimate the number of words if known or use a high number.
	const estimatedWords = 1000
	badWordsList = make([]string, 0, estimatedWords)

	// Create a buffer and attach it to scanner.
	scanner := bufio.NewScanner(file)
	const maxCapacity = 10 * 1024 // 10KB;
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)

	for scanner.Scan() {
		word := scanner.Text()
		badWordsList = append(badWordsList, word)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	// Optimize memory usage by shrinking the slice to the actual number of words.
	badWordsList = append([]string{}, badWordsList...)

	// go CompileBadWordsPattern() // Compile in a goroutine if it's safe to do asynchronously.

	// Compile the list of bad words into a trie (Aho-corasick Double-Array Trie)
	matcher = ahocorasick.CompileStrings(badWordsList)
	return nil
}

// CheckForBadWordsUsingTrie checks if the input contains any bad words using Aho-Corasick trie
func CheckForBadWordsUsingTrie(input string) (bool, error) {
	if matcher == nil {
		matcher = ahocorasick.CompileStrings(badWordsList)
		return false, os.ErrNotExist
	}
	matches := matcher.FindAllString(input)
	return len(matches) > 0, nil
}
