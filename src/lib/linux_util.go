package lib

import (
	"fmt"
	"path/filepath"
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

func Disable_jaz_server() error {
	_, err := Ssh_exec_to_str("systemctl disable jobarg-server")

	return err
}

func Start_jaz_server() error {
	_, err := Ssh_exec_to_str("systemctl start jobarg-server")

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

// Wait until it reaches a specified sleep process count
//
// Parameters:
//   - targetProcessCount: target process count to be reached
//   - timeoutDuration: timeout (minutes) for the process count checking
//   - client: ssh client
//
// Returns:
//
//	error message when the process does not reach the target count
func JobSleepProcessCountCheck(targetProcessCount int, timeoutDuration int, client *ssh.Client) error {
	// set timeout
	timeout := time.After(time.Duration(timeoutDuration) * time.Minute)

	for {
		select {
		case <-timeout:
			return fmt.Errorf("timeout after %d minutes", timeoutDuration)
		default:
			currentProcessCountStr, err := GetOutputStrFromSSHCommand(client, "ps -aux | grep sleep | grep -v grep | wc -l")

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

// ClearLogFile truncates the specified log file on a remote server via SSH.
//
// Parameters:
//   - client: An SSH client to establish the session.
//   - filePath: The path to the log file to clear.
//
// Returns:
//   - error: Returns nil if successful; otherwise, returns an error.
func ClearLogFile(client *ssh.Client, filePath string) error {
	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create SSH session for clearing log: %w", err)
	}
	defer session.Close()

	clearCommand := fmt.Sprintf("truncate -s 0 %s", filePath)
	if err := session.Run(clearCommand); err != nil {
		return fmt.Errorf("failed to clear log file %s: %w", filePath, err)
	}
	return nil
}

// CheckRemoteDirectoryExists verifies if the specified directory exists on a remote server.
// Returns true if the directory exists, false if it does not, or an error if an issue occurs (permissions, etc.).
func CheckRemoteDirectoryExists(client *ssh.Client, path string) (bool, error) {
	session, err := client.NewSession()
	if err != nil {
		return false, fmt.Errorf("failed to create session: %w", err)
	}
	defer session.Close()

	// Check if the directory exists with an SSH command
	cmd := fmt.Sprintf(`if [ -d "%s" ]; then echo "exists"; else echo "notfound"; fi`, path)
	output, err := session.Output(cmd)
	if err != nil {
		return false, fmt.Errorf("failed to execute command: %w", err)
	}

	// Check command output and trim any whitespace or newlines
	outputStr := strings.TrimSpace(string(output))
	if outputStr == "exists" {
		return true, nil
	} else if outputStr == "notfound" {
		return false, nil
	}
	return false, fmt.Errorf("unexpected output: %s", outputStr)
}

// RemoveRemoteDirectory deletes the specified directory and all its contents on a remote server.
// Returns an error if the directory could not be removed, or nil if successful.
func RemoveRemoteDirectory(client *ssh.Client, path string) error {
	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}
	defer session.Close()

	// Remove the directory remotely with SSH command
	cmd := fmt.Sprintf("rm -rf %s", path)
	if err := session.Run(cmd); err != nil {
		return fmt.Errorf("failed to remove directory %s: %w", path, err)
	}

	fmt.Printf("Directory %s successfully removed.\n", path)
	return nil
}

// RemoveAllFilesInDirectory deletes all files within the specified directory on a remote server.
// The directory itself is not removed, only the files within it.
// Returns an error if the operation fails, or nil if successful.
func RemoveAllFilesInDirectory(client *ssh.Client, path string) error {
	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}
	defer session.Close()

	// Remove all files in the directory without removing the directory itself
	cmd := fmt.Sprintf("rm -rf %s/*", path)
	if err := session.Run(cmd); err != nil {
		return fmt.Errorf("failed to remove files in directory %s: %w", path, err)
	}
	return nil
}

// CheckRemoteIndexFileExists checks if any files with a .job extension exist in the specified remote directory.
// Parameters:
//   - client: *ssh.Client, the established SSH client connection.
//   - dirPath: the directory to search for .job files (e.g., "/var/log/jobarranger/job/").
//
// Returns:
//   - bool: true if at least one .job file exists, false otherwise.
//   - error: error if there was an issue accessing the directory.
func CheckRemoteIndexFileExists(client *ssh.Client, dirPath string) (bool, error) {
	session, err := client.NewSession()
	if err != nil {
		return false, fmt.Errorf("failed to create session: %w", err)
	}
	defer session.Close()

	// Command to check for any .job files in the specified directory
	cmd := fmt.Sprintf("ls %s/*.job 1> /dev/null 2>&1 && echo found || echo notfound", dirPath)
	output, err := session.Output(cmd)
	if err != nil {
		return false, err
	}

	// Check command output
	return string(output) == "found\n", nil
}

// FindandTrimServerJobFilePrefix finds the first .job file in the specified directory and returns the name
// after removing the .job extension.
func FindandTrimServerJobFilePrefix(client *ssh.Client, dirPath string) (string, error) {
	session, err := client.NewSession()
	if err != nil {
		return "", fmt.Errorf("failed to create session: %w", err)
	}
	defer session.Close()

	// Command to list all .job files in the directory
	cmd := fmt.Sprintf("find %s -maxdepth 1 -name '*.job'", dirPath)
	output, err := session.Output(cmd)
	if err != nil {
		return "", fmt.Errorf("failed to list .job files in directory %s: %w", dirPath, err)
	}

	// Split the output into lines and return the first .job file found
	files := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(files) > 0 {
		// Extract the file name from the path
		fileName := filepath.Base(files[0])

		// Remove the .job extension to work with the file name
		fileNameWithoutExt := strings.TrimSuffix(fileName, ".job")

		// Return the file name without the .job extension
		return fileNameWithoutExt, nil
	}

	return "", nil // No .job file found
}

// FindandTrimAgentJobFilePrefix finds the first .job file in the specified directory, removes the .job extension,
// and returns the part of the file name that comes before the first underscore.
func FindandTrimAgentJobFilePrefix(client *ssh.Client, dirPath string) (string, error) {
	session, err := client.NewSession()
	if err != nil {
		return "", fmt.Errorf("failed to create session: %w", err)
	}
	defer session.Close()

	// Command to list all .job files in the directory
	cmd := fmt.Sprintf("find %s -maxdepth 1 -name '*.job'", dirPath)
	output, err := session.Output(cmd)
	if err != nil {
		return "", fmt.Errorf("failed to list .job files in directory %s: %w", dirPath, err)
	}

	// Split the output into lines and process the first .job file found
	files := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(files) > 0 && files[0] != "" {
		// Extract the file name from the path
		fileName := filepath.Base(files[0])

		// Remove the .job extension
		fileNameWithoutExt := strings.TrimSuffix(fileName, ".job")

		// Find the index of the first underscore and take only the part before it
		if underscoreIndex := strings.Index(fileNameWithoutExt, "_"); underscoreIndex != -1 {
			return fileNameWithoutExt[:underscoreIndex], nil
		}

		// If there's no underscore, return the whole name
		return fileNameWithoutExt, nil
	}

	return "", fmt.Errorf("no .job file found in directory %s", dirPath)
}

// UpdateDebugLevel updates the DebugLevel setting in jobarg_agentd.conf to the specified level on the remote server
func UpdateDebugLevel(client *ssh.Client, configFilePath string, newLevel int) error {
	// Create a new SSH session
	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}
	defer session.Close()

	// Command to replace DebugLevel value in the configuration file
	cmd := fmt.Sprintf("sudo sed -i 's/^DebugLevel=[0-4]/DebugLevel=%d/' %s", newLevel, configFilePath)

	// Run the command on the remote server
	if err := session.Run(cmd); err != nil {
		return fmt.Errorf("failed to update DebugLevel in %s: %w", configFilePath, err)
	}

	fmt.Printf("Successfully updated DebugLevel to %d in %s \n", newLevel, configFilePath)
	return nil
}

// WaitForPatternInLogFile checks a remote log file for a specific pattern within a given timeout period.
// It polls the file at regular intervals until the pattern is found or the timeout is reached.
//
// Parameters:
// - client: SSH client connection to the remote server.
// - filePath: The path to the log file on the remote server.
// - pattern: The string pattern to search for in the log file.
// - timeout: The maximum duration to wait for the pattern to appear.
// - interval: The duration between each poll of the log file.
//
// Returns:
// - An error if the pattern is not found within the timeout or if any errors occur during SSH command execution.
func WaitForPatternInLogFile(client *ssh.Client, filePath, pattern string, timeout, interval time.Duration) (string, error) {
	timeoutChan := time.After(timeout)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-timeoutChan:
			return "", fmt.Errorf("timeout reached: pattern not found in log file %s", filePath)
		case <-ticker.C:
			session, err := client.NewSession()
			if err != nil {
				return "", fmt.Errorf("failed to create SSH session: %w", err)
			}
			defer session.Close()

			command := fmt.Sprintf("cat %s", filePath)
			output, err := session.Output(command)
			if err != nil {
				return "", fmt.Errorf("failed to execute command: %w", err)
			}

			content := string(output)
			if strings.Contains(content, pattern) {
				fmt.Printf("Pattern found: %s\n", pattern)
				return content, nil // Successfully found the pattern
			}
		}
	}
}
func Delete_agent_log() error {
	_, err := Ssh_exec_to_str("> /var/log/jobarranger/jobarg_agentd.log")

	return err
}

func Delete_server_log() error {
	_, err := Ssh_exec_to_str("> /var/log/jobarranger/jobarg_server.log")

	return err
}
