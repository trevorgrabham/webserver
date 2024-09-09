package handlers

import (
	"html/template"
	"net/http"

	"github.com/trevorgrabham/webserver/webserver/html"
)

func HandleIndex(w http.ResponseWriter, r *http.Request) {
	startHyperscript := `_="on click
												send startTimer to #timer-display
												add .hidden to me
												remove .hidden from #pause-timer-button
												remove .hidden from #stop-timer-button"`
	pauseHyperscript := `_="on click
												send stopTimer to #timer-display 
												add .hidden to me
												remove .hidden from #start-timer-button"`
	stopHyperscript := `_="on click send stopTimer to #timer-display"`
	indexData := html.TimerData{
		StartButton: html.NewElementAttributes(append([]string{`id="start-timer-button"`, `class="svg-button"`}, startHyperscript)),
		PauseButton: html.NewElementAttributes(append([]string{`id="pause-timer-button"`, `class="svg-button hidden"`}, pauseHyperscript)),
		StopButton: html.NewElementAttributes(append([]string{`id="stop-timer-button"`, `class="svg-button hidden"`, `hx-get="/stopTimer"`, `hx-target="#timer-buttons-container"`, `hx-swap="outerHTML"`}, stopHyperscript)),
	}
	
	index := template.Must(template.New("home").ParseFiles(html.IncludeFiles["home"]...))
	if err := index.Execute(w, indexData); err != nil { panic(err) }
}