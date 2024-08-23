package util

import (
	"fmt"
	"strings"
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
)

var svgSource map[SVGType]string = map[SVGType]string{
	Start:
  	`<svg
			%s
			%s
			%s
    	fill="#000000"
    	height="71px"
    	width="71px"
    	version="1.1"
    	xmlns="http://www.w3.org/2000/svg"
    	xmlns:xlink="http://www.w3.org/1999/xlink"
    	viewBox="0 0 60 60"
    	xml:space="preserve"
  	>
			<!-- svg source https://www.svgrepo.com/svg/13672/play-button -->
    	<g>
      	<path
        	d="M45.563,29.174l-22-15c-0.307-0.208-0.703-0.231-1.031-0.058C22.205,14.289,22,14.629,22,15v30c0,0.371,0.205,0.711,0.533,0.884C22.679,45.962,22.84,46,23,46c0.197,0,0.394-0.059,0.563-0.174l22-15C45.836,30.64,46,30.331,46,30S45.836,29.36,45.563,29.174z M24,43.107V16.893L43.225,30L24,43.107z"
      	/>
      	<path
        	d="M30,0C13.458,0,0,13.458,0,30s13.458,30,30,30s30-13.458,30-30S46.542,0,30,0z M30,58C14.561,58,2,45.439,2,30S14.561,2,30,2s28,12.561,28,28S45.439,58,30,58z"
      	/>
    	</g>
  	</svg>`,
	Stop:
		`<svg 
			%s
			%s
			%s
			fill="#000000" 
			height="71px" 
			width="71px" 
			version="1.1" 
			xmlns="http://www.w3.org/2000/svg" 
			xmlns:xlink="http://www.w3.org/1999/xlink" 
			viewBox="0 0 60 60" 
			xml:space="preserve"
		>
			<!-- svg from https://www.svgrepo.com/svg/125999/stop -->
			<g id="SVGRepo_bgCarrier" stroke-width="0"></g>
			<g id="SVGRepo_tracerCarrier" stroke-linecap="round" stroke-linejoin="round"></g>
			<g id="SVGRepo_iconCarrier"> 
				<g> 
					<path d="M16,44h28V16H16V44z M18,18h24v24H18V18z"></path> 
					<path d="M30,0C13.458,0,0,13.458,0,30s13.458,30,30,30s30-13.458,30-30S46.542,0,30,0z M30,58C14.561,58,2,45.439,2,30 S14.561,2,30,2s28,12.561,28,28S45.439,58,30,58z"></path> 
				</g> 
			</g>
		</svg>`,
	Pause: 
		`<svg 
			%s
			%s
			%s
			fill="#000000" 
			height="71px" 
			width="71px" 
			version="1.1" 
			xmlns="http://www.w3.org/2000/svg" 
			xmlns:xlink="http://www.w3.org/1999/xlink" 
			viewBox="0 0 60 60" 
			xml:space="preserve"
		>
	 		<!-- svg source https://www.svgrepo.com/svg/83992/pause -->
			<g id="SVGRepo_bgCarrier" stroke-width="0"></g>
			<g id="SVGRepo_tracerCarrier" stroke-linecap="round" stroke-linejoin="round"></g>
			<g id="SVGRepo_iconCarrier"> 
				<g> 
					<path d="M30,0C13.458,0,0,13.458,0,30s13.458,30,30,30s30-13.458,30-30S46.542,0,30,0z M30,58C14.561,58,2,45.439,2,30 S14.561,2,30,2s28,12.561,28,28S45.439,58,30,58z"></path> 
					<path d="M33,46h8V14h-8V46z M35,16h4v28h-4V16z"></path> 
					<path d="M19,46h8V14h-8V46z M21,16h4v28h-4V16z"></path> 
				</g> 
			</g>
		</svg>`,
	Plus:
		`<svg 
			%s
			%s
			%s
			fill="#000000" 
			height="71px" 
			width="71px" 
			version="1.1" 
			xmlns="http://www.w3.org/2000/svg" 
			xmlns:xlink="http://www.w3.org/1999/xlink" 
			viewBox="0 0 52 52" 
			xml:space="preserve"
		>
			<!-- svg source https://www.svgrepo.com/svg/158372/plus -->
			<g id="SVGRepo_bgCarrier" stroke-width="0"></g>
			<g id="SVGRepo_tracerCarrier" stroke-linecap="round" stroke-linejoin="round"></g>
			<g id="SVGRepo_iconCarrier"> 
				<g> 
					<path d="M26,0C11.664,0,0,11.663,0,26s11.664,26,26,26s26-11.663,26-26S40.336,0,26,0z M26,50C12.767,50,2,39.233,2,26 S12.767,2,26,2s24,10.767,24,24S39.233,50,26,50z"></path> 
					<path d="M38.5,25H27V14c0-0.553-0.448-1-1-1s-1,0.447-1,1v11H13.5c-0.552,0-1,0.447-1,1s0.448,1,1,1H25v12c0,0.553,0.448,1,1,1 s1-0.447,1-1V27h11.5c0.552,0,1-0.447,1-1S39.052,25,38.5,25z"></path> 
				</g> 
			</g>
		</svg>`, 
	Success:
		`<svg 
		 	%s
			%s
			%s
			fill="#000000" 
			height="71px" 
			width="71px" 
			version="1.1" 
			xmlns="http://www.w3.org/2000/svg" 
			xmlns:xlink="http://www.w3.org/1999/xlink" 
			viewBox="0 0 52 52" 
			xml:space="preserve"
		>
			<!-- svg source https://www.svgrepo.com/svg/13679/success -->
			<g id="SVGRepo_bgCarrier" stroke-width="0"></g>
			<g id="SVGRepo_tracerCarrier" stroke-linecap="round" stroke-linejoin="round"></g>
			<g id="SVGRepo_iconCarrier"> 
				<g> 
					<path d="M26,0C11.664,0,0,11.663,0,26s11.664,26,26,26s26-11.663,26-26S40.336,0,26,0z M26,50C12.767,50,2,39.233,2,26 S12.767,2,26,2s24,10.767,24,24S39.233,50,26,50z"></path> 
					<path d="M38.252,15.336l-15.369,17.29l-9.259-7.407c-0.43-0.345-1.061-0.274-1.405,0.156c-0.345,0.432-0.275,1.061,0.156,1.406 l10,8C22.559,34.928,22.78,35,23,35c0.276,0,0.551-0.114,0.748-0.336l16-18c0.367-0.412,0.33-1.045-0.083-1.411 C39.251,14.885,38.62,14.922,38.252,15.336z"></path> 
				</g> 
			</g>
		</svg>`,
	Cancel:
		`<svg 
			%s
			%s
			%s
			fill="#000000" 
			height="71px" 
			width="71px" 
			version="1.1" 
			xmlns="http://www.w3.org/2000/svg" 
			xmlns:xlink="http://www.w3.org/1999/xlink" 
			viewBox="0 0 52 52" 
			xml:space="preserve"
		>
			<!-- svg source https://www.svgrepo.com/svg/138890/error -->
			<g id="SVGRepo_bgCarrier" stroke-width="0"></g>
			<g id="SVGRepo_tracerCarrier" stroke-linecap="round" stroke-linejoin="round"></g>
			<g id="SVGRepo_iconCarrier"> 
				<g> 
					<path d="M26,0C11.664,0,0,11.663,0,26s11.664,26,26,26s26-11.663,26-26S40.336,0,26,0z M26,50C12.767,50,2,39.233,2,26 S12.767,2,26,2s24,10.767,24,24S39.233,50,26,50z"></path> 
					<path d="M35.707,16.293c-0.391-0.391-1.023-0.391-1.414,0L26,24.586l-8.293-8.293c-0.391-0.391-1.023-0.391-1.414,0 s-0.391,1.023,0,1.414L24.586,26l-8.293,8.293c-0.391,0.391-0.391,1.023,0,1.414C16.488,35.902,16.744,36,17,36 s0.512-0.098,0.707-0.293L26,27.414l8.293,8.293C34.488,35.902,34.744,36,35,36s0.512-0.098,0.707-0.293 c0.391-0.391,0.391-1.023,0-1.414L27.414,26l8.293-8.293C36.098,17.316,36.098,16.684,35.707,16.293z"></path> 
				</g> 
			</g>
		</svg>`,
	Remove:
		`<svg 
		 	%s
		 	%s
		 	%s
			fill="#000000" 
			height="71px" 
			width="71px" 
			version="1.1" 
			xmlns="http://www.w3.org/2000/svg" 
			xmlns:xlink="http://www.w3.org/1999/xlink" 
			viewBox="0 0 31.112 31.112" 
			xml:space="preserve"
		>
			<!-- svg source https://www.svgrepo.com/svg/135247/multiply -->
			<g id="SVGRepo_bgCarrier" stroke-width="0"></g>
			<g id="SVGRepo_tracerCarrier" stroke-linecap="round" stroke-linejoin="round"></g>
			<g id="SVGRepo_iconCarrier"> 
				<polygon points="31.112,1.414 29.698,0 15.556,14.142 1.414,0 0,1.414 14.142,15.556 0,29.698 1.414,31.112 15.556,16.97 29.698,31.112 31.112,29.698 16.97,15.556 "></polygon> 
			</g>
		</svg>`,
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
	default:
		return ""
	}
}