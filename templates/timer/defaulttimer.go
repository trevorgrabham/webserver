package timer

import (
	"html/template"
)

var TimerTemplate = `
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
    `+StartButtonTemplate+`
	</div>`

var DefaultTimer = template.Must(template.New("default-timer-template").Parse(TimerTemplate))