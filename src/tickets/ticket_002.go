package tickets

import (
	"github.com/zukigit/remote_run-go/src/common"
	"github.com/zukigit/remote_run-go/src/lib"
)

type Ticket_002 struct {
	Ticket_no                                   uint
	Ticket_description                          string
	PASSED_count, FAILED_count, MUSTCHECK_count int
	Testcases                                   []common.TestCase
}

func (t *Ticket_002) New_testcase(testcase_id uint, testcase_description string) *common.TestCase {
	return common.New_testcase(testcase_id, testcase_description)
}

func (t *Ticket_002) Get_no() uint {
	return t.Ticket_no
}

func (t *Ticket_002) Set_PASSED_count(passed_count int) {
	t.PASSED_count = passed_count
}

func (t *Ticket_002) Set_FAILED_count(failed_count int) {
	t.FAILED_count = failed_count
}

func (t *Ticket_002) Set_MUSTCHECK_count(mustcheck_count int) {
	t.MUSTCHECK_count = mustcheck_count
}

func (t *Ticket_002) Get_dsctn() string {
	return t.Ticket_description
}

func (t *Ticket_002) Add_testcase(tc common.TestCase) {
	t.Testcases = append(t.Testcases, tc)
}

func (t *Ticket_002) Get_testcases() []common.TestCase {
	return t.Testcases
}

// Enter your ticket information here
func (t *Ticket_002) Set_values() {
	t.Ticket_no = 2 // Enter your ticket id
	t.Ticket_description = "Enter your ticket description here."
}

// Add your test case here
func (t *Ticket_002) Add_testcases() {
	// TESTCASE 001
	tc_1 := t.New_testcase(1, "Enter your test case description here.")
	tc_func := func() common.Testcase_status {

		// How to get host by type
		host, err := lib.Get_host(common.Hosts, common.WA_HOST_TYPE)
		if err != nil {
			lib.Logi(common.LOG_LEVEL_ERR, "could not get host, err: %s", err.Error())
			return FAILED
		} else {
			lib.Logi(common.LOG_LEVEL_INFO, "host_name: %s", host.Get_Host_name())
		}

		// How to run commands on host
		output, err := host.Run_cmd_str("hostname")
		if err != nil {
			lib.Logi(common.LOG_LEVEL_ERR, "could not run command on host, err: %s", err.Error())
			return FAILED
		}
		lib.Logi(common.LOG_LEVEL_INFO, "output: %s", output)

		// How to get hosts by types
		hosts, err := lib.Get_hosts(common.Hosts, common.LS_HOST_TYPE, common.WA_HOST_TYPE, common.LA_HOST_TYPE)
		if err != nil {
			lib.Logi(common.LOG_LEVEL_ERR, "could not get hosts, err: %s", err.Error())
			return FAILED
		} else {
			lib.Logi(common.LOG_LEVEL_INFO, "total hosts: %d", len(hosts))
		}

		for _, host = range hosts {
			lib.Logi(common.LOG_LEVEL_INFO, "hostname: %s", host.Get_Host_name())
			lib.Logi(common.LOG_LEVEL_INFO, "host_type: %s", host.Get_Host_type())
		}

		return PASSED
	}
	tc_1.Set_function(tc_func)
	t.Add_testcase(*tc_1)
}
