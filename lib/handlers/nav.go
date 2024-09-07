package handlers

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/trevorgrabham/webserver/webserver/lib/database"
	"github.com/trevorgrabham/webserver/webserver/lib/templateutil"
)

func HandleNav(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(ContextKey("user-id")).(int64)
	if !ok { panic(fmt.Errorf("unable to parse 'user-id' from handlenav()")) }
	user, err := database.GetUser(userID)
	if err != nil { http.Error(w, err.Error(), http.StatusBadRequest); return }
	nav := template.Must(template.New("nav").Funcs(templateutil.NavFuncMap()).ParseFiles(templateutil.ParseFiles["nav"]...))
	nav.Execute(w, user)
}