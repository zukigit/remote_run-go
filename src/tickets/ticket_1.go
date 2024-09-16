package tickets

import (
	"zukigit/remote_run-go/src/dao"
	"zukigit/remote_run-go/src/lib"
)

type Ticket_1 struct {
	ticket_no          uint
	ticket_description string
	auth               *dao.Auth
}

func (T *Ticket_1) Set_Ticket_values(auth *dao.Auth) {
	T.ticket_no = 1
	T.ticket_description = "Test 1 is goood lol"
	T.auth = auth
}

func (T *Ticket_1) Get_Ticket_no() uint {
	return T.ticket_no
}

func (T *Ticket_1) Get_Ticket_dsctn() string {
	return T.ticket_description
}

func (T *Ticket_1) Get_auth() *dao.Auth {
	return T.auth
}

func (T *Ticket_1) Run() {
	lib.Set_ticket(T)

	// Your logic from here
	// _, error := lib.Exec("pwd")

	// return error == nil
}
