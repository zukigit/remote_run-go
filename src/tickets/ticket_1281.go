package tickets

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/zukigit/remote_run-go/src/common"
	"github.com/zukigit/remote_run-go/src/dao"
	"github.com/zukigit/remote_run-go/src/lib"
	"golang.org/x/crypto/ssh"
)

type Ticket_1281 struct {
	Ticket_no                                   uint
	Ticket_description                          string
	PASSED_count, FAILED_count, MUSTCHECK_count int
	Testcases                                   []dao.TestCase
}

func (t *Ticket_1281) New_testcase(testcase_id uint, testcase_description string) *dao.TestCase {
	return dao.New_testcase(testcase_id, testcase_description)
}

func (t *Ticket_1281) Get_no() uint {
	return t.Ticket_no
}

func (t *Ticket_1281) Set_PASSED_count(passed_count int) {
	t.PASSED_count = passed_count
}

func (t *Ticket_1281) Set_FAILED_count(failed_count int) {
	t.FAILED_count = failed_count
}

func (t *Ticket_1281) Set_MUSTCHECK_count(mustcheck_count int) {
	t.MUSTCHECK_count = mustcheck_count
}

func (t *Ticket_1281) Get_dsctn() string {
	return t.Ticket_description
}

func (t *Ticket_1281) Add_testcase(tc dao.TestCase) {
	t.Testcases = append(t.Testcases, tc)
}

func (t *Ticket_1281) Get_testcases() []dao.TestCase {
	return t.Testcases
}

// Enter your ticket information here
func (t *Ticket_1281) Set_values() {
	t.Ticket_no = 1281 // Enter your ticket id
	t.Ticket_description = "Check job process start"
}

// Add your test case here
func (t *Ticket_1281) Add_testcases() {
	tc_1 := t.New_testcase(1, "Check if Check job process start or not.")
	tc_func := func() common.Testcase_status {
		return CheckJobProcessStartOrExit(tc_1, common.Client)
	}
	tc_1.Set_function(tc_func)
	t.Add_testcase(*tc_1)

	tc_2 := t.New_testcase(2, "Agentless must run successfully in Interact mode. STD_OUT should be correct.")
	tc_func = func() common.Testcase_status {
		if err := lib.Jobarg_enable_jobnet("Icon_1", "agentless_interact"); err != nil {
			lib.Logi(common.LOG_LEVEL_ERR, "Failed to enable jobnet, Error: %s", err)
			return FAILED
		}
		return AgentlessInteractCheckSTDOut("Icon_1", tc_2, common.Client)
	}
	tc_2.Set_function(tc_func)
	t.Add_testcase(*tc_2)

	tc_3 := t.New_testcase(3, "Agentless must run successfully in non-interact mode. STD_OUT should be correct.")
	tc_func = func() common.Testcase_status {
		if err := lib.Jobarg_enable_jobnet("Icon_1", "agentless_no_interact"); err != nil {
			lib.Logi(common.LOG_LEVEL_ERR, "Failed to enable jobnet, Error: %s", err)
			return FAILED
		}
		return AgentlessNoInteractCheckSTDOut("Icon_1", tc_3, common.Client)
	}
	tc_3.Set_function(tc_func)
	t.Add_testcase(*tc_3)

	tc_4 := t.New_testcase(4, "change AgentLessRegex value in jobarg-server.conf for agentless interact mode STD_OUT.")
	tc_func = func() common.Testcase_status {
		if err := lib.Jobarg_enable_jobnet("Icon_1", "agentless_no_interact"); err != nil {
			lib.Logi(common.LOG_LEVEL_ERR, "Failed to enable jobnet, Error: %s", err)
			return FAILED
		}
		return AgentlessInteractRegexCheckSTDOut("Icon_1", tc_4, common.Client)
	}
	tc_4.Set_function(tc_func)
	t.Add_testcase(*tc_4)

	tc_5 := t.New_testcase(5, "0xff character encoding issue.")
	tc_func = func() common.Testcase_status {
		if err := lib.Jobarg_enable_jobnet("Icon_1", "0xff_character_encoding"); err != nil {
			lib.Logi(common.LOG_LEVEL_ERR, "Failed to enable jobnet, Error: %s", err)
			return FAILED
		}
		return CharacterEncodingIssue("Icon_1", tc_5, common.Client)
	}
	tc_5.Set_function(tc_func)
	t.Add_testcase(*tc_5)

	tc_6 := t.New_testcase(6, "Backup thread testing.")
	tc_func = func() common.Testcase_status {
		if err := lib.Jobarg_enable_jobnet("Icon_1", "jobicon_linux"); err != nil {
			lib.Logi(common.LOG_LEVEL_ERR, "Failed to enable jobnet, Error: %s", err)
			return FAILED
		}
		return BackupThread("Icon_1", tc_6, common.Client)
	}
	tc_6.Set_function(tc_func)
	t.Add_testcase(*tc_6)

	tc_7 := t.New_testcase(7, "check icon count for sending status check request for different host through log.")
	tc_func = func() common.Testcase_status {
		return CheckIconCountForSendingStatus("Icon_1", tc_7, common.Client)
	}
	tc_7.Set_function(tc_func)
	t.Add_testcase(*tc_7)
}

