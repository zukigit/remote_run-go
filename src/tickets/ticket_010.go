package tickets

import (
	"fmt"

	"github.com/zukigit/remote_run-go/src/common"
	"github.com/zukigit/remote_run-go/src/lib"
	"golang.org/x/crypto/ssh"
)

type Ticket_010 struct {
	Ticket_no                                   uint
	Ticket_description                          string
	PASSED_count, FAILED_count, MUSTCHECK_count int
	Testcases                                   []common.TestCase
}

func (t *Ticket_010) New_testcase(testcase_id uint, testcase_description string) *common.TestCase {
	return common.New_testcase(testcase_id, testcase_description)
}

func (t *Ticket_010) Get_no() uint {
	return t.Ticket_no
}

func (t *Ticket_010) Get_dsctn() string {
	return t.Ticket_description
}
func (t *Ticket_010) Set_PASSED_count(passed_count int) {
	t.PASSED_count = passed_count
}

func (t *Ticket_010) Set_FAILED_count(failed_count int) {
	t.FAILED_count = failed_count
}

func (t *Ticket_010) Set_MUSTCHECK_count(mustcheck_count int) {
	t.MUSTCHECK_count = mustcheck_count
}
func (t *Ticket_010) Add_testcase(tc common.TestCase) {
	t.Testcases = append(t.Testcases, tc)
}

func (t *Ticket_010) Get_testcases() []common.TestCase {
	return t.Testcases
}

// Enter your ticket information here
func (t *Ticket_010) Set_values() {
	t.Ticket_no = 010 // Enter your ticket id
	t.Ticket_description = "Transaction files delete"
}

func (t *Ticket_010) Add_testcases() {
	// TESTCASE 74 (Transaction files delete)

	tc_74 := t.New_testcase(1, "Transaction file delete wile 800 parallel jobnets are running ")
	tc_func := func() common.Testcase_status {
		if err := lib.Jobarg_enable_jobnet("Icon_1", "jobicon_linux"); err != nil {
			lib.Logi(common.LOG_LEVEL_ERR, "Failed to enable jobnet, Error: %s", err)
			return FAILED
		}
		return RunJobnetAndStopTheAgentAndDeleteTransactionFileAndRestartTheAgent("Icon_800", 800, 15, tc_74, common.Client)
	}
	tc_74.Set_function(tc_func)
	t.Add_testcase(*tc_74)

	tc_78 := t.New_testcase(2, "Agent servive stop while 800  parallel jobnets are running")
	tc_func = func() common.Testcase_status {
		if err := lib.Jobarg_enable_jobnet("Icon_1", "jobicon_linux"); err != nil {
			lib.Logi(common.LOG_LEVEL_ERR, "Failed to enable jobnet, Error: %s", err)
			return FAILED
		}
		return RunJobnetAndForceStopTheAgent("Icon_800", 800, 15, tc_78, common.Client)
	}
	tc_78.Set_function(tc_func)
	t.Add_testcase(*tc_78)
}

// Run 800 Jobnets, all jobs icons are running and then stop the agent and delete all transactions file after deleting restart the agent
func RunJobnetAndStopTheAgentAndDeleteTransactionFileAndRestartTheAgent(jobnetId string, processCount int, processCheckTimeout int, testcase *common.TestCase, sshClient *ssh.Client) common.Testcase_status {

	// Clean up the agent
	err := lib.Jobarg_cleanup_linux()
	if err != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error: %s, Failed to clean up the linux agent.", err.Error()))
		return FAILED
	}
	fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "Clean up agent service success."))

	envs, _ := lib.Get_str_str_map("JA_HOSTNAME", "oss.linux", "JA_CMD", "sleep 1000")

	// Run jobnet
	run_jobnet_id, err := lib.Jobarg_exec_E(jobnetId, envs)
	if err != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error: %s, std_out: %s", err.Error(), run_jobnet_id))
		return FAILED
	}
	fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "%s has been successfully run with registry number: %s", jobnetId, run_jobnet_id))

	// Wait for all jobs to be in running state
	err = lib.JobProcessCountCheck(processCount, processCheckTimeout, sshClient)
	if err != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error: %s, Failed to get process count.", err.Error()))
		return FAILED
	}
	fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "Process count has reached %d", processCount))

	// Stop the agent
	err = lib.Stop_jaz_agent_linux()
	if err != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error : %s, Failed to stop the agent.", err.Error()))
		return FAILED
	}
	fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "Successful stop the agent."))

	// Delete all the transaction file of agent
	err = lib.Cleanup_agent_linux()
	if err != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error : %s, Failed to delete the files.", err.Error()))
		return FAILED
	}
	fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "Successful delete the files from /var/lib/jobarranger/tmp/"))

	// Restart the agent
	err = lib.Restart_jaz_agent_linux()
	if err != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error : %s, Failed to restart the agent.", err.Error()))
		return FAILED
	}
	fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "Successful restart the agent."))

	// Job Process Check in DB
	err = lib.JobProcessDBCountCheck(processCount, processCheckTimeout, run_jobnet_id, lib.CheckJobStatusCountQuery)
	if err != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error : %s, Check DB Count porcess failed.", err.Error()))
		return FAILED
	}
	fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "In Db count has reach the limit %d", processCount))

	// check zombie process
	zombieProcessCount, err := lib.CheckZombieProcess(5, sshClient)
	if err != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error checking zombie process: %s", err.Error()))
		return FAILED
	}

	if zombieProcessCount != 0 {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "There are zombie processes: %d", zombieProcessCount))
		return FAILED
	}
	fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "There is no zombie process."))

	return PASSED
}

func RunJobnetAndForceStopTheAgent(jobnetId string, processCount int, processCheckTimeout int, testcase *common.TestCase, sshClient *ssh.Client) common.Testcase_status {

	// Clean up the agent
	err := lib.Jobarg_cleanup_linux()
	if err != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error: %s, Failed to clean up the linux agent.", err.Error()))
		return FAILED
	}
	fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "Clean up agent service success."))

	envs, _ := lib.Get_str_str_map("JA_HOSTNAME", "oss.linux", "JA_CMD", "sleep 1000")

	// Run jobnet
	run_jobnet_id, err := lib.Jobarg_exec_E(jobnetId, envs)
	if err != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error: %s, std_out: %s", err.Error(), run_jobnet_id))
		return FAILED
	}
	fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "%s has been successfully run with registry number: %s", jobnetId, run_jobnet_id))

	// Wait for all jobs to be in running state
	err = lib.JobProcessCountCheck(processCount, processCheckTimeout, sshClient)
	if err != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error: %s, Failed to get process count.", err.Error()))
		return FAILED
	}
	fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "Process count has reached %d", processCount))

	// Stop the agent
	err = lib.Stop_jaz_agent_linux()
	if err != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error : %s, Failed to stop the agent.", err.Error()))
		return FAILED
	}
	fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "Successful stop the agent."))

	// Job Process Check in DB
	err = lib.JobProcessDBCountCheck(processCount, processCheckTimeout, run_jobnet_id, lib.CheckJobStatusCountQuery)
	if err != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error : %s, Check DB Count porcess failed.", err.Error()))
		return FAILED
	}
	fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "In Db count has reach the limit %d", processCount))

	// check zombie process
	zombieProcessCount, err := lib.CheckZombieProcess(1, sshClient)
	if err != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error checking zombie process: %s", err.Error()))
		return FAILED
	}

	if zombieProcessCount != 0 {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "There are zombie processes: %d", zombieProcessCount))
		return FAILED
	} else {
		fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "There is no zombie process."))
	}

	return PASSED
}
