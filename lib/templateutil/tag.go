package templateutil

import (
	"fmt"
	"html/template"
)

func normalizeCount(n int64, max int64) string {
	if max == 0 {
		return "100%"
	}
	return fmt.Sprintf("%.2f", float64(n)/float64(max)*100.0) + "%"
}

func TagFuncMap() template.FuncMap {
	return template.FuncMap{
		"normalizeCount": normalizeCount,
	}
}
