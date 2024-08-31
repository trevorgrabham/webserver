package profile

import "html/template"

var profileTemplate = `
	<div id="profile-section">
		<div id="profile-pic-container">` +
	showPicTemplate +
	`</div>
		<div id="profile-details-container">` +
	showNameTemplate +
	`</div>
		<div id="profile-details-container">` +
	showEmailTemplate +
	`</div>
	</div>
`

var Profile = template.Must(template.New("profile-template").Funcs(template.FuncMap{"defaultPicNeeded": DefaultPicNeeded}).Parse(profileTemplate))