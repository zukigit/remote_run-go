package main

import (
	"fmt"
	"os"
	"strings"
	"zukigit/remote_run-go/src/dao"
	"zukigit/remote_run-go/src/lib"
	"zukigit/remote_run-go/src/tickets"
)

var left, right, endticket, endtestcase string

func set_ticket_values(t []tickets.Ticket, auth *dao.Auth) {
	for _, ticket := range t {
		ticket.Set_values(auth)
	}
}

func logi(log string) {
	fmt.Println(log)
}

func get_logs(t tickets.Ticket) string {
	var pass_count, not_pass_count int
	var body, log_str, status string
	testcases := t.Get_testcases()

	head := fmt.Sprintf("%d - %s\n", t.Get_no(), t.Get_dsctn())

	for index, testcase := range testcases {
		is_passed := testcase.Get_is_passed()
		switch is_passed {
		case true:
			status = "Passed"
			pass_count++
		case false:
			status = "Failed"
			not_pass_count++
		}

		body = fmt.Sprintf("%s\n\nTestcase_NO: %d\nTestcase_DES: %s\nStatus: %s\nLogs:", body, testcase.Get_id(), testcase.Get_dsctn(), status)

		logs := testcase.Get_logs()
		for _, log_value := range logs {
			body = fmt.Sprintf("%s\n%s", body, log_value)
		}

		if index < len(testcases)-1 {
			body = fmt.Sprintf("%s\n%s", body, endtestcase)
		}
	}

	log_str = fmt.Sprintf("%sPassed: %d, Failed: %d", head, pass_count, not_pass_count)

	return fmt.Sprintf("%s%s\n\n%s", log_str, body, endticket)
}

func run_tc(t []tickets.Ticket) {
	for _, ticket := range t {
		ticket.Add_testcases()
		ticket.Run()
		logi(get_logs(ticket))
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

	left = strings.Repeat("|", 60)
	right = strings.Repeat("|", 60)
	endticket = fmt.Sprintf("%sX%s", left, right)

	left = strings.Repeat("-", 25)
	right = strings.Repeat("-", 25)
	endtestcase = fmt.Sprintf("%s><%s", left, right)

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
