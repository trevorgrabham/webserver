package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/trevorgrabham/webserver/webserver/lib/database"
	"github.com/trevorgrabham/webserver/webserver/lib/templateutil"
)

const (
	dateMask = "2006-01-02"
)

func HandleChart(w http.ResponseWriter, r *http.Request) {
	userID, err := CheckIDCookie(w, r)
	if err != nil { http.Error(w, fmt.Sprintf("getting 'client-id' cookie: %v", err.Error()), http.StatusBadRequest); return }

	params := r.URL.Query()
	var start, end *time.Time
	if params.Has("start") {
		t, err := time.Parse(dateMask, params.Get("start"))
		if err != nil { http.Error(w, fmt.Sprintf("parsing 'start' from url query: %v", err.Error()), http.StatusBadRequest); return }
		start = &t
	}
	if params.Has("end") {
		t, err := time.Parse(dateMask, params.Get("end"))
		if err != nil { http.Error(w, fmt.Sprintf("parsing 'end' from url query: %v", err.Error()), http.StatusBadRequest); return }
		end = &t
	}

	if start == nil || end == nil {
		s, e, err := database.GetStartEndData(userID)
		if err != nil { http.Error(w, fmt.Sprintf("getting start and end values: %v", err.Error()), http.StatusInternalServerError) }
		if start == nil { start = s }
		if end == nil { end = e }
	}

	data, err := database.GetChartData(userID, start, end)
	if err != nil { http.Error(w, fmt.Sprintf("getting relavent data from database: %v", err.Error()), http.StatusInternalServerError) }

	bars, err := templateutil.NewChartBars(data, start, end)
	if err != nil { http.Error(w, fmt.Sprintf("parsing data into bars for the chart: %v", err.Error()), http.StatusBadRequest) }

	chart := template.Must(template.New("chart").ParseFiles(templateutil.ParseFiles["chart"]...))
	err = chart.Execute(w, bars)
	if err != nil { http.Error(w, fmt.Sprintf("executing chart template: %v", err.Error()), http.StatusBadRequest) }
}