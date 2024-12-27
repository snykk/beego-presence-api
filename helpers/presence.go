package helpers

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/snykk/beego-presence-api/constants"
)

// DeterminePresenceStatus determines the presence status based on the given presence type, schedule times, current time, and late threshold.
func DeterminePresenceStatus(presenceType, scheduleInTime, scheduleOutTime string, currentTime time.Time, lateThreshold int) (string, error) {
	var scheduleTime string

	// Determine schedule time based on presence type
	switch presenceType {
	case constants.PresenceTypeIn:
		scheduleTime = scheduleInTime
	case constants.PresenceTypeOut:
		scheduleTime = scheduleOutTime
	default:
		return "", fmt.Errorf("invalid presence type: %s", presenceType)
	}

	// Parse schedule time
	timeParts := strings.Split(scheduleTime, ":")
	if len(timeParts) != 3 {
		return "", fmt.Errorf("invalid schedule time format: %s", scheduleTime)
	}

	hour, err := strconv.Atoi(timeParts[0])
	if err != nil {
		return "", fmt.Errorf("invalid hour in schedule time: %s", err)
	}
	minute, err := strconv.Atoi(timeParts[1])
	if err != nil {
		return "", fmt.Errorf("invalid minute in schedule time: %s", err)
	}
	second, err := strconv.Atoi(timeParts[2])
	if err != nil {
		return "", fmt.Errorf("invalid second in schedule time: %s", err)
	}

	// Ensure the schedule time is in the same timezone as currentTime
	location := currentTime.Location()
	scheduleTimeParsed := time.Date(
		currentTime.Year(), currentTime.Month(), currentTime.Day(),
		hour, minute, second, 0, location,
	)

	// Add late threshold only for "in" type
	var thresholdTime time.Time
	if presenceType == constants.PresenceTypeIn {
		thresholdTime = scheduleTimeParsed.Add(time.Minute * time.Duration(lateThreshold))
	} else {
		thresholdTime = scheduleTimeParsed
	}

	// Determine status
	if presenceType == constants.PresenceTypeIn && currentTime.After(thresholdTime) {
		return constants.PresenceStatusLate, nil
	} else if presenceType == constants.PresenceTypeOut && currentTime.Before(thresholdTime) {
		return constants.PresenceStatusEarly, nil
	}

	return constants.PresenceStatusOnTime, nil
}
