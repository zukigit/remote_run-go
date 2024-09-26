package main

import (
	"fmt"
	"zukigit/remote_run-go/src/dao"
	"zukigit/remote_run-go/src/lib"
	"zukigit/remote_run-go/src/tickets"
)

const INFO = 1
const ERR = 2

var log_filename string

func set_ticket_values(t []dao.Ticket, auth *dao.Auth) {
	for _, ticket := range t {
		ticket.Set_values(auth)
	}
}

func run_tc(t []dao.Ticket) {
	for _, ticket := range t {
		ticket.Add_testcases()
		ticket.Run()
		lib.Logi(dao.Get_ticket_logs(ticket), log_filename)
	}
}

func main() {
	var tickets []dao.Ticket

	add_tickets(&tickets)
	log_filename = lib.Get_log_filename()

	auth := dao.Get_auth() // Get login informations from user
	defer auth.Ssh_client.Close()

	lib.Set_common_client(auth.Ssh_client)
	dao.Set_ticket_logs_headers()
	set_ticket_values(tickets, auth)

	run_tc(tickets) // run test cases
	fmt.Println(lib.Formatted_log(INFO, "Logged Filename: %s", log_filename))
}

func add_tickets(t *[]dao.Ticket) {
	// Add your tickets here
	*t = append(*t, new(tickets.Ticket_1318))
	*t = append(*t, new(tickets.Ticket_000))
}
