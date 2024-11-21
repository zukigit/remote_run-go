package tickets

import (
	"fmt"
	"strings"
	"time"

	"github.com/zukigit/remote_run-go/src/common"
	"github.com/zukigit/remote_run-go/src/dao"
	"github.com/zukigit/remote_run-go/src/lib"
	"golang.org/x/crypto/ssh"
)

type Ticket_698 struct {
	Ticket_no                                   uint
	Ticket_description                          string
	PASSED_count, FAILED_count, MUSTCHECK_count int
	Testcases                                   []dao.TestCase
}

func (t *Ticket_698) New_testcase(testcase_id uint, testcase_description string) *dao.TestCase {
	return dao.New_testcase(testcase_id, testcase_description)
}

func (t *Ticket_698) Get_no() uint {
	return t.Ticket_no
}

func (t *Ticket_698) Set_PASSED_count(passed_count int) {
	t.PASSED_count = passed_count
}

func (t *Ticket_698) Set_FAILED_count(failed_count int) {
	t.FAILED_count = failed_count
}

func (t *Ticket_698) Set_MUSTCHECK_count(mustcheck_count int) {
	t.MUSTCHECK_count = mustcheck_count
}

func (t *Ticket_698) Get_dsctn() string {
	return t.Ticket_description
}

func (t *Ticket_698) Add_testcase(tc dao.TestCase) {
	t.Testcases = append(t.Testcases, tc)
}

func (t *Ticket_698) Get_testcases() []dao.TestCase {
	return t.Testcases
}

// Enter your ticket information here
func (t *Ticket_698) Set_values() {
	t.Ticket_no = 698 // Enter your ticket id
	t.Ticket_description = "Reconnect and retry case for connection agent(Job Icon, File Transfer, Reboot)"
}

