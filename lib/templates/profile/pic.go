package profile

import (
	"html/template"

	"github.com/trevorgrabham/webserver/webserver/lib/profile"
)

const showPicTemplate = `
	<div id="profile-pic-img-container">
		<img src="/imgs/{{defaultPicNeeded .ID}}" />
	</div>
	<div 
		id="profile-edit-pic-button"
		class="profile-edit-button"
		hx-get="/profile/editPic"
		hx-target="#profile-pic-container"
	>
		Edit
	</div>`

const editPicTemplate = `
	<input
		id="profile-pic-input"
		name="pic"
		type="file"
		accept="image/jpeg, image/jpg, image/png"
	/>
	<div 
		id="profile-edit-pic-button"
		class="profile-edit-button"
		hx-encoding="multipart/form-data"
		hx-post="/profile/savePic"
		hx-target="#profile-pic-container"
		hx-include="#profile-pic-input"
	>
		Save
	</div>
	<div id="profile-pic-error-container">
		{{range .Errors}}
			<div class="profile-error">
				{{.}}
			</div>
		{{end}}
	</div>`

var ShowPic = template.Must(template.New("show-pic-template").Funcs(template.FuncMap{"defaultPicNeeded": profile.GetProfilePic}).Parse(showPicTemplate)) 

var EditPic = template.Must(template.New("show-pic-template").Parse(editPicTemplate))