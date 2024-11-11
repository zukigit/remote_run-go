package tickets

import (
	"fmt"
	"strings"

	"github.com/zukigit/remote_run-go/src/common"
	"github.com/zukigit/remote_run-go/src/dao"
	"github.com/zukigit/remote_run-go/src/lib"
)

type Ticket_943 struct {
	Ticket_no                                   uint
	Ticket_description                          string
	PASSED_count, FAILED_count, MUSTCHECK_count int
	Testcases                                   []dao.TestCase
}

func (t *Ticket_943) New_testcase(testcase_id uint, testcase_description string) *dao.TestCase {
	return dao.New_testcase(testcase_id, testcase_description)
}

func (t *Ticket_943) Get_no() uint {
	return t.Ticket_no
}

func (t *Ticket_943) Set_PASSED_count(passed_count int) {
	t.PASSED_count = passed_count
}

func (t *Ticket_943) Set_FAILED_count(failed_count int) {
	t.FAILED_count = failed_count
}

func (t *Ticket_943) Set_MUSTCHECK_count(mustcheck_count int) {
	t.MUSTCHECK_count = mustcheck_count
}

func (t *Ticket_943) Get_dsctn() string {
	return t.Ticket_description
}

func (t *Ticket_943) Add_testcase(tc dao.TestCase) {
	t.Testcases = append(t.Testcases, tc)
}

func (t *Ticket_943) Get_testcases() []dao.TestCase {
	return t.Testcases
}

func (t *Ticket_943) Set_values() {
	t.Ticket_no = 943
	t.Ticket_description = "Issue 943 - 「cannot connect the host」error occurs"
}
func (t *Ticket_943) Add_testcases() {
	tc_1 := t.New_testcase(1, "1. Create jobnet that has job icon with HOSTNAME from checkbox and variable and sleep 120s.")

	tc_func := func() common.Testcase_status {
		job := "Icon_1"
		if err := lib.Jobarg_cleanup_windows(); err != nil {
			return t.logError(tc_1, "Error during cleanup: %s", err)
		}
		if err := lib.Jobarg_enable_jobnet("Icon_1", "943a"); err != nil {
			return t.logError(tc_1, "Error during enabling jobnet: %s", err)
		}

		status := t.runJob(tc_1, job)
		if status != PASSED {
			return status
		}
		return t.LogNotFound(tc_1)

	}

	tc_1.Set_function(tc_func)
	t.Add_testcase(*tc_1)

	tc_2 := t.New_testcase(2, "1. Create jobnet that has job icon with HOSTNAME from checkbox and variable and sleep 120s.")

	tc_func2 := func() common.Testcase_status {
		job := "Icon_1"
		if err := lib.Jobarg_cleanup_windows(); err != nil {
			return t.logError(tc_2, "Error during cleanup: %s", err)
		}
		if err := lib.Jobarg_enable_jobnet("Icon_1", "943a"); err != nil {
			return t.logError(tc_2, "Error during enabling jobnet: %s", err)
		}
		status := t.runJob2(tc_2, job)
		if status != PASSED {
			return status
		}

		t.resetHostname(tc_2)
		return PASSED

	}

	tc_2.Set_function(tc_func2)
	t.Add_testcase(*tc_2)
}

func (t *Ticket_943) runJob(tc *dao.TestCase, job string) common.Testcase_status {
	envs, err := lib.Get_str_str_map("JA_HOSTNAME", "oss.linux", "JA_CMD", "sleep 20")
	if err != nil {
		return t.logError(tc, "Error retrieving environment variables: %s", err)
	}

	run_jobnet_id, err := lib.Jobarg_exec_E(job, envs)
	if err != nil {
		return t.logError(tc, "Error executing job %s: %s", job, err)
	}
	fmt.Printf("Executed job: %s with run_jobnet_id: %s\n", job, run_jobnet_id)

	jobnet_run_info, err := lib.Jobarg_get_jobnet_run_info(run_jobnet_id)
	if err != nil {
		return t.logError(tc, "Error retrieving run info for job %s: %s", job, err)
	}
	fmt.Printf("Jobnet Run Info: %+v\n", jobnet_run_info)

	if jobnet_run_info.Jobnet_status == "END" && jobnet_run_info.Job_status == "NORMAL" {
		tc.Info_log("%s completed successfully.", job)
		return PASSED
	}

	return t.logError(tc, "%s failed. Jobnet_status: %s, Job_status: %s, Exit_cd: %d", job, jobnet_run_info.Jobnet_status, jobnet_run_info.Job_status, jobnet_run_info.Exit_cd)
}
func (t *Ticket_943) runJob2(tc *dao.TestCase, job string) common.Testcase_status {
	envs, err := lib.Get_str_str_map("JA_HOSTNAME", "oss.linux", "JA_CMD", "sleep 30")
	if err != nil {
		return t.logError(tc, "Error retrieving environment variables: %s", err)
	}

	run_jobnet_id, err := lib.Jobarg_exec_E(job, envs)
	if err != nil {
		return t.logError(tc, "Error executing job %s: %s", job, err)
	}
	fmt.Printf("Executed job: %s with run_jobnet_id: %s\n", job, run_jobnet_id)

	status := t.modifyHostnameForTestcase(tc)
	if status != PASSED {
		return status
	}

	if status := t.LogFound(tc); status != PASSED {

		t.resetHostname(tc)
		return status
	}
	t.resetHostname(tc)
	return FAILED

}

