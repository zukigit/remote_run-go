package tickets

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/zukigit/remote_run-go/src/common"
	"github.com/zukigit/remote_run-go/src/lib"
	"golang.org/x/crypto/ssh"
)

type Ticket_844 struct {
	Ticket_no                                   uint
	Ticket_description                          string
	PASSED_count, FAILED_count, MUSTCHECK_count int
	Testcases                                   []common.TestCase
}

func (t *Ticket_844) New_testcase(testcase_id uint, testcase_description string) *common.TestCase {
	return common.New_testcase(testcase_id, testcase_description)
}

func (t *Ticket_844) Get_no() uint {
	return t.Ticket_no
}

func (t *Ticket_844) Set_PASSED_count(passed_count int) {
	t.PASSED_count = passed_count
}

func (t *Ticket_844) Set_FAILED_count(failed_count int) {
	t.FAILED_count = failed_count
}

func (t *Ticket_844) Set_MUSTCHECK_count(mustcheck_count int) {
	t.MUSTCHECK_count = mustcheck_count
}

func (t *Ticket_844) Get_dsctn() string {
	return t.Ticket_description
}

func (t *Ticket_844) Add_testcase(tc common.TestCase) {
	t.Testcases = append(t.Testcases, tc)
}

func (t *Ticket_844) Get_testcases() []common.TestCase {
	return t.Testcases
}

// Enter your ticket information here
func (t *Ticket_844) Set_values() {
	t.Ticket_no = 844 // Enter your ticket id
	t.Ticket_description = "kill jobarg_command process."
}

