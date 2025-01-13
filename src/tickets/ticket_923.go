package tickets

import (
	"fmt"
	"strings"
	"time"

	"github.com/zukigit/remote_run-go/src/common"
	"github.com/zukigit/remote_run-go/src/lib"
)

type Ticket_923 struct {
	Ticket_no                                   uint
	Ticket_description                          string
	PASSED_count, FAILED_count, MUSTCHECK_count int
	Testcases                                   []common.TestCase
}

func (t *Ticket_923) Set_PASSED_count(passed_count int) {
	t.PASSED_count = passed_count
}

func (t *Ticket_923) Set_FAILED_count(failed_count int) {
	t.FAILED_count = failed_count
}

func (t *Ticket_923) Set_MUSTCHECK_count(mustcheck_count int) {
	t.MUSTCHECK_count = mustcheck_count
}

func (t *Ticket_923) New_testcase(testcase_id uint, testcase_description string) *common.TestCase {
	return common.New_testcase(testcase_id, testcase_description)
}

func (t *Ticket_923) Get_no() uint {
	return t.Ticket_no
}

func (t *Ticket_923) Get_dsctn() string {
	return t.Ticket_description
}

func (t *Ticket_923) Add_testcase(tc common.TestCase) {
	t.Testcases = append(t.Testcases, tc)
}

func (t *Ticket_923) Get_testcases() []common.TestCase {
	return t.Testcases
}

func (t *Ticket_923) Set_values() {
	t.Ticket_no = 923
	t.Ticket_description = "Check whether agentd can run with different listener process with parallel"
}

func (t *Ticket_923) Add_testcases() {
	t.addTestCase87()
}

func (t *Ticket_923) addTestCase87() {
	tc := t.New_testcase(87, "Abnormal job execution")
	tc_func := func() common.Testcase_status {
		if err := t.cleanupJobArg(); err != nil {
			lib.Logi(common.LOG_LEVEL_ERR, "Failed to cleanup jobarg, Error: %s", err)
			return FAILED
		}
		fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "Jobarg cleanup successfully."))

		logfilePath := "/var/log/jobarranger/jobarg_agentd.log"
		if err := t.clearLogFile(logfilePath); err != nil {
			lib.Logi(common.LOG_LEVEL_ERR, "Failed to clear log file, Error: %s", err)
			return FAILED
		}
		fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "Agent log file has been cleared."))

		if err := t.enableJobnet("Icon_1", "jobicon_linux"); err != nil {
			lib.Logi(common.LOG_LEVEL_ERR, "Failed to enable jobnet, Error: %s", err)
			return FAILED
		}
		fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "Jobnet 'Icon_1' enabled successfully."))

		if err := t.enableJobnet("Icon_10", "Icon_10"); err != nil {
			lib.Logi(common.LOG_LEVEL_ERR, "Failed to enable jobnet, Error: %s", err)
			return FAILED
		}
		fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "Jobnet 'Icon_10' enabled successfully."))

		run_jobnet_id, err := t.runJobnet("Icon_10", "oss.linux", "sleep 60")
		if err != nil {
			lib.Logi(common.LOG_LEVEL_ERR, "Error: %s, std_out: %s", err.Error(), run_jobnet_id)
			return FAILED
		}

		output, err := getProcessIdFromAgentLog(logfilePath)
		if err != nil {
			lib.Logi(common.LOG_LEVEL_ERR, "Failed to get process id from agent log")
			return FAILED
		}

		processIds := strings.Split(output, ",")
		processId1, processId2 := processIds[0], processIds[1]

		if processId1 != processId2 {
			fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "Info: Process IDs are different."))
			return PASSED
		}

		return FAILED
	}
	tc.Set_function(tc_func)
	t.Add_testcase(*tc)
}

func (t *Ticket_923) cleanupJobArg() error {
	return lib.Jobarg_cleanup_linux()
}

func (t *Ticket_923) clearLogFile(logfilePath string) error {
	_, err := lib.Ssh_exec("> " + logfilePath)
	if err != nil {
		return fmt.Errorf("failed to clear log file: %v", err)
	}
	fmt.Println("Agent log file has been cleared.")
	return nil
}

func (t *Ticket_923) enableJobnet(jobnet_id, jobnet_name string) error {
	if err := lib.Jobarg_enable_jobnet(jobnet_id, jobnet_name); err != nil {
		return fmt.Errorf("failed to enable jobnet: %v", err)
	}
	fmt.Printf("Jobnet '%s' enabled successfully.\n", jobnet_name)
	return nil
}

func (t *Ticket_923) runJobnet(jobnet_id, hostname, command string) (string, error) {
	envs, _ := lib.Get_str_str_map("JA_HOSTNAME", hostname, "JA_CMD", command)
	run_jobnet_id, err := lib.Jobarg_exec_E(jobnet_id, envs)
	if err != nil {
		return "", fmt.Errorf("failed to run jobnet: %v", err)
	}
	return run_jobnet_id, nil
}

func getProcessIdFromAgentLog(logFilePath string) (string, error) {
	// Wait for 5 seconds before attempting to retrieve the log data
	time.Sleep(10 * time.Second)

	// Command to extract the process ID from the log
	cmd := fmt.Sprintf("cat %s | grep 'In ja_agent_begin()' | head -n 2 | awk -F: '{print $1}'", logFilePath)

	// Execute the command
	fmt.Println(cmd)
	output, err := lib.Ssh_exec_to_str(cmd)
	if err != nil {
		return "", err
	}

	// Trim whitespace or newline characters from the output
	output = strings.Replace(strings.TrimSpace(output), "\n", ",", -1)

	// Check if output is empty after trimming
	if output == "" {
		return "", fmt.Errorf("no process ID found in log file")
	}

	// Return the trimmed output
	return output, nil
}
