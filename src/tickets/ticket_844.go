package tickets

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/zukigit/remote_run-go/src/common"
	"github.com/zukigit/remote_run-go/src/dao"
	"github.com/zukigit/remote_run-go/src/lib"
	"golang.org/x/crypto/ssh"
)

type Ticket_844 struct {
	no          uint
	description string
	testcases   []dao.TestCase
}

func (t *Ticket_844) New_testcase(testcase_id uint, testcase_description string) *dao.TestCase {
	return dao.New_testcase(testcase_id, testcase_description)
}

func (t *Ticket_844) Get_no() uint {
	return t.no
}

func (t *Ticket_844) Get_dsctn() string {
	return t.description
}

func (t *Ticket_844) Add_testcase(tc dao.TestCase) {
	t.testcases = append(t.testcases, tc)
}

func (t *Ticket_844) Get_testcases() []dao.TestCase {
	return t.testcases
}

// Enter your ticket information here
func (t *Ticket_844) Set_values() {
	t.no = 844 // Enter your ticket id
	t.description = "kill jobarg_command process."
}

// Add your test case here
func (t *Ticket_844) Add_testcases() {

	//TESTCASE 74
	tc_74 := t.New_testcase(74, "Agent Restart [Covered with Agent servive stop while 1000  parallel jobnets are running]")
	tc_func := func() common.Testcase_status {

		// Test Case: "Execute 800 jobnets simutaneously."
		//
		// Task:      "Restart the Agent"
		//
		// Result:    "The JOB returns with an error"

		var jobnet_id string = "Icon_1"   // This value must be the same value as Jobnet_ID
		var job_process_count_amt int = 1 // Enter total Job icons you want to check here. Must be same with Jobnet total jobs.
		var timeout_minute int = 10       // Enter timeout in minutes for Process count.
		var jobnet_run_manage_id string

		// 1. Jobarg_cleanup_start.
		// 2. Run Jobnet.
		// 3. Count unitl desired procces count is met.
		// 4. Restart linux Jobarranger agent.
		// 5. Counting Job Process Count again.

		if Run_Jobarg_cleanup_linux(tc_74) &&
			func() bool {
				var result bool
				result, jobnet_run_manage_id = Run_Jobnet(tc_74, jobnet_id)
				return result
			}() &&
			Run_Job_process_count(tc_74, job_process_count_amt, timeout_minute) &&
			Run_Restart_Linux_Jaz_agent(tc_74) &&
			Run_Job_process_count(tc_74, 0, 1) &&
			Tc_74_Job_Status_Check_For_Error(tc_74, jobnet_run_manage_id) {
			fmt.Println("All operations completed successfully")
			return PASSED
		} else {
			return FAILED
		}
	}

	tc_74.Set_function(tc_func)
	t.Add_testcase(*tc_74)

	//TESTCASE 75
	tc_75 := t.New_testcase(75, "Agent Restart [Covered with Agent servive stop while 1000  parallel jobnets are running]")
	tc_func = func() common.Testcase_status {

		// Test Case: "Parallel 800 jobs with loop (8hr)."
		//
		// Task:      "Execute jobnet"
		//
		// Result:    make sure that no zombies occur and there are no JOBICONs that do not end
		//             - JOB only has exit 0
		//             For linux => ps -aux | grep defunct
		//             For window=> tasklist /FI "STATUS eq NOT RESPONDING"

		var jobnet_id string = "TICKET844_TESTCASE75" // This value must be the same value as Jobnet_ID
		var jobnet_run_manage_id, last_std_out string
		var jobnet_run_info *common.Jobnet_run_info
		var num int
		var result bool

		// 1. Taking current time snapshot to calculate elapsed 8 hour time.
		// 2. Jobarg_cleanup_start
		// 3. Run Jobnet
		// 4. Get Jobnet Info
		// 5. Get jobarg_get last std out value
		// 6. Convert std out value to integer
		// 7. Check if Std out value is met or not.
		// 8. Check if Jobnet Finished successfully with no zombie process or not.

		start_time := time.Now()
		if Run_Jobarg_cleanup_linux(tc_75) &&
			func() bool {
				result, jobnet_run_manage_id = Run_Jobnet(tc_75, jobnet_id)
				return result
			}() &&
			func() bool {
				result, jobnet_run_info = Run_Jobarg_get_jobnet_run_info(tc_75, jobnet_run_manage_id)
				return result
			}() &&
			func() bool {
				result, last_std_out = Run_Jobarg_Get_LastSTDOUT(tc_75, jobnet_run_manage_id)
				return result
			}() &&
			func() bool {
				result, num = Str_To_Int(tc_75, last_std_out)
				return result
			}() &&
			num < 95 &&
			func() bool {
				// Calculating whether jobnet took actual 8 hour or not.
				end_time := time.Now()
				duration := end_time.Sub(start_time)
				if duration <= 8*time.Hour {
					fmt.Println(tc_75.Err_log("Error: Jobnet ended earlier than expected 8 hour time. Total Elapsed time: %s", duration.String()))
					return false
				}
				return true
			}() &&
			Run_Check_Jobnet_Finish_With_No_Zombie_Process(tc_75, jobnet_run_info) {
			fmt.Println("All operations completed successfully")
			return PASSED
		} else {
			return FAILED
		}
	}

	tc_75.Set_function(tc_func)
	t.Add_testcase(*tc_75)

	//TESTCASE 76
	tc_76 := t.New_testcase(76, "kill jobarg_command process")
	tc_func = func() common.Testcase_status {
		return RunJob800AndKillOneJobIconWithJobargCommand("TICKET844_TESTCASE76-77JOB800", 800, 4, tc_76, common.Client)
	}

	tc_76.Set_function(tc_func)
	t.Add_testcase(*tc_76)

	//TESTCASE 77
	tc_77 := t.New_testcase(77, "force stop  running job icon")
	tc_func = func() common.Testcase_status {
		return RunJob800AndForceStopOneJobIcon("TICKET844_TESTCASE76-77JOB800", 800, 4, tc_77, common.Client)
	}

	tc_77.Set_function(tc_func)
	t.Add_testcase(*tc_77)
}

