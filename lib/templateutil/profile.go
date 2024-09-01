package templateutil

import (
	"html/template"

	"github.com/trevorgrabham/webserver/webserver/lib/profile"
)

type profileData struct {
	profile.UserDetails
	Errors					[]string
}

func NewProfileData(user profile.UserDetails, errors []string) profileData {
	return profileData{
		UserDetails: user, 
		Errors: errors,
	}
}

func ProfileFuncMap() template.FuncMap {
	return template.FuncMap{
		"defaultPicNeeded": profile.GetProfilePic,
	}
}