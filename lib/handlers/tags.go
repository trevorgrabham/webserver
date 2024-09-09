package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/trevorgrabham/webserver/webserver/html"
	"github.com/trevorgrabham/webserver/webserver/lib/database"
)

func HandleTagSummary(w http.ResponseWriter, r *http.Request) {
	// Parse out the offset value
	if err := r.ParseForm(); err != nil { panic(fmt.Errorf("HandleTagSummary(): %v", err)) }

	userID, ok := r.Context().Value(ContextKey("user-id")).(int64)
	if !ok { panic(fmt.Errorf("unable to parse 'user-id' from handletagsummary()")) }

	offsetString := r.Form.Get("offset")
	var offset int64
	if offsetString != "" {
		var err error
		offset, err = strconv.ParseInt(offsetString, 0, 64)
		if err != nil { panic(fmt.Errorf("HandleTagSummary(): %v", err)) }
	}
	// if offset == 1 {
	// 	offset = 0
	// }

	tags, err := database.GetTagData(userID, offset) 
	if err != nil { panic(err) }

	switch r.Header.Get("Hx-Trigger") {
	// Coming from '/'
	case "tag-summary-section":
		tagsContainer := template.Must(template.New("tags-container").Funcs(html.TagFuncMap).ParseFiles(html.IncludeFiles["tags-container"]...))
		if err = tagsContainer.Execute(w, tags); err != nil { panic(fmt.Errorf("executing TagSummaryTemplate: %v", err)) }
	// Coming from the 'load more' button
	case "load-tag-summary-button":
		maxString := r.Form.Get("max")
		if maxString == "" { panic(fmt.Errorf("HandleTagSummary(): no 'max' attribute provided")) }

		max, err := strconv.ParseInt(maxString, 0, 64)
		if err != nil { panic(fmt.Errorf("HandleTagSummary(): %v", err)) }
		for i := range tags {
			tags[i].MaxCount = max
		}

		tagsTemplate := template.Must(template.New("tags").Funcs(html.TagFuncMap).ParseFiles(html.IncludeFiles["tags"]...))
		if err = tagsTemplate.Execute(w, tags); err != nil { panic(fmt.Errorf("executing TagSummaryTemplate: %v", err)) }
	}
}