package tickets

import (
	"fmt"
	"strings"
	"time"

	"github.com/zukigit/remote_run-go/src/common"
	"github.com/zukigit/remote_run-go/src/lib"
)

type Ticket_943 struct {
	Ticket_no                                   uint
	Ticket_description                          string
	PASSED_count, FAILED_count, MUSTCHECK_count int
	Testcases                                   []common.TestCase
}

func (t *Ticket_943) New_testcase(testcase_id uint, testcase_description string) *common.TestCase {
	return common.New_testcase(testcase_id, testcase_description)
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

func (t *Ticket_943) Add_testcase(tc common.TestCase) {
	t.Testcases = append(t.Testcases, tc)
}

func (t *Ticket_943) Get_testcases() []common.TestCase {
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
		if err := lib.Jobarg_cleanup_linux(); err != nil {
			return t.logError(tc_1, "Error during cleanup: %s", err)
		}
		if err := lib.Jobarg_enable_jobnet("Icon_1", "943"); err != nil {
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
		if err := lib.Jobarg_cleanup_linux(); err != nil {
			return t.logError(tc_2, "Error during cleanup: %s", err)
		}
		if err := lib.Jobarg_enable_jobnet("Icon_1", "943"); err != nil {
			return t.logError(tc_2, "Error during enabling jobnet: %s", err)
		}
		status := t.runJob2(tc_2, job)
		if status != PASSED {
			resetHostname(tc_2)
			return status
		}
		LogFound(tc_2)
		resetHostname(tc_2)
		return PASSED

	}

	tc_2.Set_function(tc_func2)
	t.Add_testcase(*tc_2)
}

func (t *Ticket_943) runJob(tc *common.TestCase, job string) common.Testcase_status {
	envs, err := lib.Get_str_str_map("JA_HOSTNAME", "oss.linux", "JA_CMD", "sleep 120")
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

	if jobnet_run_info.Jobnet_status == "END" && jobnet_run_info.Job_status == "NORMAL" {
		lib.Logi(common.LOG_LEVEL_INFO, "%s completed successfully.", job)
		return PASSED
	}

	return t.logError(tc, "%s failed. Jobnet_status: %s, Job_status: %s, Exit_cd: %d", job, jobnet_run_info.Jobnet_status, jobnet_run_info.Job_status, jobnet_run_info.Exit_cd)
}
func (t *Ticket_943) runJob2(tc *common.TestCase, job string) common.Testcase_status {
	lib.Clear_linux_jaz_server_log()
	envs, err := lib.Get_str_str_map("JA_HOSTNAME", "oss.linux", "JA_CMD", "sleep 120")
	if err != nil {
		return t.logError(tc, "Error retrieving environment variables: %s", err)
	}

	run_jobnet_id, err := lib.Jobarg_exec_E(job, envs)
	if err != nil {
		return t.logError(tc, "Error executing job %s: %s", job, err)
	}
	fmt.Printf("Executed job: %s with run_jobnet_id: %s\n", job, run_jobnet_id)

	status, err := modifyHostnameForTestcaseWhenInRunningState(tc, run_jobnet_id)

	if err != nil {
		lib.Logi(common.LOG_LEVEL_ERR, "Failed at modifyHostnameForTestcaseWhenInRunningState(tc, run_jobnet_id)")
		return FAILED
	}
	return status

}

const logFilePath = "/var/log/jobarranger/jobarg_server.log"
const logFileWarning = `In jarun_status_check()`

func (t *Ticket_943) LogNotFound(tc *common.TestCase) common.Testcase_status {
	cmd := fmt.Sprintf(`cat %s | grep "%s"`, logFilePath, logFileWarning)
	lib.Logi(common.LOG_LEVEL_INFO, "Executing command: %s", cmd)

	warningLogOutput, err := lib.Ssh_exec_to_str(cmd)
	fmt.Println(err)

	if strings.Contains(warningLogOutput, logFileWarning) {
		lib.Logi(common.LOG_LEVEL_ERR, "Warning log found, returning FAILED.")
		return FAILED
	}

	lib.Logi(common.LOG_LEVEL_INFO, "Warning log not found, returning PASSED.")
	return PASSED
}

func (t *Ticket_943) logError(tc *common.TestCase, format string, args ...interface{}) common.Testcase_status {
	fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, format, args...))
	return FAILED
}