// Add your test case here
func (t *Ticket_844) Add_testcases() {

	//TESTCASE 71
	tc_71 := t.New_testcase(71, "Normal job execution after agent restart")
	tc_func := func() common.Testcase_status {
		//Test Case : "Normal job execution after agent restart"
		//Task	    : stop the agent, delete all files in tmp, restart the agent
		//Result    : "The JOB returns with an error"

		var jobnet_id string = "Icon_1" // This value must be the same value as Jobnet_ID

		if err := lib.Jobarg_enable_jobnet(jobnet_id, "jobicon_linux"); err != nil {
			lib.Logi(common.LOG_LEVEL_ERR, "Failed to enable jobnet, Error: %s", err)
			return FAILED
		}
		fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "Info: Jobnet enabled successfully."))

		// 1. Jobarg_cleanup_start.
		// 2. Run Jobnet.
		// 3. Stop the agent
		// 4. Delete all files in tmp
		// 5. Restart the agent

		if err := lib.Jobarg_cleanup_linux(); err != nil {
			lib.Logi(common.LOG_LEVEL_ERR, "Failed to cleanup jobarg, Error: %s", err)
			return FAILED
		}
		fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "Info: Jobarg cleanup successfully."))

		envs, _ := lib.Get_str_str_map("JA_HOSTNAME", "oss.linux", "JA_CMD", "sleep 10")
		run_jobnet_id, error := lib.Jobarg_exec_E(jobnet_id, envs)
		if error != nil {
			fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error: %s, std_out: %s", error.Error(), run_jobnet_id))
			return FAILED
		}
		fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "%s has been successfully run with registry number: %s", jobnet_id, run_jobnet_id))

		err := lib.JobProcessCountCheck(1, 2, common.Client)
		if err != nil {
			fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error: %s, Failed to get process count.", err.Error()))
			return FAILED
		}
		fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "Process count has reached %d", 1))

		err = lib.Stop_jaz_agent_linux()
		if err != nil {
			lib.Logi(common.LOG_LEVEL_ERR, "Failed to stop agent, Error: %s", err)
			return FAILED
		}
		fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "Info: Agent stopped successfully."))

		if err := lib.Cleanup_agent_linux(); err != nil {
			lib.Logi(common.LOG_LEVEL_ERR, "Failed to cleanup agent, Error: %s", err)
			return FAILED
		}
		fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "Info: Agent cleanup successfully."))

		if err := lib.Restart_jaz_agent_linux(); err != nil {
			lib.Logi(common.LOG_LEVEL_ERR, "Failed to start agent, Error: %s", err)
			return FAILED
		}

		targetJobnetStatus := "RUN"
		targetJobStatus := "ERROR"
		jobnet_run_info, err := lib.Jobarg_get_jobnet_info(run_jobnet_id, targetJobnetStatus, targetJobStatus, 10)
		if err != nil {
			fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error getting jobnet info: %s", err.Error()))
			return FAILED
		}
		fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "%s with registry number %s is completed.", jobnet_id, run_jobnet_id))

		// Check jobnet run status and exit code.
		if jobnet_run_info.Jobnet_status != targetJobnetStatus && jobnet_run_info.Job_status != targetJobStatus {
			fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Unexpected Jobnet status. Jobnet_status: %s, Job_status: %s, Exit_cd: %d", jobnet_run_info.Jobnet_status, jobnet_run_info.Job_status, jobnet_run_info.Exit_cd))
			return FAILED
		}

		return PASSED
	}

	tc_71.Set_function(tc_func)
	t.Add_testcase(*tc_71)

	//TESTCASE 74
	tc_74 := t.New_testcase(100001, "Agent Restart [Covered with Agent servive stop while 1000  parallel jobnets are running]")
	tc_func = func() common.Testcase_status {

		// Test Case: "Execute 800 jobnets simutaneously."
		//
		// Task:      "Restart the Agent"
		//
		// Result:    "The JOB returns with an error"
		var jobnet_id string = "Icon_1"                 // This value must be Jobnet_ID that you want to enable
		var jobnet_name string = "TICKET844_TESTCASE74" // This value must be Jobnet_Name that you want to enable.
		var actual_jobnet_run_id string = "Icon_800"    // This is the actual parallel jobnet_id that we actually want to run.
		var job_process_count_amt int = 800             // Enter total Job icons you want to check here. Must be same with Jobnet total jobs.
		var timeout_minute int = 90                     // Enter timeout in minutes for Process count.
		var jobnet_run_manage_id string

		// 1. Enable 1 Icon jobnet.
		// 2. Enable 10 Icon Jobnet.
		// 3. Enable 100 Icon Jobnet.
		// 4. Enable 800 Icon Jobnet.
		// 5. Jobarg_cleanup_start.
		// 6. Run Jobnet.
		// 7. Count unitl desired procces count is met.
		// 8. Restart linux Jobarranger agent.
		// 9. Counting Job Process Count again.

		if lib.Run_enable_jobnet(jobnet_id, jobnet_name) &&
			lib.Run_enable_jobnet("Icon_10", "Icon_10") &&
			lib.Run_enable_jobnet("Icon_100", "Icon_100") &&
			lib.Run_enable_jobnet("Icon_800", "Icon_800") &&
			lib.Run_Jobarg_cleanup_linux() &&
			func() bool {
				var result bool
				result, jobnet_run_manage_id = lib.Run_Jobnet(actual_jobnet_run_id)
				return result
			}() &&
			lib.Run_Job_process_count(job_process_count_amt, timeout_minute) &&
			lib.Run_Restart_Linux_Jaz_agent() &&
			lib.Run_Job_process_count(0, 1) &&
			lib.Run_Job_Status_Check_For_Error(timeout_minute, jobnet_run_manage_id) {
			fmt.Println("All operations completed successfully")
			return PASSED
		} else {
			return FAILED
		}
	}

	tc_74.Set_function(tc_func)
	t.Add_testcase(*tc_74)

	//TESTCASE 75
	tc_75 := t.New_testcase(75, "Parallel 800 jobs with loop (8hr).")
	tc_func = func() common.Testcase_status {

		// Test Case: "Parallel 800 jobs with loop (8hr)."
		//
		// Task:      "Execute jobnet"
		//
		// Result:    make sure that no zombies occur and there are no JOBICONs that do not end
		//             - JOB only has exit 0
		//             For linux => ps -aux | grep defunct
		//             For window=> tasklist /FI "STATUS eq NOT RESPONDING"

		var jobnet_id string = "Icon_1"                 // This value must be Jobnet_ID that you want to run.
		var jobnet_name string = "TICKET844_TESTCASE75" // This value must be Jobnet_Name that you want to run.
		// Since these are default values, you don't need to set it, unless you want to change it for testing purpose.

		var jobnet_run_manage_id, last_std_out string
		var jobnet_run_info *common.Jobnet_run_info
		var num int
		var result bool

		// 1. Enable 1 Icon jobnet.
		// 2. Taking current time snapshot to calculate elapsed 8 hour time.
		// 3. Jobarg_cleanup_start
		// 4. Run Jobnet
		// 5. Get Jobnet Info
		// 6. Get jobarg_get last std out value
		// 7. Convert std out value to integer
		// 8. Check if Std out value is met or not.
		// 9. Check if Jobnet Finished successfully with no zombie process or not.

		start_time := time.Now()
		if lib.Run_enable_jobnet(jobnet_id, jobnet_name) &&
			lib.Run_Jobarg_cleanup_linux() &&
			func() bool {
				result, jobnet_run_manage_id = lib.Run_Jobnet(jobnet_id)
				return result
			}() &&
			func() bool {
				result, jobnet_run_info = lib.Run_Jobarg_get_jobnet_run_info(jobnet_run_manage_id)
				return result
			}() &&
			func() bool {
				result, last_std_out = lib.Run_Jobarg_Get_LastSTDOUT(jobnet_run_manage_id)
				return result
			}() &&
			func() bool {
				result, num = lib.Str_To_Int(last_std_out)
				return result
			}() &&
			func() bool {
				if num != 95 {
					fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error: Looping count (total iteration) is incorrect. Looping count: %d", num))
					return false
				}
				return true
			}() &&
			func() bool {
				// Calculating whether jobnet took actual 8 hour or not.
				end_time := time.Now()
				duration := end_time.Sub(start_time)
				if duration <= 8*time.Hour {
					fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error: Jobnet ended earlier than expected 8 hour time. Total Elapsed time: %s", duration.String()))
					return false
				}
				return true
			}() &&
			lib.Run_Check_Jobnet_Finish_With_No_Zombie_Process(jobnet_run_info) {
			fmt.Println("All operations completed successfully")
			return PASSED
		} else {
			return FAILED
		}
	}

	tc_75.Set_function(tc_func)
	t.Add_testcase(*tc_75)

	//TESTCASE 76
	tc_76 := t.New_testcase(76, "kill jobarg_command process ")
	tc_func = func() common.Testcase_status {
		// enable common jobnet
		if err := lib.Jobarg_enable_jobnet("Icon_1", "jobicon_linux"); err != nil {
			lib.Logi(common.LOG_LEVEL_ERR, "Failed to enable jobnet, Error: %s", err)
			return FAILED
		}
		return RunJob800AndKillOneJobIconWithJobargCommand("Icon_800", 800, 4, tc_76, common.Client)
	}

	tc_76.Set_function(tc_func)
	t.Add_testcase(*tc_76)

	//TESTCASE 77
	tc_77 := t.New_testcase(77, "force stop  running job icon")
	tc_func = func() common.Testcase_status {
		// enable common jobnet
		if err := lib.Jobarg_enable_jobnet("Icon_1", "jobicon_linux"); err != nil {
			lib.Logi(common.LOG_LEVEL_ERR, "Failed to enable jobnet, Error: %s", err)
			return FAILED
		}
		return RunJob800AndForceStopOneJobIcon("Icon_800", 800, 4, tc_77, common.Client)
	}

	tc_77.Set_function(tc_func)
	t.Add_testcase(*tc_77)
}

