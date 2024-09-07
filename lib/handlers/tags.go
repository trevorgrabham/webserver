package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/trevorgrabham/webserver/webserver/lib/database"
	"github.com/trevorgrabham/webserver/webserver/lib/templateutil"
)

func HandleTagSummary(w http.ResponseWriter, r *http.Request) {
	// Parse out the offset value
	if err := r.ParseForm(); err != nil { panic(fmt.Errorf("HandleTagSummary(): %v", err)) }

	userID, ok := r.Context().Value(ContextKey("user-id")).(int64)
	if !ok { panic(fmt.Errorf("unable to parse 'user-id' from handletagsummary()")) }

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
		tagsContainer := template.Must(template.New("tagscontainer").Funcs(templateutil.TagFuncMap()).ParseFiles(templateutil.ParseFiles["tagscontainer"]...))
		err = tagsContainer.Execute(w, tags)
		if err != nil { panic(fmt.Errorf("executing TagSummaryTemplate: %v", err)) }
	// Coming from the 'load more' button
	case "load-tag-summary-button":
		var max int64
		res, ok = r.Form["max"]
		if !ok { panic(fmt.Errorf("HandleTagSummary(): no 'max' attribute provided")) }

		max, err = strconv.ParseInt(res[0], 0, 64)
		if err != nil { panic(fmt.Errorf("HandleTagSummary(): %v", err)) }

		tags.MaxCount = max
		tagsTemplate := template.Must(template.New("tags").Funcs(templateutil.TagFuncMap()).ParseFiles(templateutil.ParseFiles["tags"]...))
		err = tagsTemplate.Execute(w, tags)
		if err != nil { panic(fmt.Errorf("executing TagSummaryTemplate: %v", err)) }
	}
}