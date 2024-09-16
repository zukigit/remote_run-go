package dao

type Ticket interface {
	Run()
	Set_Ticket_values(auth *Auth)
	Get_Ticket_no() uint
	Get_Ticket_dsctn() string
	Get_auth() *Auth
}
