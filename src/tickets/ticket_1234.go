package tickets

import (
	"context"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/zukigit/remote_run-go/src/common"
	"github.com/zukigit/remote_run-go/src/lib"
	"golang.org/x/crypto/ssh"
)

type Schedule struct {
	ScheduleID string
	UpdateDate int64 // Assuming update_date is of type bigint
}

type Ticket_1234 struct {
	no                                          uint
	description                                 string
	PASSED_count, FAILED_count, MUSTCHECK_count int
	testcases                                   []common.TestCase
}

func (t *Ticket_1234) New_testcase(testcase_id uint, testcase_description string) *common.TestCase {
	return common.New_testcase(testcase_id, testcase_description)
}

func (t *Ticket_1234) Get_no() uint {
	return t.no
}

func (t *Ticket_1234) Set_PASSED_count(passed_count int) {
	t.PASSED_count = passed_count
}

func (t *Ticket_1234) Set_FAILED_count(failed_count int) {
	t.FAILED_count = failed_count
}

func (t *Ticket_1234) Set_MUSTCHECK_count(mustcheck_count int) {
	t.MUSTCHECK_count = mustcheck_count
}

func (t *Ticket_1234) Get_dsctn() string {
	return t.description
}

func (t *Ticket_1234) Add_testcase(tc common.TestCase) {
	t.testcases = append(t.testcases, tc)
}

func (t *Ticket_1234) Get_testcases() []common.TestCase {
	return t.testcases
}

// Enter your ticket information here
func (t *Ticket_1234) Set_values() {
	t.no = 1234 // Enter your ticket id
	t.description = "Run Jobnet WINRM and check std out."
}

// Add your test case here
func (t *Ticket_1234) Add_testcases() {
	tc_1 := t.New_testcase(134, "Run WINRM Jobnet and check std out printed hostname.")
	tc_func := func() common.Testcase_status {
		if err := lib.Jobarg_enable_jobnet("Icon_1", "WINRM_SRV_hostname"); err != nil {
			lib.Logi(common.LOG_LEVEL_ERR, "Failed to enable jobnet, Error: %s", err)
			return FAILED
		}
		return RunJobnetWinRMJobnetHostname("Icon_1", 1600, 80, tc_1, common.Client)
	}
	tc_1.Set_function(tc_func)
	t.Add_testcase(*tc_1)

	tc_2 := t.New_testcase(135, "Run WINRM Jobnet and check std out printed hostname.")
	tc_func = func() common.Testcase_status {
		if err := lib.Jobarg_enable_jobnet("Icon_1", "WINRM_SRV_hostname"); err != nil {
			lib.Logi(common.LOG_LEVEL_ERR, "Failed to enable jobnet, Error: %s", err)
			return FAILED
		}
		return RunJobnetWinRMJobnetIp("Icon_1", 1600, 80, tc_2, common.Client)
	}
	tc_2.Set_function(tc_func)
	t.Add_testcase(*tc_2)

	tc_3 := t.New_testcase(136, "Run WINRM getHost Jobnet and check std out printed hostname.")
	tc_func = func() common.Testcase_status {
		if err := lib.Jobarg_enable_jobnet("Icon_1", "WINRM_SRV"); err != nil {
			lib.Logi(common.LOG_LEVEL_ERR, "Failed to enable jobnet, Error: %s", err)
			return FAILED
		}
		return RunJobnetWinRMJobnetHostnameAndIp("Icon_1", 1600, 80, tc_3, common.Client)
	}
	tc_3.Set_function(tc_func)
	t.Add_testcase(*tc_3)

	// tc_4 := t.New_testcase(138, "Create schedule jobnet (Jobnet load span 30 miutes)")
	// tc_func = func() common.Testcase_status {
	// 	return RunScheduleLoadSpan("Icon_1", 0, 30, tc_4, common.Client)
	// }
	// tc_4.Set_function(tc_func)
	// t.Add_testcase(*tc_4)

	// tc_5 := t.New_testcase(139, "Create schedule jobnet (Jobnet load span 60 miutes)")
	// tc_func = func() common.Testcase_status {
	// 	return RunScheduleLoadSpan("Icon_1", 0, 60, tc_5, common.Client)
	// }
	// tc_5.Set_function(tc_func)
	// t.Add_testcase(*tc_5)
}

