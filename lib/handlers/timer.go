package handlers

import (
	"fmt"
	"log"
	"net/http"
	"slices"
	"strings"
	"unicode/utf8"

	"github.com/trevorgrabham/webserver/webserver/lib/database"
	"github.com/trevorgrabham/webserver/webserver/lib/util"
)

const DEBUG = true

func HandleRemove(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintf(w, ``)
}

func HandleStartTimer(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintf(w, `
	<button
		class="svg-wrapper"
		_="on click send stopTimer to #timer-display" 
		hx-get="/pauseTimer" 
		hx-target="this"
		hx-swap="outerHTML"
	>
		%s
	</button>

	<button
		class="svg-wrapper"
		_="on click send stopTimer to #timer-display" 
		hx-get="/stopTimer" 
		hx-target="#timer-buttons-container"
		hx-swap="outerHTML"
	>
		%s
	</button>
	`, util.SVG("pause-button", []string{"timer-button","button-main"}, nil, util.Pause), util.SVG("stop-button", []string{"timer-button", "main-button"}, nil, util.Stop))
}

func HandlePauseTimer(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintf(w, `
	<button
		class="svg-wrapper"
    _="on click send startTimer to #timer-display"
    hx-get="/resumeTimer"
    hx-target="this"
    hx-swap="outerHTML"
	>
		%s
	</button>
	`, util.SVG("resume-button", []string{"timer-button", "button-main"}, nil, util.Start))
}

func HandleResumeTimer(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintf(w, `
	<button
		class="svg-wrapper"
		_="on click send stopTimer to #timer-display" 
		hx-get="/pauseTimer" 
		hx-target="this"
		hx-swap="outerHTML"
	>
		%s
	</button>
	`, util.SVG("pause-button", []string{"timer-button", "buton-main"}, nil, util.Pause))
}

func HandleStopTimer(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintf(w, `
  <form
    id="timer-form"
    hx-post="/submitActivity"
    hx-target="#timer-container"
    action=""
		autocomplete="off"
  >
    <input id="hidden-timer" name="timer" type="hidden" />
		<div id="activity-input-container" class="timer-form-input-row">

			<div id="activity-input-wrapper">
				<input type="text" 
					id="activity-input"
					class="timer-form-input"
					hx-get="activitySuggestions"
					hx-trigger="keyup changed delay:500ms"
					hx-target="#activity-suggestions"
					hx-swap="outerHTML"
					name="activity" 
					placeholder="What were you doing?" 
					required
					minlength="2"
					maxlength="255"
				/>
				<div id="activity-suggestions"></div>
			</div>
	
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
			 	%s
			</button>

		</div>

		<div id="timer-form-submit-container" class="timer-form-input-row">
			<button 
	 			id="submit-timer-form"
				class="svg-wrapper"
				type="submit"
				_="on click call resetTimer()"
			>
				%s
			</button>
	
			<button
				class="svg-wrapper"
				_="on click call resetTimer()"
				hx-get="/resetTimer"
				hx-target="#timer-container"
			>
				%s
			</button>
		</div>
	</form>
	`, util.SVG("add-tag-svg", []string{"timer-button", "button-sub-form"}, nil, util.Plus), util.SVG("", []string{"timer-button", "button-form"}, nil, util.Success), util.SVG("reset-button", []string{"timer-button", "button-form"}, nil, util.Cancel))
}

func HandleActivitySuggestions(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Fatalf("HandleActivitySuggestions(): parsing form : %v", err)
	}
	res, ok := r.Form["activity"]
	if !ok || len(res[0]) < 1 {
		fmt.Fprint(w)
		return
	}
	activityPartial := res[0]

	if DEBUG {
		fmt.Printf("Searching for previous activities matching %s\n", activityPartial)
	}

	previousActivities, err := database.GetPreviousActivities()
	if err != nil {
		log.Fatal(err)
	}

	matches := util.FilterFromPartialString(activityPartial, previousActivities)
	if matches.Length == 0 {
		fmt.Fprint(w)
		return
	}

	if DEBUG {
		fmt.Printf("Found matches %v\n", matches)
	}

	if err := util.AutocompleteTemplateReady.Execute(w, matches); err != nil {
		log.Fatalf("HandleActivitySuggestions(): %v", err)
	}
}

