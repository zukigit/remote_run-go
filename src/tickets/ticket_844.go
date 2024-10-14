package tickets

import (
	"fmt"
	"time"

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

		// Executing Test Case for 800 Job Icons
		jobnet_run_manage_id, err := lib.Jobarg_exec("Icon_100")
		if err != nil {
			fmt.Println(tc_73.Err_log("Error: Failed when trying to run the job. %s. Jobnet Management Id: %s", err.Error(), jobnet_run_manage_id))
			fmt.Println(tc_73.Err_log("Error: Did you imported the necessary xml files?(TICKET844_TESTCASE73-75)"))
			return FAILED
		} else {
			fmt.Println(tc_73.Info_log("Info: Job has start running."))
		}

		//jobnetstatus := lib.Jobarg_get_JA_JOBNETSTATUS(jobnet_run_manage_id)

		// Counting Job Process Count
		fmt.Println(tc_73.Info_log("Info: Counting Job processes."))
		err = lib.JobProcessCountCheck(100, 10, common.Client)
		if err != nil {
			fmt.Println(tc_73.Err_log("Error: Failed at counting required process count. %s Jobnet Management Id: %s ", err.Error(), jobnet_run_manage_id))
			return FAILED
		} else {
			fmt.Println(tc_73.Info_log("Info: Job process count met."))
		}

		// Restartng Jobarranger linux Agent.
		fmt.Println(tc_73.Info_log("Info: Restarting Jobarranger Agent now."))
		err = lib.Restart_jaz_agent_linux()
		if err != nil {
			fmt.Println(tc_73.Err_log("Error: Failted at restarting Jobarranger Agent. %s", err.Error()))
			return FAILED
		} else {
			fmt.Println(tc_73.Info_log("Info: Jobarranger Agent has successfully restarted."))
		}

		// Counting Job Process Count again.
		fmt.Println(tc_73.Info_log("Info: Counting Job processes again..."))
		err = lib.JobProcessCountCheck(0, 1, common.Client)
		if err != nil {
			fmt.Println(tc_73.Err_log("Error: Failed at counting required process count. %s", err.Error()))

		} else {

			// Checking 'Job_Status' to get the value 'Error' with 'Jobarg_get' command."
			fmt.Println(tc_73.Info_log("Info: Total 0 Job process count met. Waiting for Jobnet to end with RUNERR.(Running Error)"))
			timeoutDuration := 10 * time.Minute
			timeout := time.Now().Add(timeoutDuration)
			for time.Now().Before(timeout) {
				jobstatus, err := lib.Jobarg_get_JA_JOBSTATUS(jobnet_run_manage_id)
				if err == nil && jobstatus == "ERROR" {
					fmt.Println(tc_73.Info_log("Info: The job returns with Error. Success."))
					return PASSED
				}
				time.Sleep(1 * time.Second)
			}
		}

		return FAILED
	}

	tc_73.Set_function(tc_func)
	t.Add_testcase(*tc_73)
}