func extractIPAndHostname(valueTrimmed string) (string, string) {
	// Find the position where the digits end and the hostname starts
	var ip string
	var hostname string

	// Use strings.IndexFunc to find where the hostname starts
	index := strings.IndexFunc(valueTrimmed, func(r rune) bool {
		return r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' // Check for letters
	})

	if index != -1 {
		// Extract IP and hostname using the index found
		ip = valueTrimmed[:index]
		hostname = valueTrimmed[index:]
	} else {
		// Handle case where no hostname is found
		ip = valueTrimmed
	}

	return ip, hostname
}

// Run the WINRM Jobnet and check the std out
func RunJobnetWinRMJobnetHostname(jobnetId string, processCount int, processCheckTimeout int, testcase *common.TestCase, client *ssh.Client) common.Testcase_status {
	/*
		Prepare process before execute the ext jobnet
		1. cleanup data from ja_run_jobnet_table
	*/
	lib.Jobarg_cleanup_linux()

	// Get hostname
	hostname, _err := lib.GetOutputStrFromSSHCommand(client, "hostname")
	if _err != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error fetching hostname: %v", _err))
		return FAILED
	}

	// Trim any whitespace from the hostname
	hostname = strings.TrimSpace(hostname)
	fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "hostname info: %s", hostname))

	cmd := fmt.Sprintf(`bash -c 'export WINRM_SRV="%s"; 
	export HOSTNAME="%s"; 
	jobarg_exec -z %s -U Admin -P zabbix -j %s -E WINRM_SRV > /tmp/jobarg_output.txt 2>&1; 
	cat /tmp/jobarg_output.txt;'`, hostname, hostname, common.Login_info.Hostname, jobnetId)

	_run_jobnet, _err := lib.GetOutputStrFromSSHCommand(client, cmd)

	if _err != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error getting jobnet info: %v", _err))
		fmt.Println("Output from command:", _run_jobnet) // Print output for debugging
		if exitErr, ok := _err.(*exec.ExitError); ok {
			fmt.Println("Command exited with status:", exitErr.ExitCode())
		}
	} else {
		fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "Getting jobnet info: %s", _run_jobnet))

		// Adjusted regex pattern to account for spaces
		pattern := `Registry number\s*:\s*\[\s*(\d+)\s*\]`
		re := regexp.MustCompile(pattern)

		// Use FindStringSubmatch to get the ID
		matches := re.FindStringSubmatch(_run_jobnet)
		jobnetID := matches[1]
		fmt.Println("Extracted Jobnet ID:", jobnetID)
		time.Sleep(10 * time.Second) // Sleep for 2 seconds

		// cmd := fmt.Sprintf(`bash -c 'jobarg_get -z %s -U Admin -P zabbix -r %s -e 2>&1'`, common.Login_info.Hostname, jobnetID)
		cmd := fmt.Sprintf(`bash -c 'eval $(jobarg_get -z %s -U Admin -P zabbix -r %s -e) && echo -n $JA_LASTSTDOUT'`, common.Login_info.Hostname, jobnetID)
		fmt.Println("Executing command:", cmd)

		fmt.Println(cmd)
		lastSTDOut, _err_std_out := lib.GetOutputStrFromSSHCommand(client, cmd)
		fmt.Println(lastSTDOut)
		if lastSTDOut != "" {
			// Remove all spaces
			lastSTDOut = strings.ReplaceAll(lastSTDOut, " ", "")

			// Now replace with a comma
			lastSTDOut = strings.ReplaceAll(lastSTDOut, "(", ",(")
			lastSTDOut = strings.ReplaceAll(lastSTDOut, ")", "),")
			// Trim any leading/trailing whitespace
			lastSTDOut = strings.TrimSpace(lastSTDOut)

			// Split the output by commas
			values := strings.Split(lastSTDOut, ",")
			// Trim each value to remove whitespace
			for i := range values {
				values[i] = strings.TrimSpace(values[i])
			}

			// Print out the cleaned values
			for _, value := range values {
				if value != "" {
					// fmt.Println(lib.Logi(common.LOG_LEVEL_INFO,"Value: %s", value))
					// Assuming 'hostname' is already set and contains the correct value
					hostname := strings.TrimSpace(hostname)
					valueTrimmed := strings.TrimSpace(value)
					// ip := strings.TrimSpace(common.Login_info.Hostname)

					_, hostnameTrimmed := extractIPAndHostname(valueTrimmed)
					// Print for debugging
					fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "Comparing Hostname: '%s' with Value: '%s'", hostnameTrimmed, valueTrimmed))

					if strings.EqualFold(hostname, hostnameTrimmed) { // Case insensitive comparison
						return PASSED
					}
				}
			}
		} else {
			fmt.Println("No output for JA_LASTSTDOUT.")
		}

		if _err_std_out != nil {
			fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error Std out log: %s", _err_std_out))
		}
	}

	return FAILED
}

