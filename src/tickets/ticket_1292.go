package tickets

import (
	"fmt"

	"github.com/zukigit/remote_run-go/src/common"
	"github.com/zukigit/remote_run-go/src/dao"
	"github.com/zukigit/remote_run-go/src/lib"
)

type Ticket_1292 struct {
	Ticket_no                                   uint
	Ticket_description                          string
	PASSED_count, FAILED_count, MUSTCHECK_count int
	Testcases                                   []dao.TestCase
}

func (t *Ticket_1292) New_testcase(testcase_id uint, testcase_description string) *dao.TestCase {
	return dao.New_testcase(testcase_id, testcase_description)
}

func (t *Ticket_1292) Get_no() uint {
	return t.Ticket_no
}

func (t *Ticket_1292) Set_PASSED_count(passed_count int) {
	t.PASSED_count = passed_count
}

func (t *Ticket_1292) Set_FAILED_count(failed_count int) {
	t.FAILED_count = failed_count
}

func (t *Ticket_1292) Set_MUSTCHECK_count(mustcheck_count int) {
	t.MUSTCHECK_count = mustcheck_count
}

func (t *Ticket_1292) Get_dsctn() string {
	return t.Ticket_description
}

func (t *Ticket_1292) Add_testcase(tc dao.TestCase) {
	t.Testcases = append(t.Testcases, tc)
}

func (t *Ticket_1292) Get_testcases() []dao.TestCase {
	return t.Testcases
}

// Enter your ticket information here
func (t *Ticket_1292) Set_values() {
	t.Ticket_no = 1292 // Enter your ticket id
	t.Ticket_description = "Issue 1292 - Ja_purge is not working when the jobnet is aborted and it has job icons in ready state."
}

// Add your test case here
func (t *Ticket_1292) Add_testcases() {
	// TESTCASE 127

	// Issue 1292 - Ja_purge is not working when the jobnet is aborted and it has job icons in ready state.
	//
	// Default Case Check
	//
	// - Create a jobnet that has parallel 100 jobnets.
	// - Set keep-span to 5mins.
	// - Run that jobnet.
	// - After the jobnet exectuion finish, run the following query. UPDATE zabbix.ja_run_job_table SET status=1;
	//
	// - After 5 mins of the jobnet exectuion finish, run the following query. SELECT * FROM ja_run_jobnet_table WHERE inner_jobnet_id = YOUR_INNER_JOBNET_ID_FROM_MANAGEMENT_SCREEN;
	//
	// - There must be no records left.

	tc_127 := t.New_testcase(127, "Default Case Check")
	tc_func := func() common.Testcase_status {

		var jobnet_id string = "Icon_1"                   // This value must be Jobnet_ID that you want to enable
		var jobnet_name string = "TICKET1292_TESTCASE127" // This value must be Jobnet_Name that you want to enable.
		var result bool = true                            // Default value, Do not Change.
		var jobnet_manage_id string

		// 1. Run Jobarg_cleanup and enable jobnets
		// 2. Run Sql Script to set Jobnet_Keep_Span value to 5.
		// 3. Restart JAZ Server
		// 4. Run Jobnet 100 icons and retrieve Jobnet Run Info
		// 5. Run Sql Script to set ja_run_job_table set status = 1
		// 6. Wait for 5 minutes
		// 7. Check database for the empty set for the Ja Purge to work.

		if lib.Run_Jobarg_cleanup_linux() &&
			lib.Run_enable_jobnet(jobnet_id, jobnet_name) &&
			lib.Run_enable_jobnet("Icon_10", "Icon_10") &&
			lib.Run_enable_jobnet("Icon_100", "Icon_100") &&
			lib.Run_Sql_Script("UPDATE ja_parameter_table SET value = 5 WHERE parameter_name = 'JOBNET_KEEP_SPAN';") &&
			lib.Run_Restart_Linux_Jaz_server() &&
			func() bool {
				result, jobnet_manage_id = lib.Run_Jobnet("Icon_100")
				return result
			}() &&
			func() bool {
				result, _ = lib.Run_Jobarg_get_jobnet_run_info(jobnet_manage_id)
				return result
			}() &&
			lib.Run_Sql_Script("UPDATE ja_run_job_table SET status = 1;") &&
			lib.Run_Timeout(5*60) &&
			func() bool {
				var count int
				_, sql_result := lib.Run_Sql_Script_Return_Rows("SELECT count(*) FROM ja_run_jobnet_table WHERE inner_jobnet_id = '" + jobnet_manage_id + "';")
				if sql_result.Next() { // Move to the first row
					if err := sql_result.Scan(&count); err != nil {
						fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error: Error scanning result: %s", err))
						return false
					}
				}
				if count > 0 {
					fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error: Database is not empty!!!"))
					return false
				}
				fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "Info: Database is empty."))
				return true
			}() {
			fmt.Println("All operations completed successfully")
		}
		fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "Info: Resting config files back to normal parameter"))
		if lib.Run_Sql_Script("UPDATE ja_parameter_table SET value = 60 WHERE parameter_name = 'JOBNET_KEEP_SPAN';") &&
			lib.Run_Jobarg_cleanup_linux() {
			if result {
				return PASSED
			}
			return FAILED
		}
		fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "Error: Failed at restting config files back to normal settings. Please reset the config files manually."))
		return MUST_CHECK
	}
	tc_127.Set_function(tc_func)
	t.Add_testcase(*tc_127)
}
