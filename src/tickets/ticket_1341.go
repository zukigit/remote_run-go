package tickets

import (
	"bufio"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/zukigit/remote_run-go/src/common"
	"github.com/zukigit/remote_run-go/src/dao"
	"github.com/zukigit/remote_run-go/src/lib"
	"golang.org/x/crypto/ssh"
)

type Ticket_1341 struct {
	no                                          uint
	description                                 string
	PASSED_count, FAILED_count, MUSTCHECK_count int
	testcases                                   []dao.TestCase
}

func (t *Ticket_1341) New_testcase(testcase_id uint, testcase_description string) *dao.TestCase {
	return dao.New_testcase(testcase_id, testcase_description)
}

func (t *Ticket_1341) Get_no() uint {
	return t.no
}

func (t *Ticket_1341) Set_PASSED_count(passed_count int) {
	t.PASSED_count = passed_count
}

func (t *Ticket_1341) Set_FAILED_count(failed_count int) {
	t.FAILED_count = failed_count
}

func (t *Ticket_1341) Set_MUSTCHECK_count(mustcheck_count int) {
	t.MUSTCHECK_count = mustcheck_count
}

func (t *Ticket_1341) Get_dsctn() string {
	return t.description
}

func (t *Ticket_1341) Add_testcase(tc dao.TestCase) {
	t.testcases = append(t.testcases, tc)
}

func (t *Ticket_1341) Get_testcases() []dao.TestCase {
	return t.testcases
}

// Enter your ticket information here
func (t *Ticket_1341) Set_values() {
	t.no = 1341 // Enter your ticket id
	t.description = "Sleep Intervals"
}

// Add your test case here
func (t *Ticket_1341) Add_testcases() {
	tc_1 := t.New_testcase(160, "Check if Check job process start or not.")
	tc_func := func() common.Testcase_status {
		if err := lib.Jobarg_enable_jobnet("Icon_1", "two_jobicon"); err != nil {
			tc_1.Err_log("Failed to enable jobnet, Error: %s", err)
			return FAILED
		}
		return JaRunLoopNormalTest("Icon_1", 0, 30, tc_1, common.Client)
	}
	tc_1.Set_function(tc_func)
	t.Add_testcase(*tc_1)

	tc_2 := t.New_testcase(161, "Check if Check job process start or not.")
	tc_func = func() common.Testcase_status {
		if err := lib.Jobarg_enable_jobnet("Icon_1", "two_jobicon"); err != nil {
			tc_2.Err_log("Failed to enable jobnet, Error: %s", err)
			return FAILED
		}
		return JaRunLoopNormalTestWithJaRunInterval("Icon_1", 0, 30, tc_2, common.Client)
	}
	tc_2.Set_function(tc_func)
	t.Add_testcase(*tc_2)

	tc_3 := t.New_testcase(162, "JaRunInterval Test with Database. JaRun")
	tc_func = func() common.Testcase_status {
		searchTerm := "select inner_job_id, inner_jobnet_id, method_flag, job_type, test_flag, inner_jobnet_main_id from ja_run_job_table where status = 1 and method_flag <> 3"
		return JaRunIntervalTestWithDatabase("Icon_1", 0, 30, tc_3, common.Client, searchTerm, "jarun")
	}
	tc_3.Set_function(tc_func)
	t.Add_testcase(*tc_3)

	tc_4 := t.New_testcase(164, "JaRunInterval Test with Database. JaJob")
	tc_func = func() common.Testcase_status {
		searchTerm := "select inner_job_id, inner_jobnet_id, job_type, method_flag, timeout_flag, start_time from ja_run_job_table"
		return JaRunIntervalTestWithDatabase("Icon_1", 0, 30, tc_4, common.Client, searchTerm, "jajob")
	}
	tc_4.Set_function(tc_func)
	t.Add_testcase(*tc_4)
}

