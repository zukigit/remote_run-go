package tickets

import (
	"fmt"

	"github.com/zukigit/remote_run-go/src/common"
	"github.com/zukigit/remote_run-go/src/dao"
	"github.com/zukigit/remote_run-go/src/lib"
)

type Ticket_794 struct {
	no          uint
	description string
	testcases   []dao.TestCase
}

func (t *Ticket_794) New_testcase(testcase_id uint, testcase_description string) *dao.TestCase {
	return dao.New_testcase(testcase_id, testcase_description)
}

func (t *Ticket_794) Get_no() uint {
	return t.no
}

func (t *Ticket_794) Get_dsctn() string {
	return t.description
}

func (t *Ticket_794) Add_testcase(tc dao.TestCase) {
	t.testcases = append(t.testcases, tc)
}

func (t *Ticket_794) Get_testcases() []dao.TestCase {
	return t.testcases
}

// Enter your ticket information here
func (t *Ticket_794) Set_values() {
	t.no = 794 // Enter your ticket id
	t.description = "Process Check"
}
func (t *Ticket_794) Add_testcases() {

	tc_38 := t.New_testcase(38, "Normal job executions performance test")
	tc_func := func() common.Testcase_status {

		if err := lib.Jobarg_cleanup_linux(); err != nil {
			return t.logError(tc_38, "Error during cleanup: %s", err)
		}

		if err := lib.Jobarg_enable_jobnet("Icon_1", "jobicon_linux"); err != nil {
			return t.logError(tc_38, "Error during enable jobnet: %s", err)
		}

		// Define job names
		jobs := []string{
			"Icon_510",
			"Icon_1020",
			"Icon_2040",
		}

		for _, job := range jobs {
			if status := t.runJob(tc_38, job); status != PASSED {
				return status
			}
		}

		return PASSED
	}

	tc_38.Set_function(tc_func)
	t.Add_testcase(*tc_38)
}

func (t *Ticket_794) runJob(tc *dao.TestCase, job string) common.Testcase_status {
	envs, _ := lib.Get_str_str_map("JA_HOSTNAME", "oss.linux", "JA_CMD", "sleep 10")
	run_jobnet_id, err := lib.Jobarg_exec_E(job, envs)
	if err != nil {
		return t.logError(tc, "Error executing job %s: %s", job, err)
	}

	jobnet_run_info, err := lib.Jobarg_get_jobnet_run_info(run_jobnet_id)
	if err != nil {
		return t.logError(tc, "Error retrieving run info for job %s: %s", job, err)
	}

	if jobnet_run_info.Jobnet_status == "END" && jobnet_run_info.Job_status == "NORMAL" && jobnet_run_info.Exit_cd == 0 {
		fmt.Println(tc.Info_log("%s completed successfully.", job))
		return PASSED
	}

	return t.logError(tc, "%s failed. Jobnet_status: %s, Job_status: %s, Exit_cd: %d", job, jobnet_run_info.Jobnet_status, jobnet_run_info.Job_status, jobnet_run_info.Exit_cd)
}

func (t *Ticket_794) logError(tc *dao.TestCase, format string, args ...interface{}) common.Testcase_status {
	fmt.Println(tc.Err_log(format, args...))
	return FAILED
}
