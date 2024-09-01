package timer

import (
	"html/template"

	svgtemplate "github.com/trevorgrabham/webserver/webserver/lib/templates/svg"
)

var StopTemplate = `
	<button
		class="svg-wrapper"
		_="on click send stopTimer to #timer-display" 
		hx-get="/stopTimer" 
		hx-target="#timer-buttons-container"
		hx-swap="outerHTML"
	>
		`+svgtemplate.StopSVGTemplate+`
	</button>`

var StopButton = template.Must(template.New("start-button-template").Parse(StopTemplate))