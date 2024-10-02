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

var tkts, run_tickets []dao.Ticket

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

func check_duplicated_ticket() {
	seen := make(map[uint]bool)

	for _, tkt := range tkts {
		if seen[tkt.Get_no()] {
			fmt.Printf("Error: ticket[%d] is duplicated\n", tkt.Get_no())
			os.Exit(1)
		}
		seen[tkt.Get_no()] = true
	}
}

func add_run_tickets(ticket_number uint) {
	if ticket_number == 0 {
		run_tickets = tkts
	} else {
		for _, t := range tkts {
			if ticket_number == t.Get_no() {
				run_tickets = append(run_tickets, t)
				break
			}
		}
	}
}

func main() {
	common.Log_filename = lib.Get_log_filename()

	cmd.Execute()
	defer common.Client.Close()

	add_tickets(&tkts)
	set_ticket_values(tkts)
	check_duplicated_ticket()
	add_run_tickets(common.Specific_ticket_no)

	run_tc(run_tickets) // run test cases

	if len(run_tickets) > 0 {
		fmt.Println(lib.Formatted_log(common.INFO, "Logged Filename: %s", common.Log_filename))
	} else {
		fmt.Println("There is no ticket to run.")
	}
}

func add_tickets(t *[]dao.Ticket) {
	// Add your tickets here
	*t = append(*t, new(tickets.Ticket_1318))
	*t = append(*t, new(tickets.Ticket_000))
}
