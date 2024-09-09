package timer

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func ParseTimer(time string) (hours, mins int, err error) {
	splitTime := strings.Split(time, ":")
	switch len(splitTime) {
	case 2: // MM:SS
		mins, err = strconv.Atoi(splitTime[0])
		return
	case 3: // HH:MM:SS
		hours, err = strconv.Atoi(splitTime[0])
		if err != nil { return 0, 0, err }
		mins, err = strconv.Atoi(splitTime[1])
		return
	default: // Shouldn't ever hit here
		return 0, 0, fmt.Errorf("ParseTimer(%v): Unrecognized time format", time)
	}
}

func StartingTime(mins int) (startDate time.Time, err error) {
	delta, err := time.ParseDuration(fmt.Sprintf("-%dm", mins))
	if err != nil { return }
	startDate = time.Now().Add(delta)
	return
}