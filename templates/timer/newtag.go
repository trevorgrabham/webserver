package timer

import (
	"html/template"

	"github.com/trevorgrabham/webserver/webserver/lib/templates"
	svgtemplate "github.com/trevorgrabham/webserver/webserver/lib/templates/svg"
)

type TagTemplateData struct {
	Value string 
	Width int
	templates.ElementInfo 
	Remove templates.ElementInfo 
}

func NewTagTemplateData(id string, classes, htmx []string, hyperscript, value string, width int, removesvg templates.ElementInfo) TagTemplateData {
	return TagTemplateData{
		Value: value,
		Width: width, 
		ElementInfo: templates.NewElementInfo(id, classes, htmx, hyperscript),
		Remove: removesvg,
	}
}

var NewTagTemplate = `
	<input type="text"
		{{.ID}}
		{{.Classes}}
		{{.Htmx}}
		{{.Hyperscript}}
		name="temporary-tag" 
		placeholder="tags" 
		maxlength="255"
	/>

	<div class="tag-container">
		<div class="tag-wrapper">
			<input type="text" class="tag-display" name="tag" value="{{.Value}}" readonly style="width: {{.Width}}ch"/>
			{{template "removesvg" .Remove}}
		</div>
	</div>`

var main = template.Must(template.New("new-tag").Parse(NewTagTemplate))
var NewTag = template.Must(main.AddParseTree("removesvg", svgtemplate.RemoveSVG.Tree))