func RunJobnetWinRMJobnetIp(jobnetId string, processCount int, processCheckTimeout int, testcase *common.TestCase, client *ssh.Client) common.Testcase_status {
	/*
		Prepare process before execute the ext jobnet
		1. cleanup data from ja_run_jobnet_table
	*/
	lib.Jobarg_cleanup_linux()

	// Get hostname
	hostname, _err := lib.GetOutputStrFromSSHCommand(client, "hostname")
	if _err != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error fetching hostname: %v", _err))
		return FAILED
	}

	// Trim any whitespace from the hostname
	hostname = strings.TrimSpace(hostname)
	fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "hostname info: %s", hostname))

	cmd := fmt.Sprintf(`bash -c 'export WINRM_SRV="getHost(%s)"; 
	export HOSTNAME="%s"; 
	jobarg_exec -z %s -U Admin -P zabbix -j %s -e WINRM_SRV > /tmp/jobarg_output.txt 2>&1; 
	cat /tmp/jobarg_output.txt;'`, hostname, hostname, common.Login_info.Hostname, jobnetId)

	_run_jobnet, _err := lib.GetOutputStrFromSSHCommand(client, cmd)

	if _err != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error getting jobnet info: %v", _err))
		fmt.Println("Output from command:", _run_jobnet) // Print output for debugging
		if exitErr, ok := _err.(*exec.ExitError); ok {
			fmt.Println("Command exited with status:", exitErr.ExitCode())
		}
	} else {
		fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "Getting jobnet info: %s", _run_jobnet))

		// Adjusted regex pattern to account for spaces
		pattern := `Registry number\s*:\s*\[\s*(\d+)\s*\]`
		re := regexp.MustCompile(pattern)

		// Use FindStringSubmatch to get the ID
		matches := re.FindStringSubmatch(_run_jobnet)
		jobnetID := matches[1]
		fmt.Println("Extracted Jobnet ID:", jobnetID)
		time.Sleep(10 * time.Second) // Sleep for 2 seconds

		// cmd := fmt.Sprintf(`bash -c 'jobarg_get -z %s -U Admin -P zabbix -r %s -e 2>&1'`, common.Login_info.Hostname, jobnetID)
		cmd := fmt.Sprintf(`bash -c 'eval $(jobarg_get -z %s -U Admin -P zabbix -r %s -e) && echo -n $JA_LASTSTDOUT'`, common.Login_info.Hostname, jobnetID)
		fmt.Println("Executing command:", cmd)

		fmt.Println(cmd)
		lastSTDOut, _err_std_out := lib.GetOutputStrFromSSHCommand(client, cmd)
		fmt.Println(lastSTDOut)
		if lastSTDOut != "" {
			// Remove all spaces
			lastSTDOut = strings.ReplaceAll(lastSTDOut, " ", "")

			// Now replace with a comma
			lastSTDOut = strings.ReplaceAll(lastSTDOut, "(", ",(")
			lastSTDOut = strings.ReplaceAll(lastSTDOut, ")", "),")
			// Trim any leading/trailing whitespace
			lastSTDOut = strings.TrimSpace(lastSTDOut)

			// Split the output by commas
			values := strings.Split(lastSTDOut, ",")
			// Trim each value to remove whitespace
			for i := range values {
				values[i] = strings.TrimSpace(values[i])
			}

			// Print out the cleaned values
			for _, value := range values {
				if value != "" {
					// fmt.Println(lib.Logi(common.LOG_LEVEL_INFO,"Value: %s", value))
					// Assuming 'hostname' is already set and contains the correct value
					// hostname := strings.TrimSpace(hostname)
					valueTrimmed := strings.TrimSpace(value)
					ip := strings.TrimSpace(common.Login_info.Hostname)

					ipTrimed, hostnameTrimmed := extractIPAndHostname(valueTrimmed)
					// Print for debugging
					fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "Comparing Hostname: '%s' with Value: '%s'", hostnameTrimmed, valueTrimmed))

					if strings.EqualFold(ip, ipTrimed) { // Case insensitive comparison
						return PASSED
					}
				}
			}
		} else {
			fmt.Println("No output for JA_LASTSTDOUT.")
		}

		if _err_std_out != nil {
			fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error Std out log: %s", _err_std_out))
		}
	}

	return FAILED
}