func JaRunLoopNormalTest(jobnetId string, processCount int, processCheckTimeout int, testcase *dao.TestCase, client *ssh.Client) common.Testcase_status {

	_set_param_err := lib.Ja_set_server_config_linux("JaRunInterval", "1")
	if _set_param_err != nil {
		fmt.Println(testcase.Err_log("Error set params : %s", _set_param_err))
	}

	_err_restart := lib.Restart_jaz_server()
	if _err_restart != nil {
		fmt.Println(testcase.Err_log("Error Jaz server restart : %s", _err_restart))
	}

	envs, _ := lib.Get_str_str_map("JA_HOSTNAME", "oss-redhat9psql", "JA_CMD", "hostname")
	run_jobnet_id, err := lib.Jobarg_exec_E(jobnetId, envs)
	if err != nil {
		fmt.Println(testcase.Err_log("Error executing job %s: %s", jobnetId, err))
	}

	time.Sleep(30 * time.Second)
	if run_jobnet_id != "" {
		_job_status, _job_status_err := lib.Jobarg_get_JA_JOBNETSTATUS(run_jobnet_id)
		if _job_status_err != nil {
			fmt.Println(testcase.Err_log("Error Job get status : %s", _job_status_err))
		}

		if _job_status == "END" {
			fmt.Println(testcase.Info_log("Job get status is : %s", _job_status))
			return PASSED
		}
	}

	return FAILED
}

func JaRunLoopNormalTestWithJaRunInterval(jobnetId string, processCount int, processCheckTimeout int, testcase *dao.TestCase, client *ssh.Client) common.Testcase_status {

	/*
		Prepare process before execute the ext jobnet
		1. cleanup data from ja_run_jobnet_table
	*/
	lib.Jobarg_cleanup_linux()

	setIntervalCmd := "sed -i 's/^\\s*#\\?\\s*JaRunInterval=\\s*[0-9]\\+$/JaRunInterval=5/' /etc/jobarranger/jobarg_server.conf"

	_, _set_param_err := lib.GetOutputStrFromSSHCommand(client, setIntervalCmd)
	if _set_param_err != nil {
		fmt.Println(testcase.Err_log("Error set params : %s", _set_param_err))
	}

	_err_restart := lib.Restart_jaz_server()
	if _err_restart != nil {
		fmt.Println(testcase.Err_log("Error Jaz server restart : %s", _err_restart))
	}

	envs, _ := lib.Get_str_str_map("JA_HOSTNAME", "oss.linux", "JA_CMD", "hostname")
	run_jobnet_id, err := lib.Jobarg_exec_E(jobnetId, envs)
	if err != nil {
		fmt.Println(testcase.Err_log("Error executing job %s: %s", jobnetId, err))
	}

	time.Sleep(30 * time.Second)
	if run_jobnet_id != "" {
		_job_status, _job_status_err := lib.Jobarg_get_JA_JOBNETSTATUS(run_jobnet_id)
		if _job_status_err != nil {
			fmt.Println(testcase.Err_log("Error Job get status : %s", _job_status_err))
		}

		if _job_status == "END" {
			fmt.Println(testcase.Info_log("Job get status is : %s", _job_status))

			cmd := fmt.Sprintf("jobarg_joblogput -z %s -U Admin -P zabbix -r %s", common.Login_info.Hostname, run_jobnet_id)

			_job_run_info, _job_run_info_err := lib.Ssh_exec_to_str(cmd)
			if _job_run_info_err != nil {
				fmt.Println(testcase.Err_log("Error job get run info : %s", _job_run_info_err))
			}

			if _job_run_info != "" {
				// fmt.Println(_job_run_info)
				// Call the function to extract JOB-1 and JOB-2 start times
				job1Timestamp, job2Timestamp, err := ExtractJobStartTimes(_job_run_info)
				if err != nil {
					fmt.Println("Error:", err)
					return FAILED
				}

				// Print the results for JOB-1 and JOB-2
				fmt.Printf("JOB-1 start time: %s\n", job1Timestamp)
				fmt.Printf("JOB-2 start time: %s\n", job2Timestamp)

				// Define the format in which the timestamp is provided
				format := "2006/01/02 15:04:05.000"

				// Parse both timestamps into time.Time objects
				job1Time, err := time.Parse(format, job1Timestamp)
				if err != nil {
					fmt.Println("Error parsing JOB-1 timestamp:", err)
					return FAILED
				}

				job2Time, err := time.Parse(format, job2Timestamp)
				if err != nil {
					fmt.Println("Error parsing JOB-2 timestamp:", err)
					return FAILED
				}

				// Compare the times
				if job1Time.Before(job2Time) {
					fmt.Println("JOB-1 started before JOB-2")
				} else if job1Time.After(job2Time) {
					fmt.Println("JOB-1 started after JOB-2")
				} else {
					fmt.Println("JOB-1 and JOB-2 started at the same time")
				}

				// Optionally, you can also compute the difference between the two times
				duration := job2Time.Sub(job1Time)
				// Check if the difference is at least 5 seconds
				if duration >= 5*time.Second {
					fmt.Println("The time difference between JOB-1 and JOB-2 is at least 5 seconds.")
					return PASSED
				} else {
					fmt.Println("The time difference between JOB-1 and JOB-2 is less than 5 seconds.")
					return FAILED
				}
			}

			return PASSED
		} else {
			fmt.Println(testcase.Err_log("Error job status : %s", _job_status))
			return FAILED
		}
	}

	return FAILED
}

