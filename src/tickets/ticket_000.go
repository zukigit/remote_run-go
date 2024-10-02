package tickets

import (
	"fmt"
	"zukigit/remote_run-go/src/common"
	"zukigit/remote_run-go/src/dao"
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

func (t *Ticket_000) Run() {
	for _, tc := range t.testcases {
		fmt.Println(tc.Info_log("running..."))
		if !tc.Is_function_nil() {
			tc.Set_status(tc.Run_function())
		} else {
			fmt.Println(tc.Err_log("has no function. SKIPPED!"))
			tc.Set_status(FAILED)
		}

		fmt.Println(tc.Info_log("finished!"))
	}
}

// Enter your ticket information here
func (t *Ticket_000) Set_values() {
	t.no = 1318 // Enter your ticket id
	t.description = "Enter your ticket description here."
}

// Add your test case here
func (t *Ticket_000) Add_testcases() {
	// TESTCASE 000
	tc_000 := t.New_testcase(000, "Enter your test case description here.")
	tc_func := func() common.Testcase_status {
		// Enter your test case logic here
		return FAILED
	}
	tc_000.Set_function(tc_func)
	t.Add_testcase(*tc_000)
}
