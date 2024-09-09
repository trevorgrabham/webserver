package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/trevorgrabham/webserver/webserver/html"
	"github.com/trevorgrabham/webserver/webserver/lib/chart"
	"github.com/trevorgrabham/webserver/webserver/lib/database"
)

const dateMask = "2006-01-02"

func HandleChart(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(ContextKey("user-id")).(int64)
	if !ok { panic(fmt.Errorf("unable to parse 'user-id' in handlechart()")) }

	params := r.URL.Query()
	var start, end *time.Time
	if params.Has("start") {
		t, err := time.Parse(dateMask, params.Get("start"))
		if err != nil { panic(fmt.Errorf("parsing 'start' from url query: %v", err.Error())) }
		start = &t
	}
	if params.Has("end") {
		t, err := time.Parse(dateMask, params.Get("end"))
		if err != nil { panic(fmt.Errorf("parsing 'end' from url query: %v", err.Error())) }
		end = &t
	}

	if start == nil || end == nil {
		s, e, err := database.GetStartEndData(userID)
		if err != nil || s == nil || e == nil { panic(fmt.Errorf("getting start and end values: %v", err)) }
		if start == nil { start = s }
		if end == nil { end = e }
	}

	data, err := database.GetChartData(userID, start, end)
	if err != nil { panic(fmt.Errorf("getting relavent data from database: %v", err.Error())) }

	bars, err := chart.NewChartBars(data, start, end)
	if err != nil { panic(fmt.Errorf("parsing data into bars for the chart: %v", err.Error())) }

	chart := template.Must(template.New("chart").ParseFiles(html.IncludeFiles["chart"]...))
	err = chart.Execute(w, bars)
	if err != nil { panic(fmt.Errorf("executing chart template: %v", err.Error())) }
}