// Function to extract and return JOB-1 and JOB-2 start times
func ExtractJobStartTimes(logData string) (string, string, error) {
	// Regular expression to match job start lines with the format:
	// "timestamp", "job name", "status"
	// Example log entry: "2024/11/05 14:56:51.334","1600000000000000110","1600000000000000110","1","0","Icon_1","JOB-1","JC00000003","Job has started.","two_jobicon","","Admin","2024/11/05 14:54:46",""
	re := regexp.MustCompile(`"([^"]+)","[^"]*","[^"]*","[^"]*","[^"]*","[^"]*","(JOB-[^"]+)","[^"]*","Job has started\."`)

	// Split the log data into lines
	lines := strings.Split(logData, "\n")

	// Variables to hold timestamps for JOB-1 and JOB-2
	var job1Timestamp, job2Timestamp string

	// Iterate over each line and check for matches
	for _, line := range lines {
		// Find all matches for job start events
		matches := re.FindStringSubmatch(line)
		if len(matches) > 2 {
			// If there are matches, collect the timestamp and job name
			timestamp := matches[1] // timestamp (e.g., "2024/11/05 14:56:51.334")
			jobName := matches[2]   // job name (e.g., "JOB-1")

			// Store the timestamp for JOB-1 and JOB-2
			if jobName == "JOB-1" && job1Timestamp == "" {
				job1Timestamp = timestamp
			} else if jobName == "JOB-2" && job2Timestamp == "" {
				job2Timestamp = timestamp
			}
		}

		// If we found both JOB-1 and JOB-2, stop searching
		if job1Timestamp != "" && job2Timestamp != "" {
			break
		}
	}

	// Return the timestamps for JOB-1 and JOB-2 (or an error if not found)
	if job1Timestamp != "" && job2Timestamp != "" {
		return job1Timestamp, job2Timestamp, nil
	} else {
		return "", "", fmt.Errorf("could not find both JOB-1 and JOB-2 in the logs")
	}
}

