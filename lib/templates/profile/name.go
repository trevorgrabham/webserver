package profile

import (
	"html/template"
)

const showNameTemplate = `
	<div id="profile-name-row" class="profile-input-row">
		<input 
			id="profile-name-input"
			name="name"
			type="text"
			readonly="true"
			value="{{.Name}}"
		/>
		<div 
			id="profile-edit-name-button"
			class="profile-edit-button"
			hx-get="/profile/editName"
			hx-target="#profile-name-row"
		>
			Edit
		</div>
	</div>`

const editNameTemplate = `
	<div id="profile-name-row" class="profile-input-row">
		<input 
			id="profile-name-input"
			name="name"
			type="text"
			placeholder="{{.Name}}"
		/>
		<div 
			id="profile-edit-name-button"
			class="profile-edit-button"
			hx-post="/profile/saveName"
			hx-target="#profile-name-row"
			hx-include="#profile-name-input"
		>
			Save
		</div>
		<div id="profile-name-error-container">
			<div class="profile-error">
				{{range .Errors}}
					{{.}}
				{{end}}
			</div>
		</div>
	</div>`

var ShowName = template.Must(template.New("show-name-template").Parse(showNameTemplate))

var EditName = template.Must(template.New("edit-name-template").Parse(editNameTemplate))