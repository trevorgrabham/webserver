package timer

import (
	"fmt"

	"github.com/trevorgrabham/webserver/webserver/lib/svg"
)

func PlusButton(id string, classes, htmx []string) string {
	return fmt.Sprintf(`
		<button 
			id="add-tag-button"
			class="svg-wrapper"
			hx-post="/addTag"
			hx-include="#tag-input"
			hx-params="temporary-tag"
			hx-target="#tags-container"
			hx-swap="beforeend"
		>
			%s
		</button>
	`, svg.SVG(id, classes, htmx, svg.Plus))
}