func JaRunIntervalTestWithDatabase(jobnetId string, processCount int, processCheckTimeout int, testcase *dao.TestCase, client *ssh.Client, searchTerm string, jatype string) common.Testcase_status {

	var getInterval string
	var err error

	if jatype == "jarun" {
		configPath := "/etc/jobarranger/jobarg_server.conf"
		cmdGetInterval := fmt.Sprintf("grep -i 'JaRunInterval=' %s", configPath)

		// Execute the command remotely using SSH
		getInterval, err = lib.GetOutputStrFromSSHCommand(client, cmdGetInterval)
		if err != nil {
			fmt.Println(testcase.Err_log("Error get run interval value : %s", err))
		}
	} else {
		configPath := "/etc/jobarranger/jobarg_server.conf"
		cmdGetInterval := fmt.Sprintf("grep -i 'JaJobInterval=' %s", configPath)

		// Execute the command remotely using SSH
		getInterval, err = lib.GetOutputStrFromSSHCommand(client, cmdGetInterval)
		if err != nil {
			fmt.Println(testcase.Err_log("Error get job interval value : %s", err))
		}
	}

	// Trim the output to remove extra whitespace
	getInterval = strings.TrimSpace(getInterval)
	var getIntervalValue string

	// fmt.Println(getInterval)

	// Split the string by the "=" character
	parts := strings.Split(getInterval, "=")

	// Check if the split was successful and if there's a valid value
	if len(parts) > 1 && strings.TrimSpace(parts[1]) != "" {
		// If the value after '=' is not empty
		getIntervalValue = strings.TrimSpace(parts[1])
		fmt.Println("Interval Value:", getIntervalValue) // Output: 5
	} else {
		// Handle case when no value is found
		fmt.Println("No valid value found")
	}

	// fmt.Println(defaultConfig)

	if common.DB_type == "postgres" {

		defaultConfig, _def_err := findPostgresqlConf(client)

		if _def_err != nil {
			fmt.Println(testcase.Err_log("Error default config : %s", _def_err))
		}

		// Define the path to the postgresql.conf file
		newLogDirectory := "/var/log/postgresql/" // New log directory path

		// Command to create the log directory and set permissions
		createDirCmd := `bash -c 'rm -rf /var/log/postgresql/ && mkdir /var/log/postgresql/ && chown -R postgres:postgres /var/log/postgresql/'`

		// Execute the directory creation and permission command
		output, err_create_dir := lib.GetOutputStrFromSSHCommand(client, createDirCmd)
		if err_create_dir != nil {
			fmt.Println("Error creating log directory and setting permissions:", err_create_dir)
			fmt.Println("Output:", output) // Print the error and output to debug
			return FAILED
		}

		logDes_cmd := fmt.Sprintf("sed -i \"s|^#*\\(log_destination\\s*=\\s*\\).*|log_destination = 'stderr'|\" %s", defaultConfig)

		_logDes_output, _logDes_err := lib.Ssh_exec(logDes_cmd)

		if _logDes_err != nil {
			fmt.Printf("Error logDes : %v\n", _logDes_err)
			fmt.Printf("Output: %s\n", _logDes_output) // Print the output to get more details
			return FAILED
		}

		log_collector_cmd := fmt.Sprintf("sed -i 's|^#*\\(logging_collector\\s*=\\s*\\).*|\\1on|' %s", defaultConfig)

		_log_collector_output, _log_collector_err := lib.GetOutputStrFromSSHCommand(client, log_collector_cmd)

		if _log_collector_err != nil {
			fmt.Printf("Error log collector : %v\n", _log_collector_err)
			fmt.Printf("Output: %s\n", _log_collector_output) // Print the output to get more details
			return FAILED
		}

		_log_dir_cmd := fmt.Sprintf("sed -i \"s|^#*\\(log_directory\\s*=\\s*\\).*|log_directory = '/var/log/postgresql/'|\" %s", defaultConfig)

		_log_dir_output, _log_dir_err := lib.Ssh_exec(_log_dir_cmd)

		if _log_dir_err != nil {
			fmt.Printf("Error log dir : %v\n", _log_dir_err)
			fmt.Printf("Output: %s\n", _log_dir_output) // Print the output to get more details
			return FAILED
		}

		_log_statement_cmd := fmt.Sprintf("sed -i \"s|^#*\\(log_statement\\s*=\\s*\\).*|log_statement = 'all'|\" %s", defaultConfig)

		_log_statement_output, _log_statement_err := lib.GetOutputStrFromSSHCommand(client, _log_statement_cmd)

		if _log_statement_err != nil {
			fmt.Printf("Error log statement : %v\n", _log_statement_err)
			fmt.Printf("Output: %s\n", _log_statement_output) // Print the output to get more details
			return FAILED
		}

		_log_filename_cmd := fmt.Sprintf("sed -i \"s|^#*\\(log_filename\\s*=\\s*\\).*|log_filename = 'postgresql.log'|\" %s", defaultConfig)

		_log_filename_output, _log_filename_err := lib.GetOutputStrFromSSHCommand(client, _log_filename_cmd)

		if _log_filename_err != nil {
			fmt.Printf("Error log filename : %v\n", _log_filename_err)
			fmt.Printf("Output: %s\n", _log_filename_output) // Print the output to get more details
			return FAILED
		}

		// Optionally, reload PostgreSQL after the change
		reloadCmd := `bash -c 'systemctl restart postgresql'`
		_, err_restart := lib.GetOutputStrFromSSHCommand(client, reloadCmd)
		if err_restart != nil {
			fmt.Println("Error reloading PostgreSQL:", err_restart)
			return FAILED
		} else {

			// Find the most recent log file in the log directory
			latestLogFile, err := getLatestLogFileSSH(newLogDirectory, client)
			if err != nil {
				fmt.Println("Error finding the latest log file:", err)
				return FAILED
			}

			fmt.Println("Latest log file path:", latestLogFile)
			// Now monitor the latest log file
			interval := 1 * time.Second // Check every 5 seconds
			// searchTerm := "select inner_job_id, inner_jobnet_id, method_flag, job_type, test_flag, inner_jobnet_main_id from ja_run_job_table where status = 1 and method_flag <> 3"
			logIntervalValue, log_interval_err := watchLogFileSSH(latestLogFile, searchTerm, client, interval)

			if log_interval_err != nil {
				fmt.Println(testcase.Err_log("Error log interval : %s", log_interval_err))
			}

			fmt.Println(logIntervalValue)

			// Parse logIntervalValue into time.Duration
			duration, err := time.ParseDuration(logIntervalValue)
			if err != nil {
				fmt.Println("Error parsing logIntervalValue:", err)
			}

			// Parse getIntervalValue as an integer
			getInterval, err := strconv.Atoi(getIntervalValue)
			if err != nil {
				fmt.Println("Error parsing getIntervalValue:", err)
			}

			// Compare the values
			if duration >= time.Duration(getInterval)*time.Second {
				fmt.Println("logIntervalValue is equal to or greater than getIntervalValue")
				log_statement_cmd := fmt.Sprintf("sed -i 's/^log_statement\\s*=.*$/log_statement = none/' %s", defaultConfig)
				_, log_statement_err := lib.GetOutputStrFromSSHCommand(client, log_statement_cmd)
				if log_statement_err != nil {
					fmt.Println("Error setting log_statement to none:", log_statement_err)
				} else {
					fmt.Println("log_statement has been set to none successfully.")
				}

				time.Sleep(5 * time.Second)
				reloadCmd := `bash -c 'systemctl restart postgresql'`
				_, err_restart := lib.GetOutputStrFromSSHCommand(client, reloadCmd)
				if err_restart != nil {
					fmt.Println("Error reloading PostgreSQL:", err_restart)
					return FAILED
				}
				lib.Restart_jaz_server()
				return PASSED
			} else {
				fmt.Println("logIntervalValue is less than getIntervalValue")
				lib.Restart_jaz_server()
				return FAILED
			}

		}
	} else {
		fmt.Println("This is mysql")
		config := "/etc/my.cnf"
		logDir := "/var/log/mysql/general.log"

		// Command to create the log directory and set permissions
		createDirCmd := `bash -c 'cd /var/log/mysql/ && rm -rf general.log && touch general.log && chown -R mysql:mysql general.log'`

		// Execute the directory creation and permission command
		output_dir, err_create_dir := lib.GetOutputStrFromSSHCommand(client, createDirCmd)
		if err_create_dir != nil {
			fmt.Println("Error creating log directory and setting permissions:", err_create_dir)
			fmt.Println("Output:", output_dir) // Print the error and output to debug
			return FAILED
		}

		// Command to check if [mysqld] exists in the config file
		checkMysqldCmd := "grep '^\\[mysqld\\]' /etc/my.cnf "
		check_block, checkErr := lib.Ssh_exec_to_str(checkMysqldCmd)

		fmt.Println(checkErr)

		// If [mysqld] block is not found
		if check_block == "" {
			// Insert [mysqld] block at the beginning of the file
			insertMysqldCmd := fmt.Sprintf(
				"sudo sed -i '1i\\\n[mysqld]\\ngeneral_log = 1\\ngeneral_log_file = /var/log/mysql/general.log' %s",
				config,
			)
			_, insertErr := lib.Ssh_exec(insertMysqldCmd)
			if insertErr != nil {
				fmt.Println("Error inserting [mysqld] block at the beginning:", insertErr)
			} else {
				fmt.Println("[mysqld] block added successfully at the beginning of the file.")
			}
		} else {
			// [mysqld] section exists, check for general_log and general_log_file
			checkLogSettingsCmd := fmt.Sprintf(
				"grep -q '^general_log' %s && grep -q '^general_log_file' %s",
				config,
				config,
			)
			_, checkLogSettingsErr := lib.Ssh_exec(checkLogSettingsCmd)

			// If the general_log and general_log_file are not present
			if checkLogSettingsErr != nil {
				// Append general_log and general_log_file under the existing [mysqld] section
				appendMysqldCmd := fmt.Sprintf(
					"sudo sed -i '/^\\[mysqld\\]/a\\\ngeneral_log = 1\\ngeneral_log_file = /var/log/mysql/general.log' %s",
					config,
				)
				_, appendErr := lib.Ssh_exec(appendMysqldCmd)
				if appendErr != nil {
					fmt.Println("Error appending settings under the [mysqld] section:", appendErr)
				} else {
					fmt.Println("Settings successfully appended under the [mysqld] section.")
				}
			} else {
				// general_log and general_log_file already exist, so do nothing
				fmt.Println("[mysqld] section already contains general_log and general_log_file. No changes made.")
			}
		}

		// Optionally, reload PostgreSQL after the change
		reloadCmd := `bash -c 'systemctl restart mysqld'`
		_, err_restart := lib.GetOutputStrFromSSHCommand(client, reloadCmd)
		if err_restart != nil {
			fmt.Println("Error reloading mysql :", err_restart)
			return FAILED
		} else {
			// Find the most recent log file in the log directory
			// latestLogFile, err := getLatestLogFileSSH(logDir, client)
			// if err != nil {
			// 	fmt.Println("Error finding the latest log file:", err)
			// 	return FAILED
			// }

			// fmt.Println("Latest log file path:", latestLogFile)
			// Now monitor the latest log file
			interval := 1 * time.Second // Check every 5 seconds
			// searchTerm := "select inner_job_id, inner_jobnet_id, method_flag, job_type, test_flag, inner_jobnet_main_id from ja_run_job_table where status = 1 and method_flag <> 3"
			logIntervalValue, log_interval_err := watchLogFileSSH(logDir, searchTerm, client, interval)

			if log_interval_err != nil {
				fmt.Println(testcase.Err_log("Error log interval : %s", log_interval_err))
			}

			fmt.Println(logIntervalValue)

			// Parse logIntervalValue into time.Duration
			duration, err := time.ParseDuration(logIntervalValue)
			if err != nil {
				fmt.Println("Error parsing logIntervalValue:", err)
			}

			// Parse getIntervalValue as an integer
			getInterval, err := strconv.Atoi(getIntervalValue)
			if err != nil {
				fmt.Println("Error parsing getIntervalValue:", err)
			}

			// Compare the values
			if duration >= time.Duration(getInterval)*time.Second {
				fmt.Println("logIntervalValue is equal to or greater than getIntervalValue")
				closeGencmd := "sudo sed -i '/^general_log/ s/^/#/' /etc/my.cnf && sudo sed -i '/^general_log_file/ s/^/#/' /etc/my.cnf"
				_, close_gen_err := lib.Ssh_exec(closeGencmd)
				if close_gen_err != nil {
					fmt.Println(testcase.Err_log("Error close gen cmd : %s", close_gen_err))
				}

				time.Sleep(5 * time.Second)
				lib.Restart_jaz_server()
				reloadCmd := `bash -c 'systemctl restart mysqld'`
				_, err_restart := lib.GetOutputStrFromSSHCommand(client, reloadCmd)
				if err_restart != nil {
					fmt.Println("Error reloading mysql :", err_restart)
				}
				return PASSED
			} else {
				fmt.Println("logIntervalValue is less than getIntervalValue")
				lib.Restart_jaz_server()
				return FAILED
			}
		}

	}
}

