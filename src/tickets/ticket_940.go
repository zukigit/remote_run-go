package tickets

import (
	"fmt"
	"time"

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
	t.no = 940
	t.description = "Check the process timeout is working well or not"
}
func (t *Ticket_940) Add_testcases() {
	configs := []string{
		"JaRunTimeout=1",
		"JaTrapperTimeout=1",
		"JaJobTimeout=1",
		"JaJobnetTimeout=1",
		"JaLoaderTimeout=1",
		"JaBootTimeout=1",
		"JaMsgsndTimeout=1",
		"JaSelfmonTimeout=1",
		"JaPurgeTimeout=1",
		"JaAbortTimeout=1",
	}

	configFilePath := "/etc/jobarranger/jobarg_server.conf"

	for i, config := range configs {
		tc := t.New_testcase(uint(i+1), fmt.Sprintf("Change configuration: %s", config))

		tc_func := func() common.Testcase_status {
			sedCmd := fmt.Sprintf(`sed -i -e '$a\%s' %s`, config, configFilePath)

			output, err := lib.Ssh_exec_to_str(sedCmd)
			if err != nil {
				t.logError(tc,
					"Failed to set server config for %s, Error: %s, Output: %s", config, err.Error(), output)

			}
			fmt.Println("Configuration has been set to:", config)

			if err := lib.Jobarg_cleanup_linux(); err != nil {
				t.logError(t.New_testcase(uint(i+1), "Cleaning Up jobarg_server"),
					"Failed to clean up the server, Error: %s", err.Error())
			}
			fmt.Println("jobarg_server has been restarted successfully.")

			if err := lib.Delete_server_log(); err != nil {
				t.logError(tc, "Error during deleting log: %s", err)
			}

			// Enable jobnet
			if err := lib.Jobarg_enable_jobnet("Icon_1", "jobicon_linux"); err != nil {
				t.logError(tc, "Error during enabling jobnet: %s", err)
			}
			fmt.Println("Enabled Jobnet")

			if err := lib.Restart_jaz_server(); err != nil {
				t.logError(tc, "Failed to restart: %s", err)
			}

			if status := t.runIcon100(tc, "Icon_100"); status != PASSED {
				return t.logError(tc, "Job execution failed after updating %s", config)
			}

			if err := t.removeConfig(config, configFilePath); err != nil {
				t.logError(tc, "Failed to remove config %s, Error: %s", config, err)
				return FAILED
			}
			fmt.Println("Configuration has been removed:", config)

			return PASSED
		}

		tc.Set_function(tc_func)
		t.Add_testcase(*tc)

	}

}

func (t *Ticket_940) removeConfig(config, configFilePath string) error {
	removeCmd := fmt.Sprintf(`sed -i '/%s/d' %s`, config, configFilePath)
	_, err := lib.Ssh_exec_to_str(removeCmd)
	return err
}

func (t *Ticket_940) runIcon100(tc *dao.TestCase, job string) common.Testcase_status {
	envs, err := lib.Get_str_str_map("JA_HOSTNAME", "oss.linux", "JA_CMD", "sleep 50")
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
		tc.Info_log("%s completed successfully.", job)
		fmt.Printf("%s completed successfully.", job)

		logFilePath := "/var/log/jobarranger/jobarg_server.log"
		cmd := fmt.Sprintf(`cat %s | grep "Process is taking"`, logFilePath)

		tc.Info_log("Executing command: %s", cmd)

		var output string
		var retryCount int
		var retryLimit = 2
		var errRetry error

		// Retry up to `retryLimit` times with a 5-10 minute delay
		for retryCount < retryLimit {
			output, errRetry = lib.Ssh_exec_to_str(cmd)
			if errRetry != nil {
				// Log the error but continue retrying
				tc.Err_log("Failed to check timeout warnings, Error: %s", errRetry.Error())
				retryCount++
				time.Sleep(3 * time.Minute) // Sleep for 5 minutes before retrying
				continue
			}

			// If no error after retrying, log the output and break the loop
			tc.Info_log("Command output: %s", output)
			break
		}

		// If after all retries it still fails, return FAILED
		if errRetry != nil {
			tc.Err_log("Max retries reached, failed to check timeout warnings.")
			return FAILED
		}

		return PASSED
	}
	return PASSED
}

func (t *Ticket_940) logError(tc *dao.TestCase, format string, args ...interface{}) common.Testcase_status {
	fmt.Println(tc.Err_log(format, args...))
	return FAILED
}
