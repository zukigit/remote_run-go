package main

import (
	"fmt"
	"os"
	"strings"
	"zukigit/remote_run-go/src/dao"
	"zukigit/remote_run-go/src/lib"
	"zukigit/remote_run-go/src/tickets"
)

var left, right, endticket string

func set_ticket_values(t []tickets.Ticket, auth *dao.Auth) {
	for _, ticket := range t {
		ticket.Set_values(auth)
	}
}

func logi(t tickets.Ticket) {
	fmt.Printf("Ticket_NO: '%d', Ticket_DES: '%s'\n", t.Get_no(), t.Get_dsctn())
	fmt.Println("Testcases:")
	fmt.Println()
	// Print the balanced string
	testcases := t.Get_testcases()

	for _, testcase := range testcases {
		fmt.Printf("Testcase_NO: %d, Testcase_DES: %s\n", testcase.Get_id(), testcase.Get_dsctn())
	}
	fmt.Println(endticket)
}

func run_tc(t []tickets.Ticket) {
	for _, ticket := range t {
		ticket.Add_testcases()
		ticket.Run()
		logi(ticket)
	}
}

func add_tickets(t *[]tickets.Ticket) {
	// Add your tickets here
	*t = append(*t, new(tickets.Ticket_1318))
	// is it rainy day?
}

func main() {
	var tickets []tickets.Ticket
	add_tickets(&tickets)

	left = strings.Repeat("-", 60)
	right = strings.Repeat("-", 60)
	endticket = fmt.Sprintf("%s><%s", left, right)

	auth := dao.Get_auth() // Get login informations from user
	config := lib.Get_config(auth.Username, auth.Password)

	// Connect to the SSH server
	client, err := lib.Get_client(auth.Hostname, config)
	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}
	defer client.Close()

	// Create a session for the command
	session, err := lib.Get_session(client)
	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}
	defer session.Close()
	auth.Session = session

	set_ticket_values(tickets, auth)

	run_tc(tickets) // run test cases
}