func getLatestLogFileSSH(logDirectory string, client *ssh.Client) (string, error) {
	// SSH command to list files, sorted by modification time, filter for postgresql*.log
	cmd := fmt.Sprintf("ls -t %spostgresql.log 2>/dev/null", logDirectory)

	// Execute the command remotely using SSH
	output, err := lib.GetOutputStrFromSSHCommand(client, cmd)
	if err != nil {
		return "", fmt.Errorf("error executing SSH command: %v", err)
	}

	// If the result is empty, it means no log files were found
	if output == "" {
		return "", fmt.Errorf("no log files found in directory %s", logDirectory)
	}

	// The output should be a list of log files, so we return the first one (the latest)
	return output, nil
}

// Function to extract the timestamp from a PostgreSQL log line
func extractTimestamp(logLine string) (time.Time, error) {

	if common.DB_type == "postgres" {

		// PostgreSQL timestamp format: 2024-11-07 16:26:49.571 +0630
		// Use a regular expression to match the timestamp at the start of the log line
		timestampPattern := `^(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}\.\d{3})`
		re := regexp.MustCompile(timestampPattern)
		matches := re.FindStringSubmatch(logLine)

		if len(matches) > 1 {
			// Parse the timestamp into a time.Time object
			timestamp, err := time.Parse("2006-01-02 15:04:05.000", matches[1])
			if err != nil {
				return time.Time{}, err
			}
			return timestamp, nil
		}

	} else {
		// Example log line: "2024-11-08T14:59:21.888993Z        12 Query     ..."
		// Define the log time format to match the timestamp in your logs
		const logTimeFormat = "2006-01-02T15:04:05.999999Z"

		// Split the log line to get the timestamp
		parts := strings.Fields(logLine)
		if len(parts) == 0 {
			return time.Time{}, fmt.Errorf("no timestamp found in log line")
		}

		timestampStr := parts[0]
		// Parse the timestamp using the defined format
		timestamp, err := time.Parse(logTimeFormat, timestampStr)
		if err != nil {
			return time.Time{}, fmt.Errorf("error extracting timestamp: %v", err)
		}

		return timestamp, nil
	}

	return time.Time{}, fmt.Errorf("no timestamp found in log line")
}

