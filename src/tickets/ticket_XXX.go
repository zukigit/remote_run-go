package tickets

import (
	"fmt"
	"zukigit/remote_run-go/src/dao"
)

type Ticket_XXX struct {
	no          uint
	description string
	auth        *dao.Auth
	testcases   []dao.TestCase
}

func (t *Ticket_XXX) Set_values(auth *dao.Auth) {
	t.no = 000 // Enter your ticket id here
	t.description = "Enter your ticket description here."
	t.auth = auth
}

func (t *Ticket_XXX) New_testcase(testcase_id uint, testcase_description string) *dao.TestCase {
	return dao.New_testcase(testcase_id, testcase_description, t.auth)
}

func (t *Ticket_XXX) Get_no() uint {
	return t.no
}

func (t *Ticket_XXX) Get_dsctn() string {
	return t.description
}

func (t *Ticket_XXX) Add_testcase(tc dao.TestCase) {
	t.testcases = append(t.testcases, tc)
}

func (t *Ticket_XXX) Get_testcases() []dao.TestCase {
	return t.testcases
}

func (t *Ticket_XXX) Run() {
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

func (t *Ticket_XXX) Add_testcases() {
	// Add your test case here

	// ticket 000
	tc_000 := t.New_testcase(000, "Enter your test case description here.")
	tc_func := func() dao.Testcase_status {
		// Enter your test case function here
		return FAILED
	}
	tc_000.Set_function(tc_func)
	t.Add_testcase(*tc_000)
}
