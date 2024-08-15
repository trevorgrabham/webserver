package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/trevorgrabham/webserver/webserver/lib/database"
)

const cardTemplate = `{{range .}}
	{{ with $card := . }}
		<div id="date-{{.Day}}" class="card-container">
			<div class="card-header">
				<h2 class="card-date">{{.Day}}</h2>
				<h2 class="card-total-hours">{{printf "%.1fh" .TotalHours}}</h2>
			</div>
			<div class="card-data-container">
				<div class="activities-container">
					{{range .Activities}}
						<div class="duration-bar" title="{{.Description}}" data-percentage="{{fractionalTime .Duration $card.TotalHours}}"></div>
					{{end}}
				</div>
				<div class="tags-container">
					{{range .Tags}}
					<span class="tag">{{.}}</span>
					{{end}}
				</div>
			</div>
		</div>
	{{end}}
{{end}}`

var CardTemplate = template.Must(template.New("cardtemplate").Funcs(template.FuncMap{"fractionalTime": func(n int64, total float64) float64 { return float64(n)/60/total*100}}).Parse(cardTemplate))

type ActivityData struct {
	Duration		int64
	Description	string
}

type CardData struct {
	Day						string
	Activities		[]ActivityData	
	Tags					Tags					
	TotalHours		float64
}

type Tags []string

func (t Tags) contains(s string) bool {
	for _, tag := range t {
		if tag == s {
			return true
		}
	}
	return false
}

func HandleDashboard(w http.ResponseWriter, r *http.Request) {
	// Parse out maxItems parameter
	if err := r.ParseForm(); err != nil {
		log.Fatalf("Parsing form: %v", err)
	}
	
	res, ok := r.Form["maxItems"]
	if !ok {
		log.Fatal("'maxItems' was not a provided parameter")
	}
	maxItems, err := strconv.ParseInt(res[0], 10, 64)
	if err != nil {
		log.Fatalf("Parsing maxItems (%s): %v", res[0], err)
	}

	// connect to the db
	db, err := database.ConnectDb()
	if err != nil {
		log.Fatal(err)
	}

	// get maxItem (or all) days
	var dayRows *sql.Rows
	if maxItems > 0 {
		dayRows, err = db.Query("SELECT DISTINCT day FROM timer_data ORDER BY day DESC LIMIT ?", maxItems)
		if err != nil {
			log.Fatalf("Querying %d days from timer_data: %v", maxItems, err)
		}
	} else {
		dayRows, err = db.Query("SELECT DISTINCT day FROM timer_data ORDER BY day DESC")
		if err != nil {
			log.Fatalf("Querying from timer_data: %v", err)
		}
	}
	defer dayRows.Close()

	cards := make([]CardData, 0, maxItems)

	// for each day, get all sessions
	for dayRows.Next() {
		card := CardData{}
		var day string
		if err := dayRows.Scan(&day); err != nil {
			log.Fatalf("Reading from dayRow: %v", err)
		}

		// get all session ids, durations, and descriptions for each day
		dataRows, err := db.Query("SELECT id, duration, description FROM timer_data WHERE day LIKE ?", day)
		if err != nil {
			log.Fatalf("Querying duration, description WHERE day LIKE %s: %v", day, err)
		}
		for dataRows.Next() {
			var (
				desc string
				dur, id int64
			)
			err := dataRows.Scan(&id, &dur, &desc)
			if err != nil {
				dataRows.Close()
				log.Fatalf("Reading from dataRow: %v", err)
			}

			// update card data
			card.Day = day
			card.Activities = append(card.Activities, ActivityData{
				Duration: dur,
				Description: desc,
			})
			card.TotalHours += float64(dur)/60

			// get tags for each session for each day
			tagRows, err := db.Query("SELECT tag FROM activity_tag WHERE activity_id = ?", id)
			if err != nil {
				log.Fatalf("Unable to query activity_tag with id %d: %v", id, err)
			}
			for tagRows.Next() {
				var tag string
				err := tagRows.Scan(&tag)
				if err != nil {
					log.Fatalf("Scanning tagRows: %v", err)
				}

				// update card data
				if !card.Tags.contains(tag) {
					card.Tags = append(card.Tags, tag)
				}
			}
			if err := tagRows.Err(); err != nil {
				tagRows.Close()
				log.Fatalf("Querying activity_tag with id %d: %v", id, err)
			}
			tagRows.Close()
		}
		if err := dataRows.Err(); err != nil {
			dataRows.Close()
			log.Fatalf("Querying duration, description WHERE day LIKE %s: %v", day, err)
		}
		dataRows.Close()
		cards = append(cards, card)
	}
	if err := dayRows.Err(); err != nil {
			log.Fatalf("Query distinct day from timer_data: %v", err)
	}

	if err := CardTemplate.Execute(w, cards); err != nil {
		log.Fatalf("Executing template: %v", err)
	}
}	