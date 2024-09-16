package tickets

import (
	"fmt"
	"zukigit/remote_run-go/src/dao"
	"zukigit/remote_run-go/src/lib"
)

type Ticket_1318 struct {
	ticket_no          uint
	ticket_description string
	auth               *dao.Auth
}

func (T *Ticket_1318) Set_Ticket_values(auth *dao.Auth) {
	T.ticket_no = 1318
	T.ticket_description = "Fixed for negative JOB_EXT_CD return value."
	T.auth = auth
}

func (T *Ticket_1318) Get_Ticket_no() uint {
	return T.ticket_no
}

func (T *Ticket_1318) Get_Ticket_dsctn() string {
	return T.ticket_description
}

func (T *Ticket_1318) Get_auth() *dao.Auth {
	return T.auth
}

func (T *Ticket_1318) Run() {
	// Enter Your logic from here
	tc_168 := dao.New_TestCase(168, "Default Case")

	_, err := lib.Start_jaz_server()
	if err != nil {
		tc_168.Add_log(err.Error())
		tc_168.Fail()
	} else {
		// add info log
		tc_168.Succeed()
	}

	fmt.Println("Ran", tc_168)
}
