package timer

import (
	"fmt"

	"github.com/trevorgrabham/webserver/webserver/lib/svg"
)



func PauseButton(id string, classes []string, htmx []string) string {
	return fmt.Sprintf(`
	<button
		class="svg-wrapper"
		_="on click send stopTimer to #timer-display" 
		hx-get="/pauseTimer" 
		hx-target="this"
		hx-swap="outerHTML"
	>
		%s
	</button>`, svg.SVG(id, classes, htmx, svg.Pause))
}
