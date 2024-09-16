package dao

type Ticket interface {
	Run_ticket()
	Set_ticket_values(auth *Auth)
	Get_ticket_no() uint
	Get_ticket_dsctn() string
	Get_auth() *Auth
	Add_testcase(tc TestCase)
	Add_testcases()
}
