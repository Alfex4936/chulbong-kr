package util

import (
	"bufio"
	"context"
	"errors"
	"os"
	"regexp"
	"strings"
	"sync"
	"unicode/utf8"
	"unsafe"

	"github.com/rrethy/ahocorasick"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var urlRegex = regexp.MustCompile(`https?://[^\s]+`)

type BadWordUtil struct {
	BadWordsListByte [][]byte
	BadWordsList     []string
	BadWordRegex     *regexp.Regexp
	Matcher          *ahocorasick.Matcher
	ByteMatcher      *ahocorasick.Matcher
}

func NewBadWordUtil() *BadWordUtil {
	return &BadWordUtil{}
}

func RegisterBadWordUtilLifecycle(lifecycle fx.Lifecycle, badWordUtil *BadWordUtil, logger *zap.Logger) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			filePath := os.Getenv("BAD_WORDS_FILE_PATH")
			if filePath == "" {
				filePath = "badwords.txt" // Provide a default path if not set
			}
			logger.Info("Loading bad words from file", zap.String("path", filePath))
			if err := badWordUtil.LoadBadWords(filePath); err != nil {
				logger.Error("Failed to load bad words", zap.Error(err))
				return err
			}
			badWordUtil.LoadBadWordsByte(filePath)

			logger.Info("Bad words loaded successfully")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			// Cleanup if necessary
			return nil
		},
	})
}

func (b *BadWordUtil) CompileBadWordsPattern() error {
	var pattern strings.Builder
	pattern.WriteString(`(`)
	for i, word := range b.BadWordsList {
		if word == "" {
			continue
		}
		pattern.WriteString(regexp.QuoteMeta(word))
		if i < len(b.BadWordsList)-1 {
			pattern.WriteString(`|`)
		}
	}
	pattern.WriteString(`)`)

	var err error
	b.BadWordRegex, err = regexp.Compile(pattern.String())
	return err
}

func (b *BadWordUtil) CheckForBadWords(input string) (bool, error) {
	if b.BadWordRegex == nil {
		b.CompileBadWordsPattern()
		return false, errors.New("bad words pattern not compiled")
	}

	return b.BadWordRegex.MatchString(input), nil
}

