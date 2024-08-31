package nav

import (
	"html/template"

	"github.com/trevorgrabham/webserver/webserver/lib/profile"
)

const navTemplate = `
  <ul id="main-nav-list">
    <li 
      class="main-nav-list-item"
      hx-get="/profile"
      hx-target="body"
    >
      <img src="/imgs/{{defaultPicNeeded .ID}}" />
    </li>
  </ul>`

var Nav = template.Must(template.New("nav-template").Funcs(template.FuncMap{"defaultPicNeeded": profile.GetProfilePic}).Parse(navTemplate))