package fuzzy

import (
	"strings"

	fuzzypkg "github.com/sahilm/fuzzy"
)

// Match returns true if input fuzzy-matches the target string (case-insensitive).
func Match(input, target string) bool {
	if input == "" {
		return true
	}
	matches := fuzzypkg.Find(strings.ToLower(input), []string{strings.ToLower(target)})
	return len(matches) > 0
}

// Searcher returns a Searcher function compatible with promptui.Select.Searcher.
// It checks whether the input fuzzy-matches the item at the given index.
func Searcher(items []string) func(input string, index int) bool {
	return func(input string, index int) bool {
		return Match(input, items[index])
	}
}
