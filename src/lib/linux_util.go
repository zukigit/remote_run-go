package lib

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func Ja_set_config_linux(key string, value string, config_file_path string) error {
	cmd := fmt.Sprintf(`sed -i 's/^#*\(%s=\).*/\1%s/' %s`, key, value, config_file_path)

	_, err := Ssh_exec_to_str(cmd)

	return err
}

// To use this function, you must have jobarg_agentd default filepath.
func Ja_set_agent_config_linux(key string, value string) error {
	return Ja_set_config_linux(key, value, "/etc/jobarranger/jobarg_agentd.conf")
}

// To use this function, you must have jobarg_server default filepath
func Ja_set_server_config_linux(key string, value string) error {
	return Ja_set_config_linux(key, value, "/etc/jobarranger/jobarg_server.conf")
}

func Restart_jaz_agent_linux() error {
	_, err := Ssh_exec_to_str("systemctl restart jobarg-agentd")

	return err
}

func Restart_jaz_server() error {
	_, err := Ssh_exec_to_str("systemctl restart jobarg-server")

	return err
}

// Wait until it reaches a specified process count
func Job_process_count_check(targetProcessCount int) (int, error) {
	for {
		currentProcessCountStr, err := Ssh_exec_to_str("ps -aux | grep /etc/jobarranger/extendedjob/ | grep -v grep | wc -l")

		if err != nil {
			fmt.Println("Failed to obtain the process count.")
			return -1, err
		}

		// Check the current job process count if it reaches the specified count
		currentProcessCountStr = strings.TrimSpace(currentProcessCountStr)

		currentProcessCount, err := strconv.Atoi(currentProcessCountStr)

		if err != nil {
			fmt.Println("Failed to convert the process count from string to int (maybe due to invalid process count value).")
			return -2, err
		}

		if currentProcessCount == targetProcessCount {
			return currentProcessCount, nil
		}

		time.Sleep(1 * time.Second)
	}

}

// check for zombie process by finding defunct process
func CheckZombieProcess() (int, error) {
	currentProcessCountStr, err := Ssh_exec_to_str("ps -aux | grep defunct | grep -v grep | wc -l")

	if err != nil {
		fmt.Println("Failed to obtain the process count.")
		return -1, err
	}

	// Check the current job process count if it reaches the specified count
	currentProcessCountStr = strings.TrimSpace(currentProcessCountStr)

	currentProcessCount, err := strconv.Atoi(currentProcessCountStr)

	if err != nil {
		fmt.Println("Failed to convert the process count from string to int (maybe due to invalid process count value).")
		return -2, err
	}

	return currentProcessCount, nil

}
