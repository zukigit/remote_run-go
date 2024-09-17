package dao

type Ticket interface {
	Run()
	Set_values(auth *Auth)
	Get_no() uint
	Get_dsctn() string
	Get_auth() *Auth
	Add_testcase(tc TestCase)
	Add_testcases()
}
