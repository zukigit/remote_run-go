package lib

import (
	"fmt"
	"os"
	"time"

	"github.com/zukigit/remote_run-go/src/common"
)

var spinner = []rune{'|', '/', '-', '\\'}

func Get_formatted_time() string {
	currentTime := time.Now()
	return currentTime.Format("20060102150405")
}

func Formatted_log(level int, unfmt string, arg ...any) string {
	log := fmt.Sprintf(unfmt, arg...)
	formattedTime := Get_formatted_time()

	switch level {
	case common.INFO:
		log = formattedTime + " [INFO] " + log
	case common.ERR:
		log = formattedTime + " [ERROR] " + log
	default:
		log = formattedTime + " [UNKNOWN] " + log

	}

	return log
}

func Get_log_filename() string {
	return fmt.Sprintf("%s.log", Get_formatted_time())
}

// Write logs to the log file
func Logi(log string) {
	if common.Log_file == nil {
		fmt.Println("Error: Log_file is nil.")
		os.Exit(1)
	}
	if _, err := common.Log_file.WriteString(log); err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}
}

func Spinner_log(index int, log string) {
	fmt.Printf("\r%s %c", log, spinner[index%len(spinner)])
}
