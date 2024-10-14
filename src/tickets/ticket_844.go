package tickets

import (
	"fmt"

	"github.com/zukigit/remote_run-go/src/common"
	"github.com/zukigit/remote_run-go/src/dao"
	"github.com/zukigit/remote_run-go/src/lib"
)

type Ticket_844 struct {
	no          uint
	description string
	testcases   []dao.TestCase
}

func (t *Ticket_844) New_testcase(testcase_id uint, testcase_description string) *dao.TestCase {
	return dao.New_testcase(testcase_id, testcase_description)
}

func (t *Ticket_844) Get_no() uint {
	return t.no
}

func (t *Ticket_844) Get_dsctn() string {
	return t.description
}

func (t *Ticket_844) Add_testcase(tc dao.TestCase) {
	t.testcases = append(t.testcases, tc)
}

func (t *Ticket_844) Get_testcases() []dao.TestCase {
	return t.testcases
}

// Enter your ticket information here
func (t *Ticket_844) Set_values() {
	t.no = 844 // Enter your ticket id
	t.description = "Others"
}

// Add your test case here
func (t *Ticket_844) Add_testcases() {

	//TESTCASE 73
	tc_73 := t.New_testcase(73, "Agent Restart [Covered with Agent servive stop while 1000  parallel jobnets are running]")
	tc_func := func() common.Testcase_status {

		// Executing Test Case for 100 Job Icons
		run_jobnet_id, err := lib.Jobarg_exec("ICON_100")
		if err != nil {
			fmt.Println(tc_73.Err_log("Error: Failed while running job. Jobnet Management Id: %s", run_jobnet_id))
			fmt.Println(err.Error())
			return FAILED
		}

		// In case if getting jobnet info failed
		jobnet_run_info, err := lib.Jobarg_get_jobnet_run_info(run_jobnet_id)
		if err != nil {
			fmt.Println(tc_73.Err_log("Error: Failted at getting jobnet run info. Jobnet Management Id: %s ", run_jobnet_id))
			fmt.Println(err.Error())
			return FAILED
		}

		// Restarting Jobarranger agent service
		// err = lib.Restart_jaz_agent_linux()

		// if err != nil {
		// 	fmt.Println(tc_73.Err_log("Error: Failted at restarting Jobarranger Agent. %s", err.Error()))
		// }

		// Checking whether Jobnet is finished or not
		if jobnet_run_info.Job_status == "NORMAL" && jobnet_run_info.Jobnet_status == "END" && jobnet_run_info.Exit_cd == 0 {

			// Checking whether "checking at zombie process" process failed or not.
			zombieProcessCount, err := lib.CheckZombieProcess(1, common.Client)
			if err != nil {
				fmt.Println(tc_73.Err_log("Error: Failed at checking zombie Process."))
				fmt.Println(err.Error())
				// Checking whether zombie process exist or not.
				if zombieProcessCount > 0 {
					fmt.Println(tc_73.Err_log("Error: There is zombie Process left."))
					fmt.Println(err.Error())
				}

				return FAILED
			}

			return PASSED
		}

		// Checking in case Abnormal Fail.
		fmt.Println(tc_73.Err_log("Jobnet_status: %s, Job_status: %s, Exit_cd: %d", jobnet_run_info.Jobnet_status, jobnet_run_info.Job_status, jobnet_run_info.Exit_cd))
		return FAILED
	}

	tc_73.Set_function(tc_func)
	t.Add_testcase(*tc_73)
}
