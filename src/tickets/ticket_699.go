package tickets

import (
	"fmt"

	"github.com/zukigit/remote_run-go/src/common"
	"github.com/zukigit/remote_run-go/src/dao"
	"github.com/zukigit/remote_run-go/src/lib"
	"golang.org/x/crypto/ssh"
)

type Ticket_699 struct {
	no          uint
	description string
	testcases   []dao.TestCase
}

func (t *Ticket_699) New_testcase(testcase_id uint, testcase_description string) *dao.TestCase {
	return dao.New_testcase(testcase_id, testcase_description)
}

func (t *Ticket_699) Get_no() uint {
	return t.no
}

func (t *Ticket_699) Get_dsctn() string {
	return t.description
}

func (t *Ticket_699) Add_testcase(tc dao.TestCase) {
	t.testcases = append(t.testcases, tc)
}

func (t *Ticket_699) Get_testcases() []dao.TestCase {
	return t.testcases
}

// Enter your ticket information here
func (t *Ticket_699) Set_values() {
	t.no = 699 // Enter your ticket id
	t.description = "Reconnect and retry case for connection agent(Job Icon, File Transfer, Reboot)"
}

func (t *Ticket_699) Add_testcases() {

	tc_1 := t.New_testcase(1, "Check Server log output")
	tc_func := func() common.Testcase_status {
		if err := lib.Jobarg_enable_jobnet("Icon_1", "jobicon_linux"); err != nil {
			tc_1.Err_log("Failed to enable jobnet, Error: %s", err)
			return FAILED
		}
		return CheckLog("Icon_1", tc_1, common.Client)
	}
	tc_1.Set_function(tc_func)
	t.Add_testcase(*tc_1)

	tc_2 := t.New_testcase(2, "Check Server log output")
	tc_func = func() common.Testcase_status {
		if err := lib.Jobarg_enable_jobnet("Icon_1", "jobicon_linux_skip"); err != nil {
			tc_2.Err_log("Failed to enable jobnet, Error: %s", err)
			return FAILED
		}
		return checkjJobnetPurgeSuccefullyOrNot("Icon_1", tc_2, "END", "NORMAL", "END", "NORMAL", 1)
	}
	tc_2.Set_function(tc_func)
	t.Add_testcase(*tc_2)

	// tc_3 := t.New_testcase(1, "Check Server log output")
	// tc_func = func() common.Testcase_status {
	// 	if err := lib.Jobarg_enable_jobnet("Icon_1", "jobicon_linux"); err != nil {
	// 		tc_3.Err_log("Failed to enable jobnet, Error: %s", err)
	// 		return FAILED
	// 	}
	// 	return checkjJobnetPurgeSuccefullyOrNot("Icon_1", tc_2, "", "", "", "", 1)
	// }
	// tc_3.Set_function(tc_func)
	// t.Add_testcase(*tc_3)
}

func CheckLog(jobnetId string, testcase *dao.TestCase, sshClient *ssh.Client) common.Testcase_status {

	err := lib.Clean_jaz_server_log()
	if err != nil {
		fmt.Println(testcase.Err_log("Error: %s, Failed to cleanup the job arranger server log.", err))
		return FAILED
	}
	fmt.Println(testcase.Info_log("Cleanup the server log success."))

	err = lib.Restart_jaz_server()
	if err != nil {
		fmt.Println(testcase.Err_log("Error: %s, Failed to restart the jobarranger server.", err))
		return FAILED
	}
	fmt.Println(testcase.Info_log("Server restart success."))

	searchString := `'\[JASERVER000004\] server #13 started \[purge old jobnet #1\]'`
	logs, err := lib.Jobarg_server_check_log(searchString)
	if err != nil {
		fmt.Println("Error:", err)
		return FAILED
	}
	fmt.Println(testcase.Info_log("Log output result for purge old jobnet:"))
	for _, log := range logs {
		fmt.Println(log)
	}

	searchString = "[JASERVER000004] server #10 started [jobnet boot #1]"
	logs, err = lib.Jobarg_server_check_log(searchString)
	if err != nil {
		fmt.Println("Error:", err)
		return FAILED
	}

	fmt.Println(testcase.Info_log("Log output result for jobnet boot:"))
	for _, log := range logs {
		fmt.Println(log)
	}
	return PASSED
}

