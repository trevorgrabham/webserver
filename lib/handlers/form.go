package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/trevorgrabham/webserver/webserver/html"
	"github.com/trevorgrabham/webserver/webserver/lib/dashboard"
	"github.com/trevorgrabham/webserver/webserver/lib/database"
	"github.com/trevorgrabham/webserver/webserver/lib/tag"
	"github.com/trevorgrabham/webserver/webserver/lib/timer"
)

func HandleAddTag(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		panic(fmt.Errorf("parsing form: %v", err))
	}

	tagValue := r.Form.Get("temporary-tag")
	if tagValue == "" {
		panic(fmt.Errorf("'temporary-tag' was not a provided paramter"))
	}

	formTag := template.Must(template.New("form-tag").ParseFiles(html.IncludeFiles["form-tag"]...))

	formTagData := html.FormTagData{
		SVG: html.NewElementAttributes([]string{`class="tag-remove-button"`, `hx-get="/removeTag"`, `hx-target="closest .tag-container"`, `hx-swap="outerHTML"`})}
	formTagData.AddTag(tagValue)

	if err := formTag.Execute(w, formTagData); err != nil { panic(err) }
}

func HandleResetTimers(w http.ResponseWriter, _ *http.Request) {
	startHyperscript := `_="on click trigger click on #nav-start-button"`
	pauseHyperscript := `_="on click trigger click on #nav-pause-button"`
	stopHyperscript := `_="on click trigger click on #nav-stop-button"`
	timerData := html.TimerData{
		StartButton: html.NewElementAttributes(append([]string{`id="main-start-button"`, `class="svg-button start-timer-button"`}, startHyperscript)),
		PauseButton: html.NewElementAttributes(append([]string{`id="main-pause-button"`, `class="svg-button hidden pause-timer-button"`}, pauseHyperscript)),
		StopButton: html.NewElementAttributes(append([]string{`id="main-stop-button"`, `class="svg-button hidden stop-timer-button"`}, stopHyperscript)),
	}

	HandleNavTimer(w, nil)
	
	timer := template.Must(template.New("timer").ParseFiles(html.IncludeFiles["timer"]...))
	if err := timer.Execute(w, timerData); err != nil { panic(err) }
}

func HandleActivitySubmit(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil { panic(fmt.Errorf("handleActivitySubmit parsing form: %v", err)) }

	userID, ok := r.Context().Value(ContextKey("user-id")).(int64)
	if !ok { panic(fmt.Errorf("unable to parse 'user-id' from handleactivitysuggestions()")) }

	// check that all required fields are present
	if !r.Form.Has("activity") || !r.Form.Has("timer") || !r.Form.Has("tag") { panic(fmt.Errorf("missing field from %v", r.Form))}
	activity := strings.TrimSpace(r.Form.Get("activity"))
	timerString := strings.TrimSpace(r.Form.Get("timer"))
	tagStrings := r.Form["tag"]

	// format the form data
	var tags tag.Tags
	for _, t := range tagStrings {
		t = strings.TrimSpace(t)
		if !tags.Contains(t) {
			tags = append(tags, tag.TagMetaData{ID: -1, Tag: t, Count: 1})
		}
	}
	hours, mins, err := timer.ParseTimer(timerString)
	if err != nil { panic(err) }
	var totalMins int64 = int64(mins) + 60*int64(hours)

	if DEBUG {
		fmt.Printf("Split up timer: %v:%v\nActivity: %s\n", hours, mins, activity)
	}

	// grab the date when the timer started
	start, err := timer.StartingTime(int(totalMins))
	if err != nil { panic(err) }

	if DEBUG {
		fmt.Printf("Timer started %v min ago\nIt was %v\n", totalMins, start.Format("2006-01-02"))
	}

	// get the timerData ready to insert into DB
	timer := dashboard.ActivityMetaData{
		ID:          -1,
		UserID:      userID,
		Day:         start.Format("2006-01-02"),
		Duration:    totalMins,
		Description: activity,
		Tags:        tags,
	}

	// insert newest activity data into DB
	_, err = database.AddDay(timer)
	if err != nil { panic(err) }

	// grab all of the activities on this day (newly inserted included)
	activities, err := database.GetDayData(timer.UserID, timer.Day)
	if err != nil { panic(err) }

	// generate a CardMetaData for the template
	totalMins = 0
	t := make(tag.Tags, 0)
	for _, a := range activities {
		totalMins += a.Duration
		// de-duplicate any identical tags spread over different sessions on the same day
		for _, tg := range a.Tags {
			if !t.Contains(tg.Tag) {
				t = append(t, tg)
			}
		}
	}
	card := dashboard.CardMetaData{
		Activities: activities,
		Tags:       t,
		TotalMins:  totalMins,
		Day:        timer.Day,
		SwapOOB:    true,
	}

	// generate the template and respond with it
	cardTemplate := template.Must(template.New("card").Funcs(html.DashboardFuncMap).ParseFiles(html.IncludeFiles["card"]...))
	err = cardTemplate.Execute(w, card)
	if err != nil { panic(fmt.Errorf("executing template: %v", err)) }

	HandleIndex(w, r)
}