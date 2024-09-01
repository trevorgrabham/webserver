package templateutil

import (
	"text/template"

	"github.com/trevorgrabham/webserver/webserver/lib/profile"
)

func NavFuncMap() template.FuncMap {
	return template.FuncMap{
		"defaultPicNeeded": profile.GetProfilePic,
	}
}