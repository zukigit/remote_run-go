package common

type Ticket interface {
	Set_values()
	Get_no() uint
	Get_dsctn() string
	Set_PASSED_count(passed_count int)
	Set_FAILED_count(failed_count int)
	Set_MUSTCHECK_count(mustcheck_count int)
	Add_testcase(tc TestCase)
	Add_testcases()
	Get_testcases() []TestCase
	New_testcase(testcase_id uint, testcase_description string) *TestCase
}
