package timer

import (
	"html/template"

	svgtemplate "github.com/trevorgrabham/webserver/webserver/lib/templates/svg"
)

var ResumeButtonTemplate = `
	<button
		class="svg-wrapper"
    _="on click send startTimer to #timer-display"
    hx-get="/resumeTimer"
    hx-target="this"
    hx-swap="outerHTML"
	>
		`+svgtemplate.StartSVGTemplate+`
	</button>`

var ResumeButton = template.Must(template.New("start-button-template").Parse(ResumeButtonTemplate))