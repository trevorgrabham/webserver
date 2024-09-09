package handlers

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/trevorgrabham/webserver/webserver/html"
	"github.com/trevorgrabham/webserver/webserver/lib/database"
)

func HandleNav(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(ContextKey("user-id")).(int64)
	if !ok { panic(fmt.Errorf("unable to parse 'user-id' from handlenav()")) }

	user, err := database.GetUser(userID)
	if err != nil { panic(err.Error()) }

	nav := template.Must(template.New("nav").Funcs(html.NavFuncMap).ParseFiles(html.IncludeFiles["nav"]...))
	if err := nav.Execute(w, user); err != nil { panic(err) }
}