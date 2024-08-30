package webserver

import (
	"log"
	"net/http"

	"github.com/trevorgrabham/webserver/webserver/lib/handlers"
)

func StaticHandler(next http.Handler) http.Handler{
	// Makes sure that there is a `client-id` cookie set, and returns the value. Don't need the value, just need to make sure that there is a cookie present
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlers.CheckIDCookie(w, r)
		next.ServeHTTP(w, r)
	})
}


func StartServer() {
	// Timer
	http.Handle("/", StaticHandler(http.Handler(http.FileServer(http.Dir("./static")))))
	http.HandleFunc("/startTimer", handlers.HandleStartTimer)
	http.HandleFunc("/pauseTimer", handlers.HandlePauseTimer)
	http.HandleFunc("/resumeTimer", handlers.HandleResumeTimer)
	http.HandleFunc("/stopTimer", handlers.HandleStopTimer)
	http.HandleFunc("/activitySuggestions", handlers.HandleActivitySuggestions)
	http.HandleFunc("/tagSuggestions", handlers.HandleTagSuggestions)
	http.HandleFunc("/addTag", handlers.HandleAddTag)
	http.HandleFunc("/removeTag", handlers.HandleRemove)
	http.HandleFunc("/resetTimer", handlers.HandleResetTimer)
	http.HandleFunc("/submitActivity", handlers.HandleActivitySubmit)
	// Dashboard
	http.HandleFunc("/dashboard", handlers.HandleDashboard)
	// Tag Summary
	http.HandleFunc("/tagSummary", handlers.HandleTagSummary)

	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}