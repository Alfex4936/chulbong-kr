// original code: https://github.com/RRethy/ahocorasick
// edited version
package util

import (
	"fmt"
	"sort"
	"unicode/utf8"
)

// Matcher is the pattern matching state machine.
type Matcher struct {
	base        []int        // base array in the double array trie
	check       []int        // check array in the double array trie
	fail        []int        // fail function
	output      [][]int      // output function
	runeIndices map[rune]int // mapping from runes to indices
	runes       []rune       // list of unique runes
}

// Match represents a matched pattern in the text.
type Match struct {
	Word  string // the matched pattern
	Index int    // the start index of the match
}

// CompileByteSlices compiles a Matcher from a slice of byte slices. This Matcher can be
// used to find occurrences of each pattern in a text.
func CompileByteSlices(words [][]byte) *Matcher {
	wordRuneSlices := make([][]rune, len(words))
	for i, word := range words {
		runes, err := bytesToRunes(word)
		if err != nil {
			// Handle invalid UTF-8 by skipping the pattern or logging.
			// For simplicity, we'll skip invalid patterns.
			continue
		}
		wordRuneSlices[i] = runes
	}
	return compile(wordRuneSlices)
}

// CompileStrings compiles a Matcher from a slice of strings. This Matcher can
// be used to find occurrences of each pattern in a text.
func CompileStrings(words []string) *Matcher {
	wordRuneSlices := make([][]rune, len(words))
	for i, word := range words {
		wordRuneSlices[i] = []rune(word)
	}
	return compile(wordRuneSlices)
}

func compile(words [][]rune) *Matcher {
	m := &Matcher{
		base:        []int{0},
		check:       []int{0},
		fail:        []int{0},
		output:      [][]int{nil},
		runeIndices: make(map[rune]int),
	}

	// Build rune to index mapping
	for _, word := range words {
		for _, r := range word {
			if _, exists := m.runeIndices[r]; !exists {
				m.runeIndices[r] = len(m.runeIndices)
				m.runes = append(m.runes, r)
			}
		}
	}

	// Sort the words to ensure deterministic automaton construction.
	sort.Slice(words, func(i, j int) bool {
		return lessRuneSlice(words[i], words[j])
	})

	type trieNode struct {
		state int
		depth int
		start int
		end   int
	}

	queue := []trieNode{{state: 0, depth: 0, start: 0, end: len(words)}}

	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]

		if node.end <= node.start {
			continue
		}

		edges := collectEdges(words, node.depth, node.start, node.end)

		base := m.findBase(edges)
		m.base[node.state] = base

		i := node.start
		for _, edge := range edges {
			offset, exists := m.runeIndices[edge]
			if !exists {
				continue // Skip if rune not in mapping
			}

			newState := base + offset

			m.ensureStateCapacity(newState)

			m.check[newState] = node.state

			// Add fail links
			var failState int
			if node.depth == 0 {
				failState = 0
			} else {
				failState = m.getFailState(m.fail[node.state], offset)
			}
			m.fail[newState] = failState

			// Merge output functions
			if len(m.output[failState]) > 0 {
				m.output[newState] = append(m.output[newState], m.output[failState]...)
			}

			// Add output for complete words
			newNodeStart := i
			newNodeEnd := i
			for i < node.end && words[i][node.depth] == edge {
				if node.depth+1 == len(words[i]) {
					m.output[newState] = append(m.output[newState], len(words[i]))
				}
				i++
				newNodeEnd++
			}

			// Enqueue the next trie node if necessary
			if newNodeStart < newNodeEnd {
				queue = append(queue, trieNode{
					state: newState,
					depth: node.depth + 1,
					start: newNodeStart,
					end:   newNodeEnd,
				})
			}
		}
	}

	return m
}

// lessRuneSlice compares two rune slices lexicographically.
func lessRuneSlice(a, b []rune) bool {
	minLen := len(a)
	if len(b) < minLen {
		minLen = len(b)
	}
	for i := 0; i < minLen; i++ {
		if a[i] != b[i] {
			return a[i] < b[i]
		}
	}
	return len(a) < len(b)
}

// collectEdges collects the unique edges (runes) at the given depth.
func collectEdges(words [][]rune, depth, start, end int) []rune {
	edgeSet := make(map[rune]struct{})
	for i := start; i < end; i++ {
		if depth < len(words[i]) {
			edgeSet[words[i][depth]] = struct{}{}
		}
	}
	edges := make([]rune, 0, len(edgeSet))
	for edge := range edgeSet {
		edges = append(edges, edge)
	}
	sort.Slice(edges, func(i, j int) bool { return edges[i] < edges[j] })
	return edges
}

// ensureStateCapacity ensures that the state arrays have enough capacity.
func (m *Matcher) ensureStateCapacity(state int) {
	if state >= len(m.base) {
		newSize := state + 1
		m.base = append(m.base, make([]int, newSize-len(m.base))...)
		m.check = append(m.check, make([]int, newSize-len(m.check))...)
		m.fail = append(m.fail, make([]int, newSize-len(m.fail))...)
		m.output = append(m.output, make([][]int, newSize-len(m.output))...)
	}
}

// getFailState computes the fail state for a given state and offset.
func (m *Matcher) getFailState(failState, offset int) int {
	for failState != 0 && !m.hasEdge(failState, offset) {
		failState = m.fail[failState]
	}
	if m.hasEdge(failState, offset) {
		return m.base[failState] + offset
	}
	return 0
}

// findBase finds a suitable base value for the given edges.
func (m *Matcher) findBase(edges []rune) int {
	var base int
search:
	for {
		base++
		for _, edge := range edges {
			offset, exists := m.runeIndices[edge]
			if !exists {
				continue search
			}
			state := base + offset
			if state < len(m.check) && m.check[state] != 0 {
				continue search
			}
		}
		break
	}
	return base
}

// hasEdge checks if there is an edge from the given state with the given offset.
func (m *Matcher) hasEdge(state, offset int) bool {
	nextState := m.base[state] + offset
	return nextState < len(m.check) && m.check[nextState] == state
}

// FindAllString finds all instances of the patterns in the text.
func (m *Matcher) FindAllString(text string) []*Match {
	return m.FindAllRuneSlice([]rune(text))
}

// FindAllRuneSlice finds all instances of the patterns in the rune slice.
func (m *Matcher) FindAllRuneSlice(text []rune) []*Match {
	var matches []*Match
	state := 0
	for i, r := range text {
		offset, exists := m.runeIndices[r]
		if !exists {
			state = 0
			continue
		}
		for state != 0 && !m.hasEdge(state, offset) {
			state = m.fail[state]
		}
		if m.hasEdge(state, offset) {
			state = m.base[state] + offset
		} else {
			state = 0
		}
		if len(m.output[state]) > 0 {
			for _, length := range m.output[state] {
				start := i - length + 1
				if start >= 0 {
					matches = append(matches, &Match{
						Word:  string(text[start : i+1]),
						Index: start,
					})
				}
			}
		}
	}
	return matches
}

// bytesToRunes converts a byte slice to a rune slice, ensuring valid UTF-8 encoding.
func bytesToRunes(text []byte) ([]rune, error) {
	var runes []rune
	for len(text) > 0 {
		r, size := utf8.DecodeRune(text)
		if r == utf8.RuneError && size == 1 {
			return nil, fmt.Errorf("invalid UTF-8 encoding")
		}
		runes = append(runes, r)
		text = text[size:]
	}
	return runes, nil
}
