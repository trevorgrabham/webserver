package timer

import (
	"fmt"

	"github.com/trevorgrabham/webserver/webserver/lib/svg"
)

func CancelButton(id string, classes, htmx []string) string {
	return fmt.Sprintf(`
		<button
			class="svg-wrapper"
			_="on click call resetTimer()"
			hx-get="/resetTimer"
			hx-target="#timer-container"
		>
			%s
		</button>`, svg.SVG(id, classes, htmx, svg.Cancel))
}