// Run the WINRM Jobnet and check the std out
func RunJobnetWinRMJobnetHostnameAndIp(jobnetId string, processCount int, processCheckTimeout int, testcase *common.TestCase, client *ssh.Client) common.Testcase_status {
	/*
		Prepare process before execute the ext jobnet
		1. cleanup data from ja_run_jobnet_table
	*/
	lib.Jobarg_cleanup_linux()

	// Get hostname
	hostname, _err := lib.GetOutputStrFromSSHCommand(client, "hostname")
	if _err != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error fetching hostname: %v", _err))
		return FAILED
	}

	// Trim any whitespace from the hostname
	hostname = strings.TrimSpace(hostname)
	fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "hostname info: %s", hostname))

	cmd := fmt.Sprintf(`bash -c 'export WINRM_SRV="getHost(%s)"; 
	export HOSTNAME="%s"; 
	jobarg_exec -z %s -U Admin -P zabbix -j %s -e WINRM_SRV,HOSTNAME > /tmp/jobarg_output.txt 2>&1; 
	cat /tmp/jobarg_output.txt;'`, hostname, hostname, common.Login_info.Hostname, jobnetId)

	_run_jobnet, _err := lib.GetOutputStrFromSSHCommand(client, cmd)

	if _err != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error getting jobnet info: %v", _err))
		fmt.Println("Output from command:", _run_jobnet) // Print output for debugging
		if exitErr, ok := _err.(*exec.ExitError); ok {
			fmt.Println("Command exited with status:", exitErr.ExitCode())
		}
	} else {
		fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "Getting jobnet info: %s", _run_jobnet))

		// Adjusted regex pattern to account for spaces
		pattern := `Registry number\s*:\s*\[\s*(\d+)\s*\]`
		re := regexp.MustCompile(pattern)

		// Use FindStringSubmatch to get the ID
		matches := re.FindStringSubmatch(_run_jobnet)
		jobnetID := matches[1]
		fmt.Println("Extracted Jobnet ID:", jobnetID)
		time.Sleep(10 * time.Second) // Sleep for 2 seconds

		// cmd := fmt.Sprintf(`bash -c 'jobarg_get -z %s -U Admin -P zabbix -r %s -e 2>&1'`, common.Login_info.Hostname, jobnetID)
		cmd := fmt.Sprintf(`bash -c 'eval $(jobarg_get -z %s -U Admin -P zabbix -r %s -e) && echo -n $JA_LASTSTDOUT'`, common.Login_info.Hostname, jobnetID)
		fmt.Println("Executing command:", cmd)

		fmt.Println(cmd)
		lastSTDOut, _err_std_out := lib.GetOutputStrFromSSHCommand(client, cmd)
		fmt.Println(lastSTDOut)
		if lastSTDOut != "" {
			// Remove all spaces
			lastSTDOut = strings.ReplaceAll(lastSTDOut, " ", "")

			// Now replace with a comma
			lastSTDOut = strings.ReplaceAll(lastSTDOut, "(", ",(")
			lastSTDOut = strings.ReplaceAll(lastSTDOut, ")", "),")
			// Trim any leading/trailing whitespace
			lastSTDOut = strings.TrimSpace(lastSTDOut)

			// Split the output by commas
			values := strings.Split(lastSTDOut, ",")
			// Trim each value to remove whitespace
			for i := range values {
				values[i] = strings.TrimSpace(values[i])
			}

			// Print out the cleaned values
			for _, value := range values {
				if value != "" {
					// fmt.Println(lib.Logi(common.LOG_LEVEL_INFO,"Value: %s", value))
					// Assuming 'hostname' is already set and contains the correct value
					hostname := strings.TrimSpace(hostname)
					valueTrimmed := strings.TrimSpace(value)
					ip := strings.TrimSpace(common.Login_info.Hostname)

					ipTrimed, hostnameTrimmed := extractIPAndHostname(valueTrimmed)
					// Print for debugging
					fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "Comparing Hostname: '%s' with Value: '%s'", hostnameTrimmed, valueTrimmed))

					if strings.EqualFold(hostname, hostnameTrimmed) && strings.EqualFold(ip, ipTrimed) { // Case insensitive comparison
						return PASSED
					}
				}
			}
		} else {
			fmt.Println("No output for JA_LASTSTDOUT.")
		}

		if _err_std_out != nil {
			fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error Std out log: %s", _err_std_out))
		}
	}

	return FAILED

}