func checkjJobnetPurgeSuccefullyOrNot(jobnetId string, testcase *dao.TestCase, firstTargetJobnetStatus string, firstTargetJobStatus string, secondTargetJobnetStatus string, secondTargetJobStatus string, processCheckTimeout int) common.Testcase_status {

	var firstJobnetId, secondJobnetId string

	// Clean up the agent
	err := lib.Jobarg_cleanup_linux()
	if err != nil {
		fmt.Println(testcase.Err_log("Error: %s, Failed to clean up the linux agent.", err.Error()))
		return FAILED
	}
	fmt.Println(testcase.Info_log("Clean up agent service success."))

	envs, _ := lib.Get_str_str_map("JA_HOSTNAME", "oss.linux", "JA_CMD", "sleep 10")
	run_jobnet_id, err := lib.Jobarg_exec_E(jobnetId, envs)
	if err != nil {
		fmt.Println(testcase.Err_log("Error: %s, std_out: %s", err.Error(), run_jobnet_id))
		return FAILED
	}
	fmt.Println(testcase.Info_log("%s has been successfully run with registry number: %s", jobnetId, run_jobnet_id))

	firstJobnetId = run_jobnet_id

	envs, _ = lib.Get_str_str_map("JA_HOSTNAME", "oss.linux", "JA_CMD", "sleep 10")
	run_jobnet_id, err = lib.Jobarg_exec_E(jobnetId, envs)
	if err != nil {
		fmt.Println(testcase.Err_log("Error: %s, std_out: %s", err.Error(), run_jobnet_id))
		return FAILED
	}
	fmt.Println(testcase.Info_log("%s has been successfully run with registry number: %s", jobnetId, run_jobnet_id))

	secondJobnetId = run_jobnet_id

	// Wait jobnet finishes and get jobnet run info.
	jobnet_run_info, err := lib.Jobarg_get_jobnet_info(firstJobnetId, firstTargetJobnetStatus, firstTargetJobStatus, processCheckTimeout)
	if err != nil {
		fmt.Println(testcase.Err_log("Error getting jobnet info: %s", err.Error()))
		return FAILED
	}
	fmt.Println(testcase.Info_log("%s with registry number %s is completed.", jobnetId, firstJobnetId))

	// Check jobnet run status and exit code.
	if jobnet_run_info.Jobnet_status != firstTargetJobnetStatus {
		fmt.Println(testcase.Err_log("Unexpected Jobnet status. Jobnet_status: %s, Job_status: %s, Exit_cd: %d", jobnet_run_info.Jobnet_status, jobnet_run_info.Job_status, jobnet_run_info.Exit_cd))
		return FAILED
	}
	// Success (obtain the expected status, message, or exit code)
	fmt.Println(testcase.Info_log("Jobnet_status: %s, Job_status: %s, Exit_cd: %d", jobnet_run_info.Jobnet_status, jobnet_run_info.Job_status, jobnet_run_info.Exit_cd))

	// Wait jobnet finishes and get jobnet run info.
	jobnet_run_info, err = lib.Jobarg_get_jobnet_info(secondJobnetId, secondTargetJobnetStatus, secondTargetJobStatus, processCheckTimeout)
	if err != nil {
		fmt.Println(testcase.Err_log("Error getting jobnet info: %s", err.Error()))
		return FAILED
	}
	fmt.Println(testcase.Info_log("%s with registry number %s is completed.", jobnetId, secondJobnetId))

	// Check jobnet run status and exit code.
	if jobnet_run_info.Jobnet_status != secondTargetJobnetStatus {
		fmt.Println(testcase.Err_log("Unexpected Jobnet status. Jobnet_status: %s, Job_status: %s, Exit_cd: %d", jobnet_run_info.Jobnet_status, jobnet_run_info.Job_status, jobnet_run_info.Exit_cd))
		return FAILED
	}
	// Success (obtain the expected status, message, or exit code)
	fmt.Println(testcase.Info_log("Jobnet_status: %s, Job_status: %s, Exit_cd: %d", jobnet_run_info.Jobnet_status, jobnet_run_info.Job_status, jobnet_run_info.Exit_cd))

	return PASSED
}
