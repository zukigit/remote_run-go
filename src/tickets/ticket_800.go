package tickets

import (
	"fmt"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/ssh"

	"github.com/zukigit/remote_run-go/src/common"
	"github.com/zukigit/remote_run-go/src/dao"
	"github.com/zukigit/remote_run-go/src/lib"
)

type Ticket_800 struct {
	Ticket_no                                   uint
	Ticket_description                          string
	PASSED_count, FAILED_count, MUSTCHECK_count int
	Testcases                                   []dao.TestCase
}

func (t *Ticket_800) New_testcase(testcase_id uint, testcase_description string) *dao.TestCase {
	return dao.New_testcase(testcase_id, testcase_description)
}

func (t *Ticket_800) Get_no() uint {
	return t.Ticket_no
}

func (t *Ticket_800) Set_PASSED_count(passed_count int) {
	t.PASSED_count = passed_count
}

func (t *Ticket_800) Set_FAILED_count(failed_count int) {
	t.FAILED_count = failed_count
}

func (t *Ticket_800) Set_MUSTCHECK_count(mustcheck_count int) {
	t.MUSTCHECK_count = mustcheck_count
}

func (t *Ticket_800) Get_dsctn() string {
	return t.Ticket_description
}

func (t *Ticket_800) Add_testcase(tc dao.TestCase) {
	t.Testcases = append(t.Testcases, tc)
}

func (t *Ticket_800) Get_testcases() []dao.TestCase {
	return t.Testcases
}

// Enter your ticket information here
func (t *Ticket_800) Set_values() {
	t.Ticket_no = 800 // Enter your ticket id
	t.Ticket_description = "Run Jobnet, Issue 777, 783 and 800."
}

// Add your test case here
func (t *Ticket_800) Add_testcases() {
	tc_1 := t.New_testcase(1, "Stop and Start Jobar Agent Service when runing a jobnet.")
	tc_func := func() common.Testcase_status {
		if err := lib.Jobarg_enable_jobnet("Icon_1", "jobicon_linux"); err != nil {
			lib.Logi(common.LOG_LEVEL_ERR, "Failed to enable jobnet, Error: %s", err)
			return FAILED
		}
		return RunJobnet("Icon_1", 1, 1, tc_1, common.Client)
	}
	tc_1.Set_function(tc_func)
	t.Add_testcase(*tc_1)
}

// RunJobnet cleans the environment and runs the jobnet while simulating jobab-termination and restart.
func RunJobnet(jobnetId string, processCount int, processCheckTimeout int, testcase *dao.TestCase, client *ssh.Client) common.Testcase_status {
	// Cleanup jobarg server and agent
	err := lib.Jobarg_cleanup_linux()
	if err != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error cleaning up jobarg server and agent: %s", err.Error()))
		return FAILED
	}
	fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "jobarg-server and jobarg-agentd have been cleaned and restarted."))

	envs, _ := lib.Get_str_str_map("JA_HOSTNAME", "oss.linux", "JA_CMD", "sleep 50")

	// Run jobnet
	run_jobnet_id, error := lib.Jobarg_exec_E(jobnetId, envs)
	if error != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error: %s, std_out: %s", error.Error(), run_jobnet_id))
		return FAILED
	}
	fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "%s has been successfully run with registry number: %s", jobnetId, run_jobnet_id))

	// Wait for all jobs to be in running state
	err = lib.JobProcessCountCheck(processCount, processCheckTimeout, client)
	if err != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error getting process count: %s", err.Error()))
		return FAILED
	}
	fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "Process count has reached %d", processCount))

	// Simulate stopping jobarg-agentd after jobnet starts running
	fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "Stopping jobarg-agentd..."))
	err = lib.Stop_jaz_agent_linux() // Stop the agent
	if err != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error stopping jobarg-agentd: %s", err.Error()))
		return FAILED
	}
	// Verify "no connection to database" error is not in the log
	logFilePath := "/var/log/jobarranger/jobarg_server.log"
	errorMsg := "no connection to database"
	timeout := 10 * time.Second
	interval := 1 * time.Second

	// Use WaitForPatternInLogFile to check for the error message
	fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "Checking for 'no connection to database' error in log..."))
	_, err = lib.WaitForPatternInLogFile(client, logFilePath, errorMsg, timeout, interval)
	if err == nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Unexpected error message found in log: %s", errorMsg))
		return FAILED
	} else if !strings.Contains(err.Error(), "timeout reached") {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error while checking log for database connection errors: %s", err.Error()))
		return FAILED
	}
	// Wait to simulate an interruption
	time.Sleep(10 * time.Second)

	// Restart jobarg-agentd
	fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "Restarting jobarg-agentd..."))
	err = lib.Restart_jaz_agent_linux()
	if err != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error restarting jobarg-agentd: %s", err.Error()))
		return FAILED
	}

	// Wait for jobnet to finish and retrieve jobnet run info
	targetJobnetStatus := "RUN"
	targetJobStatus := "ERROR"
	jobnet_run_info, err := lib.Jobarg_get_jobnet_info(run_jobnet_id, targetJobnetStatus, targetJobStatus, 5)
	if err != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error getting jobnet info: %s", err.Error()))
		return FAILED
	}
	fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "Jobnet %s with registry number %s is completed.", jobnetId, run_jobnet_id))

	// Check jobnet run status and exit code
	// fmt.Println(jobnet_run_info.Jobnet_status)
	if jobnet_run_info.Jobnet_status != targetJobnetStatus || jobnet_run_info.Job_status != targetJobStatus {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Unexpected Jobnet or Job status. Jobnet_status: %s, Job_status: %s, Exit_cd: %d",
			jobnet_run_info.Jobnet_status, jobnet_run_info.Job_status, jobnet_run_info.Exit_cd))
		return FAILED
	}

	// Check for zombie processes after job completion
	zombieProcessCount, err := lib.CheckZombieProcess(1, client)
	if err != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error checking zombie processes: %s", err.Error()))
		return FAILED
	}
	if zombieProcessCount != 0 {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "There are zombie processes: %d", zombieProcessCount))
	} else {
		fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "No zombie process detected."))
	}

	return PASSED
}
