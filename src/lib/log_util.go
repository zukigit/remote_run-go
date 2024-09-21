package lib

import (
	"time"
)

const INFO = 1
const ERR = 2

func Formatted_log(level int, log string) string {
	currentTime := time.Now()
	formattedTime := currentTime.Format("2006-Jan-02 15:04:05.000000")

	switch level {
	case INFO:
		log = formattedTime + ", [INFO] " + log
	case ERR:
		log = formattedTime + ", [ERROR]!!! " + log
	}

	return log
}
