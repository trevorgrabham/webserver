package handlers

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/trevorgrabham/webserver/webserver/html"
	"github.com/trevorgrabham/webserver/webserver/lib/database"
)

const DEBUG = true

// TODO: add tag auto-complete

func HandleRemove(w http.ResponseWriter, _ *http.Request) { fmt.Fprint(w) }

func HandleStopTimer(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(ContextKey("user-id")).(int64)
	if !ok { panic(fmt.Errorf("handlestoptimer(): unable to parse 'user-id'")) }

	form := template.Must(template.New("form").ParseFiles(html.IncludeFiles["form"]...))

	activitySuggestions, err := database.GetPreviousActivities(userID)
	if err != nil { panic(err) }
	// TODO
	// tagSuggestions, err := database.GetPreviousTags(userID)
	// if err != nil { panic(err) }

	formData := html.FormData{
		PlusButton: html.NewElementAttributes([]string{`id="add-tag-button"`, `class="timer-button form-sub-button svg-button"`, `hx-get="/addTag"`, `hx-include="#tag-input"`, `hx-params="temporary-tag"`, `hx-target="#tags-container"`, `hx-swap="beforeend"`}),
		SuccessButton: html.NewElementAttributes([]string{`id="submit-timer-form-button"`, `class="timer-button form-button svg-button"`, `type="submit"`, `_="on click call resetTimer()"`}),
		CancelButton: html.NewElementAttributes([]string{`id="cancel-timer-form-button"`, `class="timer-button form-button svg-button"`, `hx-get="/timer"`, `hx-target="#timer-container"`, `_="on click call resetTimer()"`})}
	formData.AddSuggestions("activity", activitySuggestions)
	// TODO
	// formData.AddSuggestions("tag", tagSuggestions)

	if err := form.Execute(w, formData); err != nil { panic(err) }
}
