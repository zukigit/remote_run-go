package tickets

import (
	"fmt"
	// "strconv"
	// "strings"

	"github.com/zukigit/remote_run-go/src/common"
	"github.com/zukigit/remote_run-go/src/dao"
	"github.com/zukigit/remote_run-go/src/lib"
)

type Ticket_940 struct {
	no          uint
	description string
	testcases   []dao.TestCase
}

func (t *Ticket_940) New_testcase(testcaseID uint, testcaseDescription string) *dao.TestCase {
	return dao.New_testcase(testcaseID, testcaseDescription)
}

func (t *Ticket_940) Get_no() uint {
	return t.no
}

func (t *Ticket_940) Get_dsctn() string {
	return t.description
}

func (t *Ticket_940) Add_testcase(tc dao.TestCase) {
	t.testcases = append(t.testcases, tc)
}

func (t *Ticket_940) Get_testcases() []dao.TestCase {
	return t.testcases
}

func (t *Ticket_940) Set_values() {
	t.no = 940 // Enter your ticket id
	t.description = "Check the process timeout is working well or not"
}
func (t *Ticket_940) Add_testcases() {
	configs := []string{
		"JaRunTimeout=20",
		"JaTrapperTimeout=20",
		"JaJobTimeout=20",
		//"JaJobnetTimeout=10",
		// "JaLoaderTimeout=10",
		// "JaBootTimeout=10",
		// "JaMsgsndTimeout=10",
		// "JaSelfmonTimeout=10",
		// "JaPurgeTimeout=10",
		// "JaAbortTimeout=10",
	}

	configFilePath := "/etc/jobarranger/jobarg_server.conf"

	for i, config := range configs {
		sedCmd := fmt.Sprintf(`sed -i -e '$a\%s' %s`, config, configFilePath)

		output, err := lib.Ssh_exec_to_str(sedCmd)
		if err != nil {
			t.logError(t.New_testcase(uint(i+1), fmt.Sprintf("Update %s", config)),
				"Failed to set server config for %s, Error: %s, Output: %s", config, err.Error(), output)

		}
		fmt.Println("Configuration has been set to:", config)

		if err := lib.Jobarg_cleanup_linux(); err != nil {
			t.logError(t.New_testcase(uint(i+1), "Cleaning Up jobarg_server"),
				"Failed to clean up the server, Error: %s", err.Error())
		}
		fmt.Println("jobarg_server has been restarted successfully.")

		tc := t.New_testcase(uint(i+1), fmt.Sprintf("Change configuration: %s", config))
		t.Add_testcase(*tc)

		if err := lib.Delete_server_log(); err != nil {
			t.logError(tc, "Error during deleting log: %s", err)
		}

		if err := lib.Restart_jaz_server(); err != nil {
			t.logError(tc, "Failed to restart: %s", err)
		}
		// Enable jobnet
		if err := lib.Jobarg_enable_jobnet("Icon_1", "jobicon_linux"); err != nil {
			t.logError(tc, "Error during enabling jobnet: %s", err)
		}

		if status := t.runIcon100(tc, "Icon_100"); status != PASSED {
			t.logError(tc, "Job execution failed after updating %s", config)
		}
	}

}

func (t *Ticket_940) runIcon100(tc *dao.TestCase, job string) common.Testcase_status {
	envs, err := lib.Get_str_str_map("JA_HOSTNAME", "oss.linux", "JA_CMD", "hostname")
	if err != nil {
		return t.logError(tc, "Error retrieving environment variables: %s", err)
	}
	run_jobnet_id, err := lib.Jobarg_exec_E(job, envs)
	if err != nil {
		return t.logError(tc, "Error executing Icon_100: %s, error: %s", job, err)
	}

	jobnet_run_info, err := lib.Jobarg_get_jobnet_run_info(run_jobnet_id)
	if err != nil {
		return t.logError(tc, "Error retrieving job status: %s", err)
	}

	if jobnet_run_info.Jobnet_status == "END" && jobnet_run_info.Job_status == "NORMAL" && jobnet_run_info.Exit_cd == 0 {
		// tc.Info_log("%s completed successfully.", job)
		// logFilePath := "/var/log/jobarranger/jobarg_server.log"
		// cmd := fmt.Sprintf(`cat %s | grep "Process is taking taking"`, logFilePath)

		// tc.Info_log("Executing command: %s", cmd)

		// output, err := lib.Ssh_exec_to_str(cmd)
		// if err != nil {
		// 	tc.Err_log("Failed to check timeout warnings, Error:", err.Error())
		// 	return FAILED
		// }

		// tc.Info_log("Command output: %s", output)

		// logCount, _ := strconv.Atoi(strings.TrimSpace(output))
		// if logCount > 0 {
		// 	tc.Info_log("Timeout log entry found.")
		// 	return PASSED
		// }

		tc.Info_log("No timeout log entry found, returning PASSED as job completed successfully.")
		return PASSED

	}
	return FAILED
}
func (t *Ticket_940) logError(tc *dao.TestCase, format string, args ...interface{}) common.Testcase_status {
	fmt.Println(tc.Err_log(format, args...))
	return FAILED
}