func watchLogFileSSH(filePath, searchTerm string, client *ssh.Client, interval time.Duration) (string, error) {

	var logInterval time.Duration
	// var confInterval string

	// configPath := "/etc/jobarranger/jobarg_server.conf"

	// SSH command to tail the log file
	cmd := fmt.Sprintf("tail -f %s", filePath)

	// Execute the command remotely using SSH
	session, err := client.NewSession()
	if err != nil {
		return "", fmt.Errorf("error creating SSH session: %v", err)
	}

	defer session.Close()

	// Get the output of the command (tail -f)
	stdout, err := session.StdoutPipe()
	if err != nil {
		return "", fmt.Errorf("error getting stdout pipe: %v", err)
	}

	// Start the command execution
	err = session.Start(cmd)
	if err != nil {
		return "", fmt.Errorf("error starting command: %v", err)
	}

	// Create a scanner to read the output of the tail command line by line
	scanner := bufio.NewScanner(stdout)

	var previousTimestamp time.Time

	// Create a ticker for checking at intervals (optional)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// Handle any errors that occurred during scanning
	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading from log file %s: %v\n", filePath, err)
	}

	// Continuously read the output for new lines
	for {
		select {
		case <-ticker.C: // Triggered at every interval
			// This can be used for additional checks if needed
			// fmt.Println("Interval check triggered...") // You can comment out or remove this
			// This block is executed periodically by the ticker

		case <-time.After(60 * time.Second): // Example: Timeout after 60 seconds
			return "", fmt.Errorf("timeout reached, stopping log tailing. %v", time.Second)

		default: // Default case to read the next log line
			if scanner.Scan() {
				line := scanner.Text()

				// Filter the lines containing the searchTerm (similar to 'grep')
				if strings.Contains(line, searchTerm) {
					// Extract timestamp from the log line
					timestamp, err := extractTimestamp(line)
					if err != nil {
						fmt.Printf("Error extracting timestamp: %v\n", err)
						continue
					}

					// If this isn't the first log line, calculate the interval
					if !previousTimestamp.IsZero() {
						logInterval = timestamp.Sub(previousTimestamp)
						//fmt.Printf("Time interval between logs: %v\n", interval)
						return logInterval.String(), nil
					}

					// Update the previous timestamp to the current one
					previousTimestamp = timestamp

					// Print the matched line
					// this is log out put
					// fmt.Println(line)

				}
			}
		}
	}
}

