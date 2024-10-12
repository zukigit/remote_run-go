package tickets

import (
	"github.com/zukigit/remote_run-go/src/common"
	"github.com/zukigit/remote_run-go/src/dao"
)

type Ticket_1225 struct {
	no          uint
	description string
	testcases   []dao.TestCase
}

func (t *Ticket_1225) New_testcase(testcase_id uint, testcase_description string) *dao.TestCase {
	return dao.New_testcase(testcase_id, testcase_description)
}

func (t *Ticket_1225) Get_no() uint {
	return t.no
}

func (t *Ticket_1225) Get_dsctn() string {
	return t.description
}

func (t *Ticket_1225) Add_testcase(tc dao.TestCase) {
	t.testcases = append(t.testcases, tc)
}

func (t *Ticket_1225) Get_testcases() []dao.TestCase {
	return t.testcases
}

// Enter your ticket information here
func (t *Ticket_1225) Set_values() {
	t.no = 1225 // Enter your ticket id
	t.description = "Japurge delete records with limit 500."
}

// Add your test case here
func (t *Ticket_1225) Add_testcases() {
	// TESTCASE 001
	tc_1 := t.New_testcase(79, "Test with 1 million records.")
	tc_func := func() common.Testcase_status {
		// Enter your test case logic here
		return FAILED
	}
	tc_1.Set_function(tc_func)
	t.Add_testcase(*tc_1)
}
