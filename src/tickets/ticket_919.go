package tickets

import (
	"github.com/zukigit/remote_run-go/src/common"
	"github.com/zukigit/remote_run-go/src/dao"
)

type Ticket_919 struct {
	no          uint
	description string
	testcases   []dao.TestCase
}

func (t *Ticket_919) New_testcase(testcase_id uint, testcase_description string) *dao.TestCase {
	return dao.New_testcase(testcase_id, testcase_description)
}

func (t *Ticket_919) Get_no() uint {
	return t.no
}

func (t *Ticket_919) Get_dsctn() string {
	return t.description
}

func (t *Ticket_919) Add_testcase(tc dao.TestCase) {
	t.testcases = append(t.testcases, tc)
}

func (t *Ticket_919) Get_testcases() []dao.TestCase {
	return t.testcases
}

// Enter your ticket information here
func (t *Ticket_919) Set_values() {
	t.no = 919 // Enter your ticket id
	t.description = " Add retry process to DB deadlock situations."
}
func (t *Ticket_919) Add_testcases() {

	tc_89 := t.New_testcase(89, "Check if jobs can rerun sucessfully after deadlock.")
	tc_func := func() common.Testcase_status {

		// if err := lib.Jobarg_cleanup_linux(); err != nil {
		// 	return t.logError(tc_38, "Error during cleanup: %s", err)
		// }

		// if err := lib.Jobarg_enable_jobnet("Icon_1", "jobicon_linux"); err != nil {
		// 	return t.logError(tc_38, "Error during enable jobnet: %s", err)
		// }
		return FAILED
	}

	tc_89.Set_function(tc_func)
	t.Add_testcase(*tc_89)
}
