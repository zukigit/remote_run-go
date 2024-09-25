package dao

import (
	"fmt"
	"strings"
)

var left, right, endticket, endtestcase string

type Ticket interface {
	Run()
	Set_values(auth *Auth)
	Get_no() uint
	Get_dsctn() string
	Add_testcase(tc TestCase)
	Add_testcases()
	Get_testcases() []TestCase
	New_testcase(testcase_id uint, testcase_description string) *TestCase
}

func Set_ticket_logs_headers() {
	left = strings.Repeat("|", 60)
	right = strings.Repeat("|", 60)
	endticket = fmt.Sprintf("%sX%s", left, right)

	left = strings.Repeat("-", 25)
	right = strings.Repeat("-", 25)
	endtestcase = fmt.Sprintf("%s><%s", left, right)
}

func Get_ticket_logs(t Ticket) string {
	var pass_count, not_pass_count, must_check_count, unknown_count int
	var body, log_str string
	testcases := t.Get_testcases()

	head := fmt.Sprintf("Ticket[%d] %s\n", t.Get_no(), t.Get_dsctn())

	for index, testcase := range testcases {
		testcase_status := testcase.Get_status()
		switch testcase_status {
		case PASSED:
			pass_count++
		case FAILED:
			not_pass_count++
		case MUST_CHECK:
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
