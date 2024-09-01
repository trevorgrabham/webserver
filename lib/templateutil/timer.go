package templateutil

import (
	"fmt"
	"html/template"
	"strings"

	"github.com/trevorgrabham/webserver/webserver/lib/timer"
)

type ElementInfo struct {
	ID          template.HTMLAttr
	Classes     template.HTMLAttr
	Htmx        template.HTMLAttr
	Hyperscript template.HTMLAttr
}

type AutocompleteTemplateData struct {
	ID string
	timer.AutocompleteSuggestions 
}

type TagData struct {
	Input ElementInfo
	SVG ElementInfo
	Value string
	Width int
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

// id and classes are just the names themselves, htmx and hyperscript should have their attributes specified.
//
// Ex. NewElementInfo("foo", []string{"bar"}, []string{`hx-get="/foo"`, `hx-swap="outerHTML"`}, `_="on click set my.value to 1"`)
func NewElementInfo(id string, classes []string, htmx []string, hyperscript string) ElementInfo {
	return ElementInfo{
		ID:          template.HTMLAttr(fmt.Sprintf(`id="%s"`, id)),
		Classes:     template.HTMLAttr(fmt.Sprintf(`class="%s"`, strings.Join(classes, " "))),
		Htmx:        template.HTMLAttr(strings.Join(htmx, "\n")),
		Hyperscript: template.HTMLAttr(hyperscript),
	}
}

func NewFormData(plus, success, cancel ElementInfo) map[string]ElementInfo {
	return map[string]ElementInfo{
		"Plus": plus,
		"Success": success,
		"Cancel": cancel,
	}
}

func NewAutocompleteData(id string, suggestions timer.AutocompleteSuggestions) AutocompleteTemplateData {
	return AutocompleteTemplateData{
		ID: id, 
		AutocompleteSuggestions: suggestions,
	}
}

func NewTagData(inputElement, svgElement ElementInfo, value string, width int) TagData {
	return TagData{
		Input: inputElement,
		SVG: svgElement,
		Value: value,
		Width: width,
	}
}

func AutocompleteFuncMap() template.FuncMap {
	return template.FuncMap{
		"optionStart": optionStart,
		"optionMatch": optionMatch,
		"optionEnd": optionEnd,
	}
}