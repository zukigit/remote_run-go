package lib

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const INFO = 1
const ERR = 2

var spinner = []rune{'|', '/', '-', '\\'}

func Get_formatted_time() string {
	currentTime := time.Now()
	return currentTime.Format("20060102150405")
}

func Formatted_log(level int, unfmt string, arg ...any) string {
	log := fmt.Sprintf(unfmt, arg...)
	formattedTime := Get_formatted_time()

	switch level {
	case INFO:
		log = formattedTime + ", [INFO] " + log
	case ERR:
		log = formattedTime + ", [ERROR] " + log
	}

	return log
}

func Get_log_filename() string {
	return fmt.Sprintf("%s.log", Get_formatted_time())
}

func Logi(log string, filename string) {
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}

	sub_dir := filepath.Join(currentDir, "logs")
	file_path := filepath.Join(sub_dir, filename)

	if _, err := os.Stat(sub_dir); os.IsNotExist(err) {
		err = os.Mkdir(sub_dir, 0755) // Create the directory with read/write permissions
		if err != nil {
			fmt.Println("Error:", err.Error())
			os.Exit(1)
		}
	}

	file, err := os.OpenFile(file_path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}
	defer file.Close()

	if _, err := file.WriteString(log); err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}
}

func Spinner_log(index int, log string) {
	fmt.Printf("\r%s %c", log, spinner[index%len(spinner)])
}