// USE CheckForBadWords
func (b *BadWordUtil) CheckForBadWordsWithGo(input string) (bool, error) {
	for _, word := range b.BadWordsList {
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
func (b *BadWordUtil) CheckForBadWordsWithGoRoutine(input string) (bool, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Ensures context is canceled once we return.

	resultChan := make(chan bool)
	var wg sync.WaitGroup

	for _, word := range b.BadWordsList {
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

func (b *BadWordUtil) ReplaceBadWords(input string) (string, error) {
	if b.Matcher == nil {
		b.Matcher = ahocorasick.CompileStrings(b.BadWordsList)
		return input, errors.New("bad words matcher not initialized")
	}

	matches := b.Matcher.FindAllString(input)
	if len(matches) == 0 {
		return input, nil // No matches, return original input
	}

	runes := []rune(input)
	replaced := make([]bool, len(runes))
	runeIndices := computeRuneIndices(input) // Precompute indices

	// Apply replacements for each match
	for _, match := range matches {
		matchStart := runeIndices[match.Index]
		matchEnd := runeIndices[match.Index+len(match.Word)]

		for i := matchStart; i < matchEnd; i++ {
			if !replaced[i] {
				runes[i] = '*'
				replaced[i] = true
			}
		}
	}

	return string(runes), nil
}

func (b *BadWordUtil) ReplaceBadWordsInBytes(input []byte) ([]byte, error) {
	if b.Matcher == nil {
		b.ByteMatcher = ahocorasick.CompileByteSlices(b.BadWordsListByte)
		return input, errors.New("bad words matcher not initialized")
	}

	matches := b.Matcher.FindAllByteSlice(input)
	if len(matches) == 0 {
		return input, nil // No matches, return original input
	}

	// Create a map to mark bytes that are part of matches
	matchedBytes := make([]bool, len(input))
	for _, match := range matches {
		for i := match.Index; i < match.Index+len(match.Word); i++ {
			matchedBytes[i] = true
		}
	}

	// Process the input bytes, decoding runes, and building output bytes
	var output []byte
	i := 0
	for i < len(input) {
		if matchedBytes[i] {
			// Replace the entire rune with '*'
			_, size := utf8.DecodeRune(input[i:])
			output = append(output, '*')
			i += size
			// Skip any additional bytes that are part of the matched word
			// These are already accounted for in matchedBytes
			// No need to do anything else here
		} else {
			// Copy the rune as is
			_, size := utf8.DecodeRune(input[i:])
			output = append(output, input[i:i+size]...)
			i += size
		}
	}

	return output, nil
}

func (b *BadWordUtil) ProcessChatMessage(message []byte) ([]byte, error) {
	// Remove URLs directly from []byte
	message = RemoveURLsFromBytes(message)

	// Replace bad words directly in []byte
	return b.ReplaceBadWordsInBytes(message)
}

func RemoveURLs(input string) string {
	// Replace URLs with an empty string
	return urlRegex.ReplaceAllString(input, "")
}

func RemoveURLsFromBytes(message []byte) []byte {
	return urlRegex.ReplaceAll(message, []byte(""))
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
func (b *BadWordUtil) LoadBadWords(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Estimate the number of words if known or use a high number.
	const estimatedWords = 1000
	b.BadWordsList = make([]string, 0, estimatedWords)

	// Create a buffer and attach it to scanner.
	scanner := bufio.NewScanner(file)
	const maxCapacity = 10 * 1024 // 10KB;
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)

	for scanner.Scan() {
		word := scanner.Text()
		b.BadWordsList = append(b.BadWordsList, word)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	// Optimize memory usage by shrinking the slice to the actual number of words.
	b.BadWordsList = append([]string{}, b.BadWordsList...)

	// go CompileBadWordsPattern() // Compile in a goroutine if it's safe to do asynchronously.

	// Compile the list of bad words into a trie (Aho-corasick Double-Array Trie)
	b.Matcher = ahocorasick.CompileStrings(b.BadWordsList)
	return nil
}

func (b *BadWordUtil) LoadBadWordsByte(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	const estimatedWords = 600
	b.BadWordsListByte = make([][]byte, 0, estimatedWords)

	scanner := bufio.NewScanner(file)
	const maxCapacity = 10 * 1024 // 10KB
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)

	for scanner.Scan() {
		word := make([]byte, len(scanner.Bytes()))
		copy(word, scanner.Bytes())
		b.BadWordsListByte = append(b.BadWordsListByte, word)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	b.ByteMatcher = ahocorasick.CompileByteSlices(b.BadWordsListByte)
	return nil
}

// CheckForBadWordsUsingTrie checks if the input contains any bad words using Aho-Corasick trie
func (b *BadWordUtil) CheckForBadWordsUsingTrie(input string) (bool, error) {
	if b.Matcher == nil {
		b.Matcher = ahocorasick.CompileStrings(b.BadWordsList)
		return false, os.ErrNotExist
	}
	matches := b.Matcher.FindAllString(input)
	return len(matches) > 0, nil
}

// - fasthttp

// BytesToString converts a byte slice to a string without making a copy.
//
// Warning: This method uses unsafe operations. The conversion is safe as long as the original
// byte slice is not modified after conversion, as the resulting string will reference
// the same underlying memory.
//
// Example:
//
//	b := []byte("Hello")
//	s := BytesToString(b) // Converts []byte to string without memory allocation
//	fmt.Println(s)        // Prints: Hello
//
// Important: Do not modify the byte slice `b` after calling this function, as the string `s`
// references the same memory, and changes to `b` will lead to undefined behavior.
func BytesToString(b []byte) string {
	// if len(b) == 0 {
	// 	return ""
	// }
	return unsafe.String(&b[0], len(b))
	// return *(*string)(unsafe.Pointer(&b))
}

// StringToBytes converts a string to a byte slice without making a copy.
//
// Warning: This method uses unsafe operations. The conversion is safe as long as the
// resulting byte slice is not modified. Since strings in Go are immutable, modifying
// the byte slice will result in undefined behavior.
//
// Example:
//
//	s := "Hello"
//	b := StringToBytes(s) // Converts string to []byte without memory allocation
//	fmt.Println(b)        // Prints: [72 101 108 108 111] (ASCII values of "Hello")
//
// Important: Do not modify the byte slice `b` after calling this function, as strings in Go
// are immutable, and modifying the byte slice can lead to undefined behavior.
func StringToBytes(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

// Precompute rune indices for the whole string to avoid recalculating repeatedly
func computeRuneIndices(input string) []int {
	runeIndices := make([]int, len(input)+1)
	runeCount := 0
	for i := range input {
		runeIndices[i] = runeCount
		_, size := utf8.DecodeRuneInString(input[i:])
		runeCount += 1
		i += size - 1 // Adjust for the size of the current rune
	}
	runeIndices[len(input)] = runeCount // End of string index
	return runeIndices
}

// By official doc, Slice(ptr, len) => like (*[len]T)(ptr)[:]
func SliceFromPointer[T any](base unsafe.Pointer, length int) []T {
	return unsafe.Slice((*T)(base), length)
}
