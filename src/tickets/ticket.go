package tickets

import "zukigit/remote_run-go/src/dao"

type Ticket interface {
	Run()
	Set_values(auth *dao.Auth)
	Get_no() uint
	Get_dsctn() string
	Set_testcase(tc dao.TestCase)
	Add_testcases()
	Get_testcases() []dao.TestCase
	New_testcase(testcase_id uint, testcase_description string) *dao.TestCase
}
