package svg

import "html/template"

const RemoveSVGTemplate = `
	<svg 
		{{.ID}}
		{{.Classes}}
		{{.Htmx}}
		{{.Hyperscript}}
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
	</svg>`

var RemoveSVG = template.Must(template.New("remove-svg-template").Parse(RemoveSVGTemplate))