package tickets

import (
	"github.com/zukigit/remote_run-go/src/common"
)

type Ticket_template struct {
	Ticket_no                                   uint
	Ticket_description                          string
	PASSED_count, FAILED_count, MUSTCHECK_count int
	Testcases                                   []common.TestCase
}

func (t *Ticket_template) New_testcase(testcase_id uint, testcase_description string) *common.TestCase {
	return common.New_testcase(testcase_id, testcase_description)
}

func (t *Ticket_template) Get_ticket_no() uint {
	return t.Ticket_no
}

func (t *Ticket_template) Set_ticket_no(ticket_no uint) {
	t.Ticket_no = ticket_no
}

func (t *Ticket_template) Set_PASSED_count(passed_count int) {
	t.PASSED_count = passed_count
}

func (t *Ticket_template) Set_FAILED_count(failed_count int) {
	t.FAILED_count = failed_count
}

func (t *Ticket_template) Set_MUSTCHECK_count(mustcheck_count int) {
	t.MUSTCHECK_count = mustcheck_count
}

func (t *Ticket_template) Get_ticket_description() string {
	return t.Ticket_description
}

func (t *Ticket_template) Set_ticket_description(testcase_description string) {
	t.Ticket_description = testcase_description
}

func (t *Ticket_template) Add_testcase(tc *common.TestCase) {
	t.Testcases = append(t.Testcases, *tc)
}

func (t *Ticket_template) Get_testcases() []common.TestCase {
	return t.Testcases
}

// Add your test case here
func (t *Ticket_template) Prepare() {
	t.Set_ticket_no(1) // Enter your ticket id here
	t.Set_ticket_description("Enter your ticket description here.")

	// TESTCASE 001
	tc := t.New_testcase(3, "Enter your test case description here.")
	tc_func := func() common.Testcase_status {
		//Enter your testcase function here
		return FAILED
	}
	tc.Set_function(tc_func)
	t.Add_testcase(tc)
}
