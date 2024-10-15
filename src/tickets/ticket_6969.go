package tickets

import (
	"fmt"

	"github.com/zukigit/remote_run-go/src/common"
	"github.com/zukigit/remote_run-go/src/dao"
	"github.com/zukigit/remote_run-go/src/lib"
)

type Ticket_6969 struct {
	no          uint
	description string
	testcases   []dao.TestCase
}

func (t *Ticket_6969) New_testcase(testcase_id uint, testcase_description string) *dao.TestCase {
	return dao.New_testcase(testcase_id, testcase_description)
}

func (t *Ticket_6969) Get_no() uint {
	return t.no
}

func (t *Ticket_6969) Get_dsctn() string {
	return t.description
}

func (t *Ticket_6969) Add_testcase(tc dao.TestCase) {
	t.testcases = append(t.testcases, tc)
}

func (t *Ticket_6969) Get_testcases() []dao.TestCase {
	return t.testcases
}

// Enter your ticket information here
func (t *Ticket_6969) Set_values() {
	t.no = 6969 // Enter your ticket id
	t.description = "My first simple ticket automation"
}

// Add your test case here
func (t *Ticket_6969) Add_testcases() {

	// TESTCASE 6969
	tc_6999 := t.New_testcase(6969, "Sample testcase 1")
	tc_func := func() common.Testcase_status {

		//Jobnet Execution

		jobnet_run_manage_id, err := lib.Jobarg_exec("SIMPLE_JOB")
		if err != nil {
			fmt.Println(tc_6999.Err_log("Error: Failed when trying to run the job. %s.", err.Error()))
			fmt.Println(tc_6999.Err_log("Error: Did you imported the necessary xml files or forgot to enable jobnet?(TICKET844_TESTCASE73-75)"))
			return FAILED
		} else {
			fmt.Println(tc_6999.Info_log("Info: Job has start running."))
		}

		//Getting jobnet info
		jobnet_run_info, err := lib.Jobarg_get_jobnet_run_info(jobnet_run_manage_id)

		if err != nil {
			fmt.Println(tc_6999.Err_log("Error: %s.", err.Error()))
			return FAILED
		}

		if jobnet_run_info.Jobnet_status == "END" && jobnet_run_info.Exit_cd == 0 && jobnet_run_info.Job_status == "NORMAL" {
			return PASSED
		}

		// Checking whether Jobnet is finished or not
		// if jobnet_run_info.Job_status == "NORMAL" && jobnet_run_info.Jobnet_status == "END" && jobnet_run_info.Exit_cd == 0 {

		// 	// Checking whether "checking at zombie process" process failed or not.
		// 	zombieProcessCount, err := lib.CheckZombieProcess(1, common.Client)
		// 	if err != nil {
		// 		fmt.Println(tc_73.Err_log("Error: Failed at checking zombie Process."))
		// 		fmt.Println(tc_73.Err_log(err.Error()))
		// 		// Checking whether zombie process exist or not.
		// 		if zombieProcessCount > 0 {
		// 			fmt.Println(tc_73.Err_log("Error: There is zombie Process left."))
		// 			fmt.Println(tc_73.Err_log(err.Error()))
		// 		}

		// 		return FAILED
		// 	}

		// 	return PASSED
		// }

		// In case if getting jobnet info failed
		// jobnet_run_info, err := lib.Jobarg_get_jobnet_run_info(jobnet_run_manage_id)
		// if err != nil {
		// 	fmt.Println(tc_73.Err_log("Error: Failted at getting jobnet run info. Jobnet Management Id: %s ", jobnet_run_manage_id))
		// 	fmt.Println(tc_73.Err_log(err.Error()))
		// 	return FAILED
		// }

		// Checking in case Abnormal Fail.
		fmt.Println(tc_6999.Err_log("Jobnet_status: %s, Job_status: %s, Exit_cd: %d", jobnet_run_info.Jobnet_status, jobnet_run_info.Job_status, jobnet_run_info.Exit_cd))
		return FAILED
	}
	tc_6999.Set_function(tc_func)
	t.Add_testcase(*tc_6999)
}
