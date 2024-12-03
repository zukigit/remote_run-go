package tickets

import (
	"fmt"
	"time"

	"github.com/zukigit/remote_run-go/src/common"
	"github.com/zukigit/remote_run-go/src/dao"
	"github.com/zukigit/remote_run-go/src/lib"
)

type Ticket_962 struct {
	Ticket_no                                   uint
	Ticket_description                          string
	PASSED_count, FAILED_count, MUSTCHECK_count int
	Testcases                                   []dao.TestCase
}

func (t *Ticket_962) Set_PASSED_count(passed_count int) {
	t.PASSED_count = passed_count
}

func (t *Ticket_962) Set_FAILED_count(failed_count int) {
	t.FAILED_count = failed_count
}

func (t *Ticket_962) Set_MUSTCHECK_count(mustcheck_count int) {
	t.MUSTCHECK_count = mustcheck_count
}

func (t *Ticket_962) New_testcase(testcase_id uint, testcase_description string) *dao.TestCase {
	return dao.New_testcase(testcase_id, testcase_description)
}

func (t *Ticket_962) Get_no() uint {
	return t.Ticket_no
}

func (t *Ticket_962) Get_dsctn() string {
	return t.Ticket_description
}

func (t *Ticket_962) Add_testcase(tc dao.TestCase) {
	t.Testcases = append(t.Testcases, tc)
}

func (t *Ticket_962) Get_testcases() []dao.TestCase {
	return t.Testcases
}

// Enter your ticket information here
func (t *Ticket_962) Set_values() {
	t.Ticket_no = 962 // Enter your ticket id
	t.Ticket_description = "Check process for jobarg_session's processes"
}

// Add your test case here
func (t *Ticket_962) Add_testcases() {
	t.addTestcase85()
	t.addTestcase86()
}

// General jobnet setup function
func (t *Ticket_962) setupJobnet(jobnetID, cmd string, tc *dao.TestCase) (string, error) {
	// Cleanup jobarg
	if err := lib.Jobarg_cleanup_linux(); err != nil {
		return "", fmt.Errorf("failed to cleanup jobarg, error: %s", err)
	}
	fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "Jobarg cleanup successfully."))

	// Enable jobnet
	if err := lib.Jobarg_enable_jobnet(jobnetID, cmd); err != nil {
		return "", fmt.Errorf("failed to enable jobnet, error: %s", err)
	}
	fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "Jobnet enabled successfully."))

	// Execute jobnet
	envs, _ := lib.Get_str_str_map("JA_HOSTNAME", "oss.linux", "JA_CMD", "sleep 60")
	runJobnetID, err := lib.Jobarg_exec_E(jobnetID, envs)
	if err != nil {
		return "", fmt.Errorf("error running jobnet: %s", err)
	}
	fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "%s has been successfully run with registry number: %s\n", jobnetID, runJobnetID))

	return runJobnetID, nil
}

// Test Case 85: Abnormal job execution
func (t *Ticket_962) addTestcase85() {
	tc := t.New_testcase(85, "Abnormal job execution")
	tc_func := func() common.Testcase_status {
		jobnetID := "Icon_1"

		// Set up the jobnet
		runJobnetID, err := t.setupJobnet(jobnetID, "agentless", tc)
		if err != nil {
			lib.Logi(common.LOG_LEVEL_ERR, err.Error())
			return FAILED
		}

		// Check jobnet status
		targetJobnetStatus := "END"
		targetJobStatus := "NORMAL"
		jobnetRunInfo, err := lib.Jobarg_get_jobnet_info(runJobnetID, targetJobnetStatus, targetJobStatus, 5)
		if err != nil {
			lib.Logi(common.LOG_LEVEL_ERR, "Error getting jobnet info: %s", err.Error())
			return FAILED
		}
		fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "%s with registry number %s is completed.\n", jobnetID, runJobnetID))

		if jobnetRunInfo.Jobnet_status != targetJobnetStatus && jobnetRunInfo.Job_status != targetJobStatus {
			lib.Logi(common.LOG_LEVEL_ERR, "Unexpected Jobnet status. Jobnet_status: %s, Job_status: %s, Exit_cd: %d", jobnetRunInfo.Jobnet_status, jobnetRunInfo.Job_status, jobnetRunInfo.Exit_cd)
			return FAILED
		}

		return PASSED
	}
	tc.Set_function(tc_func)
	t.Add_testcase(*tc)
}

// Test Case 86: Abnormal job execution
func (t *Ticket_962) addTestcase86() {
	tc := t.New_testcase(86, "Abnormal job execution")
	tc_func := func() common.Testcase_status {
		jobnetID := "Icon_1"

		// Set up the jobnet
		runJobnetID, err := t.setupJobnet(jobnetID, "agentless", tc)
		if err != nil {
			lib.Logi(common.LOG_LEVEL_ERR, err.Error())
			return FAILED
		}

		// Validate process count
		if err := lib.JobProcessCountCheck(1, 2, common.Client); err != nil {
			lib.Logi(common.LOG_LEVEL_ERR, "Failed to get process count, error: %s", err.Error())
			return FAILED
		}

		// Clear server log and restart Jaz server
		logfilePath := "/var/log/jobarranger/jobarg_server.log"
		if _, err := lib.Ssh_exec("echo > " + logfilePath); err != nil {
			lib.Logi(common.LOG_LEVEL_ERR, "Failed to clear server log: %s", err.Error())
			return FAILED
		}
		fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "Server log file has been cleared."))

		if err := lib.Restart_jaz_server(); err != nil {
			lib.Logi(common.LOG_LEVEL_ERR, "Failed to restart Jaz server, error: %s", err.Error())
			return FAILED
		}
		fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "Jaz server restarted successfully."))

		// Retrieve data from the server log
		searchData := "\\[ERROR]\\ \\[JASESSIONCHCK000001]\\ The process of session_id"
		output, err := getDataFromServerLog(logfilePath, searchData)
		if err != nil {
			lib.Logi(common.LOG_LEVEL_ERR, "Failed to get data from server log: %s", err.Error())
			return FAILED
		}
		if output == "" {
			fmt.Println("Output is empty!")
		} else {
			fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, output))
		}

		// Check jobnet run status
		targetJobnetStatus := "RUN"
		targetJobStatus := "ERROR"
		jobnetRunInfo, err := lib.Jobarg_get_jobnet_info(runJobnetID, targetJobnetStatus, targetJobStatus, 5)
		if err != nil {
			lib.Logi(common.LOG_LEVEL_ERR, "Error getting jobnet info: %s", err.Error())
			return FAILED
		}

		if jobnetRunInfo.Jobnet_status != targetJobnetStatus && jobnetRunInfo.Job_status != targetJobStatus {
			lib.Logi(common.LOG_LEVEL_ERR, "Unexpected Jobnet status. Jobnet_status: %s, Job_status: %s, Exit_cd: %d", jobnetRunInfo.Jobnet_status, jobnetRunInfo.Job_status, jobnetRunInfo.Exit_cd)
			return FAILED
		}

		return PASSED
	}
	tc.Set_function(tc_func)
	t.Add_testcase(*tc)
}

// Retrieves data from the server log based on search criteria
func getDataFromServerLog(logFilePath, searchData string) (string, error) {
	time.Sleep(10 * time.Second)
	cmd := fmt.Sprintf("cat %s | grep '%s'", logFilePath, searchData)
	output, err := lib.Ssh_exec_to_str(cmd)
	if err != nil {
		return "", fmt.Errorf("failed to execute command: %s", err)
	}
	return output, nil
}
