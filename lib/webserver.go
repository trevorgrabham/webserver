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
	// http.HandleFunc("/tagInput", handlers.HandleTagInput)
	http.HandleFunc("/addTag", handlers.HandleAddTag)
	http.HandleFunc("/removeTag", handlers.HandleRemove)
	http.HandleFunc("/resetTimer", handlers.HandleResetTimer)
	http.HandleFunc("/submitActivity", handlers.HandleActivitySubmit)
	http.HandleFunc("/dashboard", handlers.HandleDashboard)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}