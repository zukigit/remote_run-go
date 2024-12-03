package dao

import (
	"fmt"
	"time"

	"github.com/zukigit/remote_run-go/src/common"
	"github.com/zukigit/remote_run-go/src/lib"
)

const (
	PASSED     common.Testcase_status = "PASSED"
	FAILED     common.Testcase_status = "FAILED"
	MUST_CHECK common.Testcase_status = "MUST_CHECK"
)

func Run_testcase(tc TestCase) {
	fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "running..."))
	common.Current_tk_no = tc.Get_ticket_no()
	common.Current_tc_no = tc.Get_no()

	if !tc.Is_function_nil() {
		// start time
		startTime := time.Now()

		tc.Set_status(tc.Run_function())

		// total elasped time or duration of testcase
		duration := time.Since(startTime)
		durationStr := fmt.Sprintf("%02d:%02d:%02d", int(duration/time.Hour), int(duration/time.Minute)%60, int(duration/time.Second)%60)

		tc.Set_duration(durationStr)
	} else {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "has no function. SKIPPED!"))
		tc.Set_status(FAILED)
	}

	common.Current_tk_no = 0
	common.Current_tc_no = 0
	fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "finished!"))
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
