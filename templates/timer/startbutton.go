package timer

import (
	"html/template"

	svgtemplate "github.com/trevorgrabham/webserver/webserver/lib/templates/svg"
)

var StartButtonTemplate = `
	<button
    class="svg-wrapper"
    _="on click send startTimer to #timer-display"
    hx-get="/startTimer"
    hx-target="#timer-buttons-container"
  >
		` + svgtemplate.StartSVGTemplate + `
	</button>`

var StartButton = template.Must(template.New("start-button-template").Parse(StartButtonTemplate))