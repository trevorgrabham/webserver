package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/trevorgrabham/webserver/webserver/lib/database"
	"github.com/trevorgrabham/webserver/webserver/lib/util"
)

func HandleTagSummary(w http.ResponseWriter, r *http.Request) {
	// Parse out the offset value
	if err := r.ParseForm(); err != nil {
		log.Fatalf("HandleTagSummary: %v", err)
	}
	res, ok := r.Form["offset"]
	var offset int64
	if ok {
		var err error
		offset, err = strconv.ParseInt(res[0], 0, 64)
		if err != nil {
			log.Fatalf("HandleTagSummary: %v", err)
		}
	}
	if offset == 1 {
		offset = 0
	}

	tags, err := database.GetTagData(offset) 
	if err != nil {
		log.Fatal(err)
	}


	switch r.Header.Get("Hx-Trigger") {
	case "tag-summary-section":
		if err := util.AllTagSummaryTemplateReady.Execute(w, tags); err != nil {
			log.Fatalf("Executing TagSummaryTemplate: %v", err)
		}
	case "load-tag-summary-button":
		var max int64
		res, ok = r.Form["max"]
		if !ok {
			log.Fatal("HandleTagSummary(): no 'max' attribute provided")
		}
		max, err = strconv.ParseInt(res[0], 0, 64)
		if err != nil {
			log.Fatalf("HandleTagSummary(): %v", err)
		}
		tags.MaxCount = max
		if err := util.TagSummaryTemplateReady.Execute(w, tags); err != nil {
			log.Fatalf("Executing TagSummaryTemplate: %v", err)
		}
	}

	if len(tags.Tags) < 10 {
		fmt.Fprint(w, util.DisabledMoreTagsButton)
	}
}