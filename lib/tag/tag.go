package tag

import "fmt"

type TagMetaData struct {
	ID       int64
	Tag      string
	Count    int64
	MaxCount int64
}

type Tags []TagMetaData

// Does not do any pre-processing, so is case-sensitive and does not trim white-space
func (t Tags) Contains(s string) bool {
	for _, tag := range t {
		if tag.Tag == s {
			return true
		}
	}
	return false
}

func NormalizeCount(n int64, max int64) string {
	if max == 0 {
		return "100%"
	}
	return fmt.Sprintf("%.2f", float64(n)/float64(max)*100.0) + "%"
}

func GetMaxCount(t Tags) int64 {
	return t[0].MaxCount
}