func CheckJobProcessStartOrExit(testcase *dao.TestCase, client *ssh.Client) common.Testcase_status {

	var ServerProcessID string
	// var ProcessID string

	_clear_err := lib.ClearLogFile(client, "/var/log/jobarranger/jobarg_server.log")
	if _clear_err != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error clear log: %s", _clear_err))
	}

	_restart_err := lib.Restart_jaz_server()
	if _restart_err != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error jaz restart: %s", _restart_err))
	}
	time.Sleep(10 * time.Second)
	pattern := "[INFO] [JASERVER000004] server #15 started [check job process #1]"

	_check_log, _check_err := lib.WaitForPatternInLogFile(client, "/var/log/jobarranger/jobarg_server.log", pattern, 30, 2)
	if _check_err != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error check log file : %s", _check_err))
	}

	if _check_log != "" {
		// fmt.Println(_check_log)
	}

	// Regular expression to match "ja_purge process execution thread : <number>"
	re := regexp.MustCompile(`ja_check_job process execution thread\s*:\s*(\d+)`) // Match "ja_purge process execution thread : <number>"

	// Use the regex to find the dynamic process ID (e.g., 250622)
	matches := re.FindStringSubmatch(_check_log)
	if len(matches) > 1 {
		ServerProcessID = matches[1] // Extract the dynamic value (e.g., 250622)
		fmt.Printf("Found process ID: %s\n", ServerProcessID)
	} else {
		fmt.Println("No matching process found in the log.")
	}

	_ps_job, _ps_job_err := lib.GetOutputStrFromSSHCommand(client, "ps -aux | grep job")
	if _ps_job_err != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error ps aux job : %s", _ps_job_err))
	}

	if _ps_job != "" {
		// Split the output by lines
		lines := strings.Split(_ps_job, "\n")
		for _, line := range lines {
			if strings.Contains(line, ServerProcessID) {
				// If the PID exists in the line, it means the process is running
				fmt.Printf("Process with PID %s is running: %s\n", ServerProcessID, line)
				return PASSED
			}
		}
		// If no matching process is found
		fmt.Printf("Process with PID %s not found in ps output.\n", ServerProcessID)
	}

	return FAILED
}

