package tickets

import (
	"fmt"

	"github.com/zukigit/remote_run-go/src/common"
	"github.com/zukigit/remote_run-go/src/dao"
	"github.com/zukigit/remote_run-go/src/lib"
	"golang.org/x/crypto/ssh"
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

// Enter your ticket information here
func (t *Ticket_811) Set_values() {
	t.no = 811 // Enter your ticket id
	t.description = "cannot stop jobnet"
}

// Add your test case here
func (t *Ticket_811) Add_testcases() {

	// TESTCASE 46 (Force stop Job job icon)
	tc_46 := t.New_testcase(46, "Check Abort process abort the Job icon (with waiting for file creation option ON) or not")
	tc_func := func() common.Testcase_status {
		return RunJobnetAndAbort("TICKET811_Jobnet", 1, 2, tc_46, common.Client)
	}
	tc_46.Set_function(tc_func)
	t.Add_testcase(*tc_46)

	//TESTCASE 70 (Force stop job icon on other agent)
	// tc_70 := t.New_testcase(70, "Check Abort process abort the job icon or not")
	// tc_func = func() common.Testcase_status {
	// 	agentSSHClient := lib.GetSSHClient("10.1.9.80", 22, "root", "123@dirace") // remote agent
	// 	return RunJobnetAndAbort("TICKET811_FileWaitJobOtherAgent1", 2, 5, tc_70, agentSSHClient)
	// }
	// tc_70.Set_function(tc_func)
	// t.Add_testcase(*tc_70)

	// // TESTCASE 47 (Force stop jobnet with icon count of 1)
	tc_47 := t.New_testcase(47, "Check Abort process abort the job icon or not")
	tc_func = func() common.Testcase_status {
		return RunJobnetAndAbort("TICKET811_Jobnet_10", 10, 5, tc_47, common.Client)
	}
	tc_47.Set_function(tc_func)
	t.Add_testcase(*tc_47)

	// // TESTCASE 48 (Force stop jobnet with icon count of 100)
	tc_48 := t.New_testcase(48, "Check Abort process abort the job icon or not")
	tc_func = func() common.Testcase_status {
		return RunJobnetAndAbort("TICKET811_Jobnet_100", 120, 15, tc_48, common.Client)
	}
	tc_48.Set_function(tc_func)
	t.Add_testcase(*tc_48)

	// TESTCASE 49 (Force stop jobnet with icon count of 800)
	// tc_49 := t.New_testcase(49, "Check Abort process abort the job icon or not")
	// tc_func = func() common.Testcase_status {
	// 	return RunJobnetAndAbort("TICKET811_Jobnet_400", 400, 20, tc_49, common.Client)
	// }
	// tc_49.Set_function(tc_func)
	// t.Add_testcase(*tc_49)

}

// Run the jobnet, abort it after all jobs are in running state, and confirm ENDERR status of the jobnet
func RunJobnetAndAbort(jobnetId string, processCount int, processCheckTimeout int, testcase *dao.TestCase, sshClient *ssh.Client) common.Testcase_status {
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