func RunJob800AndKillOneJobIconWithJobargCommand(jobnetId string, processCount int, processCheckTimeout int, testcase *common.TestCase, sshClient *ssh.Client) common.Testcase_status {

	/******************
	Pre-Operation State
	******************/
	lib.Jobarg_cleanup_linux()

	envs, _ := lib.Get_str_str_map("JA_HOSTNAME", "oss.linux", "JA_CMD", "sleep 500")
	run_jobnet_id, error := lib.Jobarg_exec_E(jobnetId, envs)
	if error != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error: %s, std_out: %s", error.Error(), run_jobnet_id))
		return FAILED
	}
	fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "%s has been successfully run with registry number: %s", jobnetId, run_jobnet_id))

	// Wait for all jobs to be in running state
	lib.JobSleepProcessCountCheck(processCount, processCheckTimeout, sshClient)

	fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "Process count has reached %d", processCount))

	/**************
	Operation State
	***************/

	// Run the command
	command := "ps -aux | grep jobarg_command | grep -v grep"
	output, err := lib.Ssh_exec(command)

	if err != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Failed to run command: %v", err))
		return FAILED
	}

	// Process the output
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	var pid string

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) > 1 {
			pid = fields[1] // Get the PID
		}
	}

	// Kill the first process if any are found
	if processCount > 0 {
		killCommand := fmt.Sprintf("kill %s", pid)
		_, err := lib.Ssh_exec(killCommand)
		if err != nil {
			fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Failed to kill process %s: %v", pid, err))
		}
		fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "Kill Job Icon using Process ID: %s", pid))
	} else {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "No processes to kill."))
		return FAILED
	}

	fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "Job icon is being aborted..."))

	/***************
	Expected Results
	****************/

	sleepDuration := 1 * time.Minute
	time.Sleep(sleepDuration)
	DBQuery := "SELECT after_value FROM ja_run_value_after_table WHERE after_value LIKE '%Check job status(end) failed.%'"

	// Execute the query
	rows, err := common.DB.Query(DBQuery)
	if err != nil {
		// If there was an error executing the query, print an error message
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error executing query: %v\n", err))
		return FAILED
	}
	// Ensure the rows are closed after processing
	defer rows.Close()

	// Variables for counting occurrences and storing result
	occurrenceCount := 0
	var afterValue string

	// Iterate over the rows
	for rows.Next() {
		// Scan the result into the variable
		err := rows.Scan(&afterValue)
		if err != nil {
			// If there was an error scanning the row, print an error message
			fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error retrieving result: %v\n", err))
			return FAILED
		}
		occurrenceCount++ // Increment the occurrence count
	}

	// Check the number of occurrences
	if occurrenceCount == 1 {
		// Exactly one occurrence found
		fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "One occurrence found. Job Icon Result: %s", afterValue))
	} else if occurrenceCount == 0 {
		// No occurrences found
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "No occurrences found with 'Check job status(end) failed.'."))
		return FAILED
	} else {
		// Multiple occurrences found
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Multiple occurrences found: %d", occurrenceCount))
		return FAILED
	}

	// Wait for all jobs to be in running state
	err = lib.JobSleepProcessCountCheck(0, 10, sshClient)

	if err != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error: %s, Failed to get process count.", err.Error()))
		return FAILED
	} else {
		fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "Process count remain %d", 0))
	}

	// check zombie process
	zombieProcessCount, err := lib.CheckZombieProcess(2, sshClient)
	if err != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error checking zombie process: %s", err.Error()))
		return FAILED
	}

	if zombieProcessCount != 0 {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "There are zombie processes: %d", zombieProcessCount))
		// return FAILED
	} else {
		fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "There is no zombie process."))
	}

	return PASSED
}

