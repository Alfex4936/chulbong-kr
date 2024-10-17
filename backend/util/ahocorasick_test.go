package util

import (
	"testing"
)

func TestNonASCIICharacters(t *testing.T) {
	patterns := []string{"안녕하세요", "안녕", "하세요"}
	matcher := CompileStrings(patterns)
	text := "안녕하세요 여러분"
	expectedMatches := []*Match{
		{Word: "안녕하세요", Index: 0},
		{Word: "안녕", Index: 0},
		{Word: "하세요", Index: 2},
	}
	matches := matcher.FindAllString(text)
	if !compareMatches(matches, expectedMatches) {
		t.Errorf("Expected matches %v, got %v", expectedMatches, matches)
	}
}

func TestLongPatterns(t *testing.T) {
	pattern := ""
	for i := 0; i < 1000; i++ {
		pattern += "a"
	}
	patterns := []string{pattern}
	matcher := CompileStrings(patterns)
	text := ""
	for i := 0; i < 1000; i++ {
		text += "a"
	}
	matches := matcher.FindAllString(text)
	if len(matches) != 1 || matches[0].Index != 0 {
		t.Errorf("Expected one match at index 0, got %v", matches)
	}
}

func TestMultipleMatchesAtSamePosition(t *testing.T) {
	patterns := []string{"he", "he", "he"}
	matcher := CompileStrings(patterns)
	text := "he"
	matches := matcher.FindAllString(text)
	if len(matches) != 3 {
		t.Errorf("Expected 3 matches, got %d", len(matches))
	}
}

func TestLongText(t *testing.T) {
	patterns := []string{"test", "long", "text"}
	matcher := CompileStrings(patterns)
	text := ""
	for i := 0; i < 10000; i++ {
		text += "This is a long text for testing."
	}
	matches := matcher.FindAllString(text)
	if len(matches) == 0 {
		t.Error("Expected matches, got none")
	}
}

func TestSpecialCharacters(t *testing.T) {
	patterns := []string{"$", "^", "*", "+", "."}
	matcher := CompileStrings(patterns)
	text := "This $ is a ^ test * with + special . characters."
	expectedMatches := []*Match{
		{Word: "$", Index: 5},
		{Word: "^", Index: 12},
		{Word: "*", Index: 19},
		{Word: "+", Index: 26},
		{Word: ".", Index: 36},
	}
	matches := matcher.FindAllString(text)
	if !compareMatches(matches, expectedMatches) {
		t.Errorf("Expected matches %v, got %v", expectedMatches, matches)
	}
}

func TestUnicodeCharacters(t *testing.T) {
	patterns := []string{"😊", "🚀", "🌟"}
	matcher := CompileStrings(patterns)
	text := "Hello 😊! Let's go to the moon 🚀 and shine like a star 🌟."
	expectedMatches := []*Match{
		{Word: "😊", Index: 6},
		{Word: "🚀", Index: 31},
		{Word: "🌟", Index: 53},
	}
	matches := matcher.FindAllString(text)
	if !compareMatches(matches, expectedMatches) {
		t.Errorf("Expected matches %v, got %v", expectedMatches, matches)
	}
}

func compareMatches(a, b []*Match) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i].Index != b[i].Index || a[i].Word != b[i].Word {
			return false
		}
	}
	return true
}