func RunJob800AndKillOneJobIconWithJobargCommand(jobnetId string, processCount int, processCheckTimeout int, testcase *dao.TestCase, sshClient *ssh.Client) common.Testcase_status {

	/******************
	Pre-Operation State
	******************/
	lib.Jobarg_cleanup_linux()

	run_jobnet_id, error := lib.Jobarg_exec(jobnetId)
	if error != nil {
		fmt.Println(testcase.Err_log("Error: %s, std_out: %s", error.Error(), run_jobnet_id))
		return FAILED
	}
	fmt.Println(testcase.Info_log("%s has been successfully run with registry number: %s", jobnetId, run_jobnet_id))

	// Wait for all jobs to be in running state
	lib.JobSleepProcessCountCheck(processCount, processCheckTimeout, sshClient)

	fmt.Println(testcase.Info_log("Process count has reached %d", processCount))

	/**************
	Operation State
	***************/

	// Run the command
	command := "ps -aux | grep jobarg_command | grep -v grep"
	output, err := lib.Ssh_exec(command)

	if err != nil {
		fmt.Println(testcase.Err_log("Failed to run command: %v", err))
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
			fmt.Println(testcase.Err_log("Failed to kill process %s: %v", pid, err))
		}
		fmt.Println(testcase.Info_log("Kill Job Icon using Process ID: %s", pid))
	} else {
		fmt.Println(testcase.Err_log("No processes to kill."))
		return FAILED
	}

	fmt.Println(testcase.Info_log("Job icon is being aborted..."))

	/***************
	Expected Results
	****************/

	sleepDuration := 2 * time.Minute
	time.Sleep(sleepDuration)
	DBQuery := "SELECT after_value FROM ja_run_value_after_table WHERE after_value LIKE '%Check job status(end) failed.%'"

	// Execute the query
	rows, err := common.DB.Query(DBQuery)
	if err != nil {
		// If there was an error executing the query, print an error message
		fmt.Println(testcase.Err_log("Error executing query: %v\n", err))
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
			fmt.Println(testcase.Err_log("Error retrieving result: %v\n", err))
			return FAILED
		}
		occurrenceCount++ // Increment the occurrence count
	}

	// Check the number of occurrences
	if occurrenceCount == 1 {
		// Exactly one occurrence found
		fmt.Println(testcase.Info_log("One occurrence found. Job Icon Result: %s", afterValue))
	} else if occurrenceCount == 0 {
		// No occurrences found
		fmt.Println(testcase.Err_log("No occurrences found with 'Check job status(end) failed.'."))
		return FAILED
	} else {
		// Multiple occurrences found
		fmt.Println(testcase.Err_log("Multiple occurrences found: %d", occurrenceCount))
		return FAILED
	}

	// Wait for all jobs to be in running state
	err = lib.JobSleepProcessCountCheck(0, 10, sshClient)

	if err != nil {
		fmt.Println(testcase.Err_log("Error: %s, Failed to get process count.", err.Error()))
		return FAILED
	} else {
		fmt.Println(testcase.Info_log("Process count remain %d", 0))
	}

	// check zombie process
	zombieProcessCount, err := lib.CheckZombieProcess(2, sshClient)
	if err != nil {
		fmt.Println(testcase.Err_log("Error checking zombie process: %s", err.Error()))
		return FAILED
	}

	if zombieProcessCount != 0 {
		fmt.Println(testcase.Err_log("There are zombie processes: %d", zombieProcessCount))
		// return FAILED
	} else {
		fmt.Println(testcase.Info_log("There is no zombie process."))
	}

	return PASSED
}

