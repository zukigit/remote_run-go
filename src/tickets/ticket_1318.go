package tickets

import (
	"fmt"
	"zukigit/remote_run-go/src/dao"
	"zukigit/remote_run-go/src/lib"
)

type Ticket_1318 struct {
	no          uint
	description string
	auth        *dao.Auth
	testcases   []dao.TestCase
}

func (t *Ticket_1318) Set_values(auth *dao.Auth) {
	t.no = 1318
	t.description = "Fixed for negative JOB_EXT_CD return value."
	t.auth = auth
}

func (t *Ticket_1318) New_testcase(testcase_id uint, testcase_description string) *dao.TestCase {
	return dao.New_testcase(testcase_id, testcase_description, t.auth)
}

func (t *Ticket_1318) Get_no() uint {
	return t.no
}

func (t *Ticket_1318) Get_dsctn() string {
	return t.description
}

func (t *Ticket_1318) Set_testcase(tc dao.TestCase) {
	t.testcases = append(t.testcases, tc)
}

func (t *Ticket_1318) Get_testcases() []dao.TestCase {
	return t.testcases
}

func (t *Ticket_1318) Run() {
	for _, tc := range t.testcases {
		fmt.Println(tc.Info_log("running..."))
		if !tc.Is_function_nil() {
			tc.Set_status(tc.Run_function())
		} else {
			fmt.Println(tc.Err_log("has no function. SKIPPED!"))
			tc.Set_status(FAILED)
		}

		fmt.Println(tc.Info_log("finished!"))
	}
}

func (t *Ticket_1318) Add_testcases() {
	// Add your test case here
	// ticket 168
	// tc_170 := t.New_testcase(170, "Default Case Check. Linux AGENT")
	// tc_func := func() dao.Testcase_status {
	// 	std_out, error := tc_170.Jobarg_exec("TICKET_1318_linux")
	// 	if error != nil {
	// 		tc_170.Err_log("Error: %s, std_out: %s", error.Error(), std_out)
	// 		return FAILED
	// 	}

	// 	jobnet_info, error := tc_170.Jobarg_get_jobnet_run_info(std_out)
	// 	if error != nil {
	// 		tc_170.Err_log("Error: %s", error.Error())
	// 		return FAILED
	// 	}

	// 	if jobnet_info.Jobnet_status == "END" && jobnet_info.Job_status == "NORMAL" && jobnet_info.Exit_cd == 96 {
	// 		return PASSED
	// 	}

	// 	tc_170.Err_log("Jobnet_status: %s, Job_status: %s, Exit_cd: %d", jobnet_info.Jobnet_status, jobnet_info.Job_status, jobnet_info.Exit_cd)
	// 	return FAILED
	// }
	// tc_170.Set_function(tc_func)
	// t.Set_testcase(*tc_170)

	// ticket 169
	tc_169 := t.New_testcase(169, "Abnormal Case") // create test case
	tc_func := func() dao.Testcase_status {
		err := lib.Restart_jaz_agent_windows()
		if err != nil {
			tc_169.Err_log("Failed to restart windows service. Error: %s", err.Error())
			return FAILED
		}

		err = lib.Stop_jaz_agent_windows()
		if err != nil {
			tc_169.Err_log("Failed to stop windows service. Error: %s", err.Error())
			return FAILED
		}

		lib.Jaz_set_agent_config_windows("ExtUnsignedFlag", "1")

		return PASSED
	}
	tc_169.Set_function(tc_func)
	t.Set_testcase(*tc_169) // Add testcase to ticket
}
