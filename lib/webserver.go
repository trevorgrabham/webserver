package webserver

import (
	"log"
	"net/http"

	"github.com/trevorgrabham/webserver/webserver/lib/handlers"
)

func StartServer() {
	// Nav 
	http.Handle("/nav", handlers.SetCookieContext(http.HandlerFunc(handlers.HandleNav)))
	// Timer
	http.Handle("/", handlers.SetCookieContext(http.Handler(http.FileServer(http.Dir("./static")))))
	http.Handle("/defaultTimer", handlers.SetCookieContext(http.HandlerFunc(handlers.HandleResetTimer)))
	http.Handle("/startTimer", handlers.SetCookieContext(http.HandlerFunc(handlers.HandleStartTimer)))
	http.Handle("/pauseTimer", handlers.SetCookieContext(http.HandlerFunc(handlers.HandlePauseTimer)))
	http.Handle("/resumeTimer", handlers.SetCookieContext(http.HandlerFunc(handlers.HandleResumeTimer)))
	http.Handle("/stopTimer", handlers.SetCookieContext(http.HandlerFunc(handlers.HandleStopTimer)))
	http.Handle("/activitySuggestions", handlers.SetCookieContext(http.HandlerFunc(handlers.HandleActivitySuggestions)))
	http.Handle("/tagSuggestions", handlers.SetCookieContext(http.HandlerFunc(handlers.HandleTagSuggestions)))
	http.Handle("/addTag", handlers.SetCookieContext(http.HandlerFunc(handlers.HandleAddTag)))
	http.Handle("/removeTag", handlers.SetCookieContext(http.HandlerFunc(handlers.HandleRemove)))
	http.Handle("/resetTimer", handlers.SetCookieContext(http.HandlerFunc(handlers.HandleResetTimer)))
	http.Handle("/submitActivity", handlers.SetCookieContext(http.HandlerFunc(handlers.HandleActivitySubmit)))
	// Dashboard
	http.Handle("/dashboard", handlers.SetCookieContext(http.HandlerFunc(handlers.HandleDashboard)))
	// Tag Summary
	http.Handle("/tagSummary", handlers.SetCookieContext(http.HandlerFunc(handlers.HandleTagSummary)))
	// Profile
	http.Handle("/profile", handlers.SetCookieContext(http.HandlerFunc(handlers.HandleProfile)))
	http.Handle("/profile/editPic", handlers.SetCookieContext(http.HandlerFunc(handlers.HandleEditPic)))
	http.Handle("/profile/savePic", handlers.SetCookieContext(http.HandlerFunc(handlers.HandleSavePic)))
	http.Handle("/profile/editName", handlers.SetCookieContext(http.HandlerFunc(handlers.HandleEditName)))
	http.Handle("/profile/saveName", handlers.SetCookieContext(http.HandlerFunc(handlers.HandleSaveName)))
	http.Handle("/profile/editEmail", handlers.SetCookieContext(http.HandlerFunc(handlers.HandleEditEmail)))
	http.Handle("/profile/saveEmail", handlers.SetCookieContext(http.HandlerFunc(handlers.HandleSaveEmail)))
	// Chart
	http.Handle("/chart", handlers.SetCookieContext(http.HandlerFunc(handlers.HandleChart)))

	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}