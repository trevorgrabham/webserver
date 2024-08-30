package timer

import (
	"fmt"

	"github.com/trevorgrabham/webserver/webserver/lib/svg"
)

func StopButton(id string, classes, htmx []string) string {
	return fmt.Sprintf(`
		<button
			class="svg-wrapper"
			_="on click send stopTimer to #timer-display" 
			hx-get="/stopTimer" 
			hx-target="#timer-buttons-container"
			hx-swap="outerHTML"
		>
			%s
		</button>`, svg.SVG(id, classes, htmx, svg.Stop))
}