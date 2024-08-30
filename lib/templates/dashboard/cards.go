package dashboard

import (
	"fmt"
	"text/template"
)

func formatTotalMin(t int64) string {
	if t > 60 {
		return fmt.Sprintf("%dh%dm", t/60, t%60)
	}
	return fmt.Sprintf("%dm", t%60)
}

func calcColor(n int64) string {
	switch {
	case n > 90:
		{
			return "great"
		}
	case n > 60:
		{
			return "good"
		}
	case n > 30:
		{
			return "ok"
		}
	default:
		return "bad"
	}
}

const cardTemplate = `%s{{ if gt .TotalMins 0 }}
			<div id="date-{{.Day}}" class="card-container"%s>
				<div class="card-header">
					<h2 class="card-date">{{.Day}}</h2>
					<h2 class="card-total-hours">{{formatTotalMin .TotalMins}}</h2>
				</div>
				<div class="card-data-container">
					<div class="activities-container">
						{{range .Activities}}
							<div class="duration-bar" style="--default-flex: {{.Duration}};">
								<div class="bar-text">{{.Description}}</div>
							</div>
						{{end}}
					</div>
					<div class="tags-container">
						{{range .Tags}}
						<div id="tag-{{.ID}}" class="tag">{{.Tag}}</div>
						{{end}}
					</div>
				</div>
			</div>
		{{end}}%s`

var singleCardTemplate = fmt.Sprintf(cardTemplate, ``, ` hx-swap-oob="true" `, ``)
var allCardsTemplate = fmt.Sprintf(cardTemplate,
	`<section id="dashboard-section" class="section-container">
		{{range .}}
		`, ``, `
		{{end}}
	</section>`)

var AllCardsTemplateReady = template.Must(
	template.New("allcardtemplate").
		Funcs(template.FuncMap{
			"formatTotalMin": formatTotalMin,
			"calcColor":      calcColor,
		}).
		Parse(allCardsTemplate))

var SingleCardTemplateReady = template.Must(
	template.New("allcardtemplate").
		Funcs(template.FuncMap{
			"formatTotalMin": formatTotalMin,
			"calcColor":      calcColor,
		}).
		Parse(singleCardTemplate))