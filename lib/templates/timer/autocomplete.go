package timer

import (
	"fmt"
	"html/template"
)

func optionStart(option string, start int) string {
	return option[:start]
}

func optionEnd(option string, end int) string {
	return option[end:]
}

func optionMatch(option string, start, end int) string {
	return option[start:end]
}

const autocompleteTemplate = `
	<div%sdata-len="{{.Length}}" data-curr="0">
		{{range .Suggestions}}
			<div class="autocomplete-option-container">
				<div class="autocomplete-option">{{optionStart .Option .MatchStart}}<span class="matched-part">{{optionMatch .Option .MatchStart .MatchEnd}}</span>{{optionEnd .Option .MatchEnd}}</div>
			</div>
		{{end}}
	</div>`

func AutocompleteTemplateReady(id string) *template.Template {
	idString := " "
	if id != "" {
		idString = fmt.Sprintf(` id="%s" `, id)
	}
	return template.Must(template.New("autocompletetemplate").
		Funcs(template.FuncMap{
			"optionStart": optionStart,
			"optionMatch": optionMatch,
			"optionEnd":   optionEnd,
		}).
		Parse(fmt.Sprintf(autocompleteTemplate, idString)))
}