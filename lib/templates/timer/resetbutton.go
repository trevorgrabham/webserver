package timer

import (
	"fmt"

	"github.com/trevorgrabham/webserver/webserver/lib/svg"
)

func ResetButton(id string, classes, htmx []string) string {
	return fmt.Sprintf(`
  	<p
    	id="timer-display"
    	_="on startTimer call startTimer() then 
        	repeat until event stopTimer
          	call updateTimer()
          	wait 1s
        	end"
    	style="font-size: 5em"
  	>
    	00:00
  	</p>
  	<div id="timer-buttons-container">
    	<button
      	class="svg-wrapper"
      	_="on click send startTimer to #timer-display"
      	hx-get="/startTimer"
      	hx-target="#timer-buttons-container"
    	>
				%s
    	</button>
  	</div>
	`, svg.SVG(id, classes, htmx, svg.Start))
}