package tickets

import (
	"fmt"
	"strings"

	"github.com/zukigit/remote_run-go/src/common"
	"github.com/zukigit/remote_run-go/src/lib"
)

type Ticket_775 struct {
	Ticket_no                                   uint
	Ticket_description                          string
	PASSED_count, FAILED_count, MUSTCHECK_count int
	Testcases                                   []common.TestCase
}

func (t *Ticket_775) New_testcase(testcase_id uint, testcase_description string) *common.TestCase {
	return common.New_testcase(testcase_id, testcase_description)
}

func (t *Ticket_775) Get_no() uint {
	return t.Ticket_no
}

func (t *Ticket_775) Set_PASSED_count(passed_count int) {
	t.PASSED_count = passed_count
}

func (t *Ticket_775) Set_FAILED_count(failed_count int) {
	t.FAILED_count = failed_count
}

func (t *Ticket_775) Set_MUSTCHECK_count(mustcheck_count int) {
	t.MUSTCHECK_count = mustcheck_count
}

func (t *Ticket_775) Get_dsctn() string {
	return t.Ticket_description
}

func (t *Ticket_775) Add_testcase(tc common.TestCase) {
	t.Testcases = append(t.Testcases, tc)
}

func (t *Ticket_775) Get_testcases() []common.TestCase {
	return t.Testcases
}

// Enter your ticket information here
func (t *Ticket_775) Set_values() {
	t.Ticket_no = 775 // Enter your ticket id
	t.Ticket_description = "Agent crash on executing 24-byte length of characters"
}

// Add your test case here
func (t *Ticket_775) Add_testcases() {
	// TESTCASE 001

	charaterString24 := "abc123$%&こんにちはＨｅｌｌｏ*&+_1"
	charaterString40 := "abc123$%&こんにちはＨｅｌｌｏ*&+_1abc123こんにちはＨｅｌｌｏ"
	charaterString72 := "abc123$%&こんにちはＨｅｌｌｏ*&+_1abc123こんにちはＨｅｌｌｏabc123$%&こんにちはＨｅｌｌｏ*&+こんにちはＨｅｌｌｏ"
	charaterString136 := "abc123$%&こんにちはＨｅｌｌｏ*&+_1abc123こんにちはＨｅｌｌｏabc123$%&こんにちはＨｅｌｌｏ*&+こんにちはＨｅｌｌｏabc123$%&こんにちはＨｅｌｌｏ*&+_1abc123こんにちはＨｅｌｌｏabc123$%&こんにちはＨｅｌｌｏ*&+_1"
	charaterString264 := "abc123$%&こんにちはＨｅｌｌｏ*&+_1abc123こんにちはＨｅｌｌｏabc123$%&こんにちはＨｅｌｌｏ*&+こんにちはＨｅｌｌｏabc123$%&こんにちはＨｅｌｌｏ*&+こんにちはＨｅｌｌｏabc123$%&こんにちはＨｅｌｌｏ*&+_1abc123$%&こんにちはＨｅｌｌｏ*&+_1abc123こんにちはＨｅｌｌｏabc123$%&こんにちはＨｅｌｌｏ*&+こんにちはＨｅｌｌｏabc123$%&こんにちはＨｅｌｌｏ*&+_1abc123こんにちはＨｅｌｌｏabc123$%&こんにちはＨｅｌｌｏ*&+_1"

	addTestCase := func(id int, description string, jobnetId string, jobname string, charcommand string) {
		tc := t.New_testcase(uint(id), description) // Convert id to uint
		tc_func := func() common.Testcase_status {

			if err := lib.Jobarg_cleanup_linux(); err != nil {
				lib.Logi(common.LOG_LEVEL_ERR, "Error during cleanup: %s", err)
				return FAILED
			}

			if err := lib.Jobarg_enable_jobnet(jobnetId, jobname); err != nil {
				lib.Logi(common.LOG_LEVEL_ERR, "Failed to enable jobnet, Error: %s", err)
				return FAILED
			}
			return executejobnet(jobnetId, charcommand, tc)
		}
		tc.Set_function(tc_func)
		t.Add_testcase(*tc)
	}

	addTestCase(1, "execute command with 24 characters", "Icon_1", "jobicon_linux", charaterString24)
	addTestCase(2, "execute command with 40 characters", "Icon_1", "jobicon_linux", charaterString40)
	addTestCase(3, "execute command with 72 characters", "Icon_1", "jobicon_linux", charaterString72)
	addTestCase(4, "execute command with 136 characters", "Icon_1", "jobicon_linux", charaterString136)
	addTestCase(5, "execute command with 264 characters", "Icon_1", "jobicon_linux", charaterString264)
}

func executejobnet(jobnetId string, charcommand string, tc *common.TestCase) common.Testcase_status {

	runcharcommand := fmt.Sprintf("echo %s", charcommand)

	envs, _ := lib.Get_str_str_map("JA_HOSTNAME", "oss.linux", "JA_CMD", runcharcommand)

	run_jobnet_id, error := lib.Jobarg_exec_E(jobnetId, envs)

	if error != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error: %s, std_out: %s", error.Error(), run_jobnet_id))
		return FAILED
	}

	fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "%s has been successfully run with registry number: %s", jobnetId, run_jobnet_id))

	jobnet_run_info, error := lib.Jobarg_get_jobnet_run_info(run_jobnet_id)

	if error != nil {
		lib.Logi(common.LOG_LEVEL_ERR, "Error: %s", error.Error())
		return FAILED
	}
	output_command := jobnet_run_info.Std_out

	if strings.Contains(output_command, charcommand) {
		fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "Characters are match in STD_ERR"))
		fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "Original Text: %q", charcommand))
		fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "Job Output: %q", output_command))
		return PASSED
	} else {
		lib.Logi(common.LOG_LEVEL_ERR, "Error: %s", "Characters are not match in STD_ERR")
		return FAILED
	}

}
