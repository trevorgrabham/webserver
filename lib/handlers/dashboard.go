package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/trevorgrabham/webserver/webserver/html"
	"github.com/trevorgrabham/webserver/webserver/lib/database"
)

func HandleDashboard(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(ContextKey("user-id")).(int64)
	if !ok { panic(fmt.Errorf("unable to parse 'user-id' from handledashboard()")) }
	
	// Parse out maxItems parameter
	if err := r.ParseForm(); err != nil { panic(fmt.Errorf("parsing form: %v", err)) }
	maxString := r.Form.Get("maxItems")
	if maxString == "" { panic(fmt.Errorf("unable to parse 'maxItems' from handledashboard()")) }
	maxItems, err := strconv.ParseInt(maxString, 10, 64)
	if err != nil { panic(fmt.Errorf("parsing maxItems (%s): %v", maxString, err)) }

	cards, err := database.GetCardData(userID, maxItems)
	if err != nil { panic(err) }

	allCards := template.Must(template.New("cards").Funcs(html.DashboardFuncMap).ParseFiles(html.IncludeFiles["cards"]...))
	if err = allCards.Execute(w, cards); err != nil { panic(fmt.Errorf("executing template: %v", err)) }
}	