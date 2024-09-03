package chart

const (
	DateMask = "2006-01-02"
	YearMask = "2006"
	MonthMask = "01"
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