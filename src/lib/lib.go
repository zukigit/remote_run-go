package lib

import "zukigit/remote_run-go/src/dao"

var ticket dao.Ticket

func Set_ticket(tc dao.Ticket) {
	ticket = tc
}
