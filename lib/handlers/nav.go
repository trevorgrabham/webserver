package handlers

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/trevorgrabham/webserver/webserver/html"
)

func HandleNav(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(ContextKey("user-id")).(int64)
	if !ok { panic(fmt.Errorf("unable to parse 'user-id' from handlenav()")) }

	startHyperscript := `_="on click 
														call startTimers()
														add .hidden to .start-timer-button
														remove .hidden from .pause-timer-button
														remove .hidden from .stop-timer-button
														repeat until event stopTimer 
															call updateTimers()
															wait 1s
														end"`
	pauseHyperscript := `_="on click
														send stopTimer to #nav-start-button
														add .hidden to .pause-timer-button
														remove .hidden from .start-timer-button"`
	stopHyperscript := `_="on click send stopTimer to #nav-start-button"`
	navData := html.NavData{
		TimerData: html.TimerData{
			StartButton: html.NewElementAttributes(append([]string{`id="nav-start-button"`, `class="svg-button start-timer-button"`}, startHyperscript)),
			PauseButton: html.NewElementAttributes(append([]string{`id="nav-pause-button"`, `class="svg-button hidden pause-timer-button"`}, pauseHyperscript)),
			StopButton: html.NewElementAttributes(append([]string{`id="nav-stop-button"`, `class="svg-button hidden stop-timer-button"`, `hx-get="/stopTimer"`, `hx-vals="js:{timer: document.querySelector('.timer-display').innerText}"`, `hx-target="#timer-container"`}, stopHyperscript)),
		},
		ID: userID}

	nav := template.Must(template.New("nav").Funcs(html.NavFuncMap).ParseFiles(html.IncludeFiles["nav"]...))
	if err := nav.Execute(w, navData); err != nil { panic(err) }
}

func HandleNavTimer(w http.ResponseWriter, _ *http.Request) {
	startHyperscript := `_="on click 
														call startTimers()
														add .hidden to .start-timer-button
														remove .hidden from .pause-timer-button
														remove .hidden from .stop-timer-button
														repeat until event stopTimer 
															call updateTimers()
															wait 1s
														end"`
	pauseHyperscript := `_="on click
														send stopTimer to #nav-start-button
														add .hidden to .pause-timer-button
														remove .hidden from .start-timer-button"`
	stopHyperscript := `_="on click send stopTimer to #nav-start-button"`
	timer := html.TimerData{
		StartButton: html.NewElementAttributes(append([]string{`id="nav-start-button"`, `class="svg-button start-timer-button"`}, startHyperscript)),
		PauseButton: html.NewElementAttributes(append([]string{`id="nav-pause-button"`, `class="svg-button hidden pause-timer-button"`}, pauseHyperscript)),
		StopButton: html.NewElementAttributes(append([]string{`id="nav-stop-button"`, `class="svg-button hidden stop-timer-button"`, `hx-get="/stopTimer"`, `hx-vals="js:{timer: document.querySelector('.timer-display').innerText}"`, `hx-target="#timer-container"`}, stopHyperscript)),
	}
	
	navTimer := template.Must(template.New("nav-timer").ParseFiles(html.IncludeFiles["nav-timer"]...))
	if err := navTimer.Execute(w, timer); err != nil { panic(err) }
}