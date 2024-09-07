package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"slices"
	"strings"
	"unicode/utf8"

	"github.com/trevorgrabham/webserver/webserver/lib/dashboard"
	"github.com/trevorgrabham/webserver/webserver/lib/database"
	tagpkg "github.com/trevorgrabham/webserver/webserver/lib/tag"
	"github.com/trevorgrabham/webserver/webserver/lib/templateutil"
	"github.com/trevorgrabham/webserver/webserver/lib/timer"
)

const DEBUG = true

func HandleRemove(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(w, ``)
}

func HandleStartTimer(w http.ResponseWriter, _ *http.Request) {
	pauseButton := template.Must(template.New("pausebutton").ParseFiles(templateutil.ParseFiles["pausebutton"]...))
	stopButton := template.Must(template.New("stopbutton").ParseFiles(templateutil.ParseFiles["stopbutton"]...))

	err := pauseButton.Execute(w, templateutil.NewElementInfo("pause-button", []string{"timer-button", "button-main"}, nil, ""))
	if err != nil { panic(err) }

	err = stopButton.Execute(w, templateutil.NewElementInfo("stop-button", []string{"timer-button", "button-main"}, nil, ""))
	if err != nil { panic(err) }
}

func HandlePauseTimer(w http.ResponseWriter, _ *http.Request) {
	startButton := template.Must(template.New("startbutton").ParseFiles(templateutil.ParseFiles["startbutton"]...))
	err := startButton.Execute(w, templateutil.NewElementInfo("resume-button", []string{"timer-button", "button-main"}, nil, ""))
	if err != nil { panic(err) }
}

func HandleResumeTimer(w http.ResponseWriter, _ *http.Request) {
	pausebutton := template.Must(template.New("pausebutton").ParseFiles(templateutil.ParseFiles["pausebutton"]...))
	err := pausebutton.Execute(w, templateutil.NewElementInfo("pause-button", []string{"timer-button", "button-main"}, nil, ""))
	if err != nil { panic(err) }
}

func HandleStopTimer(w http.ResponseWriter, _ *http.Request) {
	form := template.Must(template.New("form").ParseFiles(templateutil.ParseFiles["form"]...))
	data := templateutil.NewFormData(
		templateutil.NewElementInfo("add-tag-svg", []string{"timer-button", "button-sub-form"}, nil, ""),
		templateutil.NewElementInfo("", []string{"timer-button", "button-form"}, nil, ""),
		templateutil.NewElementInfo("reset-button", []string{"timer-button", "button-form"}, nil, ""))
	err := form.Execute(w, data)
	if err != nil { panic(err) }
}

func HandleActivitySuggestions(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil { panic(fmt.Errorf("HandleActivitySuggestions(): parsing form : %v", err)) }

	userID, ok := r.Context().Value(ContextKey("user-id")).(int64)
	if ! ok { panic(fmt.Errorf("unable to parse 'user-id' from handleactivitysuggestions()")) }

	res, ok := r.Form["activity"]
	if !ok || len(res[0]) < 1 { fmt.Fprint(w); return }
	activityPartial := res[0]

	if DEBUG {
		fmt.Printf("Searching for previous activities matching %s\n", activityPartial)
	}

	previousActivities, err := database.GetPreviousActivities(userID)
	if err != nil { panic(err) }

	matches := timer.FilterFromPartialString(activityPartial, previousActivities)
	if matches.Length == 0 { fmt.Fprint(w); return }

	if DEBUG {
		fmt.Printf("Found matches %v\n", matches)
	}

	autocomplete := template.Must(template.New("autocomplete").Funcs(templateutil.AutocompleteFuncMap()).ParseFiles(templateutil.ParseFiles["autocomplete"]...))
	err = autocomplete.Execute(w, templateutil.NewAutocompleteData("activity-suggestions", matches))
	if err != nil { panic(fmt.Errorf("HandleActivitySuggestions(): %v", err)) }
}

func HandleTagSuggestions(w http.ResponseWriter, r *http.Request) {
	
}

