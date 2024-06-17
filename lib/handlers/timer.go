package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/trevorgrabham/webserver/webserver/lib/database"
)

type timerData struct {
	day					string
	duration		int64
	description	string
}

func HandleRemove(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintf(w, ``)
}

func HandleStartTimer(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(w, `
	<div id="timer-buttons-container">

		<button id="pause-timer" 
			_="on click send stopTimer to #timer" 
			hx-get="/pauseTimer" 
			hx-swap="outerHTML"
		>
			Pause
		</button>

		<button id="stop-timer" 
			_="on click send stopTimer to #timer" 
			hx-get="/stopTimer" 
			hx-target="#timer-buttons-container"
		>
			Stop
		</button>

	</div>
	`)
}

func HandlePauseTimer(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(w, `
	<button id="resume-timer" 
		_="on click send startTimer to #timer" 
		hx-get="/resumeTimer" 
		hx-swap="outerHTML"
	>
		Resume
	</button>
	`)
}

func HandleResumeTimer(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(w, `
	<button id="pause-timer" 
		_="on click send stopTimer to #timer" 
		hx-get="/pauseTimer" 
		hx-swap="outerHTML"
	>
		Pause
	</button>
	`)
}

func HandleStopTimer(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(w, `
	<div class="form-input-row">

		<input type="text" 
			name="activity" 
			placeholder="What were you doing?" 
			required
			minlength="2"
			maxlength="255"
		/>

	</div>

	<div id="tags-container" class="form-input-row">

		<p id="add-tag-button" 
			hx-get="/tagInput" 
			hx-target="#tags-container" 
			hx-swap="beforebegin"
		>
			Add tag
		</p>

	</div>

	<button 
	 	id="submit-timer-form"
		type="submit"
		_="on click call resetTimer()"
	>
		Submit
	</button>
	`)
}

func HandleTagInput(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintf(w, `
	<span class="tag-container">

		<input type="text"
			name="tag"
		/>

		<span hx-get="/remove" 
			hx-target="closest .tag-container"
			hx-swap="outerHTML"
		>
			X
		</span>
	
	</span>
	`)
}

