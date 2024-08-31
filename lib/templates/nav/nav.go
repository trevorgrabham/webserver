package nav

import (
	"html/template"

	profiletemplate "github.com/trevorgrabham/webserver/webserver/lib/templates/profile"
)

const navTemplate = `
  <ul id="main-nav-list">
    <li 
      class="main-nav-list-item"
      hx-get="/profile"
      hx-target="body"
    >
      <img src="/imgs/{{defaultPicNeeded .ID .Ext}}" />
    </li>
  </ul>`

var Nav = template.Must(template.New("nav-template").Funcs(template.FuncMap{"defaultPicNeeded": profiletemplate.DefaultPicNeeded}).Parse(navTemplate))