package handlers

import (
	"html/template"
	"net/http"

	"github.com/trevorgrabham/webserver/webserver/lib/database"
	"github.com/trevorgrabham/webserver/webserver/lib/templateutil"
)

func HandleNav(w http.ResponseWriter, r *http.Request) {
	userID, err := CheckIDCookie(w, r)
	if err != nil { http.Error(w, err.Error(), http.StatusBadRequest); return }
	client, err := database.GetClient(userID)
	if err != nil { http.Error(w, err.Error(), http.StatusBadRequest); return }
	nav := template.Must(template.New("nav").Funcs(templateutil.NavFuncMap()).ParseFiles(templateutil.ParseFiles["nav"]...))
	nav.Execute(w, client)
}