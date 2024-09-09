package chart

import (
	"fmt"
	"time"
)

type Data struct {
	Duration float64
	Day      string
}

type ChartBar struct {
	StartDate string
	EndDate   string
	Count     float64
	Relative  float64
}

const (
	DateMask  = "2006-01-02"
	YearMask  = "2006"
	MonthMask = "01"
	hoursInFourYears = 24 * 365 * 4.0
	hoursInSixMonths = 24 * 182.5
	hoursInSixWeeks  = 24 * 7 * 6.0
)

func NewChartBars(data []Data, start, end *time.Time) (bars []ChartBar, err error) {
	if start == nil || end == nil {
		return nil, fmt.Errorf("NewChartBars(%s, %s): Nil values for 'start' or 'end' are not permitted", start.Format(DateMask), end.Format(DateMask))
	}

	switch diff := end.Sub(*start).Hours(); {
	case diff > hoursInFourYears:
		return generateBarsByYear(data, start)
	case diff > hoursInSixMonths:
		return generateBarsByMonth(data, start)
	case diff > hoursInSixWeeks:
		return generateBarsByWeek(data, start)
	case diff > 24.0:
		return generateBarsByDay(data, start)
	default:
		return generateBarsByEntry(data, start)
	}
}

func generateBarsByYear(data []Data, start *time.Time) (bars []ChartBar, err error) {
	bars = make([]ChartBar, 0)
	currDate, err := time.Parse(YearMask, start.Format(YearMask))
	if err != nil {
		return nil, fmt.Errorf("generateByYear(): %v", err)
	}

	bar := ChartBar{
		StartDate: currDate.Format(DateMask),
		EndDate:   currDate.AddDate(1, 0, -1).Format(DateMask)}
	for _, dataPoint := range data {
		t, err := time.Parse(YearMask, dataPoint.Day)
		if err != nil {
			return nil, fmt.Errorf("generateByYear(): %v", err)
		}
		year := t.Format(YearMask)
		if currDate.Format(YearMask) == year {
			bar.Count += dataPoint.Duration
			continue
		}
		bars = append(bars, bar)
		currDate = currDate.AddDate(1, 0, 0)
		for {
			if currDate.Format(YearMask) == year {
				break
			}
			bars = append(bars, ChartBar{
				StartDate: currDate.Format(DateMask),
				EndDate:   currDate.AddDate(1, 0, -1).Format(DateMask)})
			currDate = currDate.AddDate(1, 0, 0)
		}
		bar = ChartBar{
			StartDate: currDate.Format(DateMask),
			EndDate:   currDate.AddDate(1, 0, -1).Format(DateMask),
			Count:     dataPoint.Duration}
	}
	// We will always have our latest 'bar' open, since we only add it once we reach a date that is not equal to the latest bar
	bars = append(bars, bar)
	return
}

func generateBarsByMonth(data []Data, start *time.Time) (bars []ChartBar, err error) {
	bars = make([]ChartBar, 0)
	currDate, err := time.Parse(MonthMask, start.Format(MonthMask))
	if err != nil {
		return nil, fmt.Errorf("generateByMonth(): %v", err)
	}

	bar := ChartBar{
		StartDate: currDate.Format(DateMask),
		EndDate:   currDate.AddDate(0, 1, -1).Format(DateMask)}
	for _, dataPoint := range data {
		t, err := time.Parse(MonthMask, dataPoint.Day)
		if err != nil {
			return nil, fmt.Errorf("generateByMonth(): %v", err)
		}
		month := t.Format(MonthMask)
		if currDate.Format(MonthMask) == month {
			bar.Count += dataPoint.Duration
			continue
		}
		bars = append(bars, bar)
		currDate = currDate.AddDate(0, 1, 0)
		for {
			if currDate.Format(MonthMask) == month {
				break
			}
			bars = append(bars, ChartBar{
				StartDate: currDate.Format(DateMask),
				EndDate:   currDate.AddDate(0, 1, -1).Format(DateMask)})
			currDate = currDate.AddDate(0, 1, 0)
		}
		bar = ChartBar{
			StartDate: currDate.Format(DateMask),
			EndDate:   currDate.AddDate(0, 1, -1).Format(DateMask),
			Count:     dataPoint.Duration}
	}
	// We will always have our latest 'bar' open, since we only add it once we reach a date that is not equal to the latest bar
	bars = append(bars, bar)
	return
}

func generateBarsByWeek(data []Data, start *time.Time) (bars []ChartBar, err error) {
	bars = make([]ChartBar, 0)
	currDate := *start
	bar := ChartBar{
		StartDate: currDate.Format(DateMask),
		EndDate:   currDate.AddDate(0, 0, 6).Format(DateMask)}
	for _, dataPoint := range data {
		day, err := time.Parse(DateMask, dataPoint.Day)
		if err != nil {
			return nil, fmt.Errorf("generateByWeek(): %v", err)
		}
		if !day.After(currDate.AddDate(0, 0, 6)) {
			bar.Count += dataPoint.Duration
			continue
		}
		bars = append(bars, bar)
		currDate = currDate.AddDate(0, 0, 7)
		for {
			if !day.After(currDate.AddDate(0, 0, 6)) {
				break
			}
			bars = append(bars, ChartBar{StartDate: currDate.Format(DateMask), EndDate: currDate.AddDate(0, 0, 6).Format(DateMask)})
			currDate = currDate.AddDate(0, 0, 7)
		}
		bar = ChartBar{
			StartDate: currDate.Format(DateMask),
			EndDate:   currDate.AddDate(0, 0, 6).Format(DateMask),
			Count:     dataPoint.Duration}
	}
	// We will always have our latest 'bar' open, since we only add it once we reach a date that is not equal to the latest bar
	bars = append(bars, bar)
	return
}

func generateBarsByDay(data []Data, start *time.Time) (bars []ChartBar, err error) {
	bars = make([]ChartBar, 0)
	// We will always have our latest 'bar' open, since we only add it once we reach a date that is not equal to the latest bar
	currDate := *start
	bar := ChartBar{
		StartDate: currDate.Format(DateMask),
		EndDate:   currDate.Format(DateMask)}
	for _, dataPoint := range data {
		day := dataPoint.Day
		if currDate.Format(DateMask) == day {
			bar.Count += dataPoint.Duration
			continue
		}
		bars = append(bars, bar)
		currDate = currDate.AddDate(0, 0, 1)
		for {
			if currDate.Format(DateMask) == day {
				break
			}
			bars = append(bars, ChartBar{StartDate: currDate.Format(DateMask), EndDate: currDate.Format(DateMask)})
			currDate = currDate.AddDate(0, 0, 1)
		}
		bar = ChartBar{
			StartDate: currDate.Format(DateMask),
			EndDate:   currDate.Format(DateMask),
			Count:     dataPoint.Duration}
	}
	bars = append(bars, bar)
	return
}

func generateBarsByEntry(data []Data, _ *time.Time) (bars []ChartBar, err error) {
	bars = make([]ChartBar, len(data))
	for _, dataPoint := range data {
		bars = append(bars, ChartBar{Count: dataPoint.Duration})
	}
	return
}