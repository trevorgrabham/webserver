package timer

import (
	"slices"
	"strings"
)

type AutocompleteSuggestions struct {
	Suggestions 				[]AutocompleteSuggestion
	Length 							int
}

type AutocompleteSuggestion struct {
	Option 					string 
	MatchStart			int
	MatchEnd 				int
}

// Returns any members from 'list' in a sorted order, with those that match at the beginning of the member having a higher precedence than a match in the middle of the member
func FilterFromPartialString(partial string, list []string) AutocompleteSuggestions {
	beginningMatched := make([]AutocompleteSuggestion, 0)
	anywhereMatched := make([]AutocompleteSuggestion, 0)
	for i := range list {
		switch index := strings.Index(list[i], partial); {
		case index == 0:
			beginningMatched = append(beginningMatched, AutocompleteSuggestion{Option: list[i], MatchStart: index, MatchEnd: index + len(partial)})
		case index > 0:
			anywhereMatched = append(anywhereMatched, AutocompleteSuggestion{Option: list[i], MatchStart: index, MatchEnd: index + len(partial)})
		}
	}
	cmp := func(a, b AutocompleteSuggestion) int { return strings.Compare(a.Option, b.Option) }
	slices.SortFunc(beginningMatched, cmp)
	slices.SortFunc(anywhereMatched, cmp)
	return AutocompleteSuggestions{Suggestions: append(beginningMatched, anywhereMatched...), Length: len(beginningMatched) + len(anywhereMatched)}
}