func (t *Ticket_698) Add_testcases() {

	tc_1 := t.New_testcase(1, "Reconnect and retry case for connection agent(Job Icon)")
	tc_func := func() common.Testcase_status {

		if err := lib.Start_jaz_server(); err != nil {
			tc_1.Err_log("Failed to start jobarg-server, Error: %s", err.Error())
			return FAILED
		}
		fmt.Println(tc_1.Info_log("JAZ Server has been started."))

		if err := lib.Jobarg_enable_jobnet("Icon_1", "jobicon_linux"); err != nil {
			tc_1.Err_log("Failed to enable jobnet, Error: %s", err)
			return FAILED
		}
		return CheckWhetherTheJobIconRunsNormally("Icon_1", tc_1, "END", "NORMAL", 1, "", common.Client)
	}
	tc_1.Set_function(tc_func)
	t.Add_testcase(*tc_1)

	tc_2 := t.New_testcase(2, "Agent connection failure after 120 seconds.")
	tc_func = func() common.Testcase_status {

		if err := lib.Start_jaz_server(); err != nil {
			tc_2.Err_log("Failed to start jobarg-server, Error: %s", err.Error())
			return FAILED
		}
		fmt.Println(tc_2.Info_log("JAZ Server has been started."))

		if err := lib.Jobarg_enable_jobnet("Icon_1", "jobicon_linux"); err != nil {
			tc_2.Err_log("Failed to enable jobnet, Error: %s", err)
			return FAILED
		}
		return AgentConnectionReconnection("Icon_1", 180, tc_2, "RUN", "ERROR", 1, "", common.Client)
	}
	tc_2.Set_function(tc_func)
	t.Add_testcase(*tc_2)

	tc_3 := t.New_testcase(3, "Agent reconnection within 120 seconds.")
	tc_func = func() common.Testcase_status {
		if err := lib.Start_jaz_server(); err != nil {
			tc_3.Err_log("Failed to start jobarg-server, Error: %s", err.Error())
			return FAILED
		}
		fmt.Println(tc_3.Info_log("JAZ Server has been started."))

		if err := lib.Jobarg_enable_jobnet("Icon_1", "jobicon_linux"); err != nil {
			tc_3.Err_log("Failed to enable jobnet, Error: %s", err)
			return FAILED
		}
		return AgentConnectionReconnection("Icon_1", 60, tc_3, "END", "NORMAL", 1, "", common.Client)
	}
	tc_3.Set_function(tc_func)
	t.Add_testcase(*tc_3)

	tc_4 := t.New_testcase(4, "Reconnect and retry case for connection agent(FCopy)")
	tc_func = func() common.Testcase_status {
		if err := lib.Start_jaz_server(); err != nil {
			tc_4.Err_log("Failed to start jobarg-server, Error: %s", err.Error())
			return FAILED
		}
		fmt.Println(tc_4.Info_log("JAZ Server has been started."))

		if err := lib.Jobarg_enable_jobnet("Icon_1", "fcopy_overwrite"); err != nil {
			tc_4.Err_log("Failed to enable jobnet, Error: %s", err)
			return FAILED
		}
		return CheckWhetherTheJobIconRunsNormally("Icon_1", tc_4, "END", "NORMAL", 1, "rm.sh", common.Client)
	}
	tc_4.Set_function(tc_func)
	t.Add_testcase(*tc_4)

	tc_5 := t.New_testcase(5, "Agent connection failure after 120 seconds. (FCopy)")
	tc_func = func() common.Testcase_status {

		if err := lib.Start_jaz_server(); err != nil {
			tc_5.Err_log("Failed to start jobarg-server, Error: %s", err.Error())
			return FAILED
		}
		fmt.Println(tc_5.Info_log("JAZ Server has been started."))

		if err := lib.Jobarg_enable_jobnet("Icon_1", "fcopy_overwrite"); err != nil {
			tc_5.Err_log("Failed to enable jobnet, Error: %s", err)
			return FAILED
		}
		return AgentConnectionReconnection("Icon_1", 180, tc_5, "RUN", "ERROR", 1, "rm.sh", common.Client)
	}
	tc_5.Set_function(tc_func)
	t.Add_testcase(*tc_5)

	tc_6 := t.New_testcase(6, "Agent reconnection within 120 seconds. (FCopy)")
	tc_func = func() common.Testcase_status {

		if err := lib.Start_jaz_server(); err != nil {
			tc_6.Err_log("Failed to start jobarg-server, Error: %s", err.Error())
			return FAILED
		}
		fmt.Println(tc_6.Info_log("JAZ Server has been started."))

		if err := lib.Jobarg_enable_jobnet("Icon_1", "fcopy_overwrite"); err != nil {
			tc_6.Err_log("Failed to enable jobnet, Error: %s", err)
			return FAILED
		}
		return AgentConnectionReconnection("Icon_1", 60, tc_6, "END", "NORMAL", 1, "rm.sh", common.Client)
	}
	tc_6.Set_function(tc_func)
	t.Add_testcase(*tc_6)

	tc_7 := t.New_testcase(7, "Reconnect and retry case for connection agent(Reboot)")
	tc_func = func() common.Testcase_status {

		if err := lib.Start_jaz_server(); err != nil {
			tc_7.Err_log("Failed to start jobarg-server, Error: %s", err.Error())
			return FAILED
		}
		fmt.Println(tc_7.Info_log("JAZ Server has been started."))

		if err := lib.Jobarg_enable_jobnet("Icon_1", "reboot_icon"); err != nil {
			tc_7.Err_log("Failed to enable jobnet, Error: %s", err)
			return FAILED
		}
		return CheckWhetherTheJobIconRunsNormally("Icon_1", tc_7, "END", "NORMAL", 1, "", common.Client)
	}
	tc_7.Set_function(tc_func)
	t.Add_testcase(*tc_7)

	tc_8 := t.New_testcase(8, "Agent connection failure after 120 seconds. (Reboot)")
	tc_func = func() common.Testcase_status {

		sleepDuration := 1 * time.Minute
		time.Sleep(sleepDuration)

		common.Client = lib.ConnectWithRetry(common.Login_info.Hostname, common.Login_info.Port, common.Login_info.Username, common.Login_info.Password, 60)

		result2, err := lib.GetOutputStrFromSSHCommand(common.Client, "hostname")

		if err != nil {
			fmt.Println(tc_8.Err_log("error"))
			return FAILED
		}
		fmt.Print(tc_8.Info_log("Successfully rebooted hostname : %s", result2))

		if err := lib.Start_jaz_server(); err != nil {
			tc_8.Err_log("Failed to start jobarg-server, Error: %s", err.Error())
			return FAILED
		}
		fmt.Println(tc_8.Info_log("JAZ Server has been started."))

		if err := lib.Jobarg_enable_jobnet("Icon_1", "reboot_icon"); err != nil {
			tc_8.Err_log("Failed to enable jobnet, Error: %s", err)
			return FAILED
		}
		return AgentConnectionReconnection("Icon_1", 180, tc_8, "RUN", "ERROR", 1, "", common.Client)
	}
	tc_8.Set_function(tc_func)
	t.Add_testcase(*tc_8)

	tc_9 := t.New_testcase(9, "Agent reconnection within 120 seconds. (Reboot)")
	tc_func = func() common.Testcase_status {

		sleepDuration := 1 * time.Minute
		time.Sleep(sleepDuration)

		common.Client = lib.ConnectWithRetry(common.Login_info.Hostname, common.Login_info.Port, common.Login_info.Username, common.Login_info.Password, 60)

		result2, err := lib.GetOutputStrFromSSHCommand(common.Client, "hostname")

		if err != nil {
			fmt.Println(tc_9.Err_log("error"))
			return FAILED
		}
		fmt.Print(tc_9.Info_log("Successfully rebooted hostname : %s", result2))

		if err := lib.Start_jaz_server(); err != nil {
			tc_9.Err_log("Failed to start jobarg-server, Error: %s", err.Error())
			return FAILED
		}
		fmt.Println(tc_9.Info_log("JAZ Server has been started."))

		if err := lib.Jobarg_enable_jobnet("Icon_1", "reboot_icon"); err != nil {
			tc_9.Err_log("Failed to enable jobnet, Error: %s", err)
			return FAILED
		}
		return AgentConnectionReconnection("Icon_1", 60, tc_9, "END", "NORMAL", 1, "", common.Client)
	}
	tc_9.Set_function(tc_func)
	t.Add_testcase(*tc_9)

	tc_10 := t.New_testcase(10, "Reconnect and retry case for connection agent(Reboot)")
	tc_func = func() common.Testcase_status {

		sleepDuration := 1 * time.Minute
		time.Sleep(sleepDuration)

		common.Client = lib.ConnectWithRetry(common.Login_info.Hostname, common.Login_info.Port, common.Login_info.Username, common.Login_info.Password, 60)

		result2, err := lib.GetOutputStrFromSSHCommand(common.Client, "hostname")

		if err != nil {
			fmt.Println(tc_10.Err_log("error"))
			return FAILED
		}
		fmt.Print(tc_10.Info_log("Successfully rebooted hostname : %s", result2))

		if err := lib.Start_jaz_server(); err != nil {
			tc_10.Err_log("Failed to start jobarg-server, Error: %s", err.Error())
			return FAILED
		}
		fmt.Println(tc_10.Info_log("JAZ Server has been started."))

		if err := lib.Jobarg_enable_jobnet("Icon_1", "reboot_after_job_completing_jobs"); err != nil {
			tc_10.Err_log("Failed to enable jobnet, Error: %s", err)
			return FAILED
		}
		return CheckWhetherTheJobIconRunsNormally("Icon_1", tc_10, "END", "NORMAL", 1, "", common.Client)
	}
	tc_10.Set_function(tc_func)
	t.Add_testcase(*tc_10)

	tc_11 := t.New_testcase(11, "F-Transfer Icon with allow root = 0 (Root Permission File)")
	tc_func = func() common.Testcase_status {

		sleepDuration := 1 * time.Minute
		time.Sleep(sleepDuration)

		common.Client = lib.ConnectWithRetry(common.Login_info.Hostname, common.Login_info.Port, common.Login_info.Username, common.Login_info.Password, 60)

		result2, err := lib.GetOutputStrFromSSHCommand(common.Client, "hostname")

		if err != nil {
			fmt.Println(tc_11.Err_log("error"))
			return FAILED
		}
		fmt.Print(tc_11.Info_log("Successfully rebooted hostname : %s", result2))

		if err := lib.Start_jaz_server(); err != nil {
			tc_11.Err_log("Failed to start jobarg-server, Error: %s", err.Error())
			return FAILED
		}
		fmt.Println(tc_11.Info_log("JAZ Server has been started."))

		if err := lib.Jobarg_enable_jobnet("Icon_1", "fcopy_overwrite"); err != nil {
			tc_11.Err_log("Failed to enable jobnet, Error: %s", err)
			return FAILED
		}
		return AllowRootFileTransfer("Icon_1", tc_11, "RUN", "ERROR", 1, common.Client)
	}
	tc_11.Set_function(tc_func)
	t.Add_testcase(*tc_11)

	tc_12 := t.New_testcase(12, "F-Transfer Icon with allow root = 0 (Zabbix Permission File)")
	tc_func = func() common.Testcase_status {

		sleepDuration := 1 * time.Minute
		time.Sleep(sleepDuration)

		common.Client = lib.ConnectWithRetry(common.Login_info.Hostname, common.Login_info.Port, common.Login_info.Username, common.Login_info.Password, 60)

		result2, err := lib.GetOutputStrFromSSHCommand(common.Client, "hostname")

		if err != nil {
			fmt.Println(tc_12.Err_log("error"))
			return FAILED
		}
		fmt.Print(tc_12.Info_log("Successfully rebooted hostname : %s", result2))

		if err := lib.Start_jaz_server(); err != nil {
			tc_12.Err_log("Failed to start jobarg-server, Error: %s", err.Error())
			return FAILED
		}
		fmt.Println(tc_12.Info_log("JAZ Server has been started."))

		if err := lib.Jobarg_enable_jobnet("Icon_1", "fcopy_overwrite_zabbix"); err != nil {
			tc_12.Err_log("Failed to enable jobnet, Error: %s", err)
			return FAILED
		}
		return AllowRootFileTransfer("Icon_1", tc_12, "END", "NORMAL", 1, common.Client)
	}
	tc_12.Set_function(tc_func)
	t.Add_testcase(*tc_12)

	tc_13 := t.New_testcase(13, "Reboot after completing Job (AllowRoot 1) <Timeout=30>")
	tc_func = func() common.Testcase_status {

		sleepDuration := 1 * time.Minute
		time.Sleep(sleepDuration)

		common.Client = lib.ConnectWithRetry(common.Login_info.Hostname, common.Login_info.Port, common.Login_info.Username, common.Login_info.Password, 60)

		result2, err := lib.GetOutputStrFromSSHCommand(common.Client, "hostname")

		if err != nil {
			fmt.Println(tc_13.Err_log("error"))
			return FAILED
		}
		fmt.Print(tc_13.Info_log("Successfully rebooted hostname : %s", result2))

		if err := lib.Start_jaz_server(); err != nil {
			tc_13.Err_log("Failed to start jobarg-server, Error: %s", err.Error())
			return FAILED
		}
		fmt.Println(tc_13.Info_log("JAZ Server has been started."))

		if err := lib.Jobarg_enable_jobnet("Icon_1", "jobicon_linux"); err != nil {
			tc_13.Err_log("Failed to enable jobnet, Error: %s", err)
			return FAILED
		}

		return RebootAfterCompletingJob_30("Icon_1", tc_13, "RUN", "ERROR", "END", "NORMAL", 6, common.Client)

	}
	tc_13.Set_function(tc_func)
	t.Add_testcase(*tc_13)

	tc_14 := t.New_testcase(14, "Reboot after completing Job (AllowRoot 1) <Timeout=120>")
	tc_func = func() common.Testcase_status {

		sleepDuration := 1 * time.Minute
		time.Sleep(sleepDuration)

		common.Client = lib.ConnectWithRetry(common.Login_info.Hostname, common.Login_info.Port, common.Login_info.Username, common.Login_info.Password, 60)

		result2, err := lib.GetOutputStrFromSSHCommand(common.Client, "hostname")

		if err != nil {
			fmt.Println(tc_14.Err_log("error"))
			return FAILED
		}
		fmt.Print(tc_14.Info_log("Successfully rebooted hostname : %s", result2))

		if err := lib.Start_jaz_server(); err != nil {
			tc_14.Err_log("Failed to start jobarg-server, Error: %s", err.Error())
			return FAILED
		}
		fmt.Println(tc_14.Info_log("JAZ Server has been started."))

		if err := lib.Jobarg_enable_jobnet("Icon_1", "jobicon_linux"); err != nil {
			tc_14.Err_log("Failed to enable jobnet, Error: %s", err)
			return FAILED
		}

		return RebootAfterCompletingJob_120("Icon_1", tc_14, "RUN", "ERROR", "END", "TIMEOUT", 6, common.Client)

	}
	tc_14.Set_function(tc_func)
	t.Add_testcase(*tc_14)

}