func HandleAddTag(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil { panic(fmt.Errorf("parsing form: %v", err)) }

	res, ok := r.Form["temporary-tag"]
	if !ok || len(res) < 1 { panic(fmt.Errorf("'temporary-tag' was not a provided paramter")) }

	tag := res[0]

	if len(tag) <= 0 { fmt.Fprint(w); return }

	newTag := template.Must(template.New("newtag").ParseFiles(templateutil.ParseFiles["newtag"]...))
	err := newTag.Execute(w, templateutil.NewTagData(
		templateutil.NewElementInfo("tag-input", []string{"timer-form-input"}, []string{`hx-swap-oob="true"`}, ""),
		templateutil.NewElementInfo("", []string{"button-tag-remove"}, []string{`hx-get="/removeTag"`, `hx-target="closest .tag-container"`, `hx-swap="outerHTML"`}, ""),
		tag,
		utf8.RuneCountInString(tag),
	))
	if err != nil { panic(err) }
}

func HandleResetTimer(w http.ResponseWriter, _ *http.Request) {
	resetTimer(w)
}

func HandleActivitySubmit(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil { panic(fmt.Errorf("handleActivitySubmit parsing form: %v", err)) }

	if DEBUG {
		fmt.Println(r.Form)
	}

	userID, ok := r.Context().Value(ContextKey("user-id")).(int64)
	if ! ok { panic(fmt.Errorf("unable to parse 'user-id' from handleactivitysuggestions()")) }

	// check that all required fields are present
	var activityDuration, activityDesc string
	var noTags bool
	missingFields, fields, err := timer.FormHasFields(r.Form, []string{"timer", "activity", "tag"})
	if err != nil { panic(err) }
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
		panic(fmt.Errorf("ActivitySubmit: missing fields " + missingString.String()))
	}

	// format the form data
	var tags tagpkg.Tags
	if !noTags {
		for _, tag := range fields[2] {
			if !tags.Contains(tag) {
				tags = append(tags, tagpkg.TagMetaData{ID: -1, Tag: tag, Count: 1})
			}
		}
	}
	activityDuration = fields[0][0]
	activityDesc = fields[1][0]
	hours, mins, err := timer.ParseTimer(activityDuration)
	if err != nil { panic(err) }
	var totalMins int64 = int64(mins) + 60*int64(hours)

	if DEBUG {
		fmt.Printf("Split up timer: %v:%v\nActivity: %s\n", hours, mins, activityDesc)
	}

	// grab the date when the timer started
	start, err := timer.StartingTime(int(totalMins))
	if err != nil { panic(err) }

	if DEBUG {
		fmt.Printf("Timer started %v min ago\nIt was %v\n", totalMins, start.Format("2006-01-02"))
	}

	// get the timerData ready to insert into DB
	timer := dashboard.ActivityMetaData{
		ID:						-1,
		UserID: 			userID,
		Day:					start.Format("2006-01-02"),
		Duration: 		totalMins,
		Description:	activityDesc,
		Tags:					tags,
	}

	// insert newest activity data into DB
	_, err = database.AddDay(timer)
	if err != nil { panic(err) }
	
	// grab all of the activities on this day (newly inserted included)
	activities, err := database.GetDayData(timer.UserID, timer.Day)
	if err != nil { panic(err) }

	// generate a CardMetaData for the template
	totalMins = 0
	t := make(tagpkg.Tags, 0)
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
		Tags:	t,
		TotalMins: totalMins,
		Day: timer.Day,
		SwapOOB: template.HTMLAttr(`hx-swap-oob="true"`),
	}

	// generate the template and respond with it
	cardTemplate := template.Must(template.New("card").Funcs(templateutil.DashboardFuncMap()).ParseFiles(templateutil.ParseFiles["card"]...))
	err = cardTemplate.Execute(w, card)
	if err != nil { panic(fmt.Errorf("executing template: %v", err)) }

	resetTimer(w)
}

func resetTimer(w http.ResponseWriter) {
	defaultTimer := template.Must(template.New("defaulttimer").ParseFiles(templateutil.ParseFiles["defaulttimer"]...))
	err := defaultTimer.Execute(w, templateutil.NewElementInfo("start-button", []string{"timer-button", "button-main"}, nil, ""))
	if err != nil { panic(err) }
}
