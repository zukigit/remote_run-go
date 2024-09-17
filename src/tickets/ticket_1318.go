package tickets

import (
	"fmt"
	"zukigit/remote_run-go/src/dao"
)

type Ticket_1318 struct {
	no          uint
	description string
	auth        *dao.Auth
	testcases   []dao.TestCase
}

func (t *Ticket_1318) Set_values(auth *dao.Auth) {
	t.no = 1318
	t.description = "Fixed for negative JOB_EXT_CD return value."
	t.auth = auth
}

func (t *Ticket_1318) Get_no() uint {
	return t.no
}

func (t *Ticket_1318) Get_dsctn() string {
	return t.description
}

func (t *Ticket_1318) Get_auth() *dao.Auth {
	return t.auth
}

func (t *Ticket_1318) Add_testcase(tc dao.TestCase) {
	t.testcases = append(t.testcases, tc)
}

func (t *Ticket_1318) Add_testcases() {
	// Add your test case here
	tc_168 := dao.New_TestCase(168, "Normal Case")
	tc_168.Add_function(
		func() bool {
			return false
		},
	)
	t.Add_testcase(*tc_168)

	tc_169 := dao.New_TestCase(169, "Normal Case") // create test case
	tc_169.Add_function(                           // Add function to test case
		func() bool {
			return true
		},
	)
	t.Add_testcase(*tc_169) // Add testcase to ticket
}

func (t *Ticket_1318) Run() {
	t.Add_testcases()

	for _, tc := range t.testcases {
		if !tc.Is_function_nil() {
			tc.Set_is_passed(tc.Run_function())
			fmt.Println("tc", tc)
		} else {
			// will add error log
		}
	}
}
