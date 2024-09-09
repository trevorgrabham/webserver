package dashboard

import (
	"fmt"

	"github.com/trevorgrabham/webserver/webserver/lib/tag"
)

type ActivityMetaData struct {
	ID          int64
	UserID      int64
	Duration    int64
	Description string
	Day         string
	tag.Tags
}

type CardMetaData struct {
	Activities []ActivityMetaData
	tag.Tags
	TotalMins int64
	Day       string
	SwapOOB 	bool
}

func FormatTotalMin(t int64) string {
	if t > 60 {
		return fmt.Sprintf("%dh%dm", t/60, t%60)
	}
	return fmt.Sprintf("%dm", t%60)
}