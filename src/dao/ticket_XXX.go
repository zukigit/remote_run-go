package dao

import (
	"fmt"
)

type Ticket_XXX struct {
	no          uint
	description string
	auth        *Auth
	testcases   []TestCase
}

func (t *Ticket_XXX) Set_values(auth *Auth) {
	t.no = 000
	t.description = "Enter Your Ticket Description."
	t.auth = auth
}

func (t *Ticket_XXX) New_testcase(testcase_id uint, testcase_description string) *TestCase {
	return New_testcase(testcase_id, testcase_description, t.auth)
}

func (t *Ticket_XXX) Get_no() uint {
	return t.no
}

func (t *Ticket_XXX) Get_dsctn() string {
	return t.description
}

func (t *Ticket_XXX) Set_testcase(tc TestCase) {
	t.testcases = append(t.testcases, tc)
}

func (t *Ticket_XXX) Get_testcases() []TestCase {
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

	// Example adding testcase
	tc_xxx := t.New_testcase(000, "Enter your testcase Description") // Creating test case object
	tc_func := func() Testcase_status {                              // Creating test case function
		return PASSED // There are three types of status. PASSED, FAILED, MUST_CHECK
	}
	tc_xxx.Set_function(tc_func) // Adding created test case function to the test case.
	t.Set_testcase(*tc_xxx)      // Adding created test case to ticket.
}
