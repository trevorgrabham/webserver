package profile

import "html/template"

const showEmailTemplate = `
	<div id="profile-email-row" class="profile-input-row">
		<input 
			id="profile-email-input"
			name="email"
			type="text"
			readonly="true"
			value="{{.Email}}"
		/>
		<div 
			id="profile-edit-email-button"
			class="profile-edit-button"
			hx-get="/profile/editEmail"
			hx-target="#profile-email-row"
		>
			Edit
		</div>
	</div>`

const editEmailTemplate = `
	<div id="profile-email-row" class="profile-input-row">
		<input 
			id="profile-email-input"
			name="email"
			type="text"
			placeholder="{{.Email}}"
		/>
		<div 
			id="profile-edit-email-button"
			class="profile-edit-button"
			hx-post="/profile/saveEmail"
			hx-target="#profile-email-row"
			hx-include="#profile-email-input"
		>
			Save
		</div>
		<div id="profile-email-error-container">
			{{range .Errors}}
				<div class="profile-error">
					{{.}}
				</div>
			{{end}}
		</div>
	</div>`

var ShowEmail = template.Must(template.New("show-email-template").Parse(showEmailTemplate))

var EditEmail = template.Must(template.New("edit-email-template").Parse(editEmailTemplate))