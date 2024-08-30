package timer

import (
	"fmt"
	"strings"

	"github.com/trevorgrabham/webserver/webserver/lib/svg"
)

func NewTagTemplate(value string, width int, id string, classes, htmx []string) string {
	var idString, classString, htmxString string
	if id != "" { idString = fmt.Sprintf(`id="%s"`, id) }
	if classes != nil { classString = fmt.Sprintf(`class="%s"`, strings.Join(classes, " ")) }
	if htmx != nil { htmxString = strings.Join(htmx, "\n") }
	return fmt.Sprintf(`
		<input type="text"
			%s
			%s
			%s
			name="temporary-tag" 
			placeholder="tags" 
			maxlength="255"
		/>

	 	<div class="tag-container">
			<div class="tag-wrapper">
				<input type="text" class="tag-display" name="tag" value="%s" readonly style="width: %dch"/>
		 		%s
			</div>
		</div>
	`, idString, classString, htmxString, value, width, svg.SVG("", []string{"button-tag-remove"}, []string{`hx-get="/removeTag"`, `hx-target="closest .tag-container"`, `hx-swap="outerHTML"`}, svg.Remove))
}