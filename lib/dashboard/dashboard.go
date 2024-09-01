package dashboard

import (
	"html/template"

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
	SwapOOB 	template.HTMLAttr
}
