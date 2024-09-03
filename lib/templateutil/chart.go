package templateutil

import (
	"fmt"
	"time"

	"github.com/trevorgrabham/webserver/webserver/lib/chart"
)

const (
	hoursInFourYears = 24 * 365 * 4.0
	hoursInSixMonths = 24 * 182.5
	hoursInSixWeeks = 24 * 7 * 6.0
)

func NewChartBars(data []chart.Data, start, end *time.Time) (bars []chart.ChartBar, err error) {
	if start == nil || end == nil { return nil, fmt.Errorf("NewChartBars(%s, %s): Nil values for 'start' or 'end' are not permitted", start.Format(chart.DateMask), end.Format(chart.DateMask)) }

	switch diff := end.Sub(*start).Hours(); {
	case diff > hoursInFourYears:
		return generateBarsByYear(data, start, end)
	case diff > hoursInSixMonths:
		return generateBarsByMonth(data, start, end)
	case diff > hoursInSixWeeks:
		return generateBarsByWeek(data, start, end)
	case diff > 24.0:
		return generateBarsByDay(data, start, end)
	default:
		return generateBarsByEntry(data, start, end)
	}
}

func generateBarsByYear(data []chart.Data, start, end *time.Time) (bars []chart.ChartBar, err error) {
	bars = make([]chart.ChartBar, 0)
	currDate, err := time.Parse(chart.YearMask, start.Format(chart.YearMask))
	if err != nil { return nil, fmt.Errorf("generateByYear(): %v", err) }

	bar := chart.ChartBar{
		StartDate: currDate.Format(chart.DateMask),
		EndDate: currDate.AddDate(1,0,-1).Format(chart.DateMask)}
	for _, dataPoint := range data {
		t, err := time.Parse(chart.YearMask, dataPoint.Day)
		if err != nil { return nil, fmt.Errorf("generateByYear(): %v", err) }
		year := t.Format(chart.YearMask)
		if currDate.Format(chart.YearMask) == year {
			bar.Count += dataPoint.Duration
			continue
		}
		bars = append(bars, bar)
		currDate = currDate.AddDate(1,0,0)
		for {
			if currDate.Format(chart.YearMask) == year { break }
			bars = append(bars, chart.ChartBar{ 
				StartDate: currDate.Format(chart.DateMask),
				EndDate: currDate.AddDate(1,0,-1).Format(chart.DateMask)})
			currDate = currDate.AddDate(1,0,0)
		}
		bar = chart.ChartBar{ 
			StartDate: currDate.Format(chart.DateMask), 
			EndDate: currDate.AddDate(1,0,-1).Format(chart.DateMask), 
			Count: dataPoint.Duration}	
	}
	// We will always have our latest 'bar' open, since we only add it once we reach a date that is not equal to the latest bar
	bars = append(bars, bar)
	return
}

func generateBarsByMonth(data []chart.Data, start, end *time.Time) (bars []chart.ChartBar, err error) {
	bars = make([]chart.ChartBar, 0)
	currDate, err := time.Parse(chart.MonthMask, start.Format(chart.MonthMask))
	if err != nil { return nil, fmt.Errorf("generateByMonth(): %v", err) }

	bar := chart.ChartBar{
		StartDate: currDate.Format(chart.DateMask),
		EndDate: currDate.AddDate(0,1,-1).Format(chart.DateMask)}
	for _, dataPoint := range data {
		t, err := time.Parse(chart.MonthMask, dataPoint.Day)
		if err != nil { return nil, fmt.Errorf("generateByMonth(): %v", err) }
		month := t.Format(chart.MonthMask)
		if currDate.Format(chart.MonthMask) == month {
			bar.Count += dataPoint.Duration
			continue
		}
		bars = append(bars, bar)
		currDate = currDate.AddDate(0,1,0)
		for {
			if currDate.Format(chart.MonthMask) == month { break }
			bars = append(bars, chart.ChartBar{
				StartDate: currDate.Format(chart.DateMask),
				EndDate: currDate.AddDate(0,1,-1).Format(chart.DateMask)})
			currDate = currDate.AddDate(0,1,0)
		}
		bar = chart.ChartBar{
			StartDate: currDate.Format(chart.DateMask),
			EndDate: currDate.AddDate(0,1,-1).Format(chart.DateMask),
			Count: dataPoint.Duration}
	}
	// We will always have our latest 'bar' open, since we only add it once we reach a date that is not equal to the latest bar
	bars = append(bars, bar)
	return
}

func generateBarsByWeek(data []chart.Data, start, end *time.Time) (bars []chart.ChartBar, err error) {
	bars = make([]chart.ChartBar, 0)
	currDate := *start
	bar := chart.ChartBar{
		StartDate: currDate.Format(chart.DateMask),
		EndDate: currDate.AddDate(0,0,6).Format(chart.DateMask)}
	for _, dataPoint := range data {
		day, err := time.Parse(chart.DateMask, dataPoint.Day)
		if err != nil { return nil, fmt.Errorf("generateByWeek(): %v", err) }
		if !day.After(currDate.AddDate(0,0,6)) {
			bar.Count += dataPoint.Duration
			continue
		}
		bars = append(bars, bar)
		currDate = currDate.AddDate(0,0,7)
		for {
			if !day.After(currDate.AddDate(0,0,6)) { break }
			bars = append(bars, chart.ChartBar{ StartDate: currDate.Format(chart.DateMask), EndDate: currDate.AddDate(0,0,6).Format(chart.DateMask)})
			currDate = currDate.AddDate(0,0,7)
		}
		bar = chart.ChartBar{
			StartDate: currDate.Format(chart.DateMask),
			EndDate: currDate.AddDate(0,0,6).Format(chart.DateMask),
			Count: dataPoint.Duration}
	}
	// We will always have our latest 'bar' open, since we only add it once we reach a date that is not equal to the latest bar
	bars = append(bars, bar)
	return
}

func generateBarsByDay(data []chart.Data, start, end *time.Time) (bars []chart.ChartBar, err error) {
	bars = make([]chart.ChartBar, 0)
	// We will always have our latest 'bar' open, since we only add it once we reach a date that is not equal to the latest bar
	currDate := *start 
	bar := chart.ChartBar{
		StartDate: currDate.Format(chart.DateMask),
		EndDate: currDate.Format(chart.DateMask)}
	for _, dataPoint := range data {
		day := dataPoint.Day 
		if currDate.Format(chart.DateMask) == day {
			bar.Count += dataPoint.Duration
			continue
		}
		bars = append(bars, bar)
		currDate = currDate.AddDate(0,0,1)
		for {
			if currDate.Format(chart.DateMask) == day { break }
			bars = append(bars, chart.ChartBar{ StartDate: currDate.Format(chart.DateMask), EndDate: currDate.Format(chart.DateMask)})
			currDate = currDate.AddDate(0,0,1)
		}
		bar = chart.ChartBar{
			StartDate: currDate.Format(chart.DateMask),
			EndDate: currDate.Format(chart.DateMask),
			Count: dataPoint.Duration}
	}
	bars = append(bars, bar)
	return
}

func generateBarsByEntry(data []chart.Data, _, _ *time.Time) (bars []chart.ChartBar, err error) {
	bars = make([]chart.ChartBar, len(data))
	for _, dataPoint := range data {
		bars = append(bars, chart.ChartBar{ Count: dataPoint.Duration })
	}
	return
}