func RunJob800AndForceStopOneJobIcon(jobnetId string, processCount int, processCheckTimeout int, testcase *dao.TestCase, sshClient *ssh.Client) common.Testcase_status {

	/******************
	Pre-Operation State
	******************/

	lib.Jobarg_cleanup_linux()

	run_jobnet_id, error := lib.Jobarg_exec(jobnetId)

	if error != nil {
		fmt.Println(testcase.Err_log("Error: %s, std_out: %s", error.Error(), run_jobnet_id))
		return FAILED
	}
	fmt.Println(testcase.Info_log("%s has been successfully run with registry number: %s", jobnetId, run_jobnet_id))

	// Wait for all jobs to be in running state
	lib.JobSleepProcessCountCheck(processCount, processCheckTimeout, sshClient)

	fmt.Println(testcase.Info_log("Process count has reached %d", processCount))

	/**************
	Operation State
	***************/

	DBQuery := "SELECT inner_jobnet_id FROM ja_run_job_table WHERE inner_jobnet_main_id = $1 LIMIT 1"

	// Execute the query
	rows, err := common.DB.Query(DBQuery, run_jobnet_id)
	if err != nil {
		// If there was an error executing the query, print an error message
		fmt.Println(testcase.Err_log("Error executing query: %v\n", err))
		return FAILED
	}
	// Ensure the rows are closed after processing
	defer rows.Close()

	// Variable to store the result
	var afterValue string

	// Iterate over the rows
	for rows.Next() {
		// Scan the result into the variable
		err := rows.Scan(&afterValue)
		if err != nil {
			// If there was an error scanning the row, print an error message
			fmt.Println(testcase.Err_log("Error retrieving result: %v\n", err))
			return FAILED
		}
	}

	// Convert string to int64 (if needed)
	intValue, err := strconv.ParseInt(afterValue, 10, 64)
	if err != nil {
		fmt.Println(testcase.Info_log("Error converting afterValue to int64: %v", err))
	}
	fmt.Println(testcase.Info_log("Kill Job Icon using innerjobnet id %d", intValue))

	DBQuery = "UPDATE ja_run_job_table SET method_flag = 3 WHERE inner_jobnet_id = $1"

	// Execute the update query
	_, err = lib.ExecuteQuery("UPDATE ja_run_job_table SET method_flag = 3 WHERE inner_jobnet_id = $1", intValue) // Pass the necessary argument for $1
	if err != nil {
		// Log error and return FAILED if the query execution fails
		fmt.Println(testcase.Err_log("Error: %s, Failed to abort the job icon.", err.Error()))
		return FAILED
	}

	fmt.Println(testcase.Info_log("Job icon is being aborted..."))

	/***************
	Expected Results
	****************/

	sleepDuration := 2 * time.Minute
	time.Sleep(sleepDuration)

	DBQuery = "SELECT after_value FROM ja_run_value_after_table WHERE after_value LIKE '%Aborted the job%'"

	// Execute the query
	rows, err = common.DB.Query(DBQuery)
	if err != nil {
		// If there was an error executing the query, print an error message
		fmt.Println(testcase.Err_log("Error executing query: %v\n", err))
		return FAILED
	}
	// Ensure the rows are closed after processing
	defer rows.Close()

	// Variables for counting occurrences and storing result
	occurrenceCount := 0

	// Iterate over the rows
	for rows.Next() {
		// Scan the result into the variable
		err := rows.Scan(&afterValue)
		if err != nil {
			// If there was an error scanning the row, print an error message
			fmt.Println(testcase.Err_log("Error retrieving result: %v\n", err))
			return FAILED
		}
		occurrenceCount++ // Increment the occurrence count
	}

	// Check the number of occurrences
	if occurrenceCount == 1 {
		// Exactly one occurrence found
		fmt.Println(testcase.Info_log("One occurrence found. Job Icon Result: %s", afterValue))
	} else if occurrenceCount == 0 {
		// No occurrences found
		fmt.Println(testcase.Err_log("No occurrences found with 'Aborted the job'."))
		return FAILED
	} else {
		// Multiple occurrences found
		fmt.Println(testcase.Err_log("Multiple occurrences found: %d", occurrenceCount))
		return FAILED
	}

	// Wait for all jobs to be in running state
	err = lib.JobSleepProcessCountCheck(0, 10, sshClient)

	if err != nil {
		fmt.Println(testcase.Err_log("Error: %s, Failed to get process count.", err.Error()))
		return FAILED
	} else {
		fmt.Println(testcase.Info_log("Process count remain %d", 0))
	}

	zombieProcessCount, err := lib.CheckZombieProcess(2, sshClient)
	if err != nil {
		fmt.Println(testcase.Err_log("Error checking zombie process: %s", err.Error()))
		return FAILED
	}

	if zombieProcessCount != 0 {
		fmt.Println(testcase.Err_log("There are zombie processes: %d", zombieProcessCount))
		// return FAILED
	} else {
		fmt.Println(testcase.Info_log("There is no zombie process."))
	}

	return PASSED
}

