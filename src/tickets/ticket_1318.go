package tickets

import (
	"github.com/zukigit/remote_run-go/src/common"
	"github.com/zukigit/remote_run-go/src/lib"
)

type Ticket_1318 struct {
	Ticket_no                                   uint
	Ticket_description                          string
	PASSED_count, FAILED_count, MUSTCHECK_count int
	Testcases                                   []common.TestCase
}

func (t *Ticket_1318) New_testcase(testcase_id uint, testcase_description string) *common.TestCase {
	return common.New_testcase(testcase_id, testcase_description)
}

func (t *Ticket_1318) Get_no() uint {
	return t.Ticket_no
}

func (t *Ticket_1318) Set_PASSED_count(passed_count int) {
	t.PASSED_count = passed_count
}

func (t *Ticket_1318) Set_FAILED_count(failed_count int) {
	t.FAILED_count = failed_count
}

func (t *Ticket_1318) Set_MUSTCHECK_count(mustcheck_count int) {
	t.MUSTCHECK_count = mustcheck_count
}

func (t *Ticket_1318) Get_dsctn() string {
	return t.Ticket_description
}

func (t *Ticket_1318) Add_testcase(tc common.TestCase) {
	t.Testcases = append(t.Testcases, tc)
}

func (t *Ticket_1318) Get_testcases() []common.TestCase {
	return t.Testcases
}

func (t *Ticket_1318) Set_values() {
	t.Ticket_no = 1318
	t.Ticket_description = "Fixed for negative JOB_EXT_CD return value."
}