func RunJob800AndForceStopOneJobIcon(jobnetId string, processCount int, processCheckTimeout int, testcase *common.TestCase, sshClient *ssh.Client) common.Testcase_status {

	/******************
	Pre-Operation State
	******************/

	lib.Jobarg_cleanup_linux()

	envs, _ := lib.Get_str_str_map("JA_HOSTNAME", "oss.linux", "JA_CMD", "sleep 500")
	run_jobnet_id, error := lib.Jobarg_exec_E(jobnetId, envs)

	if error != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error: %s, std_out: %s", error.Error(), run_jobnet_id))
		return FAILED
	}
	fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "%s has been successfully run with registry number: %s", jobnetId, run_jobnet_id))

	// Wait for all jobs to be in running state
	lib.JobSleepProcessCountCheck(processCount, processCheckTimeout, sshClient)

	fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "Process count has reached %d", processCount))

	/**************
	Operation State
	***************/

	DBQuery := "SELECT inner_job_id FROM ja_run_job_table jr WHERE jr.job_type = 4 AND jr.status = 2"

	// Execute the query
	rows, err := common.DB.Query(DBQuery)
	if err != nil {
		// If there was an error executing the query, print an error message
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error executing query: %v\n", err))
		return FAILED
	}
	// Ensure the rows are closed after processing
	defer rows.Close()

	var inner_job_id string

	// Iterate over the rows
	for rows.Next() {
		// Scan the result into the variable
		err := rows.Scan(&inner_job_id)
		if err != nil {
			// If there was an error scanning the row, print an error message
			fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error retrieving inner_job_id: %v\n", err))
			return FAILED
		}
	}

	// Convert string to int64 (if needed)
	convert_inner_job_id, err := strconv.ParseInt(inner_job_id, 10, 64)
	if err != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "Error converting inner_jobnet_id to int64: %v", err))
	}

	fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "Kill Job Icon using inner_job_id %d", convert_inner_job_id))

	// Execute the update query
	_, err = lib.ExecuteQuery("UPDATE ja_run_job_table SET method_flag = 3 WHERE inner_job_id = $1", convert_inner_job_id)
	if err != nil {
		// Log error and return FAILED if the query execution fails
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error: %s, Failed to abort the job icon.", err.Error()))
		return FAILED
	}

	// /***************
	// Expected Results
	// ****************/

	sleepDuration := 1 * time.Minute
	time.Sleep(sleepDuration)

	DBQuery = "SELECT after_value FROM ja_run_value_after_table WHERE after_value LIKE '%Aborted the job%'"

	// Execute the query
	rows, err = common.DB.Query(DBQuery)
	if err != nil {
		// If there was an error executing the query, print an error message
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error executing query: %v\n", err))
		return FAILED
	}
	// Ensure the rows are closed after processing
	defer rows.Close()

	// Variables for counting occurrences and storing result
	occurrenceCount := 0

	var afterValue string
	// Iterate over the rows
	for rows.Next() {
		// Scan the result into the variable
		err := rows.Scan(&afterValue)
		if err != nil {
			// If there was an error scanning the row, print an error message
			fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error retrieving result: %v\n", err))
			return FAILED
		}
		occurrenceCount++ // Increment the occurrence count
	}

	// Check the number of occurrences
	if occurrenceCount == 1 {
		// Exactly one occurrence found
		fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "One occurrence found. Job Icon Result: %s", afterValue))
	} else if occurrenceCount == 0 {
		// No occurrences found
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "No occurrences found with 'Aborted the job'."))
		return FAILED
	} else {
		// Multiple occurrences found
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Multiple occurrences found: %d", occurrenceCount))
		return FAILED
	}

	// Wait for all jobs to be in running state
	err = lib.JobSleepProcessCountCheck(0, 10, sshClient)

	if err != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error: %s, Failed to get process count.", err.Error()))
		return FAILED
	} else {
		fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "Process count remain %d", 0))
	}

	zombieProcessCount, err := lib.CheckZombieProcess(2, sshClient)
	if err != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error checking zombie process: %s", err.Error()))
		return FAILED
	}

	if zombieProcessCount != 0 {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "There are zombie processes: %d", zombieProcessCount))
		// return FAILED
	} else {
		fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "There is no zombie process."))
	}

	return PASSED
}