// To run jobnet
//
// Returns
//   - True and jobnet_run_manage_id if worked.
//   - False and jobnet_run_manage_id if failed.
func Run_Jobnet(testcase *dao.TestCase, jobnet_id string) (bool, string) {
	jobnet_run_manage_id, err := lib.Jobarg_exec(jobnet_id)
	if err != nil {
		fmt.Println(testcase.Err_log("Error: Failed when trying to run the job. %s", err.Error()))
		fmt.Println(testcase.Err_log("Error: Did you imported the necessary xml files?(%s)", jobnet_id))
		return false, jobnet_run_manage_id
	} else {
		fmt.Println(testcase.Info_log("Info: Job has start running."))
		return true, jobnet_run_manage_id
	}
}

// To restart linux jobarg agentd
//
// Returns
//   - True if passed.
//   - False if error.
func Run_Restart_Linux_Jaz_agent(testcase *dao.TestCase) bool {
	// Restartng Jobarranger linux Agent.
	fmt.Println(testcase.Info_log("Info: Restarting Jobarranger Agent now."))
	err := lib.Restart_jaz_agent_linux()
	if err != nil {
		fmt.Println(testcase.Err_log("Error: Failted at restarting Jobarranger Agent. %s", err.Error()))
		return false
	} else {
		fmt.Println(testcase.Info_log("Info: Jobarranger Agent has successfully restarted."))
		return true
	}
}

// To count jobnet process
//
// Returns
//   - True if met with desired process count.
//   - False if failed to met desired process count.
func Run_Job_process_count(testcase *dao.TestCase, job_process_count_amt int, timeout_minute int) bool {
	// Counting Job Process Count
	fmt.Println(testcase.Info_log("Info: Counting Job processes."))
	err := JobProcessCountCheck_with_process_counter(job_process_count_amt, timeout_minute, common.Client)
	if err != nil {
		fmt.Println(testcase.Err_log("Error: Failed at counting or failed to meet required process count amount. %s", err.Error()))
		return false
	} else {
		fmt.Println(testcase.Info_log("Info: Job process count met."))
		return true
	}
}

// To clean up Jobarg_linux.
//
// Returns
//   - True if worked.
//   - False if failed.
func Run_Jobarg_cleanup_linux(testcase *dao.TestCase) bool {
	// Initiating Jobarg_cleanup
	fmt.Println(testcase.Info_log("Info: Starting Jobarg_cleanup"))
	err := lib.Jobarg_cleanup_linux()
	if err != nil {
		fmt.Println(testcase.Err_log("Error: %s", err.Error()))
		return false
	}
	fmt.Println(testcase.Info_log("Info: Jobarg_cleanup Success."))
	return true
}

