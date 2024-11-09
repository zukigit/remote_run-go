package tickets

import (
	"fmt"

	"github.com/zukigit/remote_run-go/src/common"
	"github.com/zukigit/remote_run-go/src/dao"
	"github.com/zukigit/remote_run-go/src/lib"
	"golang.org/x/crypto/ssh"
)

type Ticket_698 struct {
	no          uint
	description string
	testcases   []dao.TestCase
}

func (t *Ticket_698) New_testcase(testcase_id uint, testcase_description string) *dao.TestCase {
	return dao.New_testcase(testcase_id, testcase_description)
}

func (t *Ticket_698) Get_no() uint {
	return t.no
}

func (t *Ticket_698) Get_dsctn() string {
	return t.description
}

func (t *Ticket_698) Add_testcase(tc dao.TestCase) {
	t.testcases = append(t.testcases, tc)
}

func (t *Ticket_698) Get_testcases() []dao.TestCase {
	return t.testcases
}

// Enter your ticket information here
func (t *Ticket_698) Set_values() {
	t.no = 698 // Enter your ticket id
	t.description = "Reconnect and retry case for connection agent(Job Icon, File Transfer, Reboot)"
}

func (t *Ticket_698) Add_testcases() {

	tc_1 := t.New_testcase(1, "Reconnect and retry case for connection agent(Job Icon)")
	tc_func := func() common.Testcase_status {
		if err := lib.Jobarg_enable_jobnet("Icon_1", "jobicon_linux"); err != nil {
			tc_1.Err_log("Failed to enable jobnet, Error: %s", err)
			return FAILED
		}
		return CheckWhetherTheJobIconRunsNormally("Icon_1", tc_1, "END", "NORMAL", 3, common.Client)
	}
	tc_1.Set_function(tc_func)
	t.Add_testcase(*tc_1)

	tc_2 := t.New_testcase(2, "Agent connection failure after 120 seconds.")
	tc_func = func() common.Testcase_status {
		if err := lib.Jobarg_enable_jobnet("Icon_1", "jobicon_linux"); err != nil {
			tc_2.Err_log("Failed to enable jobnet, Error: %s", err)
			return FAILED
		}
		return AgentConnectionReconnection("Icon_1", 180, tc_2, "RUN", "ERROR", 2, common.Client)
	}
	tc_2.Set_function(tc_func)
	t.Add_testcase(*tc_2)

	tc_3 := t.New_testcase(3, "Agent reconnection within 120 seconds.")
	tc_func = func() common.Testcase_status {
		if err := lib.Jobarg_enable_jobnet("Icon_1", "jobicon_linux"); err != nil {
			tc_3.Err_log("Failed to enable jobnet, Error: %s", err)
			return FAILED
		}
		return AgentConnectionReconnection("Icon_1", 60, tc_3, "END", "NORMAL", 2, common.Client)
	}
	tc_3.Set_function(tc_func)
	t.Add_testcase(*tc_3)
}

func CheckWhetherTheJobIconRunsNormally(jobnetId string, testcase *dao.TestCase, targetJobnetStatus string, targetJobStatus string, processCheckTimeout int, sshClient *ssh.Client) common.Testcase_status {
	// Clean up the agent
	err := lib.Jobarg_cleanup_linux()
	if err != nil {
		fmt.Println(testcase.Err_log("Error: %s, Failed to clean up the linux agent.", err.Error()))
		return FAILED
	}
	fmt.Println(testcase.Info_log("Clean up agent service success."))

	envs, _ := lib.Get_str_str_map("JA_HOSTNAME", "oss.linux", "JA_CMD", "sleep 10")

	// Run jobnet
	run_jobnet_id, err := lib.Jobarg_exec_E(jobnetId, envs)
	if err != nil {
		fmt.Println(testcase.Err_log("Error: %s, std_out: %s", err.Error(), run_jobnet_id))
		return FAILED
	}
	fmt.Println(testcase.Info_log("%s has been successfully run with registry number: %s", jobnetId, run_jobnet_id))

	// Wait jobnet finishes and get jobnet run info.
	jobnet_run_info, err := lib.Jobarg_get_jobnet_info(run_jobnet_id, targetJobnetStatus, targetJobStatus, processCheckTimeout)
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

	// Success (obtain the expected status, message, or exit code)
	fmt.Println(testcase.Info_log("Jobnet_status: %s, Job_status: %s, Exit_cd: %d", jobnet_run_info.Jobnet_status, jobnet_run_info.Job_status, jobnet_run_info.Exit_cd))

	return PASSED
}

