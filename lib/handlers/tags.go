package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/trevorgrabham/webserver/webserver/lib/database"
	tagtemplates "github.com/trevorgrabham/webserver/webserver/lib/templates/tags"
)

func HandleTagSummary(w http.ResponseWriter, r *http.Request) {
	// Parse out the offset value
	if err := r.ParseForm(); err != nil { panic(fmt.Errorf("HandleTagSummary(): %v", err)) }

	userID, err := CheckIDCookie(w, r)
	if err != nil { panic(fmt.Errorf("HandleTagSummary(): %v", err))}

	res, ok := r.Form["offset"]
	var offset int64
	if ok {
		var err error
		offset, err = strconv.ParseInt(res[0], 0, 64)
		if err != nil { panic(fmt.Errorf("HandleTagSummary(): %v", err)) }
	}
	if offset == 1 {
		offset = 0
	}

	tags, err := database.GetTagData(userID, offset) 
	if err != nil { panic(err) }

	switch r.Header.Get("Hx-Trigger") {
	// Coming from '/'
	case "tag-summary-section":
		if err := tagtemplates.AllTagSummaryTemplateReady.Execute(w, tags); err != nil { panic(fmt.Errorf("executing TagSummaryTemplate: %v", err)) }
	// Coming from the 'load more' button
	case "load-tag-summary-button":
		var max int64
		res, ok = r.Form["max"]
		if !ok { panic(fmt.Errorf("HandleTagSummary(): no 'max' attribute provided")) }

		max, err = strconv.ParseInt(res[0], 0, 64)
		if err != nil { panic(fmt.Errorf("HandleTagSummary(): %v", err)) }

		tags.MaxCount = max
		if err := tagtemplates.TagSummaryTemplateReady.Execute(w, tags); err != nil { panic(fmt.Errorf("executing TagSummaryTemplate: %v", err)) }
	}

	if len(tags.Tags) < 10 {
		fmt.Fprint(w, tagtemplates.DisabledMoreTagsButton)
	}
}