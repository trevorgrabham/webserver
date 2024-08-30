package nav

import (
	"fmt"
	"strings"

	"github.com/trevorgrabham/webserver/webserver/lib/svg"
)

var profilePicDefault = svg.SVG("user-profile", nil, []string{`hx-get="/profilePage`, `hx-target="body"`}, svg.Profile)

func Nav(src, id string, classes, htmx []string) string {
  return fmt.Sprintf(`<nav id="main-nav" class="nav-bar">
    <ul id="main-nav-list">
      <li class="main-nav-list-item">
        %s
    	</li>
    </ul>
  </nav>`, ProfilePicCustomImg(src, id, classes, htmx))
}

func ProfilePicCustomImg(src, id string, classes, htmx []string) string {
  switch src {
  case "":
    var srcString, idString, classString, htmxString string
    srcString = fmt.Sprintf(`src="%s"`, src)
    if id != "" { idString = fmt.Sprintf(`id="%s"`, id) }
    if classes != nil { classString = `class="` + strings.Join(classes, " ") }
    if htmx != nil { htmxString = strings.Join(htmx, "\n") }
    return fmt.Sprintf(`
      <img 
        %s
        %s
        %s
        %s
      >`, srcString, idString, classString, htmxString)
  default:
    return profilePicDefault
  }
}