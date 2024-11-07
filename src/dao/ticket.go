package dao

type Ticket interface {
	Set_values()
	Get_no() uint
	Get_dsctn() string
	Add_testcase(tc TestCase)
	Add_testcases()
	Get_testcases() []TestCase
	New_testcase(testcase_id uint, testcase_description string) *TestCase
}
