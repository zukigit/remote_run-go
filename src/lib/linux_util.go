package lib

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
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

func Stop_jaz_agent_linux() error {
	_, err := Ssh_exec_to_str("systemctl stop jobarg-agentd")

	return err
}

func Restart_jaz_server() error {
	_, err := Ssh_exec_to_str("systemctl restart jobarg-server")

	return err
}

func Stop_jaz_server() error {
	_, err := Ssh_exec_to_str("systemctl stop jobarg-server")

	return err
}

// Wait until it reaches a specified process count
//
// Parameters:
//   - targetProcessCount: target process count to be reached
//   - timeoutDuration: timeout (minutes) for the process count checking
//   - client: ssh client
//
// Returns:
//
//	error message when the process does not reach the target count
func JobProcessCountCheck(targetProcessCount int, timeoutDuration int, client *ssh.Client) error {
	// set timeout
	timeout := time.After(time.Duration(timeoutDuration) * time.Minute)

	for {
		select {
		case <-timeout:
			return fmt.Errorf("timeout after %d minutes", timeoutDuration)
		default:
			currentProcessCountStr, err := GetOutputStrFromSSHCommand(client, "ps -aux | grep /etc/jobarranger/extendedjob/ | grep -v grep | wc -l")

			if err != nil {
				return fmt.Errorf("failed to obtain process count: %s", err.Error())
			}

			// Check the current job process count if it reaches the specified count
			currentProcessCountStr = strings.TrimSpace(currentProcessCountStr)

			currentProcessCount, err := strconv.Atoi(currentProcessCountStr)

			if err != nil {
				return fmt.Errorf("failed to convert the process count from string to int: %s", err.Error())
			}

			if currentProcessCount == targetProcessCount {
				return nil
			}

			time.Sleep(1 * time.Second)
		}

	}

}

// Check for zombie process by finding defunct process
// Parameters:
//   - timeoutDuration: timeout (minutes) for the process count checking
//   - client: ssh client
//
// Returns:
//
//	0 if there is no zombie process
//	zombie process count if timeout occurred
//	-1 if it times out
//	-2 if it fails to convert process count from string to int
func CheckZombieProcess(timeoutDuration int, client *ssh.Client) (int, error) {

	// set timeout
	timeout := time.After(time.Duration(timeoutDuration) * time.Minute)
	currentProcessCount := -1

	for {
		select {
		case <-timeout:
			return currentProcessCount, fmt.Errorf("timeout after %d minutes", timeoutDuration)
		default:
			currentProcessCountStr, err := GetOutputStrFromSSHCommand(client, "ps -aux | grep defunct | grep -v grep | wc -l")

			if err != nil {
				return -1, fmt.Errorf("failed to obtain the process count: %s", err.Error())
			}

			// Check the current job process count if it reaches the specified count
			currentProcessCountStr = strings.TrimSpace(currentProcessCountStr)

			currentProcessCount, err := strconv.Atoi(currentProcessCountStr)

			if err != nil {
				return -2, fmt.Errorf("failed to convert the process count from string to int: %s", err.Error())
			}

			if currentProcessCount == 0 {
				return currentProcessCount, nil
			}

			time.Sleep(1 * time.Second)

		}
	}

}

// To use this function, your jobarranger agent's TmpDir must be default(TmpDir=/var/lib/jobarranger/tmp)
func Cleanup_agent_linux() error {
	_, err := Ssh_exec_to_str("rm -rf /var/lib/jobarranger/tmp/*")

	return err
}

// To use this function, your jobarranger agent's TmpDir must be default(TmpDir=/var/lib/jobarranger/tmp).
//
// Jobarg_cleanup() cleans jobarg-server and jobarg-agentd(linux) data.
// Since this is testcase utility funtion, you must use it in testcase function.
func Jobarg_cleanup_linux() error {
	if err := Stop_jaz_server(); err != nil {
		return fmt.Errorf("failed to stop JAZ server: %w", err)
	}
	if err := Stop_jaz_agent_linux(); err != nil {
		return fmt.Errorf("failed to stop JAZ agent: %w", err)
	}
	if _, err := DBexec("delete from ja_run_jobnet_table;"); err != nil {
		return fmt.Errorf("failed to execute DB command: %w", err)
	}
	if err := Cleanup_agent_linux(); err != nil {
		return fmt.Errorf("failed to clean up agent: %w", err)
	}
	if err := Restart_jaz_server(); err != nil {
		return fmt.Errorf("failed to start JAZ server: %w", err)
	}
	if err := Restart_jaz_agent_linux(); err != nil {
		return fmt.Errorf("failed to start JAZ server: %w", err)
	}

	return nil
}