// Add your testcase here
func (t *Ticket_1318) Add_testcases() {
	// TESTCASE 168
	tc_168 := t.New_testcase(168, "Normal Case with ExtUnsignedFlag=0, Windows agent") // create test case
	tc_func := func() common.Testcase_status {

		// Set joabrg agent config value
		err := lib.Ja_set_agent_config_windows("ExtUnsignedFlag", "0")
		if err != nil {
			lib.Logi(common.LOG_LEVEL_ERR, "Failed to set joabrg agent config value. Error: %s", err.Error())
			return FAILED
		}

		// Restart jobarg-agentd
		err = lib.Restart_jaz_agent_windows()
		if err != nil {
			lib.Logi(common.LOG_LEVEL_ERR, "Failed to restart windows service. Error: %s", err.Error())
			return FAILED
		}

		// Enable jobnet
		if err := lib.Jobarg_enable_jobnet("Icon_1", "jobicon_windows"); err != nil {
			lib.Logi(common.LOG_LEVEL_ERR, "Failed to enable jobnet, Error: %s", err)
			return FAILED
		}

		envs, _ := lib.Get_str_str_map("JA_HOSTNAME", "oss.windows", "JA_CMD", "exit -100000;")

		// Run jobnet
		run_jobnet_id, error := lib.Jobarg_exec_E("Icon_1", envs)
		if error != nil {
			lib.Logi(common.LOG_LEVEL_ERR, "Error: %s, std_out: %s", error.Error(), run_jobnet_id)
			return FAILED
		}

		// Wait jobnet finishes and get jobnet run info.
		jobnet_run_info, error := lib.Jobarg_get_jobnet_run_info(run_jobnet_id)
		if error != nil {
			lib.Logi(common.LOG_LEVEL_ERR, "Error: %s", error.Error())
			return FAILED
		}

		// Check jobnet run status and exit code.
		if jobnet_run_info.Jobnet_status == "END" && jobnet_run_info.Job_status == "NORMAL" && jobnet_run_info.Exit_cd == 4294867296 {
			return PASSED
		}

		lib.Logi(common.LOG_LEVEL_ERR, "Jobnet_status: %s, Job_status: %s, Exit_cd: %d", jobnet_run_info.Jobnet_status, jobnet_run_info.Job_status, jobnet_run_info.Exit_cd)
		return FAILED
	}
	tc_168.Set_function(tc_func)
	t.Add_testcase(*tc_168) // Add testcase to ticket

	// TESTCASE 169
	tc_169 := t.New_testcase(169, "Normal Case with ExtUnsignedFlag=1, Windows agent") // create test case
	tc_func = func() common.Testcase_status {

		// Set joabrg agent config value
		tc_169.Add_doc(common.PRE_OPT, "Change jobarg-agent's config ExtUnsignedFlag=1")
		err := lib.Ja_set_agent_config_windows("ExtUnsignedFlag", "1")
		if err != nil {
			lib.Logi(common.LOG_LEVEL_ERR, "Failed to set joabrg agent config value. Error: %s", err.Error())
			return FAILED
		}

		// Restart jobarg-agentd
		tc_169.Add_doc(common.PRE_OPT, "Restart jobarg-agentd service")
		err = lib.Restart_jaz_agent_windows()
		if err != nil {
			lib.Logi(common.LOG_LEVEL_ERR, "Failed to restart windows service. Error: %s", err.Error())
			return FAILED
		}

		// Enable jobnet
		tc_169.Add_doc(common.PRE_OPT, "Enable jobnet for windows agent.")
		if err := lib.Jobarg_enable_jobnet("Icon_1", "jobicon_windows"); err != nil {
			lib.Logi(common.LOG_LEVEL_ERR, "Failed to enable jobnet, Error: %s", err)
			return FAILED
		}

		envs, _ := lib.Get_str_str_map("JA_HOSTNAME", "oss.windows", "JA_CMD", "exit -100000;")

		// Run jobnet
		tc_169.Add_doc(common.PRE_OPT, "Run jobnet. JA_HOSTNAME: oss.windows, JA_CMD: exit -100000;.")
		run_jobnet_id, error := lib.Jobarg_exec_E("Icon_1", envs)
		if error != nil {
			lib.Logi(common.LOG_LEVEL_ERR, "Error: %s, std_out: %s", error.Error(), run_jobnet_id)
			return FAILED
		}

		// Wait jobnet finishes and get jobnet run info.
		jobnet_run_info, error := lib.Jobarg_get_jobnet_run_info(run_jobnet_id)
		if error != nil {
			lib.Logi(common.LOG_LEVEL_ERR, "Error: %s", error.Error())
			return FAILED
		}

		// Check jobnet run status and exit code.
		tc_169.Add_doc(common.PRE_OPT, "JObnet status must be END, job status must be normal and exitcode must be 4294867296.")
		if jobnet_run_info.Jobnet_status == "END" && jobnet_run_info.Job_status == "NORMAL" && jobnet_run_info.Exit_cd == 4294867296 {
			return PASSED
		}

		lib.Logi(common.LOG_LEVEL_ERR, "Jobnet_status: %s, Job_status: %s, Exit_cd: %d", jobnet_run_info.Jobnet_status, jobnet_run_info.Job_status, jobnet_run_info.Exit_cd)
		return FAILED
	}
	tc_169.Set_function(tc_func)
	t.Add_testcase(*tc_169) // Add testcase to ticket

	// TESTCASE 170
	tc_170 := t.New_testcase(170, "Default Case Check. Linux AGENT")
	tc_func = func() common.Testcase_status {
		err := lib.Ja_set_agent_config_linux("ExtUnsignedFlag", "0")
		if err != nil {
			lib.Logi(common.LOG_LEVEL_ERR, "Error: %s", err.Error())
			return FAILED
		}

		err = lib.Restart_jaz_agent_linux()
		if err != nil {
			lib.Logi(common.LOG_LEVEL_ERR, "Error: %s", err.Error())
			return FAILED
		}

		// Enable jobnet
		if err := lib.Jobarg_enable_jobnet("Icon_1", "jobicon_linux"); err != nil {
			lib.Logi(common.LOG_LEVEL_ERR, "Failed to enable jobnet, Error: %s", err)
			return FAILED
		}

		envs, _ := lib.Get_str_str_map("JA_HOSTNAME", "oss.linux", "JA_CMD", "exit -100000")

		run_jobnet_id, err := lib.Jobarg_exec_E("Icon_1", envs)
		if err != nil {
			lib.Logi(common.LOG_LEVEL_ERR, "Error: %s, std_out: %s", err.Error(), run_jobnet_id)
			return FAILED
		}

		jobnet_run_info, err := lib.Jobarg_get_jobnet_run_info(run_jobnet_id)
		if err != nil {
			lib.Logi(common.LOG_LEVEL_ERR, "Error: %s", err.Error())
			return FAILED
		}

		if jobnet_run_info.Jobnet_status == "END" && jobnet_run_info.Job_status == "NORMAL" && jobnet_run_info.Exit_cd == 96 {
			return PASSED
		}

		lib.Logi(common.LOG_LEVEL_ERR, "Jobnet_status: %s, Job_status: %s, Exit_cd: %d", jobnet_run_info.Jobnet_status, jobnet_run_info.Job_status, jobnet_run_info.Exit_cd)
		return FAILED
	}
	tc_170.Set_function(tc_func)
	t.Add_testcase(*tc_170)
}
