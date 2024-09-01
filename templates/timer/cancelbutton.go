package timer

import (
	"html/template"

	svgtemplate "github.com/trevorgrabham/webserver/webserver/lib/templates/svg"
)

var CancelTemplate = `
	{{define "cancelbutton"}}
	<button
		class="svg-wrapper"
		_="on click call resetTimer()"
		hx-get="/resetTimer"
		hx-target="#timer-container"
	>
		{{template "cancelsvg" .Cancel}}
	</button>
	{{end}}`

var CancelButton = template.Must(template.Must(template.New("cancelbutton").Parse(CancelTemplate)).AddParseTree("cancelsvg", svgtemplate.CancelSVG.Tree))