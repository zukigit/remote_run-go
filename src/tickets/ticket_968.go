package tickets

import (
	"fmt"
	"strings"

	"github.com/zukigit/remote_run-go/src/common"
	"github.com/zukigit/remote_run-go/src/dao"
	"github.com/zukigit/remote_run-go/src/lib"
)

type Ticket_968 struct {
	Ticket_no                                   uint
	Ticket_description                          string
	PASSED_count, FAILED_count, MUSTCHECK_count int
	Testcases                                   []dao.TestCase
}

func (t *Ticket_968) New_testcase(testcase_id uint, testcase_description string) *dao.TestCase {
	return dao.New_testcase(testcase_id, testcase_description)
}

func (t *Ticket_968) Get_no() uint {
	return t.Ticket_no
}

func (t *Ticket_968) Set_PASSED_count(passed_count int) {
	t.PASSED_count = passed_count
}

func (t *Ticket_968) Set_FAILED_count(failed_count int) {
	t.FAILED_count = failed_count
}

func (t *Ticket_968) Set_MUSTCHECK_count(mustcheck_count int) {
	t.MUSTCHECK_count = mustcheck_count
}

func (t *Ticket_968) Get_dsctn() string {
	return t.Ticket_description
}

func (t *Ticket_968) Add_testcase(tc dao.TestCase) {
	t.Testcases = append(t.Testcases, tc)
}

func (t *Ticket_968) Get_testcases() []dao.TestCase {
	return t.Testcases
}

func (t *Ticket_968) Set_values() {
	t.Ticket_no = 968
	t.Ticket_description = "Database connection is limited with transaction file."
}
func (t *Ticket_968) Add_testcases() {

	tc_97 := t.New_testcase(97, "Normal Case Check")
	tc_func := func() common.Testcase_status {
		// Stop the jobarg_server to apply config changes

		if err := lib.Jobarg_cleanup_linux(); err != nil {
			return t.logError(tc_97, "Error during cleanup: %s", err)
		}
		if err := lib.Stop_jaz_server(); err != nil {
			return t.logError(tc_97, "Failed to stop jobarg-server, Error: %s", err.Error())
		}

		fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "JAZ Server has been stopped."))

		const DBconCount = "10"
		configFilePath := "/etc/jobarranger/jobarg_server.conf"
		sedCmd := fmt.Sprintf(`sed -i -e '$a\DBconCount=%s' %s`, DBconCount, configFilePath)

		// Execute sed command to modify DBconCount
		output, err := lib.Ssh_exec_to_str(sedCmd)
		if err != nil {
			return t.logError(tc_97, "Failed to set server config for DBconCount, Error: %s, Output: %s", err.Error(), output)
		}
		fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "DBconCount has been set to %s.", DBconCount))

		// Restart jobarg_server to apply new configuration
		if err := lib.Restart_jaz_server(); err != nil {
			return t.logError(tc_97, "Failed to restart jobarg-server, Error: %s", err.Error())
		}
		fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "JAZ Server has been restarted."))

		// Check for lock files
		status := t.CheckDBconCount(tc_97)
		if status != PASSED {
			return status
		}

		// Enable Jobnet
		if err := lib.Jobarg_enable_jobnet("Icon_1", "agentless"); err != nil {
			return t.logError(tc_97, "Error during enabling jobnet: %s", err)
		}

		agentlessStatus := t.runAgentless100(tc_97)
		if agentlessStatus != PASSED {

			t.logError(tc_97, "runAgentless100 failed.")
		}
		return PASSED
	}
	tc_97.Set_function(tc_func)
	t.Add_testcase(*tc_97)

	/////////////////////////////98//////////////////////////////

	tc_98 := t.New_testcase(98, "Abormal Case Check -Delete files")
	tc_func = func() common.Testcase_status {
		// Stop the jobarg_server to apply config changes

		if err := lib.Jobarg_cleanup_linux(); err != nil {
			return t.logError(tc_98, "Error during cleanup: %s", err)
		}

		if err := lib.Jobarg_enable_jobnet("Icon_1", "agentless"); err != nil {
			return t.logError(tc_98, "Error during enabling jobnet: %s", err)
		}

		agentlessStatus := t.runAgentless(tc_98)
		if agentlessStatus != PASSED {

			t.logError(tc_98, "Agentless Icon failed.")
			return FAILED
		}

		return PASSED
	}
	tc_98.Set_function(tc_func)
	t.Add_testcase(*tc_98)

	/////////////////////////////99//////////////////////////////

	tc_99 := t.New_testcase(99, "Abnormal Case Check - Delete Folder")
	tc_func = func() common.Testcase_status {
		if err := lib.Jobarg_cleanup_linux(); err != nil {
			return t.logError(tc_99, "Error during cleanup: %s", err)
		}
		if err := lib.Jobarg_enable_jobnet("Icon_1", "agentless"); err != nil {
			return t.logError(tc_99, "Error during enabling jobnet: %s", err)
		}

		agentlessStatus := t.runAgentless2(tc_99)
		if agentlessStatus != PASSED {
			t.logError(tc_99, "Agentless Icon failed.")
			return FAILED
		}

		return PASSED
	}
	tc_99.Set_function(tc_func)
	t.Add_testcase(*tc_99)
}
func (t *Ticket_968) logError(tc_97 *dao.TestCase, format string, args ...interface{}) common.Testcase_status {
	fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, format, args...))
	return FAILED
}

