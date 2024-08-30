package svg

import (
	"fmt"
	"strings"

	templates "github.com/trevorgrabham/webserver/webserver/lib/templates/svg"
)

type SVGType int 

const (
	Start SVGType = iota
	Stop 
	Pause 
	Plus
	Success
	Cancel
	Remove
	Profile
)

var svgSource map[SVGType]string = map[SVGType]string{
	Start: templates.StartSVG,
	Stop: templates.StopSVG,
	Pause: templates.PauseSVG,
	Plus: templates.PlusSVG,
	Success: templates.SuccessSVG,
	Cancel: templates.CancelSVG,
	Remove: templates.RemoveSVG,
	Profile: templates.ProfileSVG,
}

func SVG(id string, classes []string, htmx []string, svgType SVGType) string {
	var idString, classString, htmxString string
	if id != "" {
		idString = fmt.Sprintf(`id="%s"`, id)
	}
	if classes != nil {
		classString = fmt.Sprintf(`class="%s"`, strings.Join(classes, " "))
	}
	if htmx != nil {
		htmxString = strings.Join(htmx, "\n")
	}
	switch svgType {
	case Start:
		return fmt.Sprintf(svgSource[Start], idString, classString, htmxString)
	case Stop:
		return fmt.Sprintf(svgSource[Stop], idString, classString, htmxString)
	case Pause:
		return fmt.Sprintf(svgSource[Pause], idString, classString, htmxString)
	case Plus:
		return fmt.Sprintf(svgSource[Plus], idString, classString, htmxString)
	case Success:
		return fmt.Sprintf(svgSource[Success], idString, classString, htmxString)
	case Cancel:
		return fmt.Sprintf(svgSource[Cancel], idString, classString, htmxString)
	case Remove:
		return fmt.Sprintf(svgSource[Remove], idString, classString, htmxString)
	case Profile:
		return fmt.Sprintf(svgSource[Profile], idString, classString, htmxString)
	default:
		return ""
	}
}