func AgentlessInteractCheckSTDOut(jobnetId string, testcase *dao.TestCase, client *ssh.Client) common.Testcase_status {

	hostname_cmd := "hostname"

	hostname, hostname_err := lib.Ssh_exec_to_str(hostname_cmd)
	if hostname_err != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error hostname : %s", hostname_err))
	}
	hostname = strings.TrimSpace(hostname)
	fmt.Println(hostname)

	envs, _ := lib.Get_str_str_map("JA_HOSTNAME", "oss.linux", "JA_CMD", "hostname")
	run_jobnet_id, job_run_err := lib.Jobarg_exec_E(jobnetId, envs)

	if job_run_err != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error running job: %s", job_run_err))
	}

	fmt.Println("Run Jobnet ID:", run_jobnet_id)

	time.Sleep(20 * time.Second) // Sleep for 2 seconds

	_stdOut, _stdOut_err := lib.Jobarg_get_LASTSTDOUT(run_jobnet_id)

	if _stdOut_err != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error get std out : %s", _stdOut_err))
	}

	fmt.Println(_stdOut)

	// Regular expression to match the hostname after 'hostname'
	re := regexp.MustCompile(`hostname\s+([^\s]+)`)

	// Find the match and extract the hostname
	match := re.FindStringSubmatch(_stdOut)
	if len(match) > 1 {
		fmt.Printf("hostname: '%s'\n", hostname)
		fmt.Printf("match[1]: '%s'\n", match[1])

		// Now compare them, using strings.ToLower if you want a case-insensitive comparison
		if strings.EqualFold(hostname, match[1]) {
			return PASSED
		} else {
			fmt.Println("Hostname mismatch:", hostname, match[1])
			return FAILED
		}

	} else {
		fmt.Println("Hostname not found")
	}

	return FAILED
}

func AgentlessNoInteractCheckSTDOut(jobnetId string, testcase *dao.TestCase, client *ssh.Client) common.Testcase_status {

	hostname_cmd := "hostname"

	hostname, hostname_err := lib.Ssh_exec_to_str(hostname_cmd)
	if hostname_err != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error hostname : %s", hostname_err))
	}
	hostname = strings.TrimSpace(hostname)
	fmt.Println(hostname)

	envs, _ := lib.Get_str_str_map("JA_HOSTNAME", "oss.linux", "JA_CMD", "hostname")
	run_jobnet_id, job_run_err := lib.Jobarg_exec_E(jobnetId, envs)

	if job_run_err != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error running job: %s", job_run_err))
	}

	fmt.Println("Run Jobnet ID:", run_jobnet_id)

	time.Sleep(20 * time.Second) // Sleep for 2 seconds

	_stdOut, _stdOut_err := lib.Jobarg_get_LASTSTDOUT(run_jobnet_id)

	if _stdOut_err != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error get std out : %s", _stdOut_err))
	}

	fmt.Println(_stdOut)

	// Now compare them, using strings.ToLower if you want a case-insensitive comparison
	if strings.EqualFold(hostname, _stdOut) {
		return PASSED
	} else {
		fmt.Println("Hostname mismatch:", hostname, _stdOut)
		return FAILED
	}
}

