package dao

type Ticket interface {
	Run()
	Set_values(auth *Auth)
	Get_no() uint
	Get_dsctn() string
	Set_testcase(tc TestCase)
	Add_testcases()
	Get_testcases() []TestCase
	New_testcase(testcase_id uint, testcase_description string) *TestCase
}
