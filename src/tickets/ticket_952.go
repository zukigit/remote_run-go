package tickets

import (
	"fmt"
	"os"

	"github.com/zukigit/remote_run-go/src/common"
	"github.com/zukigit/remote_run-go/src/dao"
	"github.com/zukigit/remote_run-go/src/lib"
)

type Ticket_952 struct {
	Ticket_no                                   uint
	Ticket_description                          string
	PASSED_count, FAILED_count, MUSTCHECK_count int
	Testcases                                   []dao.TestCase
}

func (t *Ticket_952) New_testcase(testcase_id uint, testcase_description string) *dao.TestCase {
	return dao.New_testcase(testcase_id, testcase_description)
}

func (t *Ticket_952) Get_no() uint {
	return t.Ticket_no
}

func (t *Ticket_952) Set_PASSED_count(passed_count int) {
	t.PASSED_count = 43
}

func (t *Ticket_952) Set_FAILED_count(failed_count int) {
	t.FAILED_count = failed_count
}

func (t *Ticket_952) Set_MUSTCHECK_count(mustcheck_count int) {
	t.MUSTCHECK_count = mustcheck_count
}

func (t *Ticket_952) Get_dsctn() string {
	return t.Ticket_description
}

func (t *Ticket_952) Add_testcase(tc dao.TestCase) {
	t.Testcases = append(t.Testcases, tc)
}

func (t *Ticket_952) Get_testcases() []dao.TestCase {
	return t.Testcases
}

func (t *Ticket_952) Set_values() {
	t.Ticket_no = 952
	t.Ticket_description = " Commands are executed  if STD_OUT option is checked in Window agent. (Window agent only)"
}
func (t *Ticket_952) Add_testcases() {
	tc_1 := t.New_testcase(1, "Check if commands are executed when STD_OUT is checked.")
	tc_func := func() common.Testcase_status {
		status := t.commonTask(tc_1, "ticket_952a")
		if status != PASSED {
			return status
		}

		if status == PASSED {
			job := "Icon_1"
			status = t.runTestCase(tc_1, job)
		}
		return status
	}

	tc_1.Set_function(tc_func)
	t.Add_testcase(*tc_1)

	tc_2 := t.New_testcase(2, "Check if commands are executed when STD_OUT is checked and used.")
	tc_func2 := func() common.Testcase_status {
		status := t.commonTask(tc_2, "ticket_952b")
		if status != PASSED {
			return status
		}

		if status == PASSED {
			job := "Icon_1"
			status = t.runTestCase(tc_2, job)
		}
		return status
	}

	tc_2.Set_function(tc_func2)
	t.Add_testcase(*tc_2)

	tc_3 := t.New_testcase(3, "Abnormal case check")
	tc_func3 := func() common.Testcase_status {
		status := t.commonTask(tc_3, "ticket_952c")
		if status != PASSED {
			return status
		}

		if status == PASSED {
			job := "Icon_1"
			status = t.runTestCase(tc_3, job)
		}
		return status
	}

	tc_3.Set_function(tc_func3)
	t.Add_testcase(*tc_3)
}

func (t *Ticket_952) logError(tc *dao.TestCase, format string, args ...interface{}) common.Testcase_status {
	fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, format, args...))
	return FAILED
}

func (t *Ticket_952) runTestCase(tc *dao.TestCase, job string) common.Testcase_status {
	const defaultTestFolderPath = `C:\test_folder` //change default

	envs, err := lib.Get_str_str_map("JA_HOSTNAME", "oss.windows", "JA_CMD", fmt.Sprintf("echo mkdir %s", defaultTestFolderPath))
	if err != nil {
		return t.logError(tc, "Error retrieving environment variables: %s", err)
	}

	run_jobnet_id, err := lib.Jobarg_exec_E(job, envs)
	if err != nil {
		return t.logError(tc, "Error executing job %s: %s", job, err)
	}

	jobExecMessage := fmt.Sprintf("Executed job: %s with run_jobnet_id: %s", job, run_jobnet_id)
	fmt.Println(jobExecMessage)
	lib.Logi(common.LOG_LEVEL_INFO, jobExecMessage)

	jobnet_run_info, err := lib.Jobarg_get_jobnet_run_info(run_jobnet_id)
	if err != nil {
		return t.logError(tc, "Error retrieving run info for job %s: %s", job, err)
	}

	fmt.Printf("Jobnet Run Info: %+v\n", jobnet_run_info)
	lib.Logi(common.LOG_LEVEL_INFO, fmt.Sprintf("Jobnet Run Info: %+v", jobnet_run_info))

	if jobnet_run_info.Jobnet_status == "END" && jobnet_run_info.Job_status == "NORMAL" {

		_, err := os.Stat(defaultTestFolderPath)
		if err == nil {

			logMessage := fmt.Sprintf("The folder %s already exists, deleting it now.", defaultTestFolderPath)
			fmt.Println(logMessage)
			lib.Logi(common.LOG_LEVEL_INFO, logMessage)

			err = os.RemoveAll(defaultTestFolderPath)
			if err != nil {
				errorMessage := fmt.Sprintf("Error deleting folder %s: %s", defaultTestFolderPath, err)
				fmt.Println(errorMessage)
				t.logError(tc, errorMessage)
				return FAILED
			}
			successMessage := fmt.Sprintf("Successfully deleted the folder %s.", defaultTestFolderPath)
			fmt.Println(successMessage)
			lib.Logi(common.LOG_LEVEL_INFO, successMessage)
		} else if os.IsNotExist(err) {
			successMessage := fmt.Sprintf("Folder %s does not exist, which is as expected. Test Passed.", defaultTestFolderPath)
			fmt.Println(successMessage)
			lib.Logi(common.LOG_LEVEL_INFO, successMessage)
		} else {
			errorMessage := fmt.Sprintf("Error checking folder %s: %s", defaultTestFolderPath, err)
			fmt.Println(errorMessage)
			t.logError(tc, errorMessage)
			return FAILED
		}

		successMessage := fmt.Sprintf("%s completed successfully.", job)
		fmt.Println(successMessage)
		lib.Logi(common.LOG_LEVEL_INFO, successMessage)
		return PASSED
	} else {

		failMessage := fmt.Sprintf("%s failed. Jobnet_status: %s, Job_status: %s, Exit_cd: %d", job, jobnet_run_info.Jobnet_status, jobnet_run_info.Job_status, jobnet_run_info.Exit_cd)
		fmt.Println(failMessage)
		lib.Logi(common.LOG_LEVEL_INFO, failMessage)
		return FAILED
	}
}
func (t *Ticket_952) commonTask(tc *dao.TestCase, ticket string) common.Testcase_status {

	if err := lib.Jobarg_cleanup_windows(); err != nil {
		return t.logError(tc, "Error during cleanup: %s", err)
	}
	switch ticket {
	case "ticket_952a":
		if err := lib.Jobarg_enable_jobnet("Icon_1", "ticket_952a"); err != nil {
			return t.logError(tc, "Failed to enable jobnet, Error: %s", err)
		}
	case "ticket_952b":
		if err := lib.Jobarg_enable_jobnet("Icon_1", "ticket_952b"); err != nil {
			return t.logError(tc, "Failed to enable jobnet, Error: %s", err)
		}
	case "ticket_952c":
		if err := lib.Jobarg_enable_jobnet("Icon_1", "ticket_952c"); err != nil {
			return t.logError(tc, "Failed to enable jobnet, Error: %s", err)
		}
	default:
		return t.logError(tc, "Invalid ticket provided: %s", ticket)
	}

	return PASSED
}
