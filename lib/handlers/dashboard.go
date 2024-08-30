package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/trevorgrabham/webserver/webserver/lib/database"
	dashboardtemplates "github.com/trevorgrabham/webserver/webserver/lib/templates/dashboard"
)

func HandleDashboard(w http.ResponseWriter, r *http.Request) {
	// Parse out maxItems parameter
	if err := r.ParseForm(); err != nil { panic(fmt.Errorf("Parsing form: %v", err)) }

	id, err := CheckIDCookie(w, r)
	if err != nil { panic(err) }
	
	res, ok := r.Form["maxItems"]
	var maxItems int64
	if ok {
		var err error
		maxItems, err = strconv.ParseInt(res[0], 10, 64)
		if err != nil { panic(fmt.Errorf("Parsing maxItems (%s): %v", res[0], err)) }
	}

	cards, err := database.GetCardData(id, maxItems)
	if err != nil { panic(err) }

	if err := dashboardtemplates.AllCardsTemplateReady.Execute(w, cards); err != nil { panic(fmt.Errorf("Executing template: %v", err)) }
}	