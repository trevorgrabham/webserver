package timer

import (
	"fmt"

	"github.com/trevorgrabham/webserver/webserver/lib/svg"
)

func ResumeButton(id string, classes, htmx []string) string {
	return fmt.Sprintf(`
		<button
			class="svg-wrapper"
    	_="on click send startTimer to #timer-display"
    	hx-get="/resumeTimer"
    	hx-target="this"
    	hx-swap="outerHTML"
		>
			%s
		</button>`, svg.SVG(id, classes, htmx, svg.Start))
}