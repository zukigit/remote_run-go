package main

import (
	"fmt"
	"os"
	"zukigit/remote_run-go/src/dao"
	"zukigit/remote_run-go/src/lib"
	"zukigit/remote_run-go/src/tickets"
)

func set_ticket_values(t []dao.Ticket, auth *dao.Auth) {
	for _, ticket := range t {
		ticket.Set_values(auth)
	}
}

func run_tc(t []dao.Ticket) {
	for _, ticket := range t {
		ticket.Add_testcases()
		ticket.Run()
		lib.Logi()
	}
}

func add_tickets(t *[]dao.Ticket) {
	// Add your tickets here
	*t = append(*t, new(tickets.Ticket_1318))
	// is it rainy day?
}

func main() {
	var tickets []dao.Ticket
	add_tickets(&tickets)

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
