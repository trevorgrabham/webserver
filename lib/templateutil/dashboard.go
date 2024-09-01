package templateutil

import (
	"fmt"
	"html/template"
)

func DashboardFuncMap() template.FuncMap {
	return template.FuncMap{
		"formatTotalMin": formatTotalMin,
	}
}

func formatTotalMin(t int64) string {
	if t > 60 {
		return fmt.Sprintf("%dh%dm", t/60, t%60)
	}
	return fmt.Sprintf("%dm", t%60)
}