// To get jobnet_run_info
//
// Returns
//   - True if worked.
//   - False if failed.
func Run_Jobarg_get_jobnet_run_info(testcase *dao.TestCase, jobnet_run_manage_id string) (bool, *common.Jobnet_run_info) {
	//In case if getting jobnet info failed
	fmt.Println(testcase.Info_log("Info: Getting Jobnet Info."))
	jobnet_run_info, err := Jobarg_get_jobnet_run_info_with_timer(jobnet_run_manage_id)
	if err != nil {
		fmt.Println(testcase.Err_log("Error: Failted at getting jobnet run info. %s. Jobnet Management Id: %s ", err.Error(), jobnet_run_manage_id))
		_, lastSTDErr := lib.Jobarg_get_LASTSTDERR(jobnet_run_manage_id)
		_, lastSTDOut := lib.Jobarg_get_LASTSTDOUT(jobnet_run_manage_id)
		fmt.Println(testcase.Err_log("Error: LAST STD ERR : %s. LAST STD OUT : %s.", lastSTDErr, lastSTDOut))
		if lastSTDErr == nil && lastSTDOut == nil {
			fmt.Println(testcase.Err_log("Error: Maybe hostname in one of your job is incorrect or not found?"))
		}
		return false, jobnet_run_info
	} else {
		fmt.Println(testcase.Info_log("Info: Successfully received getting Jobnet Info."))
		return true, jobnet_run_info
	}
}

// To check if jobnet finished with no zombie process
//
// Returns
//   - True if worked.
//   - False if failed.
func Run_Check_Jobnet_Finish_With_No_Zombie_Process(testcase *dao.TestCase, jobnet_run_info *common.Jobnet_run_info) bool {

	//Checking whether Jobnet is finished or not
	fmt.Println(testcase.Info_log("Info: Checking for Job Status, Jobnet Status and Exit CD."))
	if jobnet_run_info.Job_status == "NORMAL" && jobnet_run_info.Jobnet_status == "END" && jobnet_run_info.Exit_cd == 0 {

		// Checking whether "checking at zombie process" process failed or not.
		fmt.Println(testcase.Info_log("Info: Success. Job Status = %s, Jobnet Status = %s, Exit_CD = %d", jobnet_run_info.Job_status, jobnet_run_info.Jobnet_status, jobnet_run_info.Exit_cd))
		fmt.Println(testcase.Info_log("Info: Checking for Zombie Process."))
		zombieProcessCount, err := lib.CheckZombieProcess(1, common.Client)
		if err != nil {
			fmt.Println(testcase.Err_log("Error: Failed at checking zombie Process."))
			fmt.Println(testcase.Err_log(err.Error()))
			// Checking whether zombie process exist or not.
			if zombieProcessCount > 0 {
				fmt.Println(testcase.Err_log("Error: There is zombie Process left."))
				fmt.Println(testcase.Err_log(err.Error()))
			}

			return false
		}
		fmt.Println(testcase.Info_log("Info: There's no Zombie Processes. Success."))
		return true
	}

	fmt.Println(testcase.Err_log("Error: Either Job Status or Jobnet Status or Exit CD doesn't meet required contidions.(Normal, End, 0)"))
	fmt.Println(testcase.Err_log("Jobnet_status: %s, Job_status: %s, Exit_cd: %d", jobnet_run_info.Jobnet_status, jobnet_run_info.Job_status, jobnet_run_info.Exit_cd))
	return false
}

// To retrieve value from Jobarg_get_LASTSTDOUT
//
// Returns
//   - True and LASTSTDOUT if worked.
//   - False and nil if failed.
func Run_Jobarg_Get_LastSTDOUT(testcase *dao.TestCase, jobnet_run_manage_id string) (bool, string) {
	last_std_out, err := lib.Jobarg_get_LASTSTDOUT(jobnet_run_manage_id)
	if err != nil {
		fmt.Println(testcase.Err_log("Error: Failed to get LASTSTDOUT value. %s", err))
		return false, ""
	}
	fmt.Println(testcase.Info_log("Info: LASTSTDOUT value: %s", last_std_out))
	return true, last_std_out
}

// To convert from string to int
//
// Returns
//   - True and number if worked.
//   - False and 0 if failed.

func Str_To_Int(testcase *dao.TestCase, numberStr string) (bool, int) {
	numberInt, err := strconv.Atoi(numberStr)
	if err != nil {
		fmt.Println(testcase.Err_log("Failed to conver String to integer. %s", err.Error()))
		return false, 0
	}
	return true, numberInt
}