func RunScheduleLoadSpan(jobnetId string, processCount int, processCheckTimeout int, testcase *common.TestCase, client *ssh.Client) common.Testcase_status {
	/*
		Prepare process before execute the ext jobnet
		1. cleanup data from ja_run_jobnet_table
	*/
	lib.Jobarg_cleanup_linux()

	currentDate := time.Now()

	// Format the date to YYYYMMDD
	formattedDate := currentDate.Format("20060102")

	// Print the formatted date
	// fmt.Println(formattedDate)

	setCurrentDateQuery := fmt.Sprintf("UPDATE ja_calendar_detail_table SET operating_date = '%s' WHERE calendar_id = 'COMMON_CALENDAR';", formattedDate)
	_, current_date_error := lib.ExecuteQuery(lib.DBQuery(setCurrentDateQuery))
	if current_date_error != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error set calendar detail : %s", current_date_error))
	}

	setPublicFlagQuery := "UPDATE ja_calendar_control_table SET valid_flag  = 1 WHERE calendar_id = 'COMMON_CALENDAR';"
	_, set_public_err := lib.ExecuteQuery(lib.DBQuery(setPublicFlagQuery))
	if set_public_err != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error set public flag : %s", set_public_err))
	}

	/*
		enable the jobnet
	*/
	if _enable_err := lib.Jobarg_enable_jobnet("Icon_1", "hostname"); _enable_err != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Failed to enable jobnet, Error: %s", _enable_err))
		return FAILED
	}

	/*
		set the standard time to server time

	*/

	setStandardTimeQuery := `UPDATE ja_parameter_table SET value = 1 WHERE parameter_name = 'MANAGER_TIME_SYNC';`
	_, standard_err := lib.ExecuteQuery(lib.DBQuery(setStandardTimeQuery))

	if standard_err != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Set Standard time to server time: %s", standard_err))
	}

	/*
		select schedule control table from update date

	*/
	query := `SELECT schedule_id, update_date FROM ja_schedule_control_table WHERE schedule_id = 'COMMON_SCHEDULE';`
	rows, _query_err := lib.GetData(lib.DBQuery(query))
	if _query_err != nil {
		fmt.Println("Select query error: ", _query_err)
		return FAILED
	}
	defer rows.Close()

	var schedules []Schedule
	var scheduleUpdateDate int64
	for rows.Next() {
		var schedule Schedule
		err := rows.Scan(&schedule.ScheduleID, &schedule.UpdateDate)
		if err != nil {
			fmt.Println("Error scanning row: ", err)
			return FAILED
		}
		schedules = append(schedules, schedule)

		scheduleUpdateDate = schedule.UpdateDate
	}

	// Check for errors during iteration
	if err := rows.Err(); err != nil {
		fmt.Println("Error occurred during row iteration: ", err)
		return FAILED
	}

	if len(schedules) > 0 {
		fmt.Println("Update Date Rows:")
		for _, schedule := range schedules {
			fmt.Printf("Schedule ID: %s, Update Date: %d\n", schedule.ScheduleID, schedule.UpdateDate)
		}
	} else {
		fmt.Println("No rows returned.")
	}

	fmt.Printf("Update Date: %d\n", scheduleUpdateDate)

	/*
	   get current time from db
	*/
	currentTimeQuery := `SELECT now()`
	current_rows, _current_time_err := lib.GetData(lib.DBQuery(currentTimeQuery))

	if _current_time_err != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error current time: %s", _current_time_err))
		return FAILED // Consider returning to avoid further processing
	}
	defer current_rows.Close() // Ensure rows are closed after processing

	var currentTime time.Time // Assuming the expected type is time.Time
	if current_rows.Next() {
		if err := current_rows.Scan(&currentTime); err != nil { // Use current_rows here
			fmt.Println("Error scanning row:", err)
			return FAILED // Handle the error as needed
		}
	} else {
		fmt.Println("No rows returned for current time query.")
	}

	fmt.Printf("Current time: %s\n", currentTime.Format("15:04"))

	// Calculate boot time (current time + 10 minutes)
	bootTime := currentTime.Add(5 * time.Minute)
	bootTimeFormatted := bootTime.Format("1504") // Format as HHMM
	fmt.Printf("Starting boot time: %s\n", bootTimeFormatted)

	// Calculate the end time (boot time + 30 minutes)
	endTime := bootTime.Add(time.Duration(processCheckTimeout) * time.Minute)
	endTimeFormatted := endTime.Format("1504") // Format as HHMM
	fmt.Printf("Ending time: %s\n", endTimeFormatted)

	// Generate boot times
	var bootTimes []string // Store boot times as strings to preserve leading zero
	for t := bootTime; t.Before(endTime) || t.Equal(endTime); t = t.Add(2 * time.Minute) {
		bootTimeValue := t.Format("1504")            // Format as HHMM
		bootTimes = append(bootTimes, bootTimeValue) // Store the formatted boot time as string
	}

	// Assuming you want to remove existing entries before inserting new ones
	removeQuery := `DELETE FROM ja_schedule_detail_table;`
	_, removeErr := lib.ExecuteQuery(lib.DBQuery(removeQuery))
	if removeErr != nil {
		fmt.Println("Error removing the schedule detail:", removeErr)
		return FAILED
	}

	// Output the generated boot times
	fmt.Println("Generated boot times:")
	for _, bt := range bootTimes {
		fmt.Println(bt)

		// Insert each boot time into the table as a string
		insertQuery := fmt.Sprintf(`
			INSERT INTO ja_schedule_detail_table (
				schedule_id, 
				calendar_id, 
				boot_time, 
				update_date, 
				created_date, 
				object_flag
			) VALUES (
				'COMMON_SCHEDULE', 
				'COMMON_CALENDAR', 
				'%s',  -- Boot time in HHMM format (e.g., "0842")
				%d,  -- Assuming scheduleUpdateDate is available
				CURRENT_TIMESTAMP, 
				0
			);
		`, bt, scheduleUpdateDate)

		_, insertErr := lib.ExecuteQuery(lib.DBQuery(insertQuery))
		if insertErr != nil {
			fmt.Println("Error inserting into ja_schedule_detail_table:", insertErr)
			return FAILED
		}
	}

	time.Sleep(30 * time.Second)

	//schedule enable query
	enableQuery := `UPDATE ja_schedule_control_table SET valid_flag = 1 WHERE schedule_id  = 'COMMON_SCHEDULE';`
	_, enable_err := lib.ExecuteQuery(lib.DBQuery(enableQuery))

	if enable_err != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error schedule enable: %s", enable_err))
	}

	time.Sleep(30 * time.Second)
	maxCount := 16
	if processCheckTimeout == 60 {
		maxCount = maxCount * 2
	}
	jobDoneCount, errJobCountWithDone := runStatusProcess(nil, processCheckTimeout*10, &maxCount) // Correctly capturing count and error
	if errJobCountWithDone != nil {
		fmt.Println(errJobCountWithDone) // Handle the error
		return FAILED
	}

	if jobDoneCount == maxCount {
		fmt.Println("Success all the jobnet are done status")
		return PASSED
	}

	return FAILED
}

