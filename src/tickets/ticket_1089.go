package tickets

import (
	"fmt"

	"github.com/zukigit/remote_run-go/src/common"
	"github.com/zukigit/remote_run-go/src/dao"
)

type Ticket_1089 struct {
	no          uint
	description string
	testcases   []dao.TestCase
}

func (t *Ticket_1089) New_testcase(testcase_id uint, testcase_description string) *dao.TestCase {
	return dao.New_testcase(testcase_id, testcase_description)
}

func (t *Ticket_1089) Get_no() uint {
	return t.no
}

func (t *Ticket_1089) Get_dsctn() string {
	return t.description
}

func (t *Ticket_1089) Add_testcase(tc dao.TestCase) {
	t.testcases = append(t.testcases, tc)
}

func (t *Ticket_1089) Get_testcases() []dao.TestCase {
	return t.testcases
}

// Enter your ticket information here
func (t *Ticket_1089) Set_values() {
	t.no = 1089 // Enter your ticket id
	t.description = "Enter your ticket description here."
}

// Add your test case here
func (t *Ticket_1089) Add_testcases() {
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

		if Run_Jobarg_cleanup_linux(tc_101) &&
			Run_enable_jobnet(tc_101, jobnet_id, jobnet_name_1) &&
			func() bool {
				result, jobnet_run_manage_id = Run_Jobnet_Exec(tc_101, jobnet_id, execute_command)
				return result
			}() &&
			func() bool {
				result, jobnet_run_info = Run_Jobarg_get_jobnet_run_info(tc_101, jobnet_run_manage_id)
				return result
			}() &&
			func() bool {
				fmt.Print(tc_101.Info_log("Info: Job Info: %s. Jobnet Info: %s", jobnet_run_info.Job_status, jobnet_run_info.Jobnet_status))
				if jobnet_run_info.Job_status == "TIMEOUT" {
					fmt.Println(tc_101.Info_log("Info: Job Icon end up with Timeout."))
					return true
				}
				fmt.Println(tc_101.Info_log("Error: Job Icon doesn't end up with Timeout."))
				return false
			}() &&
			Run_Jobarg_cleanup_linux(tc_101) &&
			Run_enable_jobnet(tc_101, jobnet_id, jobnet_name_2) &&
			func() bool {
				result, jobnet_run_manage_id = Run_Jobnet_Exec(tc_101, jobnet_id, execute_command)
				return result
			}() &&
			func() bool {
				result, jobnet_run_info = Run_Jobarg_get_jobnet_run_info(tc_101, jobnet_run_manage_id)
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
