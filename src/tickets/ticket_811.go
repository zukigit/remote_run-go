package tickets

import (
	"fmt"
	"zukigit/remote_run-go/src/common"
	"zukigit/remote_run-go/src/dao"
	"zukigit/remote_run-go/src/lib"
)

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

func (t *Ticket_811) Run() {
	for _, tc := range t.testcases {
		fmt.Println(tc.Info_log("running..."))
		if !tc.Is_function_nil() {
			tc.Set_status(tc.Run_function())
		} else {
			fmt.Println(tc.Err_log("has no function. SKIPPED!"))
			tc.Set_status(FAILED)
		}

		fmt.Println(tc.Info_log("finished!"))
	}
}

// Enter your ticket information here
func (t *Ticket_811) Set_values() {
	t.no = 811 // Enter your ticket id
	t.description = "cannot stop jobnet"
}

// Add your test case here
func (t *Ticket_811) Add_testcases() {

	// TESTCASE 69 (Force stop FWait job icon)
	tc_69 := t.New_testcase(69, "Check Abort process abort the fwait icon (with waiting for file creation option ON) or not")
	tc_func := func() common.Testcase_status {
		return RunJobnetAndAbortFwaitJobIcon("10.1.9.199", "Jobnet1", 2, tc_69)
	}
	tc_69.Set_function(tc_func)
	t.Add_testcase(*tc_69)

	// TESTCASE 70 (Force stop FWait job icon on other agent)
	tc_70 := t.New_testcase(70, "Check Abort process abort the fwait icon (with waiting for file creation option ON) or not")
	tc_func = func() common.Testcase_status {
		return RunJobnetAndAbort("10.1.9.198", "FWAIT_JOBNET", 2, tc_70)
	}
	tc_70.Set_function(tc_func)
	t.Add_testcase(*tc_70)

	// TESTCASE 71 (Force stop FWait jobnet with icon count of 1)
	tc_71 := t.New_testcase(71, "Check Abort process abort the fwait icon (with waiting for file creation option ON) or not")
	tc_func = func() common.Testcase_status {
		return RunJobnetAndAbort("10.1.9.199", "Jobnet1", 2, tc_71)
	}
	tc_71.Set_function(tc_func)
	t.Add_testcase(*tc_71)

	// TESTCASE 72 (Force stop FWait jobnet with icon count of 100)
	tc_72 := t.New_testcase(72, "Check Abort process abort the fwait icon (with waiting for file creation option ON) or not")
	tc_func = func() common.Testcase_status {
		return RunJobnetAndAbort("10.1.9.199", "Jobnet100", 200, tc_72)
	}
	tc_72.Set_function(tc_func)
	t.Add_testcase(*tc_72)

	// TESTCASE 73 (Force stop FWait jobnet with icon count of 800)
	tc_73 := t.New_testcase(73, "Check Abort process abort the fwait icon (with waiting for file creation option ON) or not")
	tc_func = func() common.Testcase_status {
		return RunJobnetAndAbort("10.1.9.199", "Jobnet800", 1600, tc_73)
	}
	tc_73.Set_function(tc_func)
	// t.Add_testcase(*tc_73)

}

// Run the jobnet, abort it after all jobs are in running state, and confirm ENDERR status of the jobnet
func RunJobnetAndAbort(agentIP string, jobnetId string, processCount int, testcase *dao.TestCase) common.Testcase_status {
	// Get DB Connection as testcases inside this ticket require DB operation
	db, err := lib.ConnectDB(lib.Postgres, agentIP, "5432", "zabbix", "zabbix", "zabbix")

	if err != nil {
		fmt.Println(testcase.Err_log("Error connecting to the database."))
	}

	// Clean the ja_run_jobnet_table
	_, err = db.UpdateData(lib.DeleteRunJobnetQuery)
	if err != nil {
		testcase.Err_log("Error: %s, Failed to clean the ja_run_jobnet_table.", err.Error())
		return FAILED
	}
	fmt.Println(testcase.Info_log("`ja_run_jobnet_table` is cleaned."))

	// Run jobnet
	run_jobnet_id, error := lib.Jobarg_exec(jobnetId)
	if error != nil {
		testcase.Err_log("Error: %s, std_out: %s", error.Error(), run_jobnet_id)
		return FAILED
	}
	fmt.Println(testcase.Info_log("%s has been successfully run with registry number: %s", jobnetId, run_jobnet_id))

	// Wait for all jobs to be in running state
	count, err := lib.Job_process_count_check(processCount)

	if err != nil {
		testcase.Err_log("Error: %s, Failed to get process count.", err.Error())
		return FAILED
	}
	fmt.Println(testcase.Info_log("Process count has reached %d", count))

	// Abort the jobnet
	_, err = db.UpdateData(lib.AbortJobnetQuery, run_jobnet_id)
	if err != nil {
		testcase.Err_log("Error: %s, Failed to abort the jobnet.", err.Error())
		return FAILED
	}
	fmt.Println(testcase.Info_log("Jobnet is being aborted..."))

	// Wait for all jobs to be purged
	processCount = 0
	count, err = lib.Job_process_count_check(processCount)

	if err != nil {
		testcase.Err_log("Error: %s, Failed to get process count.", err.Error())
		return FAILED
	}
	fmt.Println(testcase.Info_log("Process count has reached %d", count))

	// Wait jobnet finishes and get jobnet run info.
	targetJobnetStatus := "ENDERR"
	targetJobStatus := "ERROR"
	jobnet_run_info, err := lib.Jobarg_get_jobnet_info(run_jobnet_id, targetJobnetStatus, targetJobStatus)
	if err != nil {
		testcase.Err_log("Error: %s", err.Error())
		return FAILED
	}
	fmt.Println(testcase.Info_log("%s with registry number %s is completed.", jobnetId, run_jobnet_id))

	// Check jobnet run status and exit code.
	if jobnet_run_info.Jobnet_status != targetJobnetStatus {
		testcase.Err_log("Unexpected Jobnet status. Jobnet_status: %s, Job_status: %s, Exit_cd: %d", jobnet_run_info.Jobnet_status, jobnet_run_info.Job_status, jobnet_run_info.Exit_cd)
		return FAILED
	}

	// check zombie process
	zombieProcessCount, err := lib.CheckZombieProcess()
	if err != nil {
		testcase.Err_log("Error: %s", err.Error())
		return FAILED
	}

	if zombieProcessCount != 0 {
		testcase.Err_log("There are zombie processes: %d", zombieProcessCount)
		// return FAILED
	} else {
		testcase.Info_log("There is no zombie process.")
	}

	// Success (obtain the expected status, message, or exit code)
	testcase.Info_log("Jobnet_status: %s, Job_status: %s, Exit_cd: %d", jobnet_run_info.Jobnet_status, jobnet_run_info.Job_status, jobnet_run_info.Exit_cd)

	return PASSED
}

