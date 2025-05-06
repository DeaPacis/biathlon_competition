package lib

import (
	"log"
	"time"
)

func ParseStringToDuration(inputString string) time.Duration {
	parsedString, err := time.Parse("15:04:05", inputString)
	if err != nil {
		log.Println(err)
		return 0
	}

	duration := time.Duration(
		parsedString.Hour()*int(time.Hour) + parsedString.Minute()*int(time.Minute) +
			parsedString.Second()*int(time.Second) + parsedString.Nanosecond(),
	)
	return duration
}