const logFilePath = "/var/log/jobarranger/jobarg_server.log"
const logFileWarning = `In jarun_status_check()`

func (t *Ticket_943) LogNotFound(tc *dao.TestCase) common.Testcase_status {
	cmd := fmt.Sprintf(`cat %s | grep "%s"`, logFilePath, logFileWarning)
	tc.Info_log("Executing command: %s", cmd)

	warningLogOutput, err := lib.Ssh_exec_to_str(cmd)
	fmt.Println(err)

	if strings.Contains(warningLogOutput, logFileWarning) {
		tc.Err_log("Warning log found, returning FAILED.")
		return FAILED
	}

	tc.Info_log("Warning log not found, returning PASSED.")
	return PASSED
}

func (t *Ticket_943) logError(tc *dao.TestCase, format string, args ...interface{}) common.Testcase_status {
	fmt.Println(tc.Err_log(format, args...))
	return FAILED
}

func (t *Ticket_943) LogFound(tc *dao.TestCase) common.Testcase_status {

	cmd := fmt.Sprintf(`cat %s | grep "%s"`, logFilePath, logFileWarning)
	tc.Info_log("Executing command: %s", cmd)

	warningLogOutput, err := lib.Ssh_exec_to_str(cmd)
	fmt.Println(err)

	if strings.Contains(warningLogOutput, logFileWarning) {
		tc.Info_log("Warning log found, returning PASSED.")
		return PASSED

	}
	return FAILED
}

func (t *Ticket_943) modifyHostnameForTestcase(tc *dao.TestCase) common.Testcase_status {
	var changeHostCmd string

	// Check if it's MySQL
	if common.Is_mysql {
		changeHostCmd = `UPDATE hosts SET host = 'new.hostname' WHERE host = 'oss.linux';`
		_, err := common.DB.Exec(changeHostCmd)
		if err != nil {
			return t.logError(tc, "Failed to change hostname in MySQL database: %s", err)
		}
		tc.Info_log("Hostname successfully changed in MySQL database.")

		// Check if it's PostgreSQL
	} else if common.Is_psql {
		changeHostCmd = `UPDATE hosts SET host = 'new.hostname' WHERE host = 'oss.linux';`
		_, err := common.DB.Exec(changeHostCmd)
		if err != nil {
			return t.logError(tc, "Failed to change hostname in PostgreSQL database: %s", err)
		}
		tc.Info_log("Hostname successfully changed in PostgreSQL database.")
	} else {
		// Handle unsupported database types
		return t.logError(tc, "Unsupported database type.")
	}

	return PASSED
}

func (t *Ticket_943) resetHostname(tc *dao.TestCase) common.Testcase_status {
	var resetHostCmd string

	if common.Is_mysql {

		resetHostCmd = `UPDATE hosts SET host = 'oss.linux' WHERE host = 'new.hostname';`
		_, err := common.DB.Exec(resetHostCmd)
		if err != nil {
			return t.logError(tc, "Failed to reset hostname in MySQL database: %s", err)
		}
		tc.Info_log("Hostname successfully reset in MySQL database.")

	} else if common.Is_psql {
		resetHostCmd = `UPDATE hosts SET host = 'oss.linux' WHERE host = 'new.hostname';`
		_, err := common.DB.Exec(resetHostCmd)
		if err != nil {
			return t.logError(tc, "Failed to reset hostname in PostgreSQL database: %s", err)
		}
		tc.Info_log("Hostname successfully reset in PostgreSQL database.")
	} else {
		return t.logError(tc, "Unsupported database type for resetting hostname.")
	}

	return PASSED
}