func AgentConnectionReconnection(jobnetId string, sleepCount int, testcase *dao.TestCase, targetJobnetStatus string, targetJobStatus string, processCheckTimeout int, sshClient *ssh.Client) common.Testcase_status {

	// Clean up the agent
	err := lib.Jobarg_cleanup_linux()
	if err != nil {
		fmt.Println(testcase.Err_log("Error: %s, Failed to clean up the linux agent.", err.Error()))
		return FAILED
	}
	fmt.Println(testcase.Info_log("Clean up agent service success."))

	// Stop the agent
	err = lib.Stop_jaz_agent_linux()
	if err != nil {
		fmt.Println(testcase.Err_log("Error : %s, Failed to stop the agent.", err.Error()))
		return FAILED
	}
	fmt.Println(testcase.Info_log("Successful stop the agent."))

	envs, _ := lib.Get_str_str_map("JA_HOSTNAME", "oss.linux", "JA_CMD", "sleep 10")

	// Run jobnet
	run_jobnet_id, err := lib.Jobarg_exec_E(jobnetId, envs)
	if err != nil {
		fmt.Println(testcase.Err_log("Error: %s, std_out: %s", err.Error(), run_jobnet_id))
		return FAILED
	}
	fmt.Println(testcase.Info_log("%s has been successfully run with registry number: %s", jobnetId, run_jobnet_id))

	// Sleep the virtual machine
	err = lib.Sleep_linux(sleepCount)
	if err != nil {
		fmt.Println(testcase.Err_log("Error: %s, Failed to sleep your linux os. Sleep count: %d", err.Error(), sleepCount))
		return FAILED
	}
	fmt.Println(testcase.Info_log("Successfully sleep your linux os in bg. Sleep count: %d", sleepCount))

	// Restart the agent
	err = lib.Restart_jaz_agent_linux()
	if err != nil {
		fmt.Println(testcase.Err_log("Error : %s, Failed to restart the agent.", err.Error()))
		return FAILED
	}
	fmt.Println(testcase.Info_log("Successful restart the agent."))

	// Wait jobnet finishes and get jobnet run info.
	jobnet_run_info, err := lib.Jobarg_get_jobnet_info(run_jobnet_id, targetJobnetStatus, targetJobStatus, processCheckTimeout)
	if err != nil || jobnet_run_info == nil {
		fmt.Println(testcase.Err_log("Error getting jobnet info: %s", err.Error()))
	}
	fmt.Println(testcase.Info_log("%s with registry number %s is completed.", jobnetId, run_jobnet_id))

	// Check jobnet run status and exit code.
	if jobnet_run_info.Jobnet_status != targetJobnetStatus {
		fmt.Println(testcase.Err_log("Unexpected Jobnet status. Jobnet_status: %s, Job_status: %s, Exit_cd: %d", jobnet_run_info.Jobnet_status, jobnet_run_info.Job_status, jobnet_run_info.Exit_cd))
		return FAILED
	}

	// Success (obtain the expected status, message, or exit code)
	fmt.Println(testcase.Info_log("Jobnet_status: %s, Job_status: %s, Exit_cd: %d", jobnet_run_info.Jobnet_status, jobnet_run_info.Job_status, jobnet_run_info.Exit_cd))

	return PASSED
}
