package timer

import (
	"html/template"

	svgtemplate "github.com/trevorgrabham/webserver/webserver/lib/templates/svg"
)

var PlusTemplate = `
	{{define "plusbutton"}}
	<button 
		id="add-tag-button"
		class="svg-wrapper"
		hx-post="/addTag"
		hx-include="#tag-input"
		hx-params="temporary-tag"
		hx-target="#tags-container"
		hx-swap="beforeend"
	>
	{{template "plusbutton" .Plus}}
	</button>
	{{end}}`

var PlusButton = template.Must(template.Must(template.New("plusbutton").Parse(PlusTemplate)).AddParseTree("plusbutton", svgtemplate.PlusSVG.Tree))