func CheckWhetherTheJobIconRunsNormally(jobnetId string, testcase *dao.TestCase, targetJobnetStatus string, targetJobStatus string, processCheckTimeout int, fileName string, sshClient *ssh.Client) common.Testcase_status {
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
	fmt.Println(testcase.Info_log("Jobnet_status: %s, Job_status: %s, Exit_cd: %d", jobnet_run_info.Jobnet_status, jobnet_run_info.Job_status, jobnet_run_info.Exit_cd))

	if testcase.Get_ticket_no() == 4 {
		folderPath := "/tmp/"
		configFilePath := folderPath
		if !strings.HasSuffix(folderPath, "/") {
			configFilePath += "/"
		}
		configFilePath += fileName

		checkFileCmd := fmt.Sprintf(`[ -f %s ] && echo "File exists and is $(du -m %s | cut -f1) MB." || echo "File does not exist."`, configFilePath, configFilePath)

		output, err := lib.Ssh_exec_to_str(checkFileCmd)
		if err != nil {
			fmt.Println(testcase.Err_log("Error: %s Failed to check file existence and size.", err.Error()))
			return FAILED
		}
		fmt.Println(testcase.Info_log(output))
	}

	return PASSED
}

func AgentConnectionReconnection(jobnetId string, sleepCount int, testcase *dao.TestCase, targetJobnetStatus string, targetJobStatus string, processCheckTimeout int, fileName string, sshClient *ssh.Client) common.Testcase_status {

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
	fmt.Println(testcase.Info_log("Jobnet_status: %s, Job_status: %s, Exit_cd: %d", jobnet_run_info.Jobnet_status, jobnet_run_info.Job_status, jobnet_run_info.Exit_cd))

	if testcase.Get_ticket_no() == 5 || testcase.Get_ticket_no() == 6 {
		folderPath := "/tmp/"
		configFilePath := folderPath
		if !strings.HasSuffix(folderPath, "/") {
			configFilePath += "/"
		}
		configFilePath += fileName

		checkFileCmd := fmt.Sprintf(`[ -f %s ] && echo "File exists and is $(du -m %s | cut -f1) MB." || echo "File does not exist."`, configFilePath, configFilePath)

		output, err := lib.Ssh_exec_to_str(checkFileCmd)
		if err != nil {
			fmt.Println(testcase.Err_log("Error: %s Failed to check file existence and size.", err.Error()))
			return FAILED
		}
		fmt.Println(testcase.Info_log(output))
	}

	return PASSED
}

