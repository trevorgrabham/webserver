package webserver

import (
	"fmt"
	"log"
	"net/http"
)

func handleStartTimer(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(w, `
	<div id="timer-buttons-container">
		<button id="pause-timer" _="on click send stopTimer to #timer" hx-get="/pauseTimer" hx-swap="outerHTML">Pause</button>
		<button id="stop-timer" _="on click send stopTimer to #timer" hx-get="/stopTimer" hx-target="#timer-buttons-container">Stop</button>
	</div>
	`)
}

func handlePauseTimer(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(w, `
	<button id="resume-timer" _="on click send startTimer to #timer" hx-get="/resumeTimer" hx-swap="outerHTML">Resume</button>
	`)
}

func handleResumeTimer(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(w, `
	<button id="pause-timer" _="on click send stopTimer to #timer" hx-get="/pauseTimer" hx-swap="outerHTML">Pause</button>
	`)
}

func handleStopTimer(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(w, `
	<input type="text" name="activity" placeholder="What were you doing?" />
	<button hx-post="/submitActivity" hx-target="#timer-container" _="on click call resetTimer()">Submit</button>
	`)
}

func handleActivitySubmit(w http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		log.Fatalf("handleActivitySubmit parsing form: %v", err)
	}

	log.Println(req.Form["activity"])
	log.Println(req.Form["timer"])

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
      <input id="hidden-timer" name="timer" type="hidden" value="0:00:00" />
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

func StartServer() {
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/startTimer", handleStartTimer)
	http.HandleFunc("/pauseTimer", handlePauseTimer)
	http.HandleFunc("/resumeTimer", handleResumeTimer)
	http.HandleFunc("/stopTimer", handleStopTimer)
	http.HandleFunc("/submitActivity", handleActivitySubmit)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}