// GetCountFromDB executes a count query and returns the count or an error.
func GetStatusFromDB(query string) (int, error) {
	dbQuery := lib.DBQuery(query) // Ensure to use the correct DBQuery from lib

	rows, err := lib.GetData(dbQuery) // Ensure GetData accepts this type
	if err != nil {
		return 0, fmt.Errorf("error fetching count: %w", err)
	}
	defer rows.Close()

	var count int
	if rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			return 0, fmt.Errorf("error scanning count: %w", err)
		}
	} else {
		return 0, fmt.Errorf("no rows found")
	}

	return count, nil
}

// runProcess monitors the count and aborts if it exceeds a threshold
func runStatusProcess(query *string, processCheckTimeout int, maxCount *int) (int, error) {
	// Use default query if none provided
	defaultQuery := "SELECT COUNT(*) FROM ja_run_jobnet_table WHERE status = 3;"
	actualQuery := defaultQuery
	if query != nil {
		actualQuery = *query
	}

	// Use default maxCount of 0 if none provided
	actualMaxCount := 0
	if maxCount != nil {
		actualMaxCount = *maxCount
	}

	timeoutDuration := time.Duration(processCheckTimeout) * time.Minute

	ctx, cancel := context.WithTimeout(context.Background(), timeoutDuration)
	defer cancel()

	count := 0
	var err error

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return 0, fmt.Errorf("error: timeout reached, exiting loop")

		case <-ticker.C:
			count, err = GetStatusFromDB(actualQuery)
			if err != nil {
				fmt.Println(err) // Log and continue
				continue
			}

			valueCopy := uint(count)
			if count == actualMaxCount {
				fmt.Printf("\rCount has reached or exceeded %d, stopping the loop.\n", valueCopy)
				return count, nil
			} else {
				// fmt.Printf("Max count %d\n", actualMaxCount)
				fmt.Printf("\rCount has not reached or exceeded %d, continuing to poll...", valueCopy)
				time.Sleep(500 * time.Millisecond)
				continue
			}
		}
	}
}
