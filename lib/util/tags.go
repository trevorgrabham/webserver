package util

import (
	"fmt"
	"html/template"
)

const DisabledMoreTagsButton = `
	<div 
		id="load-tag-summary-button"
		hx-swap-oob="true"
	>No More Data<div>`

func normalizeCount(n int64, max int64) string {
	if max == 0 {
		return "100%"
	}
	return fmt.Sprintf("%.2f", float64(n)/float64(max) * 100.0) + "%"
}

const tagDataTitle = "Tags"

const tagSummaryTemplate = `
	{{range .Tags}}
		<div class="tag-summary-row">
			<div 
				class="tag-summary-tag-title"
			>{{.Tag}}</div>
			<div 
				class="tag-summary-bar" 
				style="--default-width: {{normalizeCount .Count $.MaxCount}}" 
			>
				{{.Count}}
			</div>
		</div>
	{{end}}`

var allTagSummaryTemplate = fmt.Sprintf(`
	<div id="tag-summary-container">
		<h2 id="tag-summary-title">%s</h2>
		<div id="tag-summary-data-container">
				{{range .Tags}}
					<div class="tag-summary-row">
						<div 
							class="tag-summary-tag-title"
						>{{.Tag}}</div>
						<div 
							class="tag-summary-bar" 
							style="--default-width: {{normalizeCount .Count $.MaxCount}}" 
						>
							{{.Count}}
						</div>
					</div>
				{{end}}
		</div>
		<div 
			id="load-tag-summary-button"
			_="on click
					get @data-offset as Int + 10
					set @data-offset to it"
			hx-get="/tagSummary"
			hx-target="#tag-summary-data-container"
			hx-swap="beforeend"
			hx-vals='js:{offset: event.target.getAttribute("data-offset"), max: event.target.getAttribute("data-max-count")}'
			data-offset="11"
			data-max-count="{{$.MaxCount}}"
		>Load More<div>
	</div>
`, tagDataTitle)

var TagSummaryTemplateReady = template.Must(
	template.New("tagsummarytemplate").
	Funcs(template.FuncMap{
		"normalizeCount": normalizeCount,
	}).
	Parse(tagSummaryTemplate))

var AllTagSummaryTemplateReady = template.Must(
	template.New("tagsummarytemplate").
	Funcs(template.FuncMap{
		"normalizeCount": normalizeCount,
	}).
	Parse(allTagSummaryTemplate))