func AgentlessInteractRegexCheckSTDOut(jobnetId string, testcase *dao.TestCase, client *ssh.Client) common.Testcase_status {

	config_cmd := "sudo sed -i -e '$a AgentLessRegex=.' /etc/jobarranger/jobarg_server.conf"

	_, _config_err := lib.Ssh_exec(config_cmd)
	if _config_err != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error set config : %s", _config_err))
	}

	jaz_server_restart_err := lib.Restart_jaz_server()

	if jaz_server_restart_err != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error jaz restart : %s", jaz_server_restart_err))
	}

	hostname_cmd := "hostname"

	hostname, hostname_err := lib.Ssh_exec_to_str(hostname_cmd)
	if hostname_err != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error hostname : %s", hostname_err))
	}
	hostname = strings.TrimSpace(hostname)
	fmt.Println(hostname)

	envs, _ := lib.Get_str_str_map("JA_HOSTNAME", "oss.linux", "JA_CMD", "hostname")
	run_jobnet_id, job_run_err := lib.Jobarg_exec_E(jobnetId, envs)

	if job_run_err != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error running job: %s", job_run_err))
	}

	fmt.Println("Run Jobnet ID:", run_jobnet_id)

	time.Sleep(20 * time.Second) // Sleep for 2 seconds

	_stdOut, _stdOut_err := lib.Jobarg_get_LASTSTDOUT(run_jobnet_id)

	if _stdOut_err != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error get std out : %s", _stdOut_err))
	}

	fmt.Println(_stdOut)

	// Regular expression to match the hostname after 'hostname'
	re := regexp.MustCompile(`hostname\s+([^\s]+)`)

	// Find the match and extract the hostname
	match := re.FindStringSubmatch(_stdOut)
	if len(match) > 1 {
		fmt.Printf("hostname: '%s'\n", hostname)
		fmt.Printf("match[1]: '%s'\n", match[1])

		// Now compare them, using strings.ToLower if you want a case-insensitive comparison
		if strings.EqualFold(hostname, match[1]) {
			return PASSED
		} else {
			fmt.Println("Hostname mismatch:", hostname, match[1])
			return FAILED
		}

	} else {
		fmt.Println("Hostname not found")
	}

	return FAILED
}

func CharacterEncodingIssue(jobnetId string, testcase *dao.TestCase, client *ssh.Client) common.Testcase_status {

	// envs, _ := lib.Get_str_str_map("JA_HOSTNAME", "oss.linux", "JA_CMD", `echo -e "\xFFF"`)
	run_jobnet_id, job_run_err := lib.Jobarg_exec(jobnetId)

	if job_run_err != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error running job: %s", job_run_err))
	}

	time.Sleep(20 * time.Second) // Sleep for 20 seconds
	job_status, job_status_err := lib.Jobarg_get_JA_JOBSTATUS(run_jobnet_id)

	if job_status_err != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error job status : %s", job_status_err))
	}
	fmt.Println(job_status)

	if job_status == "NORMAL" {
		return PASSED
	}

	return FAILED
}

// Function to execute the ls command and extract the day and hour
func extractDayAndHour() (int, int, error) {
	// Run the 'ls -l --time-style=full-iso' command via SSH
	cmd := "cd /var/lib/jobarranger/tmp/close && ls -l --time-style=full-iso | awk '{print $6, $7}'"
	output, err := lib.Ssh_exec(cmd)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to execute ls command: %w", err)
	}

	// Print the raw output for debugging
	fmt.Println("Raw Output:\n", string(output))

	// Split the output into lines
	lines := strings.Split(string(output), "\n")

	// Loop through each line
	for _, line := range lines {
		// Skip empty lines
		if strings.TrimSpace(line) == "" {
			continue
		}

		// Extract the date and time (6th and 7th columns)
		columns := strings.Fields(line)
		if len(columns) < 2 {
			// Only 2 columns expected (date and time)
			continue
		}

		// The 1st column is the date, and the 2nd column is the time
		timestamp := columns[0] + " " + columns[1]

		// Print the timestamp for debugging
		fmt.Println("Parsed Timestamp:", timestamp)

		// Parse the timestamp (date and time in 'YYYY-MM-DD HH:MM:SS' format)
		layout := "2006-01-02 15:04:05"
		parsedTime, err := time.Parse(layout, timestamp)
		if err != nil {
			fmt.Printf("Failed to parse time: %v\n", err)
			continue
		}

		// Extract the day and hour
		day := parsedTime.Day()   // Day of the month
		hour := parsedTime.Hour() // Hour of the day

		// Return the day and hour
		return day, hour, nil
	}

	// If no valid timestamp was found, return an error
	return 0, 0, fmt.Errorf("no valid timestamp found")
}

