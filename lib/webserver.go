package webserver

import (
	"log"
	"net/http"

	"github.com/trevorgrabham/webserver/webserver/lib/handlers"
)


func StartServer() {
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/startTimer", handlers.HandleStartTimer)
	http.HandleFunc("/pauseTimer", handlers.HandlePauseTimer)
	http.HandleFunc("/resumeTimer", handlers.HandleResumeTimer)
	http.HandleFunc("/stopTimer", handlers.HandleStopTimer)
	http.HandleFunc("/submitActivity", handlers.HandleActivitySubmit)
	http.HandleFunc("/tagInput", handlers.HandleTagInput)
	http.HandleFunc("/dashboard", handlers.HandleDashboard)
	http.HandleFunc("/remove", handlers.HandleRemove)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}