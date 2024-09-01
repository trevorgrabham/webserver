package timer

import (
	"html/template"

	svgtemplate "github.com/trevorgrabham/webserver/webserver/lib/templates/svg"
)

var PauseTemplate = `
	<button
		class="svg-wrapper"
		_="on click send stopTimer to #timer-display" 
		hx-get="/pauseTimer" 
		hx-target="this"
		hx-swap="outerHTML"
	>
		`+svgtemplate.PauseSVGTemplate+`
	</button>`

var PauseButton = template.Must(template.New("pause-button-template").Parse(PauseTemplate))