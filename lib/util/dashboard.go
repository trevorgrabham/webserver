package util

import (
	"fmt"
	"text/template"
)

type Tags []TagMetaData

// Does not do any pre-processing, so is case-sensitive and does not trim white-space
func (t Tags) Contains(s string) bool {
	for _, tag := range t {
		if tag.Tag == s {
			return true
		}
	}
	return false
}

type TagMetaData struct {
	Id  int64
	Tag string
}

type ActivityMetaData struct {
	Id          int64
	Duration    int64
	Description string
	Day         string
	Tags
}

type CardMetaData struct {
	Activities 			[]ActivityMetaData
	Tags						
	TotalMins				int64 
	Day							string
}

func formatTotalMin(t int64) string {
	if t > 60 { return fmt.Sprintf("%dh%dm", t/60, t%60) }
	return fmt.Sprintf("%dm", t%60)
}

func calcColor(n int64) string {
	switch {
	case n > 90: { return "great" }
	case n > 60: { return "good" }
	case n > 30: { return "ok" }
	default: return "bad"
	}
}

const CardTemplate = `%s{{ if gt .TotalMins 0 }}
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
						<div id="tag-{{.Id}}" class="tag">{{.Tag}}</div>
						{{end}}
					</div>
				</div>
			</div>
		{{end}}%s`
var SingleCardTemplate = fmt.Sprintf(CardTemplate, ``, ` hx-swap-oob="true" `, ``)
var AllCardsTemplate = fmt.Sprintf(CardTemplate, `{{range .}}
	`, ``, `
	{{end}}`)

var AllCardsTemplateReady = template.Must(
	template.New("allcardtemplate").
	Funcs(template.FuncMap{
		"formatTotalMin": formatTotalMin,
		"calcColor": calcColor,
		}).
	Parse(AllCardsTemplate))

var SingleCardTemplateReady= template.Must(
	template.New("allcardtemplate").
	Funcs(template.FuncMap{
		"formatTotalMin": formatTotalMin,
		"calcColor": calcColor,
		}).
	Parse(SingleCardTemplate))
