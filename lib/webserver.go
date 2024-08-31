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
	// Nav 
	http.HandleFunc("/nav", handlers.HandleNav)
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
	// Profile
	http.HandleFunc("/profile", handlers.HandleProfile)
	http.HandleFunc("/profile/editPic", handlers.HandleEditPic)
	http.HandleFunc("/profile/savePic", handlers.HandleSavePic)
	http.HandleFunc("/profile/editName", handlers.HandleEditName)
	http.HandleFunc("/profile/saveName", handlers.HandleSaveName)
	http.HandleFunc("/profile/editEmail", handlers.HandleEditEmail)
	http.HandleFunc("/profile/saveEmail", handlers.HandleSaveEmail)

	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}