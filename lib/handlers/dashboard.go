package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/trevorgrabham/webserver/webserver/lib/database"
	"github.com/trevorgrabham/webserver/webserver/lib/util"
)

func HandleDashboard(w http.ResponseWriter, r *http.Request) {
	// Parse out maxItems parameter
	if err := r.ParseForm(); err != nil {
		log.Fatalf("Parsing form: %v", err)
	}
	
	res, ok := r.Form["maxItems"]
	var maxItems int64
	if ok {
		var err error
		maxItems, err = strconv.ParseInt(res[0], 10, 64)
		if err != nil {
			log.Fatalf("Parsing maxItems (%s): %v", res[0], err)
		}
	}

	cards, err := database.GetCardData(maxItems)
	if err != nil {
		log.Fatal(err)
	}

	if err := util.AllCardsTemplateReady.Execute(w, cards); err != nil {
		log.Fatalf("Executing template: %v", err)
	}
}	