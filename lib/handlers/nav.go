package handlers

import (
	"net/http"

	"github.com/trevorgrabham/webserver/webserver/lib/database"
	navtemplate "github.com/trevorgrabham/webserver/webserver/lib/templates/nav"
)

func HandleNav(w http.ResponseWriter, r *http.Request) {
	userID, err := CheckIDCookie(w, r)
	if err != nil { http.Error(w, err.Error(), http.StatusBadRequest); return }
	client, err := database.GetClient(userID)
	if err != nil { http.Error(w, err.Error(), http.StatusBadRequest); return }
	navtemplate.Nav.Execute(w, client)
}