func LogFound(tc_1 *common.TestCase) bool {
	const logFilePath = "/var/log/jobarranger/jobarg_server.log"
	const logFileWarning = `"jarun_status_check()"`
	const maxRetries = 10
	const retryInterval = 10 * time.Second

	for i := 0; i < maxRetries; i++ {
		cmd := fmt.Sprintf(`cat %s | grep %s`, logFilePath, logFileWarning)
		fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "Executing command: %s", cmd))

		warningLogOutput, err := lib.Ssh_exec_to_str(cmd)

		fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "checking for warning log: %s", err))

		if warningLogOutput != "" || warningLogOutput == " " {
			fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "Warning log found, returning Passed."))
			return true
		}

		fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "Warning log not found. Retrying in %v...", retryInterval))
		time.Sleep(retryInterval)
	}

	fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Warning log not found after retries, returning FAILED."))
	return false
}

func modifyHostnameForTestcase(tc *common.TestCase) common.Testcase_status {
	var changeHostCmd string

	// Check if it's MySQL
	if common.DB_type == common.MYSQL {
		changeHostCmd = `UPDATE hosts SET host = 'new.hostname' WHERE host = 'oss.linux';`
		_, err := common.DB.Exec(changeHostCmd)
		if err != nil {
			lib.Logi(common.LOG_LEVEL_ERR, "Failed to change hostname in MySQL database: %s", err)
			return FAILED
		}
		lib.Logi(common.LOG_LEVEL_INFO, "Hostname successfully changed in MySQL database.")

		// Check if it's PostgreSQL
	} else if common.DB_type == common.PSQL {
		changeHostCmd = `UPDATE hosts SET host = 'new.hostname' WHERE host = 'oss.linux';`
		_, err := common.DB.Exec(changeHostCmd)
		if err != nil {
			lib.Logi(common.LOG_LEVEL_ERR, "Failed to change hostname in PostgreSQL database: %s", err)
			return FAILED
		}
		lib.Logi(common.LOG_LEVEL_INFO, "Hostname successfully changed in PostgreSQL database.")
	} else {
		// Handle unsupported database types
		lib.Logi(common.LOG_LEVEL_ERR, "Unsupported database type.")
		return FAILED
	}

	return PASSED
}

func resetHostname(tc *common.TestCase) common.Testcase_status {
	var resetHostCmd string

	if common.DB_type == common.MYSQL {

		resetHostCmd = `UPDATE hosts SET host = 'oss.linux' WHERE host = 'new.hostname';`
		_, err := common.DB.Exec(resetHostCmd)
		if err != nil {
			lib.Logi(common.LOG_LEVEL_ERR, "Failed to reset hostname in MySQL database: %s", err.Error())
			return FAILED
		}
		lib.Logi(common.LOG_LEVEL_INFO, "Hostname successfully reset in MySQL database.")

	} else if common.DB_type == common.PSQL {
		resetHostCmd = `UPDATE hosts SET host = 'oss.linux' WHERE host = 'new.hostname';`
		_, err := common.DB.Exec(resetHostCmd)
		if err != nil {
			lib.Logi(common.LOG_LEVEL_ERR, "Failed to reset hostname in PostgreSQL database: %s", err.Error())
			return FAILED
		}
		lib.Logi(common.LOG_LEVEL_INFO, "Hostname successfully reset in PostgreSQL database.")
	} else {
		lib.Logi(common.LOG_LEVEL_ERR, "Unsupported database type for resetting hostname.")
		return FAILED
	}

	return PASSED
}

func modifyHostnameForTestcaseWhenInRunningState(tc *common.TestCase, registry_number string) (common.Testcase_status, error) {
	var jobnet_status, job_status string
	var err error
	var index int

	for {
		jobnet_status, err = lib.Jobarg_get_JA_JOBNETSTATUS(registry_number)
		if err != nil {
			return FAILED, err
		}

		job_status, err = lib.Jobarg_get_JA_JOBSTATUS(registry_number)
		if err != nil {
			return FAILED, err
		}

		if jobnet_status == common.RUN && job_status == common.NORMAL {
			modifyHostnameForTestcase(tc)
			if LogFound(tc) {
				resetHostname(tc)
				return PASSED, nil
			}
			break
		}
		lib.Spinner_log(index, lib.Formatted_log(common.LOG_LEVEL_INFO, "Getting jobnet[%s] run info but jobnet is not in running state. jobnet_status: %s, job_status: %s", registry_number, jobnet_status, job_status))
		time.Sleep(1 * time.Second)
		index++
	}

	fmt.Println()
	return FAILED, nil
}
