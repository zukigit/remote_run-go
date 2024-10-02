package main

import (
	"fmt"
	"zukigit/remote_run-go/cmd"
	"zukigit/remote_run-go/src/common"
	"zukigit/remote_run-go/src/dao"
	"zukigit/remote_run-go/src/lib"
	"zukigit/remote_run-go/src/tickets"
)

const INFO = 1
const ERR = 2

func set_ticket_values(t []dao.Ticket) {
	for _, ticket := range t {
		ticket.Set_values()
	}
}

func run_tc(t []dao.Ticket) {
	for _, ticket := range t {
		ticket.Add_testcases()
		ticket.Run()
		lib.Logi(dao.Get_ticket_logs(ticket), common.Log_filename)
	}
}

func main() {
	var tickets []dao.Ticket

	common.Log_filename = lib.Get_log_filename()

	cmd.Execute()
	defer common.Client.Close()

	add_tickets(&tickets)
	set_ticket_values(tickets)

	run_tc(tickets) // run test cases
	fmt.Println(lib.Formatted_log(INFO, "Logged Filename: %s", common.Log_filename))
}

func add_tickets(t *[]dao.Ticket) {
	// Add your tickets here
	*t = append(*t, new(tickets.Ticket_1318))
	*t = append(*t, new(tickets.Ticket_000))
}