// Function to get directory contents (files and folders)
func getDirectoryContents() ([]string, error) {
	// Run the 'ls' command to list the files and directories
	cmd := "cd /var/lib/jobarranger/tmp/close && ls"
	output, err := lib.Ssh_exec(cmd) // Use CombinedOutput to capture both stdout and stderr
	if err != nil {
		return nil, fmt.Errorf("failed to execute ls command: %w", err)
	}

	// Convert output to string and split by newline
	contents := strings.Split(string(output), "\n")

	// Filter out any empty strings (to handle trailing newlines)
	var result []string
	for _, item := range contents {
		if item != "" {
			result = append(result, item)
		}
	}

	return result, nil
}

func BackupThread(jobnetId string, testcase *dao.TestCase, client *ssh.Client) common.Testcase_status {

	rm_cmd := "sudo rm -rf /var/lib/jobarranger/tmp/close/*"
	_, rm_err := lib.Ssh_exec(rm_cmd)
	if rm_err != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error remove close dir : %s", rm_err))
	}

	time.Sleep(5 * time.Second) // Sleep for 2 seconds
	envs, _ := lib.Get_str_str_map("JA_HOSTNAME", "oss.linux", "JA_CMD", "hostname")
	_, job_run_err := lib.Jobarg_exec_E(jobnetId, envs)

	if job_run_err != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error running job: %s", job_run_err))
	}

	envs, _ = lib.Get_str_str_map("JA_HOSTNAME", "oss.linux", "JA_CMD", "hostname")
	_, job_run_err = lib.Jobarg_exec_E(jobnetId, envs)

	if job_run_err != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error running job: %s", job_run_err))
	}

	time.Sleep(20 * time.Second) // Sleep for 2 seconds
	day, hr, err := extractDayAndHour()
	if err != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error: %s", err))
	}

	fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "Day : %d", day))
	fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "Hr : %d", hr))

	contents, get_jobfile_and_folder_err := getDirectoryContents()
	if get_jobfile_and_folder_err != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error before job file and folder : %s", get_jobfile_and_folder_err))
	}

	fmt.Println(contents)

	rm_jobFile_cmd := fmt.Sprintf(`cd /var/lib/jobarranger/tmp/close && rm -rf %s`, contents[1])
	fmt.Println(rm_jobFile_cmd)

	_, rm_jobFile_err := lib.Ssh_exec(rm_jobFile_cmd)
	if rm_jobFile_err != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error remove .job file : %s", rm_jobFile_err))
	}

	rm_jobFolder_cmd := fmt.Sprintf(`cd /var/lib/jobarranger/tmp/close && rm -rf %s`, contents[2])
	fmt.Println(rm_jobFolder_cmd)

	_, rm_jobFolder_err := lib.Ssh_exec(rm_jobFolder_cmd)
	if rm_jobFolder_err != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error remove job folder : %s", rm_jobFolder_err))
	}

	time.Sleep(10 * time.Second) // Sleep for 10 seconds

	contents, get_jobfile_and_folder_err = getDirectoryContents()
	if get_jobfile_and_folder_err != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error after job file and folder : %s", get_jobfile_and_folder_err))
	}

	fmt.Println(contents)

	// Check if JaBackupTime already exists in the file
	backupTimeCheckCmd := "grep -q '^\\s*JaBackupTime=' /etc/jobarranger/jobarg_agentd.conf && echo 'Found' || echo 'Not found'"
	backupTimeCheckOutput, err := lib.Ssh_exec_to_str(backupTimeCheckCmd)
	fmt.Println(backupTimeCheckOutput)

	if err != nil || strings.TrimSpace(backupTimeCheckOutput) == "Not found" {
		// If JaBackupTime is not set, run the sed command to set it
		backupTimeCmd := "sudo sed -i '/^# JaBackupTime=24/a JaBackupTime=1' /etc/jobarranger/jobarg_agentd.conf"
		_, backupTimeErr := lib.Ssh_exec(backupTimeCmd)
		if backupTimeErr != nil {
			fmt.Println("Error ja backup time: ", backupTimeErr)
		}
	} else {
		fmt.Println("JaBackupTime already set, skipping command.")
	}

	// Check if JaBackupRunTime already exists in the file
	backupRunTimeCheckCmd := "grep -q '^\\s*JaBackupRunTime=' /etc/jobarranger/jobarg_agentd.conf && echo 'Found' || echo 'Not found'"
	backupRunTimeCheckOutput, err := lib.Ssh_exec_to_str(backupRunTimeCheckCmd)

	fmt.Println(backupRunTimeCheckOutput)

	if err != nil || strings.TrimSpace(backupRunTimeCheckOutput) == "Not found" {
		// If JaBackupRunTime is not set, run the sed command to set it
		backupRunTimeCmd := fmt.Sprintf("sudo sed -i '/^# JaBackupRunTime=24/a JaBackupRunTime=%d' /etc/jobarranger/jobarg_agentd.conf", hr+1)
		_, backupRunTimeErr := lib.Ssh_exec(backupRunTimeCmd)
		if backupRunTimeErr != nil {
			fmt.Println("Error ja backup run time: ", backupRunTimeErr)
		}
	} else {
		fmt.Println("JaBackupRunTime already set, skipping command.")
	}
	agent_err := lib.Restart_jaz_agent_linux()
	if agent_err != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error agentd restart : %s", agent_err))
	}

	set_ntp_cmd := "sudo timedatectl set-ntp false"
	_, set_ntp_err := lib.Ssh_exec(set_ntp_cmd)
	if set_ntp_err != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error set ntp false : %s", set_ntp_err))
	}

	year := time.Now().Year()
	minutesBefore := 4

	// Calculate the time for the previous hour minus the minutes
	// Construct the current time to manipulate
	currentTime := time.Date(year, time.November, day+1, hr+1, 0, 0, 0, time.UTC)

	// Subtract the minutes (4 minutes before the hour)
	adjustedTime := currentTime.Add(-time.Minute * time.Duration(minutesBefore))

	// Format the adjusted time to match the 'timedatectl' format
	formattedTime := adjustedTime.Format("2006-01-02 15:04:05")

	// Create the timedatectl command
	//setDatetimeCmd := fmt.Sprintf("sudo timedatectl set-time '%s'", formattedTime)
	setDatetimeCmd := "sudo timedatectl set-time '" + formattedTime + "'"

	// Output the command to verify
	fmt.Println("Generated Command:", setDatetimeCmd)

	// set_datetime_cmd := fmt.Sprintf(`sudo timedatectl set-time '2024-11-07 15:30:00'`)
	_, set_dateTime_err := lib.Ssh_exec(setDatetimeCmd)
	if set_dateTime_err != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error set date time : %s", set_dateTime_err))
	}

	set_ntp_cmd = "sudo timedatectl set-ntp true"
	_, set_ntp_err = lib.Ssh_exec(set_ntp_cmd)
	if set_ntp_err != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error set ntp true : %s", set_ntp_err))
	}

	//time.Sleep(5 * time.Minute)
	lib.Run_Timeout(5 * 60)
	contents, get_jobfile_and_folder_err = getDirectoryContents()
	if get_jobfile_and_folder_err != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error after backup thread start job file and folder : %s", get_jobfile_and_folder_err))
	}

	if len(contents) == 0 {

		loc, err := time.LoadLocation("Asia/Yangon")
		if err != nil {
			fmt.Println("Error loading timezone:", err)
		}

		// Get the current time in Asia/Yangon timezone
		currentTimeInZone := time.Now().In(loc)

		// Format the current time to '2024-11-09 22:10:55'
		formattedCurrentTime := currentTimeInZone.Format("2006-01-02 15:04:05")

		// Print the formatted time (you can use this in your 'timedatectl' command)
		fmt.Println("Formatted time:", formattedCurrentTime)

		set_ntp_cmd = "sudo timedatectl set-ntp false"
		_, set_ntp_err = lib.Ssh_exec(set_ntp_cmd)
		if set_ntp_err != nil {
			fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error set ntp false : %s", set_ntp_err))
		}

		// Example: Set the time using the formatted string in a command
		setCurrentTimeZone := fmt.Sprintf("sudo timedatectl set-time '%s'", formattedCurrentTime)
		fmt.Println("Command to set time:", setCurrentTimeZone)

		_, set_curr_time_err := lib.Ssh_exec(setCurrentTimeZone)
		if set_curr_time_err != nil {
			fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error set current time : %s", set_curr_time_err))
		}

		set_ntp_cmd = "sudo timedatectl set-ntp true"
		_, set_ntp_err = lib.Ssh_exec(set_ntp_cmd)
		if set_ntp_err != nil {
			fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error set ntp true : %s", set_ntp_err))
		}

		_set_current_backup_time_cmd := "sudo sed -i '/^JaBackupTime=1$/d' /etc/jobarranger/jobarg_agentd.conf"
		_, _set_current_backup_time_err := lib.Ssh_exec(_set_current_backup_time_cmd)
		if _set_current_backup_time_err != nil {
			fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error set current backup time : %s", _set_current_backup_time_err))
		}

		_set_current_backup_run_time_cmd := fmt.Sprintf("sudo sed -i '/^JaBackupRunTime=%d$/d' /etc/jobarranger/jobarg_agentd.conf", hr+1)
		_, _set_current_backup_run_time_err := lib.Ssh_exec(_set_current_backup_run_time_cmd)
		if _set_current_backup_run_time_err != nil {
			fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error set current backup time : %s", _set_current_backup_run_time_err))
		}

		agent_err = lib.Restart_jaz_agent_linux()
		if agent_err != nil {
			fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error agentd restart : %s", agent_err))
		}

		return PASSED
	}

	return FAILED
}