func AllowRootFileTransfer(jobnetId string, testcase *dao.TestCase, targetJobnetStatus string, targetJobStatus string, processCheckTimeout int, sshClient *ssh.Client) common.Testcase_status {

	// Clean up the agent
	err := lib.Jobarg_cleanup_linux()
	if err != nil {
		fmt.Println(testcase.Err_log("Error: %s, Failed to clean up the linux agent.", err.Error()))
		return FAILED
	}
	fmt.Println(testcase.Info_log("Clean up agent service success."))

	configFilePath := "/etc/jobarranger/jobarg_agentd.conf"
	sedCmd := fmt.Sprintf(`sed -i -e '$a\AllowRoot=%d' %s`, 0, configFilePath)

	// Execute sed command to modify AllowRoot
	_, err = lib.Ssh_exec_to_str(sedCmd)
	if err != nil {
		fmt.Println(testcase.Err_log("Error: %s Failed to set jobarg agent config for AllowRoot 0.", err.Error()))
		return FAILED
	}
	fmt.Println(testcase.Info_log("Jobarg agent config change for parameter AllowRoot 0 is success."))

	// Attempt to restart the agent initially
	err = lib.Restart_jaz_agent_linux()
	if err == nil {
		fmt.Println(testcase.Info_log("Agent restarted successfully."))
	} else {
		// Log initial restart failure
		fmt.Println(testcase.Err_log("Error: %s, Failed to restart the agent. Proceeding with cleanup...", err.Error()))

		// Define cleanup commands
		cleanupCommands := []string{
			"rm -rf /var/lib/jobarranger/tmp/*",
			"rm -rf /var/log/jobarranger/jobarg_agentd.log",
		}

		// Run cleanup commands
		for _, cmd := range cleanupCommands {
			_, err := lib.Ssh_exec_to_str(cmd)
			if err != nil {
				fmt.Println(testcase.Err_log("Error: %s, Failed to execute cleanup command: %s", err.Error(), cmd))
				return FAILED
			}
			fmt.Println(testcase.Info_log("Cleanup successful for command: %s", cmd))
		}

		// Attempt to restart the agent again after cleanup
		err = lib.Restart_jaz_agent_linux()
		if err != nil {
			fmt.Println(testcase.Err_log("Error: %s, Failed to restart the agent after cleanup.", err.Error()))
			return FAILED
		}
		fmt.Println(testcase.Info_log("Agent restarted successfully after cleanup."))
	}

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
	if err != nil || jobnet_run_info == nil {
		fmt.Println(testcase.Err_log("Error getting jobnet info: %s", err.Error()))
	}
	fmt.Println(testcase.Info_log("%s with registry number %s is completed.", jobnetId, run_jobnet_id))

	// Check jobnet run status and exit code.
	if jobnet_run_info.Jobnet_status != targetJobnetStatus {
		fmt.Println(testcase.Err_log("Unexpected Jobnet status. Jobnet_status: %s, Job_status: %s, Exit_cd: %d", jobnet_run_info.Jobnet_status, jobnet_run_info.Job_status, jobnet_run_info.Exit_cd))
		return FAILED
	}
	fmt.Println(testcase.Info_log("Jobnet_status: %s, Job_status: %s, Exit_cd: %d", jobnet_run_info.Jobnet_status, jobnet_run_info.Job_status, jobnet_run_info.Exit_cd))

	sedCmd = fmt.Sprintf(`sed -i -e '$a\AllowRoot=%d' %s`, 1, configFilePath)

	// Execute sed command to modify AllowRoot
	_, err = lib.Ssh_exec_to_str(sedCmd)
	if err != nil {
		fmt.Println(testcase.Err_log("Error: %s Failed to set jobarg agent config for AllowRoot 1.", err.Error()))
		return FAILED
	}
	fmt.Println(testcase.Info_log("Jobarg agent config change for parameter AllowRoot 1 is success."))

	return PASSED
}

