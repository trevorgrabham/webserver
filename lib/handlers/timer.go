package handlers

import (
	"fmt"
	"net/http"
	"slices"
	"strings"
	"unicode/utf8"

	"github.com/trevorgrabham/webserver/webserver/lib/dashboard"
	"github.com/trevorgrabham/webserver/webserver/lib/database"
	tagpkg "github.com/trevorgrabham/webserver/webserver/lib/tag"
	dashboardtemplate "github.com/trevorgrabham/webserver/webserver/lib/templates/dashboard"
	timertemplate "github.com/trevorgrabham/webserver/webserver/lib/templates/timer"
	"github.com/trevorgrabham/webserver/webserver/lib/timer"
)

const DEBUG = true

func HandleRemove(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(w, ``)
}

func HandleStartTimer(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(w, 
		timertemplate.PauseButton("pause-button", []string{"timer-button", "button-main"}, nil) + 
		timertemplate.StopButton("stop-button", []string{"timer-button", "main-button"}, nil))
}

func HandlePauseTimer(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(w, timertemplate.ResumeButton("resume-button", []string{"timer-button", "button-main"}, nil))
}

func HandleResumeTimer(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(w, timertemplate.PauseButton("pause-button", []string{"timer-button", "buton-main"}, nil))
}

func HandleStopTimer(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(w, timertemplate.Form())
}

func HandleActivitySuggestions(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil { panic(fmt.Errorf("HandleActivitySuggestions(): parsing form : %v", err)) }

	userID, err := CheckIDCookie(w, r)
	if err != nil { panic(err) }

	res, ok := r.Form["activity"]
	if !ok || len(res[0]) < 1 {
		fmt.Fprint(w)
		return
	}
	activityPartial := res[0]

	if DEBUG {
		fmt.Printf("Searching for previous activities matching %s\n", activityPartial)
	}

	previousActivities, err := database.GetPreviousActivities(userID)
	if err != nil { panic(err) }

	matches := timer.FilterFromPartialString(activityPartial, previousActivities)
	if matches.Length == 0 {
		fmt.Fprint(w)
		return
	}

	if DEBUG {
		fmt.Printf("Found matches %v\n", matches)
	}

	if err := timertemplate.AutocompleteTemplateReady("activity-suggestions").Execute(w, matches); err != nil { panic(fmt.Errorf("HandleActivitySuggestions(): %v", err)) }
}

func HandleTagSuggestions(w http.ResponseWriter, r *http.Request) {
	
}

func HandleAddTag(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil { panic(fmt.Errorf("parsing form: %v", err)) }

	res, ok := r.Form["temporary-tag"]
	if !ok || len(res) < 1 { panic(fmt.Errorf("'temporary-tag' was not a provided paramter")) }

	tag := res[0]

	if len(tag) <= 0 {
		fmt.Fprint(w)
		return
	}

	fmt.Fprint(w, 
		timertemplate.NewTagTemplate(tag, utf8.RuneCountInString(tag), "tag-input", []string{"timer-form-input"}, []string{`hx-swap-oob="true"`})) 
}

func HandleResetTimer(w http.ResponseWriter, _ *http.Request) {
	resetTimer(w)
}

func HandleActivitySubmit(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil { panic(fmt.Errorf("handleActivitySubmit parsing form: %v", err)) }

	if DEBUG {
		fmt.Println(r.Form)
	}

	userID, err := CheckIDCookie(w, r)
	if err != nil { panic(err) }

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
	}

	// generate the template and respond with it
	if err := dashboardtemplate.SingleCardTemplateReady.Execute(w, card); err != nil { panic(fmt.Errorf("executing template: %v", err)) }

	resetTimer(w)
}

func resetTimer(w http.ResponseWriter) {
	fmt.Fprint(w, timertemplate.ResetButton("start-button", []string{"timer-button", "button-main"}, nil))
}