// Checking 'Job_Status' to get the value 'Error' with 'Jobarg_get' command."
//
// Returns
//   - True if Passed.
//   - False if Not.
func Tc_74_Job_Status_Check_For_Error(testcase *dao.TestCase, jobnet_run_manage_id string) bool {
	// Checking 'Job_Status' to get the value 'Error' with 'Jobarg_get' command."
	fmt.Println(testcase.Info_log("Info: Total 0 Job process count met. Waiting for Jobnet to end with RUNERR.(Running Error)"))

	// Setting timeout duration
	timeoutDuration := 30 * time.Minute
	timeout := time.Now().Add(timeoutDuration)

	// taking current time snapshot
	start := time.Now()

	for time.Now().Before(timeout) {

		// Calculating elapsed time.
		elapsed := time.Since(start)

		// Extract hours, minutes, and seconds from elapsed time
		hours := int(elapsed.Hours())
		minutes := int(elapsed.Minutes()) % 60
		seconds := int(elapsed.Seconds()) % 60

		jobstatus, err := lib.Jobarg_get_JA_JOBSTATUS(jobnet_run_manage_id)
		if err == nil && jobstatus == "ERROR" {
			fmt.Println(testcase.Info_log("Info: The job returns with Error. Success."))

			return true
		}
		// Print in HH:MM:SS format with \r to overwrite the line
		fmt.Printf("\rElapsed Time: %02d:%02d:%02d", hours, minutes, seconds)
		time.Sleep(1 * time.Millisecond)
	}
	return false
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
//
// Fancyyyyyyy. slay. Sigma. Chad.
// Ctl + C programmer.
func JobProcessCountCheck_with_process_counter(targetProcessCount int, timeoutDuration int, client *ssh.Client) error {
	// set timeout
	timeout := time.After(time.Duration(timeoutDuration) * time.Minute)
	var index int
	for {
		select {
		case <-timeout:
			return fmt.Errorf("timeout after %d minutes", timeoutDuration)
		default:
			currentProcessCountStr, err := lib.GetOutputStrFromSSHCommand(client, "ps -aux | grep /etc/jobarranger/extendedjob/ | grep -v grep | wc -l")

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
			fmt.Printf("\rProcess Count: %d", currentProcessCount)
			time.Sleep(1 * time.Second)
			index++
		}

	}
}

// Jobarg_get_jobnet_run_info waits util the jobnet is done or get error and returns Jobnet run info but with timer.
// Fancyyyyyyy. slay. Sigma. Chad.
// Ctl + C programmer.
func Jobarg_get_jobnet_run_info_with_timer(registry_number string) (*common.Jobnet_run_info, error) {
	var jobnet_status, job_status, std_out, std_error string
	var err error
	var index int
	var exit_cd int64

	// taking current time snapshot
	start := time.Now()

	for {
		jobnet_status, err = lib.Jobarg_get_JA_JOBNETSTATUS(registry_number)
		if err != nil {
			return nil, err
		}

		job_status, err = lib.Jobarg_get_JA_JOBSTATUS(registry_number)
		if err != nil {
			return nil, err
		}

		if jobnet_status == "END" || (jobnet_status == "RUN" && job_status == "ERROR") {
			break
		}

		// Calculating elapsed time.
		elapsed := time.Since(start)
		// Extract hours, minutes, and seconds from elapsed time
		hours := int(elapsed.Hours())
		minutes := int(elapsed.Minutes()) % 60
		seconds := int(elapsed.Seconds()) % 60

		// Print in HH:MM:SS format with \r to overwrite the line
		//fmt.Printf("\r%02d:%02d:%02d", hours, minutes, seconds)

		lib.Spinner_log(index, lib.Formatted_log(common.INFO, "Getting jobnet[%s] run info but jobnet is not finished yet. Elapsed Time: %02d:%02d:%02d", registry_number, hours, minutes, seconds))
		time.Sleep(1 * time.Millisecond)
		index++
	}

	exit_cd, err = lib.Jobarg_get_LASTEXITCD(registry_number)
	if err != nil {
		return nil, err
	}

	std_out, err = lib.Jobarg_get_LASTSTDOUT(registry_number)
	if err != nil {
		return nil, err
	}

	std_error, err = lib.Jobarg_get_LASTSTDERR(registry_number)
	if err != nil {
		return nil, err
	}

	fmt.Println()
	return common.New_jobnet_run_info(jobnet_status, job_status, std_out, std_error, exit_cd), nil
}
