package main

import (
	"fmt"
	"os"
	"zukigit/remote_run-go/cmd"
	"zukigit/remote_run-go/src/common"
	"zukigit/remote_run-go/src/dao"
	"zukigit/remote_run-go/src/lib"
	"zukigit/remote_run-go/src/tickets"
)

var tkts []dao.Ticket

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

func check_ticket_duplicate() error {
	seen := make(map[uint]bool)

	for _, tkt := range tkts {
		if seen[tkt.Get_no()] {
			return fmt.Errorf("ticket[%d] is duplicated", tkt.Get_no())
		}
		seen[tkt.Get_no()] = true
	}
	return nil
}

func main() {
	common.Log_filename = lib.Get_log_filename()

	cmd.Execute()
	defer common.Client.Close()

	add_tickets(&tkts)
	set_ticket_values(tkts)

	if err := check_ticket_duplicate(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	run_tc(tkts) // run test cases
	fmt.Println(lib.Formatted_log(common.INFO, "Logged Filename: %s", common.Log_filename))
}

func add_tickets(t *[]dao.Ticket) {
	// Add your tickets here
	*t = append(*t, new(tickets.Ticket_1318))
	*t = append(*t, new(tickets.Ticket_000))
}
