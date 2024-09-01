package timer

import (
	"html/template"

	svgtemplate "github.com/trevorgrabham/webserver/webserver/lib/templates/svg"
)

var SuccessTemplate = `
	{{define "successbutton"}}
	<button 
	 	id="submit-timer-form"
		class="svg-wrapper"
		type="submit"
		_="on click call resetTimer()"
	>
		{{template "successsvg" .Success}}
	</button>
	{{end}}`

var SuccessButton = template.Must(template.Must(template.New("successbutton").Parse(SuccessTemplate)).AddParseTree("successsvg", svgtemplate.SuccessSVG.Tree))