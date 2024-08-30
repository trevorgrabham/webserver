package timer

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"
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

func FormHasFields(form map[string][]string, fields []string) (missingFields []string, values [][]string, err error) {
	if form == nil {
		return nil, nil, fmt.Errorf("FormHasFields(): map was nil")
	}
	if len(form) <= 0 {
		return fields, nil, nil
	}
	for _, field := range fields {
		value, ok := form[field]
		if !ok {
			missingFields = append(missingFields, field)
			continue
		}
		for i := range value {
			value[i] = strings.TrimSpace(value[i])
		}
		values = append(values, value)
	}
	return
}

func ParseTimer(time string) (hours, mins int, err error) {
	splitTime := strings.Split(time, ":")
	switch len(splitTime) {
	case 2:		// MM:SS
		mins, err = strconv.Atoi(splitTime[0])
		return
	case 3:		// HH:MM:SS
		hours, err = strconv.Atoi(splitTime[0])
		if err != nil {
			return 0, 0, err
		}
		mins, err = strconv.Atoi(splitTime[1])
		return
	default:	// Shouldn't ever hit here 
		return 0, 0, fmt.Errorf("ParseTimer(%v): Unrecognized time format", time)
	}
}

func StartingTime(mins int) (startDate time.Time, err error) {
	delta, err := time.ParseDuration(fmt.Sprintf("-%dm", mins))
	if err != nil {
		return 
	}
	startDate = time.Now().Add(delta)
	return
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