func HandleTagSuggestions(w http.ResponseWriter, r *http.Request) {
	
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

	if len(tag) <= 0 {
		fmt.Fprint(w)
		return
	}

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
		 		%s
			</div>
		</div>
	`, tag, utf8.RuneCountInString(tag), util.SVG("", []string{"button-tag-remove"}, []string{`hx-get="/removeTag"`, `hx-target="closest .tag-container`, `hx-swap="outerHTML"`}, util.Remove))
}

func HandleResetTimer(w http.ResponseWriter, _ *http.Request) {
	resetTimer(w)
}

func HandleActivitySubmit(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Fatalf("handleActivitySubmit parsing form: %v", err)
	}

	if DEBUG {
		fmt.Println(r.Form)
	}

	// check that all required fields are present
	var activityDuration, activityDesc string
	var noTags bool
	missingFields, fields, err := util.FormHasFields(r.Form, []string{"timer", "activity", "tag"})
	if err != nil {
		log.Fatal(err)
	}
	// Since tags are optional, remove "tag" from missingFields so we don't log.Fatal() on no tags
	if index := slices.Index(missingFields, "tag"); index != -1 {
		if DEBUG {
			fmt.Println("No tags were present in the form")
		}
		noTags = true
		if len(missingFields) <= 1 {
			missingFields = nil
		} else {
			missingFields = append(missingFields[:index], missingFields[index+1:]...)
		}
	}
	if len(missingFields) > 0 {
		var missingString strings.Builder
		for i, field := range missingFields {
			missingString.WriteString(field)
			if i < len(missingFields)-1 {
				missingString.WriteString(", ")
			}
		}
		log.Fatal("ActivitySubmit: missing fields " + missingString.String())
	}

	// format the form data
	var tags util.Tags
	if !noTags {
		for _, tag := range fields[2] {
			if !tags.Contains(tag) {
				tags = append(tags, util.TagMetaData{Id: -1, Tag: tag, Count: 1})
			}
		}
	}
	activityDuration = fields[0][0]
	activityDesc = fields[1][0]
	hours, mins, err := util.ParseTimer(activityDuration)
	if err != nil {
		log.Fatal(err)
	}
	var totalMins int64 = int64(mins) + 60*int64(hours)

	if DEBUG {
		fmt.Printf("Split up timer: %v:%v\nActivity: %s\n", hours, mins, activityDesc)
	}

	// grab the date when the timer started
	start, err := util.StartingTime(int(totalMins))
	if err != nil {
		log.Fatal(err)
	}

	if DEBUG {
		fmt.Printf("Timer started %v min ago\nIt was %v\n", totalMins, start.Format("2006-01-02"))
	}

	// get the timerData ready to insert into DB
	timer := util.ActivityMetaData{
		Id:						-1,
		Day:					start.Format("2006-01-02"),
		Duration: 		totalMins,
		Description:	activityDesc,
		Tags:					tags,
	}

	// insert newest activity data into DB
	_, err = database.AddDay(timer)
	if err != nil {
		log.Fatal(err)
	}
	
	// grab all of the activities on this day (newly inserted included)
	activities, err := database.GetDayData(timer.Day)
	if err != nil {
		log.Fatal(err)
	}

	// generate a CardMetaData for the template
	totalMins = 0
	t := make(util.Tags, 0)
	for _, a := range activities {
		totalMins += a.Duration
		// de-duplicate any identical tags spread over different sessions on the same day
		for _, tg := range a.Tags {
			if !t.Contains(tg.Tag) {
				t = append(t, tg)
			}
		}
	}
	card := util.CardMetaData{
		Activities: activities,
		Tags:	t,
		TotalMins: totalMins,
		Day: timer.Day,
	}

	// generate the template and respond with it
	if err := util.SingleCardTemplateReady.Execute(w, card); err != nil {
		log.Fatalf("executing template: %v", err)
	}

	resetTimer(w)
}

func resetTimer(w http.ResponseWriter) {
	fmt.Fprintf(w, `
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
				%s
      </button>
    </div>
	`, util.SVG("start-button", []string{"timer-button", "button-main"}, nil, util.Start))
}