// Run the jobnet, abort the fwait job icon after all jobs are in running state, and confirm ENDERR status of the jobnet
func RunJobnetAndAbortFwaitJobIcon(agentIP string, jobnetId string, processCount int, testcase *dao.TestCase) common.Testcase_status {
	// Get DB Connection as testcases inside this ticket require DB operation
	db, err := lib.ConnectDB(lib.Postgres, agentIP, "5432", "zabbix", "zabbix", "zabbix")

	if err != nil {
		fmt.Println(testcase.Err_log("Error connecting to the database."))
	}

	// Clean the ja_run_jobnet_table
	_, err = db.UpdateData(lib.DeleteRunJobnetQuery)
	if err != nil {
		testcase.Err_log("Error: %s, Failed to clean the ja_run_jobnet_table.", err.Error())
		return FAILED
	}
	fmt.Println(testcase.Info_log("`ja_run_jobnet_table` is cleaned."))

	// Run jobnet
	run_jobnet_id, error := lib.Jobarg_exec(jobnetId)
	if error != nil {
		testcase.Err_log("Error: %s, std_out: %s", error.Error(), run_jobnet_id)
		return FAILED
	}
	fmt.Println(testcase.Info_log("%s has been successfully run with registry number: %s", jobnetId, run_jobnet_id))

	// Wait for all jobs to be in running state
	count, err := lib.Job_process_count_check(processCount)

	if err != nil {
		testcase.Err_log("Error: %s, Failed to get process count.", err.Error())
		return FAILED
	}
	fmt.Println(testcase.Info_log("Process count has reached %d", count))

	// Abort the jobnet
	_, err = db.UpdateData(lib.AbortSingleFWaitJobQuery, run_jobnet_id)
	if err != nil {
		testcase.Err_log("Error: %s, Failed to abort the fwait job icon.", err.Error())
		return FAILED
	}
	fmt.Println(testcase.Info_log("Fwait job icon is being aborted..."))

	// Wait for all jobs to be purged
	processCount = 0
	count, err = lib.Job_process_count_check(processCount)

	if err != nil {
		testcase.Err_log("Error: %s, Failed to get process count.", err.Error())
		return FAILED
	}
	fmt.Println(testcase.Info_log("Process count has reached %d", count))

	// Wait jobnet finishes and get jobnet run info.
	targetJobnetStatus := "RUN"
	targetJobStatus := "ERROR"
	jobnet_run_info, err := lib.Jobarg_get_jobnet_info(run_jobnet_id, targetJobnetStatus, targetJobStatus)
	if err != nil {
		testcase.Err_log("Error: %s", err.Error())
		return FAILED
	}
	fmt.Println(testcase.Info_log("%s with registry number %s is completed.", jobnetId, run_jobnet_id))

	// Check jobnet run status and exit code.
	if jobnet_run_info.Jobnet_status != targetJobnetStatus {
		testcase.Err_log("Unexpected Jobnet status. Jobnet_status: %s, Job_status: %s, Exit_cd: %d", jobnet_run_info.Jobnet_status, jobnet_run_info.Job_status, jobnet_run_info.Exit_cd)
		return FAILED
	}

	// check zombie process
	zombieProcessCount, err := lib.CheckZombieProcess()
	if err != nil {
		testcase.Err_log("Error: %s", err.Error())
		return FAILED
	}

	if zombieProcessCount != 0 {
		testcase.Err_log("There are zombie processes: %d", zombieProcessCount)
		// return FAILED
	} else {
		testcase.Info_log("There is no zombie process.")
	}

	// Success (obtain the expected status, message, or exit code)
	testcase.Info_log("Jobnet_status: %s, Job_status: %s, Exit_cd: %d", jobnet_run_info.Jobnet_status, jobnet_run_info.Job_status, jobnet_run_info.Exit_cd)

	return PASSED
}