func RebootAfterCompletingJob_30(jobnetId string, testcase *dao.TestCase, normalJobnetStatus string, normalJobStatus string, rebootJobnetStatus string, rebootJobStatus string, processCheckTimeout int, sshClient *ssh.Client) common.Testcase_status {

	// Clean up the agent
	err := lib.Jobarg_cleanup_linux()
	if err != nil {
		fmt.Println(testcase.Err_log("Error: %s, Failed to clean up the linux agent.", err.Error()))
		return FAILED
	}
	fmt.Println(testcase.Info_log("Clean up agent service success."))

	envs, _ := lib.Get_str_str_map("JA_HOSTNAME", "oss.linux", "JA_CMD", "sleep 200")

	// Run normal jobnet
	normalJobnetId, err := lib.Jobarg_exec_E(jobnetId, envs)
	if err != nil {
		fmt.Println(testcase.Err_log("Error: %s, std_out: %s", err.Error(), normalJobnetId))
		return FAILED
	}
	fmt.Println(testcase.Info_log("%s has been successfully run with registry number: %s", jobnetId, normalJobnetId))

	// Enable the Reboot Jobnet
	if err := lib.Jobarg_enable_jobnet("Icon_1", "reboot_after_job_completing_jobs_timeout_30"); err != nil {
		fmt.Println(testcase.Err_log("Failed to enable reboot jobnet, Error: %s", err))
		return FAILED
	}
	fmt.Println(testcase.Info_log("Enable the reboot jobnet success."))

	// Run reboot jobnet
	rebootJobnetId, err := lib.Jobarg_exec_E(jobnetId, envs)
	if err != nil {
		fmt.Println(testcase.Err_log("Error: %s, std_out: %s", err.Error(), rebootJobnetId))
		return FAILED
	}
	fmt.Println(testcase.Info_log("%s has been successfully run with registry number: %s", jobnetId, rebootJobnetId))

	// Sleep the virtual machine
	err = lib.Sleep_linux(30)
	if err != nil {
		fmt.Println(testcase.Err_log("Error: %s, Failed to sleep your linux os. Sleep count: %d", err.Error(), 30))
		return FAILED
	}
	fmt.Println(testcase.Info_log("Successfully sleep your linux os in bg. Sleep count: %d", 30))

	// Check Reboot Icon is run or not
	var CheckJobStatusQuery lib.DBQuery = `SELECT status FROM ja_run_job_id where job_id = "REBOOT" and inner_jobnet_main_id = $1`
	err = lib.JobProcessDBCountCheck(4, 1, rebootJobnetId, CheckJobStatusQuery)
	if err != nil {
		fmt.Println(testcase.Info_log("Reboot Icon not run in 30s."))
	} else {
		fmt.Println(testcase.Err_log("Reboot Icon run in 30s."))
	}

	// After Reboot Connect with ssh again
	sleepDuration := 1 * time.Minute
	time.Sleep(sleepDuration)

	common.Client = lib.ConnectWithRetry(common.Login_info.Hostname, common.Login_info.Port, common.Login_info.Username, common.Login_info.Password, 60)

	result2, err := lib.GetOutputStrFromSSHCommand(common.Client, "hostname")

	if err != nil {
		fmt.Println(testcase.Err_log("error"))
		return FAILED
	}
	fmt.Print(testcase.Info_log("Successfully rebooted hostname : %s", result2))

	// Wait jobnet finishes and get reboot jobnet run info.
	reboot_jobnet_run_info, err := lib.Jobarg_get_jobnet_info(rebootJobnetId, rebootJobnetStatus, rebootJobStatus, processCheckTimeout)
	if err != nil || reboot_jobnet_run_info == nil {
		fmt.Println(testcase.Err_log("Error getting reboot jobnet info: %s", err.Error()))
	}
	fmt.Println(testcase.Info_log("%s with registry number %s is completed.", jobnetId, rebootJobnetId))

	// Check jobnet run status and exit code.
	if reboot_jobnet_run_info.Jobnet_status != rebootJobnetStatus {
		fmt.Println(testcase.Err_log("Unexpected Reboot Jobnet status. Reboot Jobnet_status: %s, Job_status: %s, Exit_cd: %d", reboot_jobnet_run_info.Jobnet_status, reboot_jobnet_run_info.Job_status, reboot_jobnet_run_info.Exit_cd))
		return FAILED
	}
	fmt.Println(testcase.Info_log("Reboot_Jobnet_status: %s, Job_status: %s, Exit_cd: %d", reboot_jobnet_run_info.Jobnet_status, reboot_jobnet_run_info.Job_status, reboot_jobnet_run_info.Exit_cd))

	// Wait jobnet finishes and get jobnet run info.
	normal_jobnet_run_info, err := lib.Jobarg_get_jobnet_info(normalJobnetId, normalJobnetStatus, normalJobStatus, processCheckTimeout)

	if err != nil || normal_jobnet_run_info == nil {
		fmt.Println(testcase.Err_log("Error getting jobnet info: %s", err.Error()))
	}
	fmt.Println(testcase.Info_log("%s with registry number %s is completed.", jobnetId, normalJobnetId))

	// Check jobnet run status and exit code.
	if normal_jobnet_run_info.Jobnet_status != normalJobnetStatus {
		fmt.Println(testcase.Err_log("Unexpected Jobnet status. Jobnet_status: %s, Job_status: %s, Exit_cd: %d", normal_jobnet_run_info.Jobnet_status, normal_jobnet_run_info.Job_status, normal_jobnet_run_info.Exit_cd))
		return FAILED
	}
	fmt.Println(testcase.Info_log("Jobnet_status: %s, Job_status: %s, Exit_cd: %d", normal_jobnet_run_info.Jobnet_status, normal_jobnet_run_info.Job_status, normal_jobnet_run_info.Exit_cd))

	return PASSED

}

