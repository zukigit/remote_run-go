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
	//TESTCASE 001
	tc_76 := t.New_testcase(76, "kill jobarg_command process")
	tc_func := func() common.Testcase_status {
		return RunJob800AndKillOneJobIconWithJobargCommand("TICKET844_TESTCASE76-77JOB800", 800, 4, tc_76, common.Client)
	}
	tc_76.Set_function(tc_func)
	t.Add_testcase(*tc_76)

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