func HandleActivitySubmit(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Fatalf("handleActivitySubmit parsing form: %v", err)
	}

	fmt.Println(r.Form)

	// check that all required fields are present
	var dur, desc string
	res, ok := r.Form["timer"]
	if !ok {
		log.Fatal("ActivitySubmit: no timer data submitted")
	}
	dur = res[0]
	res, ok = r.Form["activity"]
	if !ok {
		log.Fatal("ActivitySubmit: no description data submitted")
	}
	desc = res[0]
	durSplit := strings.Split(dur, ":")

	fmt.Printf("split up duration: %v\n", durSplit)

	// parse the timer into a duration
	var mins int64
	if len(durSplit) == 3 {
		hours, err := strconv.ParseInt(durSplit[0], 10, 64)
		if err != nil {
			log.Fatalf("Parse hours(dur): %v", err)
		}
		mins, err = strconv.ParseInt(durSplit[1], 10, 64)
		if err != nil {
			log.Fatalf("Parse mins(dur): %v", err)
		}

		fmt.Printf("%d hours\n", hours)

		mins += hours * 60 
	} else {
		var err error
		mins, err = strconv.ParseInt(durSplit[0], 10, 64)
		if err != nil {
			log.Fatalf("Parse mins(dur): %v", err)
		}
	}

	fmt.Printf("%d min total\n", mins)


	// grab the date when the timer started
	now := time.Now()
	timerDuration, err := time.ParseDuration(fmt.Sprintf("-%dm",mins))
	if err != nil {
		log.Fatalf("Creating timer duration: %v", err)
	}
	start := now.Add(timerDuration)

	fmt.Printf("Timer started %v min ago\nIt was %v\n", timerDuration, start.Format("2006-01-02"))

	// grab the tags 
	var tags []string
	data, ok := r.Form["tag"]
	if ok {
		for _, tag := range data {
			tag = strings.TrimSpace(tag)
			if len(tag) > 0 {
				tags = append(tags, tag)
			}
		}
	}

	timer := timerData{
		day:					start.Format("2006-01-02"),
		duration: 		mins,
		description:	desc,
	}

	// generate the dashboard card data
	
	// before insert, check to see if there was already a session today
	db, err := database.ConnectDb()
	if err != nil {
		log.Fatal(err)
	}
	dataRows, err := db.Query("SELECT timer_data.id AS id, duration, description, tag FROM timer_data JOIN activity_tag ON timer_data.id = activity_tag.activity_id WHERE day LIKE ?", timer.day)
	if err != nil {
		log.Fatalf("checking if entry for today already exists in timer_data: %v", err)
	}

	// add all pre-existing data for this date
	var (
		prevId, currId, duration int64
		description, t string
	)
	card := CardData{Day: timer.day}
	for dataRows.Next() {
		err := dataRows.Scan(&currId, &duration, &description, &t)
		if err != nil {
			log.Fatalf("reading from a row for the card data: %v", err)
		}
		card.Tags = append(card.Tags, t)
		if prevId != currId {
			card.Activities = append(card.Activities, ActivityData{Duration: duration, Description: description})
		}
		prevId = currId
	}
	if dataRows.Err() != nil {
		log.Fatalf("error querying if today exists: %v", err)
	}

	// add the new data for today
	card.Tags = append(card.Tags, tags...)
	card.Activities = append(card.Activities, ActivityData{Duration: timer.duration, Description: timer.description})

	// insert into timer_data db
	result, err := db.Exec("INSERT INTO timer_data (day, duration, description) VALUES (?, ?, ?)", timer.day, timer.duration, timer.description)
	if err != nil {
		log.Fatalf("Inserting into timer_data (%s, %d, %s): %v", timer.day, timer.duration, timer.description, err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		log.Fatalf("Retrieving inserted id: %v", err)
	}

	// insert into activity_tag db
	for _, tag := range tags {

		fmt.Printf("Adding tag: %v\n", tag)

		if len(tag) > 0 {
			result, err := db.Exec("INSERT INTO activity_tag (tag, activity_id) VALUES (?, ?)", tag, id)
			if err != nil {
				log.Fatalf("Inserting into activity_tag (%s, %d): %v", tag, id, err)
			}
			_, err = result.LastInsertId()
			if err != nil {
				log.Fatalf("Retrieving inserted id: %v", err)
			}
		}
	}

	templateString := `<div id="date-{{.Day}}" class="card-container" hx-swap-oob="true">
		<h2 class="card-date">{{.Day}}</h2>
		<div class="activities-container">
			{{range .Activities}}
			<span class="duration-bar" data-dur="{{.Duration}}" data-desc="{{.Description}}"></span>
			{{end}}
		</div>
		<div class="tags-container">
			{{range .Tags}}
			<span class="tag">{{.}}</span>
			{{end}}
		</div>
	</div>`

	cardTemplate := template.Must(template.New("cardtemplate").Parse(templateString))
	if err := cardTemplate.Execute(w, card); err != nil {
		log.Fatalf("executing template: %v", err)
	}

	fmt.Fprint(w, `
	<div
    id="timer"
    _="on startTimer call startTimer() then 
      repeat until event stopTimer
        call updateTimer()
        wait 1s
      end"
  >
    00:00
  </div>
  <form
    id="timer-form"
    hx-post="/submitActivity"
    hx-target="#timer-container"
    action=""
  >
    <input id="hidden-timer" name="timer" type="hidden" value="0:00:00" />
    <button
      _="on click send startTimer to #timer"
      hx-get="/startTimer"
      hx-target="this"
      hx-swap="outerHTML"
    >
      Start
    </button>
  </form>`)
	
}