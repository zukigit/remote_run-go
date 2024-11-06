package tickets

import (
	"fmt"
	"strings"

	"github.com/zukigit/remote_run-go/src/common"
	"github.com/zukigit/remote_run-go/src/dao"
	"github.com/zukigit/remote_run-go/src/lib"
)

type Ticket_968 struct {
	no          uint
	description string
	testcases   []dao.TestCase
}

func (t *Ticket_968) New_testcase(testcase_id uint, testcase_description string) *dao.TestCase {
	return dao.New_testcase(testcase_id, testcase_description)
}

func (t *Ticket_968) Get_no() uint {
	return t.no
}

func (t *Ticket_968) Get_dsctn() string {
	return t.description
}

func (t *Ticket_968) Add_testcase(tc dao.TestCase) {
	t.testcases = append(t.testcases, tc)
}

func (t *Ticket_968) Get_testcases() []dao.TestCase {
	return t.testcases
}

func (t *Ticket_968) Set_values() {
	t.no = 968
	t.description = "Database connection is limited with transaction file."
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

		fmt.Println(tc_97.Info_log("JAZ Server has been stopped."))

		const DBconCount = "10"
		configFilePath := "/etc/jobarranger/jobarg_server.conf"
		sedCmd := fmt.Sprintf(`sed -i -e '$a\DBconCount=%s' %s`, DBconCount, configFilePath)

		// Execute sed command to modify DBconCount
		output, err := lib.Ssh_exec_to_str(sedCmd)
		if err != nil {
			return t.logError(tc_97, "Failed to set server config for DBconCount, Error: %s, Output: %s", err.Error(), output)
		}
		fmt.Println(tc_97.Info_log("DBconCount has been set to %s.", DBconCount))

		// Restart jobarg_server to apply new configuration
		if err := lib.Restart_jaz_server(); err != nil {
			return t.logError(tc_97, "Failed to restart jobarg-server, Error: %s", err.Error())
		}
		fmt.Println(tc_97.Info_log("JAZ Server has been restarted."))

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
}

func (t *Ticket_968) logError(tc_97 *dao.TestCase, format string, args ...interface{}) common.Testcase_status {
	fmt.Println(tc_97.Err_log(format, args...))
	return FAILED
}

// checks for the presence of 10 lock files
func (t *Ticket_968) CheckDBconCount(tc_97 *dao.TestCase) common.Testcase_status {

	output, err := lib.Ssh_exec("ls /var/log/jobarranger/session_dbc_locks")
	if err != nil {
		return t.logError(tc_97, "Failed to list lock files: %s", err)
	}

	files := strings.Split(strings.TrimSpace(string(output)), "\n")
	fmt.Printf("Lock files: %s\n", output)

	if len(files) != 10 {
		return t.logError(tc_97, "Expected 10 lock files, but found %d", len(files))
	}

	return PASSED
}

func (t *Ticket_968) runAgentless100(tc_97 *dao.TestCase) common.Testcase_status {
	envs, err := t.setupEnv()
	if err != nil {
		return t.logError(tc_97, "Error getting environment variables: %s", err)
	}

	// Execute job
	jobID, status := t.executeJob("Icon_100", envs, tc_97)
	if status != PASSED {
		return status
	}

	// Check job status
	return t.checkJobStatus(jobID, tc_97)
}

func (t *Ticket_968) setupEnv() (map[string]string, error) {
	return lib.Get_str_str_map("JA_HOSTNAME", "oss.linux", "JA_CMD", "sleep 30")
}

func (t *Ticket_968) executeJob(job string, envs map[string]string, tc_97 *dao.TestCase) (string, common.Testcase_status) {
	fmt.Printf("Starting job execution for: %s with envs: %v\n", job, envs)

	runJobnetID, err := lib.Jobarg_exec_E(job, envs)
	if err != nil {
		return "", t.logError(tc_97, "Error executing job %s: %s", job, err)
	}

	fmt.Printf("Received jobnet ID: %s\n", runJobnetID)
	return runJobnetID, PASSED // Return the job ID with a success status
}

func (t *Ticket_968) checkJobStatus(runJobnetID string, tc_97 *dao.TestCase) common.Testcase_status {
	jobnetRunInfo, err := lib.Jobarg_get_jobnet_run_info(runJobnetID)
	if err != nil || jobnetRunInfo == nil {
		return t.logError(tc_97, "Error retrieving jobnet run info or received nil jobnet_run_info.")
	}

	if jobnetRunInfo.Job_status == "NORMAL" && jobnetRunInfo.Exit_cd == 0 {
		fmt.Println(tc_97.Info_log("Job %s completed successfully.", runJobnetID))
		return PASSED
	}

	return t.logError(tc_97, "Job %s failed. Jobnet_status: %s, Job_status: %s, Exit_cd: %d", runJobnetID, jobnetRunInfo.Jobnet_status, jobnetRunInfo.Job_status, jobnetRunInfo.Exit_cd)
}
