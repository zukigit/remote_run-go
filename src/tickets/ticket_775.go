package tickets

import (
	"fmt"
	"strings"

	"github.com/zukigit/remote_run-go/src/common"
	"github.com/zukigit/remote_run-go/src/dao"
	"github.com/zukigit/remote_run-go/src/lib"
)

type Ticket_775 struct {
	no          uint
	description string
	testcases   []dao.TestCase
}

func (t *Ticket_775) New_testcase(testcase_id uint, testcase_description string) *dao.TestCase {
	return dao.New_testcase(testcase_id, testcase_description)
}

func (t *Ticket_775) Get_no() uint {
	return t.no
}

func (t *Ticket_775) Get_dsctn() string {
	return t.description
}

func (t *Ticket_775) Add_testcase(tc dao.TestCase) {
	t.testcases = append(t.testcases, tc)
}

func (t *Ticket_775) Get_testcases() []dao.TestCase {
	return t.testcases
}

// Enter your ticket information here
func (t *Ticket_775) Set_values() {
	t.no = 775 // Enter your ticket id
	t.description = "Agent crash on executing 24-byte length of characters"
}

// Add your test case here
func (t *Ticket_775) Add_testcases() {
	// TESTCASE 001

	charaterString24 := "1233ffweldfdkl12sderfdes"
	charaterString40 := "1233ffweldfdkl12sderfdes1233ffweldfdkl12"
	charaterString72 := "1233ffweldfdkl12sderfdes1233ffweldfdkl121233ffweldfdkl12sderfdes1233ffwe"
	charaterString136 := "1233ffweldfdkl12sderfdes1233ffweldfdkl121233ffweldfdkl12sderfdes1233ffwe1233ffweldfdkl12sderfdes1233ffweldfdkl121233ffweldfdkl12sderfdes"
	charaterString264 := "1233ffweldfdkl12sderfdes1233ffweldfdkl121233ffweldfdkl12sderfdes1233ffwe1233ffweldfdkl12sderfdes1233ffweldfdkl121233ffweldfdkl12sderfdes1233ffweldfdkl12sderfdes1233ffweldfdkl121233ffweldfdkl12sderfdes1233ffwe1233ffweldfdkl12sderfdes1233ffweldfdkl121233ffweldfdkl12"

	addTestCase := func(id int, description string, jobnetId string, jobname string, charcommand string) {
		tc := t.New_testcase(uint(id), description) // Convert id to uint
		tc_func := func() common.Testcase_status {

			if err := lib.Jobarg_cleanup_linux(); err != nil {
				tc.Err_log("Error during cleanup: %s", err)
				return FAILED
			}

			if err := lib.Jobarg_enable_jobnet(jobnetId, jobname); err != nil {
				tc.Err_log("Failed to enable jobnet, Error: %s", err)
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

func executejobnet(jobnetId string, charcommand string, tc *dao.TestCase) common.Testcase_status {

	envs, _ := lib.Get_str_str_map("JA_HOSTNAME", "oss.linux", "JA_CMD", charcommand)

	run_jobnet_id, error := lib.Jobarg_exec_E(jobnetId, envs)

	if error != nil {
		fmt.Println(tc.Err_log("Error: %s, std_out: %s", error.Error(), run_jobnet_id))
		return FAILED
	}

	fmt.Println(tc.Info_log("%s has been successfully run with registry number: %s", jobnetId, run_jobnet_id))

	jobnet_run_info, error := lib.Jobarg_get_jobnet_run_info(run_jobnet_id)

	if error != nil {
		tc.Err_log("Error: %s", error.Error())
		return FAILED
	}
	output_command := jobnet_run_info.Std_error

	if strings.Contains(output_command, charcommand) {
		fmt.Println(tc.Info_log("Characters are match in STD_ERR"))
		return PASSED
	} else {
		tc.Err_log("Error: %s", "Characters are not match in STD_ERR")
		return FAILED
	}

}
