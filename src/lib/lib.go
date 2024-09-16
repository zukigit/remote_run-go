package lib

import "zukigit/remote_run-go/src/dao"

var test_case dao.Ticket

func Set_ticket(tc dao.Ticket) {
	test_case = tc
}
