package tickets

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/zukigit/remote_run-go/src/common"
	"github.com/zukigit/remote_run-go/src/dao"
	"github.com/zukigit/remote_run-go/src/lib"
	"golang.org/x/crypto/ssh"
)

type Ticket_1281 struct {
	no          uint
	description string
	testcases   []dao.TestCase
}

func (t *Ticket_1281) New_testcase(testcase_id uint, testcase_description string) *dao.TestCase {
	return dao.New_testcase(testcase_id, testcase_description)
}

func (t *Ticket_1281) Get_no() uint {
	return t.no
}

func (t *Ticket_1281) Get_dsctn() string {
	return t.description
}

func (t *Ticket_1281) Add_testcase(tc dao.TestCase) {
	t.testcases = append(t.testcases, tc)
}

func (t *Ticket_1281) Get_testcases() []dao.TestCase {
	return t.testcases
}

// Enter your ticket information here
func (t *Ticket_1281) Set_values() {
	t.no = 1281 // Enter your ticket id
	t.description = "Check job process start"
}

// Add your test case here
func (t *Ticket_1281) Add_testcases() {
	tc_1 := t.New_testcase(140, "Check if Check job process start or not.")
	tc_func := func() common.Testcase_status {
		return CheckJobProcessStartOrExit(tc_1, common.Client)
	}
	tc_1.Set_function(tc_func)
	t.Add_testcase(*tc_1)
}

func CheckJobProcessStartOrExit(testcase *dao.TestCase, client *ssh.Client) common.Testcase_status {

	var ServerProcessID string
	// var ProcessID string

	_clear_err := lib.ClearLogFile(client, "/var/log/jobarranger/jobarg_server.log")
	if _clear_err != nil {
		fmt.Println(testcase.Err_log("Error clear log: %s", _clear_err))
	}

	_restart_err := lib.Restart_jaz_server()
	if _restart_err != nil {
		fmt.Println(testcase.Err_log("Error jaz restart: %s", _restart_err))
	}
	time.Sleep(10 * time.Second)
	pattern := "[INFO] [JASERVER000004] server #15 started [check job process #1]"

	_check_log, _check_err := lib.WaitForPatternInLogFile(client, "/var/log/jobarranger/jobarg_server.log", pattern, 30, 2)
	if _check_err != nil {
		fmt.Println(testcase.Err_log("Error check log file : %s", _check_err))
	}

	if _check_log != "" {
		// fmt.Println(_check_log)
	}

	// Regular expression to match "ja_purge process execution thread : <number>"
	re := regexp.MustCompile(`ja_check_job process execution thread\s*:\s*(\d+)`) // Match "ja_purge process execution thread : <number>"

	// Use the regex to find the dynamic process ID (e.g., 250622)
	matches := re.FindStringSubmatch(_check_log)
	if len(matches) > 1 {
		ServerProcessID = matches[1] // Extract the dynamic value (e.g., 250622)
		fmt.Printf("Found process ID: %s\n", ServerProcessID)
	} else {
		fmt.Println("No matching process found in the log.")
	}

	_ps_job, _ps_job_err := lib.GetOutputStrFromSSHCommand(client, "ps -aux | grep job")
	if _ps_job_err != nil {
		fmt.Println(testcase.Err_log("Error ps aux job : %s", _ps_job_err))
	}

	if _ps_job != "" {
		// Split the output by lines
		lines := strings.Split(_ps_job, "\n")
		for _, line := range lines {
			if strings.Contains(line, ServerProcessID) {
				// If the PID exists in the line, it means the process is running
				fmt.Printf("Process with PID %s is running: %s\n", ServerProcessID, line)
				return PASSED
			}
		}
		// If no matching process is found
		fmt.Printf("Process with PID %s not found in ps output.\n", ServerProcessID)
	}

	return FAILED
}

func CheckIconCountForSendingStatus(jobnetId string, processCount int, processCheckTimeout int, testcase *dao.TestCase, client *ssh.Client) common.Testcase_status {

	return FAILED
}