func WatchLog(client *ssh.Client, pattern string, checkInterval time.Duration) (bool, error) {
	// Start an interactive session to tail the log file
	session, err := client.NewSession()
	if err != nil {
		return false, fmt.Errorf("failed to create session: %s", err)
	}
	defer session.Close()

	// Run the 'tail -f' command on the remote server
	cmd := "tail -f /var/log/jobarranger/jobarg_server.log"
	stdout, err := session.StdoutPipe()
	if err != nil {
		return false, fmt.Errorf("failed to get stdout pipe: %s", err)
	}

	// Start the session
	err = session.Start(cmd)
	if err != nil {
		return false, fmt.Errorf("failed to start session: %s", err)
	}

	// Create a scanner to read the output from the remote command
	scanner := bufio.NewScanner(stdout)

	// Define the regular expression pattern to match the desired log line
	re := regexp.MustCompile(pattern)

	// Continuously read the log output, checking every second
	for {
		// Sleep for the specified check interval (e.g., 1 second)
		time.Sleep(checkInterval)

		// Check if there's new output
		if scanner.Scan() {
			line := scanner.Text()
			fmt.Println("Scanned line:", line) // Optional: debug output to see scanned lines

			// Check if the line matches the pattern
			if re.MatchString(line) {
				// If a match is found, print the line and return true
				fmt.Println("Matched log line:", line)
				return true, nil
			}
		} else {
			if err := scanner.Err(); err != nil {
				return false, fmt.Errorf("error reading from remote log: %v", err)
			}
		}
	}
}

