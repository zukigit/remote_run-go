package tickets

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/zukigit/remote_run-go/src/common"
	"github.com/zukigit/remote_run-go/src/dao"
	"github.com/zukigit/remote_run-go/src/lib"
	"golang.org/x/crypto/ssh"
)

type Ticket_1341 struct {
	no          uint
	description string
	testcases   []dao.TestCase
}

func (t *Ticket_1341) New_testcase(testcase_id uint, testcase_description string) *dao.TestCase {
	return dao.New_testcase(testcase_id, testcase_description)
}

func (t *Ticket_1341) Get_no() uint {
	return t.no
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

	tc_3 := t.New_testcase(162, "JaRunInterval Test with Database.")
	tc_func = func() common.Testcase_status {
		return JaRunIntervalTestWithDatabase("Icon_1", 0, 30, tc_3, common.Client)
	}
	tc_3.Set_function(tc_func)
	t.Add_testcase(*tc_3)
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

	_set_param_err := lib.Ja_set_server_config_linux("JaRunInterval", "5")
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

func JaRunIntervalTestWithDatabase(jobnetId string, processCount int, processCheckTimeout int, testcase *dao.TestCase, client *ssh.Client) common.Testcase_status {
	// Define the path to the postgresql.conf file
	newLogDirectory := "/var/log/postgresql/" // New log directory path

	// Command to create the log directory and set permissions
	createDirCmd := fmt.Sprintf(`bash -c 'rm -rf /var/log/postgresql/ && mkdir /var/log/postgresql/ && chown -R postgres:postgres %s'`, newLogDirectory)

	// Execute the directory creation and permission command
	output, err_create_dir := lib.Ssh_exec(createDirCmd)
	if err_create_dir != nil {
		fmt.Println("Error creating log directory and setting permissions:", err_create_dir)
		fmt.Println("Output:", output) // Print the error and output to debug
		return FAILED
	}

	logDes_cmd := "sed -i \"s|^#*\\(log_destination\\s*=\\s*\\).*|log_destination = 'stderr'|\" /var/lib/pgsql/data/postgresql.conf"

	_logDes_output, _logDes_err := lib.Ssh_exec(logDes_cmd)

	if _logDes_err != nil {
		fmt.Printf("Error logDes : %v\n", _logDes_err)
		fmt.Printf("Output: %s\n", _logDes_output) // Print the output to get more details
		return FAILED
	}

	log_collector_cmd := "sed -i 's|^#*\\(logging_collector\\s*=\\s*\\).*|\\1on|' /var/lib/pgsql/data/postgresql.conf"

	_log_collector_output, _log_collector_err := lib.Ssh_exec(log_collector_cmd)

	if _log_collector_err != nil {
		fmt.Printf("Error log collector : %v\n", _log_collector_err)
		fmt.Printf("Output: %s\n", _log_collector_output) // Print the output to get more details
		return FAILED
	}

	_log_dir_cmd := "sed -i \"s|^#*\\(log_directory\\s*=\\s*\\).*|log_directory = '/var/log/postgresql/'|\" /var/lib/pgsql/data/postgresql.conf"

	_log_dir_output, _log_dir_err := lib.Ssh_exec(_log_dir_cmd)

	if _log_dir_err != nil {
		fmt.Printf("Error log dir : %v\n", _log_dir_err)
		fmt.Printf("Output: %s\n", _log_dir_output) // Print the output to get more details
		return FAILED
	}

	_log_statement_cmd := "sed -i \"s|^#*\\(log_statement\\s*=\\s*\\).*|log_statement = 'all'|\" /var/lib/pgsql/data/postgresql.conf"

	_log_statement_output, _log_statement_err := lib.Ssh_exec(_log_statement_cmd)

	if _log_statement_err != nil {
		fmt.Printf("Error log statement : %v\n", _log_statement_err)
		fmt.Printf("Output: %s\n", _log_statement_output) // Print the output to get more details
		return FAILED
	}

	// Optionally, reload PostgreSQL after the change
	reloadCmd := `bash -c 'systemctl restart postgresql'`
	_, err := lib.Ssh_exec(reloadCmd)
	if err != nil {
		fmt.Println("Error reloading PostgreSQL:", err)
		return FAILED
	}
	return FAILED
}