// check for the presence of 10 lock files
func (t *Ticket_968) CheckDBconCount(tc_97 *dao.TestCase) common.Testcase_status {

	output, err := lib.Ssh_exec("ls /var/log/jobarranger/session_dbc_locks")
	if err != nil {
		return t.logError(tc_97, "Failed to list lock files: %s", err)
	}

	files := strings.Split(strings.TrimSpace(string(output)), "\n")
	fmt.Printf("Lock files: %s\n  ", output)

	if len(files) != 10 {
		return t.logError(tc_97, "Expected 10 lock files, but found %d", len(files))
	}

	return PASSED
}

// /////////////////ENV_SETUP/////////////////////
func (t *Ticket_968) setupEnv() (map[string]string, error) {
	return lib.Get_str_str_map("JA_HOSTNAME", "oss.linux", "JA_CMD", "sleep 30")
}
func (t *Ticket_968) setupEnv2() (map[string]string, error) {
	return lib.Get_str_str_map("JA_HOSTNAME", "oss.linux", "JA_CMD", "sleep 3")
}

func (t *Ticket_968) executeJob(job string, envs map[string]string, tc_97 *dao.TestCase) (string, common.Testcase_status) {
	fmt.Printf("Starting job execution for: %s with envs: %v\n", job, envs)

	runJobnetID, err := lib.Jobarg_exec_E(job, envs)
	if err != nil {
		return "", t.logError(tc_97, "Error executing job %s: %s", job, err)
	}

	fmt.Printf("Received jobnet ID: %s\n", runJobnetID)
	return runJobnetID, PASSED
}

func (t *Ticket_968) runAgentless100(tc_97 *dao.TestCase) common.Testcase_status {
	envs, err := t.setupEnv()
	if err != nil {
		return t.logError(tc_97, "Error getting environment variables: %s", err)
	}

	jobID, status := t.executeJob("Icon_10", envs, tc_97)
	if status != PASSED {
		return status
	}
	return t.checkJobStatus(jobID, tc_97)
}

func (t *Ticket_968) runAgentless(tc_98 *dao.TestCase) common.Testcase_status {
	envs, err := t.setupEnv2()
	if err != nil {
		return t.logError(tc_98, "Error getting environment variables: %s", err)
	}
	jobID, status := t.executeJob("Icon_1", envs, tc_98)
	if status != PASSED {
		return status
	}

	if status := t.deleteLockFiles(tc_98); status != PASSED {
		return status
	}

	if status := t.checkLog(tc_98); status != PASSED {
		return status
	}
	t.checkJobStatus(jobID, tc_98)
	return PASSED
}

func (t *Ticket_968) runAgentless2(tc_98 *dao.TestCase) common.Testcase_status {
	envs, err := t.setupEnv2()
	if err != nil {
		return t.logError(tc_98, "Error getting environment variables: %s", err)
	}
	jobID, status := t.executeJob("Icon_1", envs, tc_98)
	if status != PASSED {
		return status
	}

	if status := t.deleteLockFiles2(tc_98); status != PASSED {
		return status
	}

	if status := t.checkLog(tc_98); status != PASSED {
		return status
	}
	t.checkJobStatus(jobID, tc_98)
	return PASSED
}

func (t *Ticket_968) deleteLockFiles(tc_98 *dao.TestCase) common.Testcase_status {
	cmd := "rm -rf /var/log/jobarranger/session_dbc_locks/*" // Delete all files in the folder
	_, err := lib.Ssh_exec(cmd)
	if err != nil {
		return t.logError(tc_98, "Failed to delete lock files: %s", err)
	}
	return PASSED
}

func (t *Ticket_968) deleteLockFiles2(tc_99 *dao.TestCase) common.Testcase_status {
	cmd := "rm -rf /var/log/jobarranger/session_dbc_locks/" // Delete all files in the folder
	_, err := lib.Ssh_exec(cmd)
	if err != nil {
		return t.logError(tc_99, "Failed to delete lock files: %s", err) // Fix the error argument to tc_99 instead of tc_98
	}
	return PASSED
}

func (t *Ticket_968) checkLog(tc_98 *dao.TestCase) common.Testcase_status {
	const logFilePath = "/var/log/jobarranger/jobarg_server.log"
	const logFileWarning = `In ja_session_dbc_lock(), /var/log/jobarranger/session_dbc_locks does not exist`

	// Prepare the command to search for the warning log in the jobarg_server.log file
	cmd := fmt.Sprintf(`cat %s | grep "%s"`, logFilePath, logFileWarning)
	lib.Logi(common.LOG_LEVEL_INFO, "Executing command: %s", cmd)

	warningLogOutput, err := lib.Ssh_exec_to_str(cmd)
	if err != nil {
		lib.Logi(common.LOG_LEVEL_INFO, "checking for warning log: %s", err)
	}

	if strings.Contains(warningLogOutput, logFileWarning) {
		lib.Logi(common.LOG_LEVEL_INFO, "Warning log found, returning Passed.")
		return PASSED
	}

	lib.Logi(common.LOG_LEVEL_ERR, "Warning log not found, returning FAILED.")
	return FAILED
}

func (t *Ticket_968) checkJobStatus(runJobnetID string, tc_97 *dao.TestCase) common.Testcase_status {
	jobnetRunInfo, err := lib.Jobarg_get_jobnet_run_info(runJobnetID)
	if err != nil {
		return t.logError(tc_97, "Error retrieving jobnet run info or received nil jobnet_run_info: %s", err.Error())
	}

	if jobnetRunInfo.Job_status == "NORMAL" && jobnetRunInfo.Exit_cd == 0 {
		fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "Job %s completed successfully.", runJobnetID))
		return PASSED
	}

	return t.logError(tc_97, "Job %s failed. Jobnet_status: %s, Job_status: %s, Exit_cd: %d", runJobnetID, jobnetRunInfo.Jobnet_status, jobnetRunInfo.Job_status, jobnetRunInfo.Exit_cd)
}