func findPostgresqlConf(client *ssh.Client) (string, error) {
	// Define the command to search for the postgresql.conf file
	// cmd := exec.Command("bash", "-c", "sudo find /var/lib/pgsql /etc/postgresql /var/lib/postgresql -name postgresql.conf 2>/dev/null | head -n 1")

	cmd := "bash -c 'sudo find /var/lib/pgsql /etc/postgresql /var/lib/postgresql -name postgresql.conf 2>/dev/null | head -n 1'"

	// Get the output of the command
	output, err := lib.GetOutputStrFromSSHCommand(client, cmd)
	if err != nil {
		return "", fmt.Errorf("error executing find command: %v", err)
	}

	// Clean up the output by trimming any leading/trailing whitespace
	result := strings.TrimSpace(string(output))

	// If the result is empty, it means no postgresql.conf file was found
	if result == "" {
		return "", fmt.Errorf("postgresql.conf not found")
	}

	// Return the path to the found postgresql.conf file
	return result, nil
}

// GetJaRunInterval fetches the JaRunInterval parameter from the jobarg_server.conf file
func GetJaRunInterval(client *ssh.Client, configPath string) (int, error) {
	// SSH command to find the JaRunInterval line in the config file, whether commented or not
	cmd := fmt.Sprintf("grep -E '^#?JaRunInterval=' %s", configPath)

	// Execute the command remotely using SSH
	output, err := lib.GetOutputStrFromSSHCommand(client, cmd)
	if err != nil {
		return 0, fmt.Errorf("error executing SSH command: %v", err)
	}

	// Trim the output to remove extra whitespace
	output = strings.TrimSpace(output)

	// If output is empty, JaRunInterval is not set
	if output == "" {
		return 0, fmt.Errorf("JaRunInterval not found in config file")
	}

	// Now remove the comment if present and extract the value
	re := regexp.MustCompile(`^#?\s*JaRunInterval\s*=\s*(\d+)`)
	matches := re.FindStringSubmatch(output)

	// If the regular expression matches, extract the value
	if len(matches) > 1 {
		// Convert the string value to an integer
		interval, err := strconv.Atoi(matches[1])
		if err != nil {
			return 0, fmt.Errorf("error converting JaRunInterval to integer: %v", err)
		}
		return interval, nil
	}

	return 0, fmt.Errorf("JaRunInterval line format is incorrect")
}