// Function to check if the desired string exists in the file
func lineExistsInFile(filePath, searchString string) (bool, error) {
	// Use 'grep' to search for the string in the file
	cmd := fmt.Sprintf("grep -q '%s' %s", searchString, filePath)
	output, err := lib.Ssh_exec(cmd)

	// If grep finds the string, it will return exit code 0, otherwise exit code 1
	if err != nil {
		// 'grep -q' returns a non-zero exit code if the string is not found
		if strings.Contains(string(output), "no such file") {
			return false, fmt.Errorf("file not found: %s", filePath) // lowercase error message
		}
		// If error is not related to string not being found, report it
		return false, fmt.Errorf("failed to search the file: %v", err) // lowercase error message
	}
	// If exit code 0, string exists, return true
	return true, nil
}

func CheckIconCountForSendingStatus(jobnetId string, testcase *dao.TestCase, client *ssh.Client) common.Testcase_status {

	//clean up before job count
	lib.Jobarg_cleanup_linux()

	//clean up log file
	clean_log_cmd := "> /var/log/jobarranger/jobarg_server.log"
	_, clean_log_err := lib.Ssh_exec(clean_log_cmd)
	if clean_log_err != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error clean log : %s", clean_log_err))
	}

	filePath := "/etc/jobarranger/locale/logmessage_64BIT.txt"
	searchString := "JACHECKJOB000002,0,0," // The string we're looking to ensure exists

	// Check if the desired line is already in the file
	exists, err := lineExistsInFile(filePath, searchString)
	if err != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error checking if line exists: %v", err))
	}

	// If the line doesn't exist, apply the modification with sed
	if !exists {
		// Construct the sed command to replace the line
		configCmd := fmt.Sprintf("sudo sed -i 's/JACHECKJOB000002,4,0,/JACHECKJOB000002,0,0,/' %s", filePath)
		_, configErr := lib.Ssh_exec(configCmd)
		if configErr != nil {
			fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error applying sed command: %s", configErr))
		}
	}

	// config_cmd := "sudo sed -i 's/JACHECKJOB000002,4,0,/JACHECKJOB000002,0,0,/' /etc/jobarranger/locale/logmessage_64BIT.txt"
	// _, config_err := lib.Ssh_exec(config_cmd)
	// if config_err != nil {
	// 	fmt.Println(lib.Logi(common.LOG_LEVEL_ERR,"Error config : %s", config_err))
	// }

	// start noamal job sleep 1000 host 1
	if err := lib.Jobarg_enable_jobnet("Icon_1", "sleep 1000"); err != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Failed to enable jobnet, Error: %s", err))
	}
	_, job_run_err := lib.Jobarg_exec(jobnetId)

	if job_run_err != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error running job: %s", job_run_err))
	}
	// End noraml job sleep 1000 host 1

	time.Sleep(10 * time.Second)
	// Start reboot after host 1
	if err := lib.Jobarg_enable_jobnet("Icon_1", "Reboot_after_completing_jobs_linux"); err != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Failed to enable jobnet, Error: %s", err))
	}

	_, job_reboot_after_err := lib.Jobarg_exec(jobnetId)
	if job_reboot_after_err != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error reboot after job : %s", job_reboot_after_err))
	}

	// watch log job count host1
	pattern := `\[\s*JACHECKJOB000002\s*\]  In check_job_process\(\), job count: (\d+), reboot job count: (\d+)`
	normal_job_count, normal_job_count_err := WatchLog(client, pattern, 1*time.Second)
	if normal_job_count_err != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error normal job count : %s", normal_job_count_err))
	}

	fmt.Println(normal_job_count)

	if normal_job_count {
		abortRebootIconQuery := "UPDATE ja_run_job_table SET status = 4 WHERE job_id  = 'REBOOT-1';"
		_, abort_reboot_err := lib.ExecuteQuery(lib.DBQuery(abortRebootIconQuery))
		if abort_reboot_err != nil {
			fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error abort reboot icon : %s", abort_reboot_err))
		}

		time.Sleep(20 * time.Second)

		abortRebootJobnetQuery := "UPDATE ja_run_jobnet_summary_table SET jobnet_abort_flag = 1 WHERE jobnet_name  = 'Reboot_after_completing_jobs_linux';"
		_, abort_reboot_err = lib.ExecuteQuery(lib.DBQuery(abortRebootJobnetQuery))
		if abort_reboot_err != nil {
			fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error abort reboot jobnet : %s", abort_reboot_err))
		}
	}

	return FAILED
}
