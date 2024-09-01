package svg

import "html/template"

const StopSVGTemplate = `
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
	</svg>`

var StopSVG = template.Must(template.New("stop-svg-template").Parse(StopSVGTemplate))