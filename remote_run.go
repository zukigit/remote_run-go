package main

import (
	"fmt"
	"os"
	"strings"
	"zukigit/remote_run-go/src/dao"
	"zukigit/remote_run-go/src/lib"
)

var left, right, endticket, endtestcase, log_filename string

const INFO = 1
const ERR = 2

func set_ticket_values(t []dao.Ticket, auth *dao.Auth) {
	for _, ticket := range t {
		ticket.Set_values(auth)
	}
}

func get_logs(t dao.Ticket) string {
	var pass_count, not_pass_count, must_check_count, unknown_count int
	var body, log_str string
	testcases := t.Get_testcases()

	head := fmt.Sprintf("Ticket[%d] %s\n", t.Get_no(), t.Get_dsctn())

	for index, testcase := range testcases {
		testcase_status := testcase.Get_status()
		switch testcase_status {
		case dao.PASSED:
			pass_count++
		case dao.FAILED:
			not_pass_count++
		case dao.MUST_CHECK:
			must_check_count++
		default:
			unknown_count++
		}

		body = fmt.Sprintf("%s\nTestcase_NO: %d\nTestcase_DES: %s\nStatus: %s\nLogs:", body, testcase.Get_id(), testcase.Get_dsctn(), testcase_status)

		logs := testcase.Get_logs()
		for _, log_value := range logs {
			body = fmt.Sprintf("%s\n%s", body, log_value)
		}

		if index < len(testcases)-1 {
			body = fmt.Sprintf("%s\n%s", body, endtestcase)
		}
	}

	if unknown_count > 0 {
		log_str = fmt.Sprintf("%sPASSED: %d, FAILED: %d, MUST_CHECK: %d, UNKNOWN: %d\n\n%s", head, pass_count, not_pass_count, must_check_count, unknown_count, endtestcase)
	} else {
		log_str = fmt.Sprintf("%sPASSED: %d, FAILED: %d, MUST_CHECK: %d\n\n%s", head, pass_count, not_pass_count, must_check_count, endtestcase)
	}

	return fmt.Sprintf("%s%s\n\n%s", log_str, body, endticket)
}

func run_tc(t []dao.Ticket) {
	for _, ticket := range t {
		ticket.Add_testcases()
		ticket.Run()
		lib.Logi(get_logs(ticket), log_filename)
	}
}

func add_tickets(t *[]dao.Ticket) {
	// Add your tickets here
	*t = append(*t, new(dao.Ticket_1318))
	// is it rainy day?
}

func main() {
	var tickets []dao.Ticket
	add_tickets(&tickets)
	log_filename = lib.Get_log_filename()

	left = strings.Repeat("|", 60)
	right = strings.Repeat("|", 60)
	endticket = fmt.Sprintf("%sX%s", left, right)

	left = strings.Repeat("-", 25)
	right = strings.Repeat("-", 25)
	endtestcase = fmt.Sprintf("%s><%s", left, right)

	auth := dao.Get_auth() // Get login informations from user
	config := lib.Get_config(auth.Username, auth.Password)

	// Connect to the SSH server
	client, err := lib.Get_client(auth.Hostname+":22", config)
	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}
	defer client.Close()

	auth.Ssh_client = client

	set_ticket_values(tickets, auth)
	run_tc(tickets) // run test cases
	fmt.Println(lib.Formatted_log(1, "Logged Filename: %s", log_filename))
}
