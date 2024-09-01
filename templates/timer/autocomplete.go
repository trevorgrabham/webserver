package timer

import (
	"html/template"

	"github.com/trevorgrabham/webserver/webserver/lib/timer"
)

type AutocompleteTemplateData struct {
	ID string
	timer.AutocompleteSuggestions 
}

func optionStart(option string, start int) string {
	return option[:start]
}

func optionEnd(option string, end int) string {
	return option[end:]
}

func optionMatch(option string, start, end int) string {
	return option[start:end]
}

const AutocompleteTemplate = `
	<div id="{{.ID}}" data-len="{{.Length}}" data-curr="0">
		{{range .Suggestions}}
			<div class="autocomplete-option-container">
				<div class="autocomplete-option">{{optionStart .Option .MatchStart}}<span class="matched-part">{{optionMatch .Option .MatchStart .MatchEnd}}</span>{{optionEnd .Option .MatchEnd}}</div>
			</div>
		{{end}}
	</div>`

var AutocompleteTemplateReady = 
template.Must(template.New("autocompletetemplate").
	Funcs(template.FuncMap{ 
		"optionStart": optionStart, 
		"optionMatch": optionMatch, 
		"optionEnd":   optionEnd, }).  
	Parse(AutocompleteTemplate))