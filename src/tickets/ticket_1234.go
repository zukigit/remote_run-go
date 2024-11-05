package tickets

import (
	"github.com/zukigit/remote_run-go/src/common"
	"github.com/zukigit/remote_run-go/src/dao"
)

type Ticket_1234 struct {
	no          uint
	description string
	testcases   []dao.TestCase
}

func (t *Ticket_1234) New_testcase(testcase_id uint, testcase_description string) *dao.TestCase {
	return dao.New_testcase(testcase_id, testcase_description)
}

func (t *Ticket_1234) Get_no() uint {
	return t.no
}

func (t *Ticket_1234) Get_dsctn() string {
	return t.description
}

func (t *Ticket_1234) Add_testcase(tc dao.TestCase) {
	t.testcases = append(t.testcases, tc)
}

func (t *Ticket_1234) Get_testcases() []dao.TestCase {
	return t.testcases
}

// Enter your ticket information here
func (t *Ticket_1234) Set_values() {
	t.no = 1234 // Enter your ticket id
	t.description = "Enter your ticket description here."
}

// Add your test case here
func (t *Ticket_1234) Add_testcases() {
	// TESTCASE 001
	tc_1234 := t.New_testcase(12345, "Enter your test case description here.")
	tc_func := func() common.Testcase_status {
		if Run_enable_jobnet(tc_1234, "Icon_1", "fcopy") {
			return PASSED
		}
		return FAILED
	}
	tc_1234.Set_function(tc_func)
	t.Add_testcase(*tc_1234)
}
