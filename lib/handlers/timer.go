package handlers

import (
	"fmt"
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


	timer := timerData{
		day:					start.Format("2006-01-02"),
		duration: 		mins,
		description:	desc,
	}

	// insert into timer_data db
	db, err := database.ConnectDb()
	if err != nil {
		log.Fatal(err)
	}
	result, err := db.Exec("INSERT INTO timer_data (day, duration, description) VALUES (?, ?, ?)", timer.day, timer.duration, timer.description)
	if err != nil {
		log.Fatalf("Inserting into timer_data (%s, %d, %s): %v", timer.day, timer.duration, timer.description, err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		log.Fatalf("Retrieving inserted id: %v", err)
	}

	// parse tags and insert into activity_tag db
	if tags, ok := r.Form["tag"]; ok && len(tags) > 0 {
		for _, tag := range tags {

			fmt.Printf("Adding tag: %v\n", tag)

			tag = strings.TrimSpace(tag)
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
      0:00:00
    </div>

    <form>

      <input id="hidden-timer" 
				name="timer" 
				type="hidden" 
				value="0:00:00" 
			/>

      <button
        _="on click send startTimer to #timer"
        hx-get="/startTimer"
        hx-swap="outerHTML"
      >
        Start
      </button>
    </form>
	`)
}