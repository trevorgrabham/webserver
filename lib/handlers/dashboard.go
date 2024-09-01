package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/trevorgrabham/webserver/webserver/lib/database"
	"github.com/trevorgrabham/webserver/webserver/lib/templateutil"
)

func HandleDashboard(w http.ResponseWriter, r *http.Request) {
	// Parse out maxItems parameter
	if err := r.ParseForm(); err != nil { panic(fmt.Errorf("parsing form: %v", err)) }

	id, err := CheckIDCookie(w, r)
	if err != nil { panic(err) }
	
	res, ok := r.Form["maxItems"]
	var maxItems int64
	if ok {
		var err error
		maxItems, err = strconv.ParseInt(res[0], 10, 64)
		if err != nil { panic(fmt.Errorf("parsing maxItems (%s): %v", res[0], err)) }
	}

	cards, err := database.GetCardData(id, maxItems)
	if err != nil { panic(err) }

	allCards := template.Must(template.New("cards").Funcs(templateutil.DashboardFuncMap()).ParseFiles(templateutil.ParseFiles["cards"]...))
	err = allCards.Execute(w, cards)
	if err != nil { panic(fmt.Errorf("executing template: %v", err)) }
}	