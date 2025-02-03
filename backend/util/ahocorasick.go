package util

import (
	"sort"
	"unsafe"
)

//TODO: it won't find anything

// Matcher is the pattern matching state machine.
type Matcher struct {
	base        []int        // base array in the double-array trie
	check       []int        // check array in the double-array trie
	fail        []int        // fail function
	output      [][]int      // output function
	runeIndices map[rune]int // mapping from runes to compact indices
	runes       []rune       // list of unique runes
	patterns    [][]rune     // store the patterns
}

// Match represents a matched pattern in the text.
type Match struct {
	Word  string // the matched pattern
	Index int    // the start index of the match
}

// CompileStrings compiles a Matcher from a slice of strings.
func CompileStrings(words []string) *Matcher {
	wordRuneSlices := make([][]rune, len(words))
	for i, word := range words {
		wordRuneSlices[i] = []rune(word)
	}
	return compile(wordRuneSlices)
}

// compile compiles the patterns into the Matcher.
func compile(words [][]rune) *Matcher {
	m := &Matcher{
		runeIndices: make(map[rune]int),
	}

	// Build rune to index mapping
	for _, word := range words {
		for _, r := range word {
			if _, exists := m.runeIndices[r]; !exists {
				m.runeIndices[r] = len(m.runes)
				m.runes = append(m.runes, r)
			}
		}
	}

	m.patterns = words

	// Initialize the base, check, fail, and output arrays
	initialSize := 1
	m.base = make([]int, initialSize)
	m.check = make([]int, initialSize)
	m.fail = make([]int, initialSize)
	m.output = make([][]int, initialSize)

	// Sort the words to ensure deterministic automaton construction
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
		m.ensureStateCapacity(node.state)
		m.base[node.state] = base

		i := node.start
		for _, edge := range edges {
			offset := m.runeIndices[edge]
			newState := base + offset

			m.ensureStateCapacity(newState)

			m.check[newState] = node.state

			// Set fail link
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
			for i < node.end && node.depth < len(words[i]) && words[i][node.depth] == edge {
				if node.depth+1 == len(words[i]) {
					m.output[newState] = append(m.output[newState], i) // Store pattern index
				}
				i++
				newNodeEnd++
			}

			// Enqueue the next trie node
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
	base := 1 // Start from 1 to avoid conflicts with initial state
search:
	for {
		for _, edge := range edges {
			offset := m.runeIndices[edge]
			state := base + offset
			if state < len(m.check) && m.check[state] != 0 {
				base++
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
	// Convert string to rune slice without allocation
	textRunes := *(*[]rune)(unsafe.Pointer(&text))
	return m.findAll(textRunes)
}

// findAll finds all matches in the given text.
func (m *Matcher) findAll(text []rune) []*Match {
	matches := make([]*Match, 0)
	state := 0
	for i, r := range text {
		offset, exists := m.runeIndices[r]
		if !exists {
			// Follow fail links until we reach the root or find a match
			for state != 0 && !exists {
				state = m.fail[state]
				exists = m.hasEdge(state, offset)
			}
			if !exists {
				state = 0
				continue
			}
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
			for _, idx := range m.output[state] {
				pattern := m.patterns[idx]
				length := len(pattern)
				start := i - length + 1
				if start >= 0 {
					matches = append(matches, &Match{
						Word:  string(pattern),
						Index: start,
					})
				}
			}
		}
	}
	return matches
}
