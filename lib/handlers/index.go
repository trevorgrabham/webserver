package handlers

import (
	"html/template"
	"net/http"

	"github.com/trevorgrabham/webserver/webserver/html"
)

func HandleIndex(w http.ResponseWriter, _ *http.Request) {
	startHyperscript := `_="on click trigger click on #nav-start-button"`
	pauseHyperscript := `_="on click trigger click on #nav-pause-button"`
	stopHyperscript := `_="on click trigger click on #nav-stop-button"`
	indexData := html.TimerData{
		ButtonContainer: html.NewElementAttributes([]string{`id="main-timer-buttons-container"`, `class="timer-buttons-container"`}),
		StartButton: html.NewElementAttributes(append([]string{`id="main-start-button"`, `class="svg-button start-timer-button"`}, startHyperscript)),
		PauseButton: html.NewElementAttributes(append([]string{`id="main-pause-button"`, `class="svg-button hidden pause-timer-button"`}, pauseHyperscript)),
		StopButton: html.NewElementAttributes(append([]string{`id="main-stop-button"`, `class="svg-button hidden stop-timer-button"`}, stopHyperscript)),
	}
	
	index := template.Must(template.New("home").ParseFiles(html.IncludeFiles["home"]...))
	if err := index.Execute(w, indexData); err != nil { panic(err) }
}