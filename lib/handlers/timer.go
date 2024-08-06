package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

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
	<button
		class="svg-wrapper"
		_="on click send stopTimer to #timer-display" 
		hx-get="/pauseTimer" 
		hx-target="this"
		hx-swap="outerHTML"
	>
		<svg 
			id="pause-button" 
			class="timer-button button-main"
			fill="#000000" 
			height="71px" 
			width="71px" 
			version="1.1" 
			xmlns="http://www.w3.org/2000/svg" 
			xmlns:xlink="http://www.w3.org/1999/xlink" 
			viewBox="0 0 60 60" 
			xml:space="preserve"
		>
	 		<!-- svg source https://www.svgrepo.com/svg/83992/pause -->
			<g id="SVGRepo_bgCarrier" stroke-width="0"></g>
			<g id="SVGRepo_tracerCarrier" stroke-linecap="round" stroke-linejoin="round"></g>
			<g id="SVGRepo_iconCarrier"> 
				<g> 
					<path d="M30,0C13.458,0,0,13.458,0,30s13.458,30,30,30s30-13.458,30-30S46.542,0,30,0z M30,58C14.561,58,2,45.439,2,30 S14.561,2,30,2s28,12.561,28,28S45.439,58,30,58z"></path> 
					<path d="M33,46h8V14h-8V46z M35,16h4v28h-4V16z"></path> 
					<path d="M19,46h8V14h-8V46z M21,16h4v28h-4V16z"></path> 
				</g> 
			</g>
		</svg>	
	</button>

	<button
		class="svg-wrapper"
		_="on click send stopTimer to #timer-display" 
		hx-get="/stopTimer" 
		hx-target="#timer-buttons-container"
		hx-swap="outerHTML"
	>
		<svg 
			id="stop-button" 
			class="timer-button button-main"
			fill="#000000" 
			height="71px" 
			width="71px" 
			version="1.1" 
			xmlns="http://www.w3.org/2000/svg" 
			xmlns:xlink="http://www.w3.org/1999/xlink" 
			viewBox="0 0 60 60" 
			xml:space="preserve"
		>
			<!-- svg from https://www.svgrepo.com/svg/125999/stop -->
			<g id="SVGRepo_bgCarrier" stroke-width="0"></g>
			<g id="SVGRepo_tracerCarrier" stroke-linecap="round" stroke-linejoin="round"></g>
			<g id="SVGRepo_iconCarrier"> 
				<g> 
					<path d="M16,44h28V16H16V44z M18,18h24v24H18V18z"></path> 
					<path d="M30,0C13.458,0,0,13.458,0,30s13.458,30,30,30s30-13.458,30-30S46.542,0,30,0z M30,58C14.561,58,2,45.439,2,30 S14.561,2,30,2s28,12.561,28,28S45.439,58,30,58z"></path> 
				</g> 
			</g>
		</svg>
	</button>
	`)
}

func HandlePauseTimer(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(w, `
	<button
		class="svg-wrapper"
    _="on click send startTimer to #timer-display"
    hx-get="/resumeTimer"
    hx-target="this"
    hx-swap="outerHTML"
	>
  	<svg
    	id="resume-button"
    	class="timer-button button-main"
    	fill="#000000"
    	height="71px"
    	width="71px"
    	version="1.1"
    	xmlns="http://www.w3.org/2000/svg"
    	xmlns:xlink="http://www.w3.org/1999/xlink"
    	viewBox="0 0 60 60"
    	xml:space="preserve"
  	>
			<!-- svg source https://www.svgrepo.com/svg/13672/play-button -->
    	<g>
      	<path
        	d="M45.563,29.174l-22-15c-0.307-0.208-0.703-0.231-1.031-0.058C22.205,14.289,22,14.629,22,15v30c0,0.371,0.205,0.711,0.533,0.884C22.679,45.962,22.84,46,23,46c0.197,0,0.394-0.059,0.563-0.174l22-15C45.836,30.64,46,30.331,46,30S45.836,29.36,45.563,29.174z M24,43.107V16.893L43.225,30L24,43.107z"
      	/>
      	<path
        	d="M30,0C13.458,0,0,13.458,0,30s13.458,30,30,30s30-13.458,30-30S46.542,0,30,0z M30,58C14.561,58,2,45.439,2,30S14.561,2,30,2s28,12.561,28,28S45.439,58,30,58z"
      	/>
    	</g>
  	</svg>
	</button>
	`)
}

func HandleResumeTimer(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(w, `
	<button
		class="svg-wrapper"
		_="on click send stopTimer to #timer-display" 
		hx-get="/pauseTimer" 
		hx-target="this"
		hx-swap="outerHTML"
	>
		<svg 
			id="pause-button" 
			class="timer-button button-main"
			fill="#000000" 
			height="71px" 
			width="71px" 
			version="1.1" 
			xmlns="http://www.w3.org/2000/svg" 
			xmlns:xlink="http://www.w3.org/1999/xlink" 
			viewBox="0 0 60 60" 
			xml:space="preserve"
		>
			<!-- svg source https://www.svgrepo.com/svg/83992/pause -->
			<g id="SVGRepo_bgCarrier" stroke-width="0"></g>
			<g id="SVGRepo_tracerCarrier" stroke-linecap="round" stroke-linejoin="round"></g>
			<g id="SVGRepo_iconCarrier"> 
				<g> 
					<path d="M30,0C13.458,0,0,13.458,0,30s13.458,30,30,30s30-13.458,30-30S46.542,0,30,0z M30,58C14.561,58,2,45.439,2,30 S14.561,2,30,2s28,12.561,28,28S45.439,58,30,58z"></path> 
					<path d="M33,46h8V14h-8V46z M35,16h4v28h-4V16z"></path> 
					<path d="M19,46h8V14h-8V46z M21,16h4v28h-4V16z"></path> 
				</g> 
			</g>
		</svg>	
	</button>
	`)
}

func HandleStopTimer(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(w, `
  <form
    id="timer-form"
    hx-post="/submitActivity"
    hx-target="#timer-container"
    action=""
  >
    <input id="hidden-timer" name="timer" type="hidden" />
		<div id="activity-input-container" class="timer-form-input-row">

			<input type="text" 
				id="activity-input"
				class="timer-form-input"
				name="activity" 
				placeholder="What were you doing?" 
				required
				minlength="2"
				maxlength="255"
			/>
	
			<div id="tags-wrapper">
				<div id="tags-container"></div>
			</div>

		</div>


		<div 
			id="tag-input-container" 
			class="timer-form-input-row"
		>

			<input type="text"
				id="tag-input" 
				class="timer-form-input"
				name="temporary-tag" 
				placeholder="tags" 
				maxlength="255"
			/>

			<button 
			 	id="add-tag-button"
				class="svg-wrapper"
				hx-post="/addTag"
				hx-include="#tag-input"
				hx-params="temporary-tag"
				hx-target="#tags-container"
				hx-swap="beforeend"
			>
				<svg 
					id="add-tag-svg"
					class="timer-button button-sub-form"
					fill="#000000" 
					height="71px" 
					width="71px" 
					version="1.1" 
					xmlns="http://www.w3.org/2000/svg" 
					xmlns:xlink="http://www.w3.org/1999/xlink" 
					viewBox="0 0 52 52" 
					xml:space="preserve"
				>
			 		<!-- svg source https://www.svgrepo.com/svg/158372/plus -->
					<g id="SVGRepo_bgCarrier" stroke-width="0"></g>
					<g id="SVGRepo_tracerCarrier" stroke-linecap="round" stroke-linejoin="round"></g>
					<g id="SVGRepo_iconCarrier"> 
						<g> 
							<path d="M26,0C11.664,0,0,11.663,0,26s11.664,26,26,26s26-11.663,26-26S40.336,0,26,0z M26,50C12.767,50,2,39.233,2,26 S12.767,2,26,2s24,10.767,24,24S39.233,50,26,50z"></path> 
							<path d="M38.5,25H27V14c0-0.553-0.448-1-1-1s-1,0.447-1,1v11H13.5c-0.552,0-1,0.447-1,1s0.448,1,1,1H25v12c0,0.553,0.448,1,1,1 s1-0.447,1-1V27h11.5c0.552,0,1-0.447,1-1S39.052,25,38.5,25z"></path> 
						</g> 
					</g>
				</svg>
			</button>

		</div>

		<div id="timer-form-submit-container" class="timer-form-input-row">
			<button 
	 			id="submit-timer-form"
				class="svg-wrapper"
				type="submit"
				_="on click call resetTimer()"
			>
				<svg 
			 		class="timer-button button-form"
					fill="#000000" 
					height="71px" 
					width="71px" 
					version="1.1" 
					xmlns="http://www.w3.org/2000/svg" 
					xmlns:xlink="http://www.w3.org/1999/xlink" 
					viewBox="0 0 52 52" 
					xml:space="preserve"
				>
					<g id="SVGRepo_bgCarrier" stroke-width="0"></g>
					<g id="SVGRepo_tracerCarrier" stroke-linecap="round" stroke-linejoin="round"></g>
					<g id="SVGRepo_iconCarrier"> 
						<g> 
							<path d="M26,0C11.664,0,0,11.663,0,26s11.664,26,26,26s26-11.663,26-26S40.336,0,26,0z M26,50C12.767,50,2,39.233,2,26 S12.767,2,26,2s24,10.767,24,24S39.233,50,26,50z"></path> 
							<path d="M38.252,15.336l-15.369,17.29l-9.259-7.407c-0.43-0.345-1.061-0.274-1.405,0.156c-0.345,0.432-0.275,1.061,0.156,1.406 l10,8C22.559,34.928,22.78,35,23,35c0.276,0,0.551-0.114,0.748-0.336l16-18c0.367-0.412,0.33-1.045-0.083-1.411 C39.251,14.885,38.62,14.922,38.252,15.336z"></path> 
						</g> 
					</g>
				</svg>	
			</button>
	
			<button
				class="svg-wrapper"
				_="on click call resetTimer()"
				hx-get="/resetTimer"
				hx-target="#timer-container"
			>
				<svg 
					id="reset-button"
			 		class="timer-button button-form"
					fill="#000000" 
					height="71px" 
					width="71px" 
					version="1.1" 
					xmlns="http://www.w3.org/2000/svg" 
					xmlns:xlink="http://www.w3.org/1999/xlink" 
					viewBox="0 0 52 52" 
					xml:space="preserve"
				>
					<!-- svg source https://www.svgrepo.com/svg/138890/error -->
					<g id="SVGRepo_bgCarrier" stroke-width="0"></g>
					<g id="SVGRepo_tracerCarrier" stroke-linecap="round" stroke-linejoin="round"></g>
					<g id="SVGRepo_iconCarrier"> 
						<g> 
							<path d="M26,0C11.664,0,0,11.663,0,26s11.664,26,26,26s26-11.663,26-26S40.336,0,26,0z M26,50C12.767,50,2,39.233,2,26 S12.767,2,26,2s24,10.767,24,24S39.233,50,26,50z"></path> 
							<path d="M35.707,16.293c-0.391-0.391-1.023-0.391-1.414,0L26,24.586l-8.293-8.293c-0.391-0.391-1.023-0.391-1.414,0 s-0.391,1.023,0,1.414L24.586,26l-8.293,8.293c-0.391,0.391-0.391,1.023,0,1.414C16.488,35.902,16.744,36,17,36 s0.512-0.098,0.707-0.293L26,27.414l8.293,8.293C34.488,35.902,34.744,36,35,36s0.512-0.098,0.707-0.293 c0.391-0.391,0.391-1.023,0-1.414L27.414,26l8.293-8.293C36.098,17.316,36.098,16.684,35.707,16.293z"></path> 
						</g> 
					</g>
				</svg>
			</button>
		</div>
	</form>
	`)
}

func HandleAddTag(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Fatalf("Parsing form: %v", err)
	}

	res, ok := r.Form["temporary-tag"]
	if !ok || len(res) < 1 {
		log.Fatal("'temporary-tag' was not a provided paramter")
	}

	tag := res[0]

	fmt.Fprintf(w, `
		<input type="text"
		 	hx-swap-oob="true"
			id="tag-input" 
			class="timer-form-input"
			name="temporary-tag" 
			placeholder="tags" 
			maxlength="255"
		/>

	 	<div class="tag-container">
			<div class="tag-wrapper">
				<input type="text" class="tag-display" name="tag" value="%s" readonly style="width: %dch"/>
				<svg 
					class="button-tag-remove"
			 		hx-get="/removeTag"
					hx-target="closest .tag-container"
					hx-swap="outerHTML"
					fill="#000000" 
					height="71px" 
					width="71px" 
					version="1.1" 
					xmlns="http://www.w3.org/2000/svg" 
					xmlns:xlink="http://www.w3.org/1999/xlink" 
					viewBox="0 0 31.112 31.112" 
					xml:space="preserve"
				>
					<g id="SVGRepo_bgCarrier" stroke-width="0"></g>
					<g id="SVGRepo_tracerCarrier" stroke-linecap="round" stroke-linejoin="round"></g>
					<g id="SVGRepo_iconCarrier"> 
						<polygon points="31.112,1.414 29.698,0 15.556,14.142 1.414,0 0,1.414 14.142,15.556 0,29.698 1.414,31.112 15.556,16.97 29.698,31.112 31.112,29.698 16.97,15.556 "></polygon> 
					</g>
				</svg>
			</div>
		</div>
	`, tag, utf8.RuneCountInString(tag))
}

func HandleResetTimer(w http.ResponseWriter, _ *http.Request) {
	resetTimer(w)
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

	resetTimer(w)
	
}

func resetTimer(w http.ResponseWriter) {
	fmt.Fprint(w, `
    <p
      id="timer-display"
      _="on startTimer call startTimer() then 
          repeat until event stopTimer
            call updateTimer()
            wait 1s
          end"
      style="font-size: 5em"
    >
      00:00
    </p>
    <div id="timer-buttons-container">
      <button
        class="svg-wrapper"
        _="on click send startTimer to #timer-display"
        hx-get="/startTimer"
        hx-target="#timer-buttons-container"
      >
        <svg
          id="start-button"
          class="timer-button button-main"
          fill="#000000"
          height="71px"
          width="71px"
          version="1.1"
          xmlns="http://www.w3.org/2000/svg"
          xmlns:xlink="http://www.w3.org/1999/xlink"
          viewBox="0 0 60 60"
          xml:space="preserve"
        >
          <!-- svg source https://www.svgrepo.com/svg/13672/play-button -->
          <g>
            <path
              d="M45.563,29.174l-22-15c-0.307-0.208-0.703-0.231-1.031-0.058C22.205,14.289,22,14.629,22,15v30
		  c0,0.371,0.205,0.711,0.533,0.884C22.679,45.962,22.84,46,23,46c0.197,0,0.394-0.059,0.563-0.174l22-15
		  C45.836,30.64,46,30.331,46,30S45.836,29.36,45.563,29.174z M24,43.107V16.893L43.225,30L24,43.107z"
            />
            <path
              d="M30,0C13.458,0,0,13.458,0,30s13.458,30,30,30s30-13.458,30-30S46.542,0,30,0z M30,58C14.561,58,2,45.439,2,30
		  S14.561,2,30,2s28,12.561,28,28S45.439,58,30,58z"
            />
          </g>
        </svg>
      </button>
    </div>
	`) 
}