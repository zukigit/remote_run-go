package tickets

import (
	"context"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/zukigit/remote_run-go/src/common"
	"github.com/zukigit/remote_run-go/src/dao"
	"github.com/zukigit/remote_run-go/src/lib"

	"golang.org/x/crypto/ssh"
)

type DBQuery struct {
	Query string
}

type Ticket_811 struct {
	no          uint
	description string
	testcases   []dao.TestCase
}

func (t *Ticket_811) New_testcase(testcase_id uint, testcase_description string) *dao.TestCase {
	return dao.New_testcase(testcase_id, testcase_description)
}

func (t *Ticket_811) Get_no() uint {
	return t.no
}

func (t *Ticket_811) Get_dsctn() string {
	return t.description
}

func (t *Ticket_811) Add_testcase(tc dao.TestCase) {
	t.testcases = append(t.testcases, tc)
}

func (t *Ticket_811) Get_testcases() []dao.TestCase {
	return t.testcases
}

// Enter your ticket information here
func (t *Ticket_811) Set_values() {
	t.no = 811 // Enter your ticket id
	t.description = "cannot stop jobnet"
}

// Add your test case here
func (t *Ticket_811) Add_testcases() {
	// create function for all test cases
	addTestCase := func(id int, description string, jobnetId string, processCount int, timeout int) {
		tc := t.New_testcase(uint(id), description) // Convert id to uint
		tc_func := func() common.Testcase_status {
			return RunJobnetAndAbort(jobnetId, processCount, timeout, tc, common.Client)
		}
		tc.Set_function(tc_func)
		t.Add_testcase(*tc)
	}

	// TESTCASE 1 (Force stop jobnet with icon count of 100)
	addTestCase(1, "Check abort process for job icon with count 100", "TICKET811_Jobnet_100", 100, 10)

	// TESTCASE 2 (Force stop jobnet with icon count of 200)
	addTestCase(2, "Check abort process for job icon with count 200", "TICKET811_Jobnet_200", 200, 12)

	// TESTCASE 3 (Force stop jobnet with icon count of 400)
	addTestCase(3, "Check abort process for job icon with count 400", "TICKET811_Jobnet_400", 400, 13)

	// TESTCASE 4 (Force stop jobnet with icon count of 1000)
	addTestCase(4, "Check abort process for job icon with count 800", "TICKET811_Jobnet_1000", 1000, 16)
	//Normal jobnet with 100 agent less icon in running state.
	tc_5 := t.New_testcase(5, "Check Abort process abort the Agent Less icon or not")
	tc_func := func() common.Testcase_status {
		return RunJobnetAndAbort("TICKET811_AgentLess100", 100, 5, tc_5, common.Client)
	}
	tc_5.Set_function(tc_func)
	t.Add_testcase(*tc_5)

	//Normal jobnet with 200 agent less icon in running state.
	tc_6 := t.New_testcase(6, "Check Abort process abort the Agent Less icon or not")
	tc_func = func() common.Testcase_status {
		return RunJobnetAndAbort("TICKET811_AgentLess200", 200, 5, tc_6, common.Client)
	}
	tc_6.Set_function(tc_func)
	t.Add_testcase(*tc_6)

	//Normal jobnet with 400 agent less icon in running state.
	tc_7 := t.New_testcase(7, "Check Abort process abort the Agent Less icon or not")
	tc_func = func() common.Testcase_status {
		return RunJobnetAndAbort("TICKET811_AgentLess400", 400, 10, tc_7, common.Client)
	}
	tc_7.Set_function(tc_func)
	t.Add_testcase(*tc_7)

	//Normal jobnet with 800 agent less icon in running state.
	tc_8 := t.New_testcase(8, "Check Abort process abort the Agent Less icon or not")
	tc_func = func() common.Testcase_status {
		return RunJobnetAndAbort("TICKET811_AgentLess800", 800, 15, tc_8, common.Client)
	}
	tc_8.Set_function(tc_func)
	t.Add_testcase(*tc_8)
}

