package tickets

import (
	"fmt"
	"strings"
	"time"

	"github.com/zukigit/remote_run-go/src/common"
	"github.com/zukigit/remote_run-go/src/dao"
	"github.com/zukigit/remote_run-go/src/lib"
)

const (
	timeoutDuration = 60 * time.Second
	pollInterval    = 10 * time.Second
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

// GetDescription returns the ticket description.
func (t *Ticket_940) Get_dsctn() string {
	return t.description
}

// AddTestCase adds a test case to the ticket.
func (t *Ticket_940) Add_testcase(tc dao.TestCase) {
	t.testcases = append(t.testcases, tc)
}

// GetTestCases returns all test cases associated with the ticket.
func (t *Ticket_940) Get_testcases() []dao.TestCase {
	return t.testcases
}

func (t *Ticket_940) Set_values() {
	t.no = 940
	t.description = "Check the process timeout is working well or not"
}

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
			return
		}
		configKey := configParts[0]
		configValue := configParts[1]

		tc := t.New_testcase(uint(i+1), fmt.Sprintf("Change configuration: %s and run Icon_100", config))

		// change the configuration
		if err := lib.Ja_set_server_config_linux(configKey, configValue); err != nil {
			t.logError(tc, "failed to update configuration %s to %s: %s", configKey, configValue, err)
			return
		}
		tc.Info_log("Configuration %s updated to %s successfully.", configKey, configValue)

		// Clean up and restart the server
		if err := lib.Jobarg_cleanup_linux(); err != nil {
			t.logError(tc, "failed to clean up the server: %s", err)
			return
		}
		tc.Info_log("Server cleaned up successfully.")

		// Enable the jobnet
		if err := lib.Jobarg_enable_jobnet("Icon_1", "jobicon_linux"); err != nil {
			t.logError(tc, "failed to enable jobnet: %s", err)
			return
		}

		// Run only Icon_100 and check its status
		if status := t.runIcon100(tc); status != PASSED {
			return
		}

		t.Add_testcase(*tc)
	}
}

// runIcon100 executes the Icon_100 job and monitors its completion status.
func (t *Ticket_940) runIcon100(tc *dao.TestCase) common.Testcase_status {
	envs, err := lib.Get_str_str_map("JA_HOSTNAME", "oss.linux", "JA_CMD", "sleep 10")
	if err != nil {
		return t.logError(tc, "error retrieving environment variables: %s", err)
	}

	runJobnetID, err := lib.Jobarg_exec_E("Icon_100", envs)
	if err != nil {
		return t.logError(tc, "error executing Icon_100: %s", err)
	}

	// Poll for job completion with a timeout
	timeout := time.After(timeoutDuration)
	tick := time.Tick(pollInterval)

	for {
		select {
		case <-timeout:
			return t.logError(tc, "job Icon_100 did not complete within the expected time")
		case <-tick:
			jobnetRunInfo, err := lib.Jobarg_get_jobnet_run_info(runJobnetID)
			if err != nil {
				return t.logError(tc, "error retrieving run info for Icon_100: %s", err)
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

// logError logs an error message and handles the failure without returning a status.
func (t *Ticket_940) logError(tc *dao.TestCase, format string, args ...interface{}) common.Testcase_status {
	fmt.Println(tc.Err_log(format, args...))
	return FAILED
}
