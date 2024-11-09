package tickets

import (
	"fmt"

	"github.com/zukigit/remote_run-go/src/common"
	"github.com/zukigit/remote_run-go/src/dao"
	"github.com/zukigit/remote_run-go/src/lib"
	"golang.org/x/crypto/ssh"
)

type Ticket_699 struct {
	no          uint
	description string
	testcases   []dao.TestCase
}

func (t *Ticket_699) New_testcase(testcase_id uint, testcase_description string) *dao.TestCase {
	return dao.New_testcase(testcase_id, testcase_description)
}

func (t *Ticket_699) Get_no() uint {
	return t.no
}

func (t *Ticket_699) Get_dsctn() string {
	return t.description
}

func (t *Ticket_699) Add_testcase(tc dao.TestCase) {
	t.testcases = append(t.testcases, tc)
}

func (t *Ticket_699) Get_testcases() []dao.TestCase {
	return t.testcases
}

// Enter your ticket information here
func (t *Ticket_699) Set_values() {
	t.no = 699 // Enter your ticket id
	t.description = "Reconnect and retry case for connection agent(Job Icon, File Transfer, Reboot)"
}

func (t *Ticket_699) Add_testcases() {
	tc_1 := t.New_testcase(1, "Check Server log output")
	tc_func := func() common.Testcase_status {
		if err := lib.Jobarg_enable_jobnet("Icon_1", "jobicon_linux"); err != nil {
			tc_1.Err_log("Failed to enable jobnet, Error: %s", err)
			return FAILED
		}
		return CheckLog("Icon_1", tc_1, common.Client)
	}
	tc_1.Set_function(tc_func)
	t.Add_testcase(*tc_1)
}

func CheckLog(jobnetId string, testcase *dao.TestCase, sshClient *ssh.Client) common.Testcase_status {

	err := lib.Clean_jaz_server_log()
	if err != nil {
		fmt.Println(testcase.Err_log("Error: %s, Failed to cleanup the job arranger server log.", err))
		return FAILED
	}
	fmt.Println(testcase.Info_log("Cleanup the server log success."))

	err = lib.Restart_jaz_server()
	if err != nil {
		fmt.Println(testcase.Err_log("Error: %s, Failed to restart the jobarranger server.", err))
		return FAILED
	}
	fmt.Println(testcase.Info_log("Server restart success."))

	searchString := `'\[JASERVER000004\] server #13 started \[purge old jobnet #1\]'`
	logs, err := lib.Jobarg_server_check_log(searchString)
	if err != nil {
		fmt.Println("Error:", err)
		return FAILED
	}
	fmt.Println(testcase.Info_log("Log output result for purge old jobnet:"))
	for _, log := range logs {
		fmt.Println(log)
	}

	searchString = "[JASERVER000004] server #10 started [jobnet boot #1]"
	logs, err = lib.Jobarg_server_check_log(searchString)
	if err != nil {
		fmt.Println("Error:", err)
		return FAILED
	}

	fmt.Println(testcase.Info_log("Log output result for jobnet boot:"))
	for _, log := range logs {
		fmt.Println(log)
	}
	return PASSED
}
