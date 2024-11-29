package tickets

import (
	"fmt"

	"github.com/zukigit/remote_run-go/src/common"
	"github.com/zukigit/remote_run-go/src/dao"
	"github.com/zukigit/remote_run-go/src/lib"
)

type Ticket_1089 struct {
	Ticket_no                                   uint
	Ticket_description                          string
	PASSED_count, FAILED_count, MUSTCHECK_count int
	Testcases                                   []dao.TestCase
}

func (t *Ticket_1089) New_testcase(testcase_id uint, testcase_description string) *dao.TestCase {
	return dao.New_testcase(testcase_id, testcase_description)
}

func (t *Ticket_1089) Get_no() uint {
	return t.Ticket_no
}

func (t *Ticket_1089) Set_PASSED_count(passed_count int) {
	t.PASSED_count = passed_count
}

func (t *Ticket_1089) Set_FAILED_count(failed_count int) {
	t.FAILED_count = failed_count
}

func (t *Ticket_1089) Set_MUSTCHECK_count(mustcheck_count int) {
	t.MUSTCHECK_count = mustcheck_count
}

func (t *Ticket_1089) Get_dsctn() string {
	return t.Ticket_description
}

func (t *Ticket_1089) Add_testcase(tc dao.TestCase) {
	t.Testcases = append(t.Testcases, tc)
}

func (t *Ticket_1089) Get_testcases() []dao.TestCase {
	return t.Testcases
}

// Enter your ticket information here
func (t *Ticket_1089) Set_values() {
	t.Ticket_no = 1089 // Enter your ticket id
	t.Ticket_description = "Enter your ticket description here."
}

// Add your test case here
func (t *Ticket_1089) Add_testcases() {
	// Issue 1089 - Jobnet Timeout
	//
	// Normal Case Check
	//
	// - Create a jobnet with a timeout warning flag set to 1 minute.
	// - Create a jobnet with a timeout jobnet stop flag set to 1 minute.
	//
	// - Run those jobnet and wait 1 min at job execution management screen.
	//
	// - Jobnet's status background color should be changed from yellow to orange after 1 min..
	// - Jobnet's status background color should be changed from yellow to red after 1 min.

	// TESTCASE 101
	tc_101 := t.New_testcase(101, "Enter your test case description here.")
	tc_func := func() common.Testcase_status {

		var jobnet_id = "Icon_1"                       // This value must be Jobnet_ID that you want to run.
		var jobnet_name_1 = "TICKET1089_TESTCASE101_1" // This value must be Jobnet_Name that you want to run for first task.
		var jobnet_name_2 = "TICKET1089_TESTCASE101_2" // This value must be Jobnet_Name that you want to run for second task.
		var execute_command = "sleep 120"              // Enter command you want to execute for this jobnet
		var jobnet_run_manage_id string
		var jobnet_run_info *common.Jobnet_run_info
		var result bool

		if lib.Run_Jobarg_cleanup_linux() &&
			lib.Run_enable_jobnet(jobnet_id, jobnet_name_1) &&
			func() bool {
				result, jobnet_run_manage_id = lib.Run_Jobnet_Exec(jobnet_id, execute_command)
				return result
			}() &&
			func() bool {
				result, jobnet_run_info = lib.Run_Jobarg_get_jobnet_run_info(jobnet_run_manage_id)
				return result
			}() &&
			func() bool {
				fmt.Print(tc_101.Info_log("Info: Job Info: %s. Jobnet Info: %s", jobnet_run_info.Job_status, jobnet_run_info.Jobnet_status))
				if jobnet_run_info.Job_status == "TIMEOUT" {
					var count int
					_, sql_result := lib.Run_Sql_Script_Return_Rows("SELECT count(*) from ja_run_job_table WHERE inner_jobnet_main_id = '" + jobnet_run_manage_id + "' AND status = 3 AND timeout_flag = 1;")
					if sql_result.Next() { // Move to the first row
						if err := sql_result.Scan(&count); err != nil {
							fmt.Println(tc_101.Err_log("Error: Error scanning result: %s", err))
							return false
						}
					}
					if count != 0 {
						fmt.Println(tc_101.Info_log("Info: Jobnet End with Timeout. Job Icon end up with Timeout."))
						return true
					}
					fmt.Println(tc_101.Info_log("Info: Jobnet End with Timeout."))
					return true
				}
				fmt.Println(tc_101.Info_log("Error: Jobnet doesn't end up with Timeout."))
				return false
			}() &&
			lib.Run_Jobarg_cleanup_linux() &&
			lib.Run_enable_jobnet(jobnet_id, jobnet_name_2) &&
			func() bool {
				result, jobnet_run_manage_id = lib.Run_Jobnet_Exec(jobnet_id, execute_command)
				return result
			}() &&
			func() bool {
				result, jobnet_run_info = lib.Run_Jobarg_get_jobnet_run_info(jobnet_run_manage_id)
				if jobnet_run_info == nil {
					fmt.Println(tc_101.Info_log("Info: Job Icon end up with Error."))
				} else {
					fmt.Println(tc_101.Info_log("Info: Job Icon doesn't end up with Error."))
				}
				return !result
			}() {
			fmt.Println("All operations completed successfully")
			return PASSED
		}
		return FAILED
	}
	tc_101.Set_function(tc_func)
	t.Add_testcase(*tc_101)
}