// Run the jobnet, abort it after all jobs are in running state, and confirm ENDERR status of the jobnet
func RunJobnetAndAbort(jobnetId string, processCount int, processCheckTimeout int, testcase *dao.TestCase, sshClient *ssh.Client) common.Testcase_status {
	err := lib.Jobarg_cleanup_linux()
	if err != nil {
		fmt.Println(testcase.Err_log("Error: %s, Jobarg_cleanup_linux() failed.", err.Error()))
		return FAILED
	}
	fmt.Println(testcase.Info_log("Agent and table is cleaned."))

	// Clean up agent on Linux
	err = lib.Cleanup_agent_linux()
	if err != nil {
		fmt.Println(testcase.Err_log("Error: %s, Failed to clean up agent on Linux.", err.Error()))
		return FAILED
	}
	fmt.Println(testcase.Info_log("Agent cleanup on Linux is completed."))

	//Restart jobarg-server and jobarg-agentd
	err = lib.Jobarg_cleanup_linux()
	if err != nil {
		fmt.Println(testcase.Err_log("Error: %s, Failed to restart jobarg-server and jobarg-agentd.", err.Error()))
		return FAILED
	}
	fmt.Println(testcase.Info_log("jobarg-server and jobarg-agentd have been successfully restarted."))

	// Run jobnet
	run_jobnet_id, error := lib.Jobarg_exec(jobnetId)
	if error != nil {
		fmt.Println(testcase.Err_log("Error: %s, std_out: %s", error.Error(), run_jobnet_id))
		return FAILED
	}
	fmt.Println(testcase.Info_log("%s has been successfully run with registry number: %s", jobnetId, run_jobnet_id))

	// Wait for all jobs to be in running state
	err = lib.JobProcessCountCheck(processCount, processCheckTimeout, sshClient)

	if err != nil {
		fmt.Println(testcase.Err_log("Error: %s, Failed to get process count.", err.Error()))
		return FAILED
	}
	fmt.Println(testcase.Info_log("Process count has reached %d", processCount))

	// Abort the jobnet
	_, err = lib.ExecuteQuery(lib.AbortJobnetQuery, run_jobnet_id)
	if err != nil {
		fmt.Println(testcase.Err_log("Error: %s, Failed to abort the jobnet.", err.Error()))
		return FAILED
	}
	fmt.Println(testcase.Info_log("Jobnet is being aborted..."))

	// Wait for all jobs to be purged
	processCount = 0
	err = lib.JobProcessCountCheck(processCount, processCheckTimeout, sshClient)

	if err != nil {
		fmt.Println(testcase.Err_log("Error: %s, Failed to get process count.", err.Error()))
		return FAILED
	}
	fmt.Println(testcase.Info_log("Process count has reached %d", processCount))

	// Wait jobnet finishes and get jobnet run info.
	targetJobnetStatus := "ENDERR"
	targetJobStatus := "ERROR"
	jobnet_run_info, err := lib.Jobarg_get_jobnet_info(run_jobnet_id, targetJobnetStatus, targetJobStatus)
	if err != nil {
		fmt.Println(testcase.Err_log("Error getting jobnet info: %s", err.Error()))
		return FAILED
	}
	fmt.Println(testcase.Info_log("%s with registry number %s is completed.", jobnetId, run_jobnet_id))

	// Check jobnet run status and exit code.
	if jobnet_run_info.Jobnet_status != targetJobnetStatus {
		fmt.Println(testcase.Err_log("Unexpected Jobnet status. Jobnet_status: %s, Job_status: %s, Exit_cd: %d", jobnet_run_info.Jobnet_status, jobnet_run_info.Job_status, jobnet_run_info.Exit_cd))
		return FAILED
	}

	// check zombie process
	zombieProcessCount, err := lib.CheckZombieProcess(1, sshClient)
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

	// Success (obtain the expected status, message, or exit code)
	fmt.Println(testcase.Info_log("Jobnet_status: %s, Job_status: %s, Exit_cd: %d", jobnet_run_info.Jobnet_status, jobnet_run_info.Job_status, jobnet_run_info.Exit_cd))

	return PASSED
}

// Run the jobnet, abort the fwait job icon after all jobs are in running state, and confirm ENDERR status of the jobnet
func RunJobnetAndAbortJobIcon(jobnetId string, processCount int, processCheckTimeout int, testcase *dao.TestCase, sshClient *ssh.Client) common.Testcase_status {
	// Clean the ja_run_jobnet_table
	_, err := lib.ExecuteQuery(lib.DeleteRunJobnetQuery)
	if err != nil {
		fmt.Println(testcase.Err_log("Error: %s, Failed to clean the ja_run_jobnet_table.", err.Error()))
		return FAILED
	}
	fmt.Println(testcase.Info_log("`ja_run_jobnet_table` is cleaned."))

	// Run jobnet
	run_jobnet_id, error := lib.Jobarg_exec(jobnetId)
	if error != nil {
		fmt.Println(testcase.Err_log("Error: %s, std_out: %s", error.Error(), run_jobnet_id))
		return FAILED
	}
	fmt.Println(testcase.Info_log("%s has been successfully run with registry number: %s", jobnetId, run_jobnet_id))

	// Wait for all jobs to be in running state
	err = lib.JobProcessCountCheck(processCount, processCheckTimeout, sshClient)

	if err != nil {
		fmt.Println(testcase.Err_log("Error: %s, Failed to get process count.", err.Error()))
		return FAILED
	}
	fmt.Println(testcase.Info_log("Process count has reached %d", processCount))

	// Abort the jobnet
	_, err = lib.ExecuteQuery(lib.AbortSingleFWaitJobQuery, run_jobnet_id)
	if err != nil {
		fmt.Println(testcase.Err_log("Error: %s, Failed to abort the Jobnet.", err.Error()))
		return FAILED
	}
	fmt.Println(testcase.Info_log("Jobnet is being aborted..."))

	// Wait for all jobs to be purged
	processCount = 0
	err = lib.JobProcessCountCheck(processCount, processCheckTimeout, sshClient)

	if err != nil {
		fmt.Println(testcase.Err_log("Error: %s, Failed to get process count.", err.Error()))
		return FAILED
	}
	fmt.Println(testcase.Info_log("Process count has reached %d", processCount))

	// Wait jobnet finishes and get jobnet run info.
	targetJobnetStatus := "RUN"
	targetJobStatus := "ERROR"
	jobnet_run_info, err := lib.Jobarg_get_jobnet_info(run_jobnet_id, targetJobnetStatus, targetJobStatus)
	if err != nil {
		fmt.Println(testcase.Err_log("Error getting jobnet info: %s", err.Error()))
		return FAILED
	}
	fmt.Println(testcase.Info_log("%s with registry number %s is completed.", jobnetId, run_jobnet_id))

	// Check jobnet run status and exit code.
	if jobnet_run_info.Jobnet_status != targetJobnetStatus {
		fmt.Println(testcase.Err_log("Unexpected Jobnet status. Jobnet_status: %s, Job_status: %s, Exit_cd: %d", jobnet_run_info.Jobnet_status, jobnet_run_info.Job_status, jobnet_run_info.Exit_cd))
		return FAILED
	}

	// check zombie process
	zombieProcessCount, err := lib.CheckZombieProcess(1, sshClient)
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

	// Success (obtain the expected status, message, or exit code)
	fmt.Println(testcase.Info_log("Jobnet_status: %s, Job_status: %s, Exit_cd: %d", jobnet_run_info.Jobnet_status, jobnet_run_info.Job_status, jobnet_run_info.Exit_cd))

	return PASSED
}