func RebootAfterCompletingJob_120(jobnetId string, testcase *dao.TestCase, normalJobnetStatus string, normalJobStatus string, rebootJobnetStatus string, rebootJobStatus string, processCheckTimeout int, sshClient *ssh.Client) common.Testcase_status {

	// Clean up the agent
	err := lib.Jobarg_cleanup_linux()
	if err != nil {
		fmt.Println(testcase.Err_log("Error: %s, Failed to clean up the linux agent.", err.Error()))
		return FAILED
	}
	fmt.Println(testcase.Info_log("Clean up agent service success."))

	envs, _ := lib.Get_str_str_map("JA_HOSTNAME", "oss.linux", "JA_CMD", "sleep 200")

	// Run normal jobnet
	normalJobnetId, err := lib.Jobarg_exec_E(jobnetId, envs)
	if err != nil {
		fmt.Println(testcase.Err_log("Error: %s, std_out: %s", err.Error(), normalJobnetId))
		return FAILED
	}
	fmt.Println(testcase.Info_log("%s has been successfully run with registry number: %s", jobnetId, normalJobnetId))

	// Enable the Reboot Jobnet
	if err := lib.Jobarg_enable_jobnet("Icon_1", "reboot_after_job_completing_jobs_timeout_120"); err != nil {
		fmt.Println(testcase.Err_log("Failed to enable reboot jobnet, Error: %s", err))
		return FAILED
	}
	fmt.Println(testcase.Info_log("Enable the reboot jobnet success."))

	// Run reboot jobnet
	rebootJobnetId, err := lib.Jobarg_exec_E(jobnetId, envs)
	if err != nil {
		fmt.Println(testcase.Err_log("Error: %s, std_out: %s", err.Error(), rebootJobnetId))
		return FAILED
	}
	fmt.Println(testcase.Info_log("%s has been successfully run with registry number: %s", jobnetId, rebootJobnetId))

	// Sleep the virtual machine
	err = lib.Sleep_linux(55)
	if err != nil {
		fmt.Println(testcase.Err_log("Error: %s, Failed to sleep your linux os. Sleep count: %d", err.Error(), 55))
		return FAILED
	}
	fmt.Println(testcase.Info_log("Successfully sleep your linux os in bg. Sleep count: %d", 55))

	// Check Reboot Icon is run or not
	timeoutDuration := 1.0 / 6.0
	var CheckRebootJobStatusQuery lib.DBQuery = `SELECT count(*) FROM ja_run_job_id where job_id = "REBOOT" AND timeout_flag = 1 AND inner_jobnet_main_id = $1`
	err = lib.JobProcessDBCountCheck(1, int(timeoutDuration), rebootJobnetId, CheckRebootJobStatusQuery)
	if err != nil {
		fmt.Println(testcase.Info_log("Reboot icon become warning state with orange color after 1 minute."))
	} else {
		fmt.Println(testcase.Err_log("Reboot icon become warning state with orange color before 1 minute."))
	}

	// Sleep the virtual machine
	err = lib.Sleep_linux(55)
	if err != nil {
		fmt.Println(testcase.Err_log("Error: %s, Failed to sleep your linux os. Sleep count: %d", err.Error(), 55))
		return FAILED
	}
	fmt.Println(testcase.Info_log("Successfully sleep your linux os in bg. Sleep count: %d", 55))

	// Check Reboot Icon is run or not
	var CheckJobStatusQuery lib.DBQuery = `SELECT count(*) FROM ja_run_job_id where job_id = "REBOOT" AND status = 4 AND inner_jobnet_main_id = $1`
	err = lib.JobProcessDBCountCheck(1, 1, rebootJobnetId, CheckJobStatusQuery)
	if err != nil {
		fmt.Println(testcase.Info_log("Reboot Icon not run in 120s."))
	} else {
		fmt.Println(testcase.Err_log("Reboot Icon run in "))
	}

	// After Reboot Connect with ssh again
	sleepDuration := 1 * time.Minute
	time.Sleep(sleepDuration)

	common.Client = lib.ConnectWithRetry(common.Login_info.Hostname, common.Login_info.Port, common.Login_info.Username, common.Login_info.Password, 60)

	result2, err := lib.GetOutputStrFromSSHCommand(common.Client, "hostname")

	if err != nil {
		fmt.Println(testcase.Err_log("error"))
		return FAILED
	}
	fmt.Print(testcase.Info_log("Successfully rebooted hostname : %s", result2))

	// Wait jobnet finishes and get reboot jobnet run info.
	reboot_jobnet_run_info, err := lib.Jobarg_get_jobnet_info(rebootJobnetId, rebootJobnetStatus, rebootJobStatus, processCheckTimeout)
	if err != nil || reboot_jobnet_run_info == nil {
		fmt.Println(testcase.Err_log("Error getting reboot jobnet info: %s", err.Error()))
	}
	fmt.Println(testcase.Info_log("%s with registry number %s is completed.", jobnetId, rebootJobnetId))

	// Check jobnet run status and exit code.
	if reboot_jobnet_run_info.Jobnet_status != rebootJobnetStatus {
		fmt.Println(testcase.Err_log("Unexpected Reboot Jobnet status. Reboot Jobnet_status: %s, Job_status: %s, Exit_cd: %d", reboot_jobnet_run_info.Jobnet_status, reboot_jobnet_run_info.Job_status, reboot_jobnet_run_info.Exit_cd))
		return FAILED
	}
	fmt.Println(testcase.Info_log("Reboot_Jobnet_status: %s, Job_status: %s, Exit_cd: %d", reboot_jobnet_run_info.Jobnet_status, reboot_jobnet_run_info.Job_status, reboot_jobnet_run_info.Exit_cd))

	// Wait jobnet finishes and get jobnet run info.
	normal_jobnet_run_info, err := lib.Jobarg_get_jobnet_info(normalJobnetId, normalJobnetStatus, normalJobStatus, processCheckTimeout)

	if err != nil || normal_jobnet_run_info == nil {
		fmt.Println(testcase.Err_log("Error getting jobnet info: %s", err.Error()))
	}
	fmt.Println(testcase.Info_log("%s with registry number %s is completed.", jobnetId, normalJobnetId))

	// Check jobnet run status and exit code.
	if normal_jobnet_run_info.Jobnet_status != normalJobnetStatus {
		fmt.Println(testcase.Err_log("Unexpected Jobnet status. Jobnet_status: %s, Job_status: %s, Exit_cd: %d", normal_jobnet_run_info.Jobnet_status, normal_jobnet_run_info.Job_status, normal_jobnet_run_info.Exit_cd))
		return FAILED
	}
	fmt.Println(testcase.Info_log("Jobnet_status: %s, Job_status: %s, Exit_cd: %d", normal_jobnet_run_info.Jobnet_status, normal_jobnet_run_info.Job_status, normal_jobnet_run_info.Exit_cd))

	return PASSED

}
