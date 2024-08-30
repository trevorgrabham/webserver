package timer

import (
	"fmt"

	"github.com/trevorgrabham/webserver/webserver/lib/svg"
)

func SuccessButton(id string, classes, htmx []string) string {

	return fmt.Sprintf(`
		<button 
	 		id="submit-timer-form"
			class="svg-wrapper"
			type="submit"
			_="on click call resetTimer()"
		>
			%s
		</button>`, svg.SVG(id, classes, htmx, svg.Success))
}