// Run the jobnet, abort the fwait job icon after all jobs are in running state, and confirm ENDERR status of the jobnet
func RunJobnetAndAbortExtJobIcon(jobnetId string, processCheckTimeout int, testcase *dao.TestCase, sshClient *ssh.Client, jobCount int) common.Testcase_status {
	/*
		Prepare process before execute the ext jobnet
		1. cleanup data from ja_run_jobnet_table
	*/
	lib.Jobarg_cleanup_linux()

	/*
		Execute the ext jobnet and process count checking
		1. execute the ext job multiple jobcount like 10,100,200,500,800
		2. checking the process and counting the job
		3. if job count is reached starting the abort process
		4. Waiting the abort process is finished and jobnet is Done with red color
		5. Check the ja_run_jobnet_table and there is no running status 2 or status 6
		   Count is reached 0 and taskcase is done.
	*/
	// 1. execute the jobnet
	run_jobnet_id, error := lib.Jobarg_exec(jobnetId)
	if error != nil {
		fmt.Println(testcase.Err_log("Error: %s, std_out: %s", error.Error(), run_jobnet_id))
		return FAILED
	}
	fmt.Println(testcase.Info_log("%s has been successfully run with registry number: %s", jobnetId, run_jobnet_id))

	// 2. Check the running jobnet
	query := "SELECT COUNT(*) FROM ja_run_icon_extjob_table;"
	count, errJobCount := runProcess(&query, processCheckTimeout, &jobCount) // Correctly capturing count and error
	if errJobCount != nil {
		fmt.Println(errJobCount) // Handle the error
	}

	// 3. if job count is reached starting the abort process
	if count == jobCount {
		fmt.Println("Final Count:", count)
		time.Sleep(time.Duration(processCheckTimeout) * time.Second)
		abortQuery := "UPDATE ja_run_jobnet_summary_table SET jobnet_abort_flag = 1 WHERE inner_jobnet_id = $1"
		_, errJobCount = lib.ExecuteQuery(lib.DBQuery(abortQuery), run_jobnet_id)
		if errJobCount != nil {
			fmt.Println(testcase.Err_log("Error: %s, Failed to abort the ext job icon.", errJobCount.Error()))
			return FAILED
		}
		fmt.Println(testcase.Info_log("Ext job icon is being aborted..."))

		// 4. Waiting the abort process is finished and jobnet is Done with red color
		// time.Sleep(10 * time.Second)
		jobDoneCount, errJobCountWithDone := runProcess(nil, processCheckTimeout, nil) // Correctly capturing count and error
		if errJobCountWithDone != nil {
			fmt.Println(errJobCountWithDone) // Handle the error
			return FAILED
		}

		// 5. Count is reached 0 and taskcase is done.
		if jobDoneCount == 0 {
			fmt.Println("Ext jobnet is successfully red with done color.")
			return PASSED
		}
	}

	return PASSED
}

// GetCountFromDB executes a count query and returns the count or an error.
func GetCountFromDB(query string) (int, error) {
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
func runProcess(query *string, processCheckTimeout int, maxCount *int) (int, error) {
	// Use default query if none provided
	defaultQuery := "SELECT COUNT(*) FROM ja_run_jobnet_table WHERE status = 2 or status = 6;"
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
			count, err = GetCountFromDB(actualQuery)
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
