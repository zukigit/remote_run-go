package dao

import (
	"fmt"

	"github.com/zukigit/remote_run-go/src/common"
)

const (
	PASSED     common.Testcase_status = "PASSED"
	FAILED     common.Testcase_status = "FAILED"
	MUST_CHECK common.Testcase_status = "MUST_CHECK"
)

func Run_testcase(tc TestCase) {
	fmt.Println(tc.Info_log("running..."))
	if !tc.Is_function_nil() {
		tc.Set_status(tc.Run_function())
	} else {
		fmt.Println(tc.Err_log("has no function. SKIPPED!"))
		tc.Set_status(FAILED)
	}
	fmt.Println(tc.Info_log("finished!"))
}

func Update_testcase_results_in_tickets(tks []Ticket) {
	var passed_count, failed_count, mustcheck_count int
	for _, ticket := range tks {
		for _, testcase := range ticket.Get_testcases() {
			switch testcase.Get_status() {
			case PASSED:
				passed_count++
			case FAILED:
				failed_count++
			default:
				mustcheck_count++
			}
		}
		ticket.Set_PASSED_count(passed_count)
		ticket.Set_FAILED_count(failed_count)
		ticket.Set_MUSTCHECK_count(mustcheck_count)
	}
}
