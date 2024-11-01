package tickets

import (
	"fmt"
	"strings"
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

func (t *Ticket_940) New_testcase(testcase_id uint, testcase_description string) *dao.TestCase {
	return dao.New_testcase(testcase_id, testcase_description)
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

// Add_testcases sets up a test case that configures and runs only the Icon_100 job
func (t *Ticket_940) Add_testcases() {
	configs := []string{
		"JaRunTimeout=5s",
		"JaTrapperTimeout=5s",
		"JaJobTimeout=5s",
		"JaJobnetTimeout=5s",
		"JaLoaderTimeout=5s",
		"JaBootTimeout=5s",
		"JaMsgsndTimeout=5s",
		"JaSelfmonTimeout=5s",
		"JaPurgeTimeout=5s",
		"JaAbortTimeout=5s",
	}

	for i, config := range configs {
		configParts := strings.Split(config, "=")
		if len(configParts) != 2 {
			fmt.Printf("Invalid configuration format: %s\n", config)
			continue
		}
		configKey := configParts[0]
		configValue := configParts[1]

		tc_88 := t.New_testcase(uint(i+1), fmt.Sprintf("Change configuration: %s and run Icon_100", config))

		tcFunc := func() common.Testcase_status {
			// Enable the jobnet
			if err := lib.Jobarg_enable_jobnet("Icon_1", "jobicon_linux"); err != nil {
				return t.logError(tc_88, "Failed to enable jobnet, Error: %s", err)
			}

			// Clean up and restart the server
			if err := lib.Jobarg_cleanup_linux(); err != nil {
				return t.logError(tc_88, "Failed to clean up the server, Error: %s", err)
			}
			tc_88.Info_log("Server cleaned up successfully.")

			// Update the configuration
			if err := lib.Ja_set_server_config_linux(configKey, configValue); err != nil {
				return t.logError(tc_88, "Failed to update configuration %s to %s, Error: %s", configKey, configValue, err)
			}
			tc_88.Info_log("Configuration %s updated to %s successfully.", configKey, configValue)

			// Run only Icon_100 and check its status
			return t.runIcon100(tc_88)
		}
		tc_88.Set_function(tcFunc)
		t.Add_testcase(*tc_88)
	}
}

// runIcon100 executes the Icon_100 job and monitors its completion status
func (t *Ticket_940) runIcon100(tc *dao.TestCase) common.Testcase_status {
	envs, _ := lib.Get_str_str_map("JA_HOSTNAME", "oss.linux", "JA_CMD", "sleep 10")
	runJobnetID, err := lib.Jobarg_exec_E("Icon_100", envs)
	if err != nil {
		return t.logError(tc, "Error executing Icon_100: %s", err)
	}

	// Poll for job completion with a timeout
	timeout := time.After(60 * time.Second)
	tick := time.Tick(10 * time.Second)

	for {
		select {
		case <-timeout:
			return t.logError(tc, "Job Icon_100 did not complete within the expected time")
		case <-tick:
			jobnetRunInfo, err := lib.Jobarg_get_jobnet_run_info(runJobnetID)
			if err != nil {
				return t.logError(tc, "Error retrieving run info for Icon_100: %s", err)
			}

			// Check if the job has completed successfully
			if jobnetRunInfo.Jobnet_status == "END" && jobnetRunInfo.Job_status == "NORMAL" && jobnetRunInfo.Exit_cd == 0 {
				tc.Info_log("Icon_100 completed successfully.")
				return PASSED
			}

			// Log the current job status and continue polling
			tc.Info_log("Icon_100 is still running... Current status: %s, Job_status: %s, Exit_cd: %d",
				jobnetRunInfo.Jobnet_status, jobnetRunInfo.Job_status, jobnetRunInfo.Exit_cd)
		}
	}
}

// logError logs an error message and returns a FAILED status
func (t *Ticket_940) logError(tc *dao.TestCase, format string, args ...interface{}) common.Testcase_status {
	fmt.Println(tc.Err_log(format, args...))
	return FAILED
}
