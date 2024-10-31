package tickets

import (
	"github.com/zukigit/remote_run-go/src/common"
	"github.com/zukigit/remote_run-go/src/dao"
	"github.com/zukigit/remote_run-go/src/lib"
)

type Ticket_1021 struct {
	no          uint
	description string
	testcases   []dao.TestCase
}

func (t *Ticket_1021) New_testcase(testcase_id uint, testcase_description string) *dao.TestCase {
	return dao.New_testcase(testcase_id, testcase_description)
}

func (t *Ticket_1021) Get_no() uint {
	return t.no
}

func (t *Ticket_1021) Get_dsctn() string {
	return t.description
}

func (t *Ticket_1021) Add_testcase(tc dao.TestCase) {
	t.testcases = append(t.testcases, tc)
}

func (t *Ticket_1021) Get_testcases() []dao.TestCase {
	return t.testcases
}

// Enter your ticket information here
func (t *Ticket_1021) Set_values() {
	t.no = 1021 // Enter your ticket id
	t.description = "Simple Testing."
}

// Add your test case here
func (t *Ticket_1021) Add_testcases() {
	// TESTCASE 001
	tc_104 := t.New_testcase(1021, "Simple Description.")

	tc_func := func() common.Testcase_status {

		// var jobnet_id string = "Icon_1"          // This value must be Jobnet_ID that you want to enable
		// var jobnet_name string = "jobicon_linux" // This value must be Jobnet_name that you want to enable.
		// var jobnet_run_manage_id string
		// var jobnet_run_info *common.Jobnet_run_info
		// var result bool

		// if Run_enable_jobnet(tc_104, jobnet_id, jobnet_name) &&
		// 	func() bool {
		// 		result, jobnet_run_manage_id = Run_Jobnet(tc_104, jobnet_id)
		// 		return result
		// 	}() &&
		// 	func() bool {
		// 		result, jobnet_run_info = Run_Jobarg_get_jobnet_run_info(tc_104, jobnet_run_manage_id)
		// 		return result
		// 	}() &&
		// 	Run_Check_Jobnet_Finish_With_No_Zombie_Process(tc_104, jobnet_run_info) {
		// 	fmt.Println("All operations completed successfully")
		// 	return PASSED
		// } else {
		// 	return FAILED
		// }

		var systemEncode string
		var result bool

		if Run_Clear_Agent_log(tc_104) &&
			func() bool {
				result, systemEncode = Run_Linux_Command(tc_104, "locale | grep 'LC_CTYPE='")
				return result
			}() &&
			func() bool {
				lib.Ssh_exec_to_str("cat /var/log/jobarranger/jobarg_agentd.log | grep 'LC_CTYPE :'")
				return true
			}() {
			return PASSED
		} else {
			return FAILED
		}
	}
	tc_104.Set_function(tc_func)
	t.Add_testcase(*tc_104)
}
