package tickets

import (
	"github.com/zukigit/remote_run-go/src/common"
	"github.com/zukigit/remote_run-go/src/dao"
)

type Ticket_000 struct {
	no          uint
	description string
	testcases   []dao.TestCase
}

func (t *Ticket_000) New_testcase(testcase_id uint, testcase_description string) *dao.TestCase {
	return dao.New_testcase(testcase_id, testcase_description)
}

func (t *Ticket_000) Get_no() uint {
	return t.no
}

func (t *Ticket_000) Get_dsctn() string {
	return t.description
}

func (t *Ticket_000) Add_testcase(tc dao.TestCase) {
	t.testcases = append(t.testcases, tc)
}

func (t *Ticket_000) Get_testcases() []dao.TestCase {
	return t.testcases
}

// Enter your ticket information here
func (t *Ticket_000) Set_values() {
	t.no = 1 // Enter your ticket id
	t.description = "Enter your ticket description here."
}

// Add your test case here
func (t *Ticket_000) Add_testcases() {
	// TESTCASE 001
	tc_1 := t.New_testcase(1, "Enter your test case description here.")
	tc_func := func() common.Testcase_status {
		// err := lib.Jobarg_enable_jobnet("Icon_1", "jobicon_linux")
		// if err != nil {
		// 	fmt.Println("err in enable jobnet", err.Error())
		// 	return FAILED
		// }
		// envs, _ := lib.Get_str_str_map("JA_HOSTNAME", "oss.linux", "JA_CMD", "hostname")
		// if _, err := lib.Jobarg_exec_E("Icon_1", envs); err != nil {
		// 	fmt.Println("err", err.Error())
		// 	return FAILED
		// } else {
		// 	return FAILED
		// }
		return FAILED
	}
	tc_1.Set_function(tc_func)
	t.Add_testcase(*tc_1)
}
