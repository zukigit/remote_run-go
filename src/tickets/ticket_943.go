package tickets

import (
	"fmt"
	"strings"

	"github.com/zukigit/remote_run-go/src/common"
	"github.com/zukigit/remote_run-go/src/dao"
	"github.com/zukigit/remote_run-go/src/lib"
)

type Ticket_943 struct {
	no          uint
	description string
	testcases   []dao.TestCase
}

func (t *Ticket_943) New_testcase(testcase_id uint, testcase_description string) *dao.TestCase {
	return dao.New_testcase(testcase_id, testcase_description)
}

func (t *Ticket_943) Get_no() uint {
	return t.no
}

func (t *Ticket_943) Get_dsctn() string {
	return t.description
}

func (t *Ticket_943) Add_testcase(tc dao.TestCase) {
	t.testcases = append(t.testcases, tc)
}

func (t *Ticket_943) Get_testcases() []dao.TestCase {
	return t.testcases
}

func (t *Ticket_943) Set_values() {
	t.no = 943
	t.description = "Issue 943 - 「cannot connect the host」error occurs"
}

func (t *Ticket_943) Add_testcases() {
	t.addTestcase(1, "Run Job Icon without updating HOSTNAME", "Icon_1", false)
	t.addTestcase(2, "Run Job Icon with updated HOSTNAME", "Icon_1", true)
}

func (t *Ticket_943) addTestcase(id uint, description, job string, shouldBeFound bool) {
	tc := t.New_testcase(id, description)
	tcFunc := func() common.Testcase_status {
		return t.runJob(tc, job, shouldBeFound)
	}
	tc.Set_function(tcFunc)
	t.Add_testcase(*tc)
}

func (t *Ticket_943) runJob(tc *dao.TestCase, job string, shouldBeFound bool) common.Testcase_status {
	if err := lib.Jobarg_cleanup_linux(); err != nil {
		return t.logError(tc, "Error during cleanup: %s", err)
	}

	envs, err := lib.Get_str_str_map("JA_HOSTNAME", "oss.linux", "JA_CMD", "sleep 120")
	if err != nil {
		return t.logError(tc, "Error retrieving environment variables: %s", err)
	}
	if err := lib.Jobarg_enable_jobnet("Icon_1", "jobicon_linux"); err != nil {
		return t.logError(tc, "Error during enable jobnet: %s", err)
	}

	runJobnetID, err := lib.Jobarg_exec_E(job, envs)
	if err != nil {
		return t.logError(tc, "Error executing job %s: %s", job, err)
	}
	fmt.Printf("Executed job: %s with runJobnetID: %s\n", job, runJobnetID)

	if shouldBeFound {
		t.updateHostname(tc, "/etc/jobarranger/jobarg_agentd.conf", "oss.linux33")
		if err := lib.Restart_jaz_server(); err != nil {
			return t.logError(tc, "Failed to restart jobarg-server, Error: %s", err.Error())
		}
	}

	jobnetRunInfo, err := lib.Jobarg_get_jobnet_run_info(runJobnetID)
	if err != nil {
		return t.logError(tc, "Error retrieving run info for job %s: %s", job, err)
	}

	if jobnetRunInfo.Jobnet_status == "END" && (jobnetRunInfo.Job_status == "NORMAL" || jobnetRunInfo.Job_status == "ERROR") {
		tc.Info_log("%s completed successfully.", job)
		return t.checkLog(tc, "/var/log/jobarranger/jobarg_server.log", "Can not get host id for hostname", shouldBeFound)
	}

	return t.logError(tc, "%s failed. Jobnet_status: %s, Job_status: %s, Exit_cd: %d", job, jobnetRunInfo.Jobnet_status, jobnetRunInfo.Job_status, jobnetRunInfo.Exit_cd)
}

func (t *Ticket_943) updateHostname(tc *dao.TestCase, configFilePath, newHostname string) error {
	updateCmd := fmt.Sprintf(`sed -i 's/^Hostname=.*/Hostname=%s/' %s`, newHostname, configFilePath)
	_, err := lib.Ssh_exec_to_str(updateCmd)
	if err != nil {
		tc.Err_log("Error updating hostname in config file: %s", err.Error())
		return err
	}
	tc.Info_log("Successfully updated Hostname to '%s' in config file", newHostname)
	return nil
}

func (t *Ticket_943) checkLog(tc *dao.TestCase, logFilePath, logFileWarning string, shouldBeFound bool) common.Testcase_status {
	cmd := fmt.Sprintf(`cat %s | grep "%s"`, logFilePath, logFileWarning)
	tc.Info_log("Executing command: %s", cmd)

	logOutput, err := lib.Ssh_exec_to_str(cmd)
	if err != nil {
		tc.Info_log(" checking for log: %s", err)

	}

	found := strings.Contains(logOutput, logFileWarning)
	if found != shouldBeFound {
		tc.Err_log("Log state mismatch: expected found=%t, but got found=%t. Returning FAILED.", shouldBeFound, found)
		return FAILED
	}

	if shouldBeFound {
		tc.Info_log("Log found as expected. Returning PASSED.")
	} else {
		tc.Info_log("Log not found as expected. Returning PASSED.")
	}
	return PASSED
}

func (t *Ticket_943) logError(tc *dao.TestCase, format string, args ...interface{}) common.Testcase_status {
	tc.Err_log(format, args...)
	return FAILED
}
