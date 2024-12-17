package tickets

import (
	"fmt"

	"github.com/zukigit/remote_run-go/src/common"
	"github.com/zukigit/remote_run-go/src/lib"
)

type Ticket_940 struct {
	Ticket_no                                   uint
	Ticket_description                          string
	PASSED_count, FAILED_count, MUSTCHECK_count int
	Testcases                                   []common.TestCase
}

func (t *Ticket_940) New_testcase(testcaseID uint, testcaseDescription string) *common.TestCase {
	return common.New_testcase(testcaseID, testcaseDescription)
}

func (t *Ticket_940) Get_no() uint {
	return t.Ticket_no
}

func (t *Ticket_940) Set_PASSED_count(passed_count int) {
	t.PASSED_count = passed_count
}

func (t *Ticket_940) Set_FAILED_count(failed_count int) {
	t.FAILED_count = failed_count
}

func (t *Ticket_940) Set_MUSTCHECK_count(mustcheck_count int) {
	t.MUSTCHECK_count = mustcheck_count
}

func (t *Ticket_940) Get_dsctn() string {
	return t.Ticket_description
}

func (t *Ticket_940) Add_testcase(tc common.TestCase) {
	t.Testcases = append(t.Testcases, tc)
}

func (t *Ticket_940) Get_testcases() []common.TestCase {
	return t.Testcases
}

func (t *Ticket_940) Set_values() {
	t.Ticket_no = 940
	t.Ticket_description = "Check the process timeout is working well or not"
}

func (t *Ticket_940) Add_testcases() {
	// All configurations to be added at once
	configs := []string{
		"JaRunTimeout=5",
		"JaTrapperTimeout=5",
		"JaJobTimeout=5",
		"JaJobnetTimeout=5",
		"JaLoaderTimeout=5",
		"JaBootTimeout=5",
		"JaMsgsndTimeout=5",
		"JaSelfmonTimeout=5",
		"JaPurgeTimeout=5",
		"JaAbortTimeout=5",
	}

	configFilePath := "/etc/jobarranger/jobarg_server.conf"

	// Create a test case for applying all configurations
	tc := t.New_testcase(1, "Apply all configurations and check logs")

	tc_func := func() common.Testcase_status {
		if status := t.applyConfigAndRunTests(tc, configs, configFilePath); status != PASSED {
			return status
		}
		return PASSED
	}

	tc.Set_function(tc_func)
	t.Add_testcase(*tc)
}

// Consolidated function for applying all configurations at once, running Icon_100, and checking logs
func (t *Ticket_940) applyConfigAndRunTests(tc *common.TestCase, configs []string, configFilePath string) common.Testcase_status {
	// Apply all configurations at once
	for _, config := range configs {
		sedCmd := fmt.Sprintf(`sed -i -e '$a\%s' %s`, config, configFilePath)
		output, err := lib.Ssh_exec_to_str(sedCmd)
		if err != nil {
			return t.logError(tc, "Failed to set server config for %s, Error: %s, Output: %s", config, err.Error(), output)
		}
		fmt.Println("Configuration has been set to:", config)
	}

	// Delete logs
	if err := lib.Clear_linux_jaz_server_log(); err != nil {
		return t.logError(tc, "Error during deleting log: %s", err)
	}

	// Cleanup server
	if err := lib.Jobarg_cleanup_linux(); err != nil {
		return t.logError(tc, "Failed to clean up the server, Error: %s", err.Error())
	}
	fmt.Println("jobarg_server has been restarted successfully.")

	// Enable jobnet
	if err := lib.Jobarg_enable_jobnet("Icon_1", "jobicon_linux"); err != nil {
		return t.logError(tc, "Error during enabling jobnet: %s", err)
	}

	// Run Icon_100 after applying configurations
	if status := t.runIcon100(tc, "Icon_100"); status != PASSED {
		return status
	}

	if err := lib.Restart_jaz_server(); err != nil {
		return t.logError(tc, "Failed to restart: %s", err)
	}

	lib.Run_Timeout(10)

	logFilePath := "/var/log/jobarranger/jobarg_server.log"
	cmd := fmt.Sprintf(`cat %s | grep "Process is taking"`, logFilePath)

	lib.Logi(common.LOG_LEVEL_INFO, "Executing command: %s", cmd)

	output, err := lib.Ssh_exec_to_str(cmd)
	if err != nil {
		return t.logError(tc, "Failed to check timeout warnings, Error: %s", err.Error())
	}

	lib.Logi(common.LOG_LEVEL_INFO, "Command output: %s", output)

	logCount := countOccurrences(output, "Process is taking")

	if logCount > 0 {
		for _, config := range configs {
			if err := t.removeConfig(config, configFilePath); err != nil {
				return t.logError(tc, "Failed to remove config %s, Error: %s", config, err)
			}
			fmt.Println("Configuration has been removed:", config)
		}
		return PASSED
	}

	return FAILED
}

func countOccurrences(str, substr string) int {
	count := 0
	for i := 0; i+len(substr) <= len(str); i++ {
		if str[i:i+len(substr)] == substr {
			count++
		}
	}
	return count
}

func (t *Ticket_940) logError(tc *common.TestCase, format string, args ...interface{}) common.Testcase_status {
	fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, format, args...))
	return FAILED
}

func (t *Ticket_940) removeConfig(config, configFilePath string) error {
	removeCmd := fmt.Sprintf(`sed -i '/%s/d' %s`, config, configFilePath)
	_, err := lib.Ssh_exec_to_str(removeCmd)
	return err
}

func (t *Ticket_940) runIcon100(tc *common.TestCase, job string) common.Testcase_status {
	envs, err := lib.Get_str_str_map("JA_HOSTNAME", "oss.linux", "JA_CMD", "sleep 30")
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

	fmt.Printf("jobnet status: %s, job status: %s\n", jobnet_run_info.Jobnet_status, jobnet_run_info.Job_status)

	if jobnet_run_info.Jobnet_status == "END" && jobnet_run_info.Job_status == "NORMAL" {
		lib.Logi(common.LOG_LEVEL_INFO, "%s completed successfully.", job)
		fmt.Printf("%s completed successfully.", job)
		return PASSED
	}
	return PASSED
}
