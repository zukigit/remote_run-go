package tickets

import (
	"github.com/zukigit/remote_run-go/src/common"
	"github.com/zukigit/remote_run-go/src/dao"
	"github.com/zukigit/remote_run-go/src/lib"
)

type Ticket_775 struct {
	no          uint
	description string
	testcases   []dao.TestCase
}

func (t *Ticket_775) New_testcase(testcase_id uint, testcase_description string) *dao.TestCase {
	return dao.New_testcase(testcase_id, testcase_description)
}

func (t *Ticket_775) Get_no() uint {
	return t.no
}

func (t *Ticket_775) Get_dsctn() string {
	return t.description
}

func (t *Ticket_775) Add_testcase(tc dao.TestCase) {
	t.testcases = append(t.testcases, tc)
}

func (t *Ticket_775) Get_testcases() []dao.TestCase {
	return t.testcases
}

// Enter your ticket information here
func (t *Ticket_775) Set_values() {
	t.no = 1 // Enter your ticket id
	t.description = "Enter your ticket description here."
}

// Add your test case here
func (t *Ticket_775) Add_testcases() {
	// TESTCASE 001
	tc_1 := t.New_testcase(1, "Agent wiht error prone string length of command input.")
	tc_func := func() common.Testcase_status {

		// if err := lib.Jobarg_cleanup_linux(); err != nil {
		// 	return t.logError(tc_1, "Error during cleanup: %s", err)
		// }

		return FAILED
	}
	tc_1.Set_function(tc_func)
	t.Add_testcase(*tc_1)
}
