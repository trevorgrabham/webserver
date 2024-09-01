package tag

type TagMetaData struct {
	ID       int64
	Tag      string
	Count    int64
	MaxCount int64
}

type Tags []TagMetaData

type TagSummaryData struct {
	Tags
	TotalCount int64
	MaxCount   int64
}

// Does not do any pre-processing, so is case-sensitive and does not trim white-space
func (t Tags) Contains(s string) bool {
	for _, tag := range t {
		if tag.Tag == s {
			return true
		}
	}
	return false
}