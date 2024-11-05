package tickets

import (
	"fmt"
	"strconv"
	"time"

	"github.com/zukigit/remote_run-go/src/common"
	"github.com/zukigit/remote_run-go/src/dao"
	"github.com/zukigit/remote_run-go/src/lib"
)

type Ticket_1091 struct {
	no          uint
	description string
	testcases   []dao.TestCase
}

func (t *Ticket_1091) New_testcase(testcase_id uint, testcase_description string) *dao.TestCase {
	return dao.New_testcase(testcase_id, testcase_description)
}

func (t *Ticket_1091) Get_no() uint {
	return t.no
}

func (t *Ticket_1091) Get_dsctn() string {
	return t.description
}

func (t *Ticket_1091) Add_testcase(tc dao.TestCase) {
	t.testcases = append(t.testcases, tc)
}

func (t *Ticket_1091) Get_testcases() []dao.TestCase {
	return t.testcases
}

// Enter your ticket information here
func (t *Ticket_1091) Set_values() {
	t.no = 1091 // Enter your ticket id
	t.description = "Issue 1091 - Jobarg-server crash when it's restarted while it has running jobs"
}

// Add your test case here
func (t *Ticket_1091) Add_testcases() {
	// TESTCASE 99

	tc_99 := t.New_testcase(99, "Normal Case Check")
	tc_func := func() common.Testcase_status {

		// Issue 1091 - Jobarg-server crash when it's restarted while it has running jobs
		//
		// Normal Case Check
		//
		//  - Create parallel jobs with different hosts as many as you can use. Command -> sleep 120
		//
		//  - Run the jobnet.
		//  - Make sure all job icons are running.
		//  - Restart the jobarg-server
		//
		//  - The following log must be printed and no CRIT logs and no error logs should be printed.
		//  - "[INFO] In ja_send_ipchange_request(), Server Ip is up to date in host_name."

		var jobnet_id string = "Icon_1"                      // This value must be Jobnet_ID that you want to run.
		var jobnet_name string = "jobicon_linux_and_windows" // This value must be Jobnet_Name that you want to run.
		var linux_hostname string = "oss.linux"              // This value must be hostname of linux agent you want to run. Default is oss.linux
		var window_hostname string = "DESKTOP-8KPOJCT"       // This value must be hostname of window agent you want to run. Default is oss.linux2
		var timeout_time int = 60                            // This is timeout value in seconds for waiting for the log to appear.
		var result bool

		envs, err := lib.Get_str_str_map("JA_HOSTNAME1", linux_hostname, "JA_HOSTNAME2", window_hostname, "JA_CMD1", "sleep 120", "JA_CMD2", "timeout /t 120")

		if err != nil {
			fmt.Print("Environments values that is being set are wrong. Please reset them again.")
			return FAILED
		}

		// 1. Run Jobarg_Cleanup_linux
		// 2. Run Jobarg_Cleanup_window
		// 3. Clear Linux Agent log
		// 4. Clear Linux Server log
		// 5. Enable Jobnet ID "Icon_1". Jobnet Name "jobicon_linux_and_windows"
		// 6. Run Jobnet with Parallel Window and Linux agent.
		// 7. Sleep 10 seconds for all Jobicon to start running.
		// 8. Restart linux Jobarg_Server
		// 9. Looping until timeout for the log to write.

		if Run_Jobarg_cleanup_linux(tc_99) &&
			Run_Jobarg_cleanup_window(tc_99) &&
			Run_Clear_Linux_Agent_log(tc_99) &&
			Run_Clear_Linux_Server_log(tc_99) &&
			Run_enable_jobnet(tc_99, jobnet_id, jobnet_name) &&
			func() bool {
				result, _ = Run_Jobnet_Envs_And_Exec(tc_99, jobnet_id, envs)
				return result
			}() &&
			//Run_Job_process_count(tc_99, 1, 10) && // Process count is commented out since it's impossible/hard to count window process.
			Run_sleep() && //Instead, a simple sleep is here.
			Run_Restart_Linux_Jaz_server(tc_99) &&
			func() bool {
				var executeResult string
				timeoutDuration := time.Duration(timeout_time) * time.Second
				timeout := time.Now().Add(timeoutDuration)
				for time.Now().Before(timeout) {
					result, executeResult = Run_Linux_Command_Str(tc_99, "cat /var/log/jobarranger/jobarg_server.log | grep 'In ja_send_ipchange_request(), Server Ip is up to date in'")
					if executeResult != "" && executeResult != " " {
						break
					}
					fmt.Print(tc_99.Info_log("Info: Executing the command Again."))
					time.Sleep(1 * time.Second)
				}
				if executeResult == "" || executeResult == " " {
					fmt.Print(tc_99.Err_log("Error: There's no 'In ja_send_ipchange_request(), Server Ip is up to date in' log jobarranger_server log."))
					return false
				}
				fmt.Print(tc_99.Info_log("Info: Execution result: %s", executeResult))
				return result
			}() {
			fmt.Println("All operations completed successfully")
			return PASSED
		}
		return FAILED
	}
	tc_99.Set_function(tc_func)
	t.Add_testcase(*tc_99)

	tc_100 := t.New_testcase(100, "Abnormal Case Check")
	tc_func = func() common.Testcase_status {

		// Issue 1091 - Jobarg-server crash when it's restarted while it has running jobs
		//
		// Abnormal Case Check
		//
		//  - Create parallel jobs with different hosts as many as you can use. Command -> sleep 120
		//
		//  - Run the jobnet.
		//  - Make sure all job icons are running.
		//  - Delete all files under /var/lib/jobarranger/tmp/serverIPs/serverIPs.json
		//  - Restart the jobarg-server
		//
		//  - The jobnet should be ended with green and the following log must be printed and no CRIT logs and no error logs should be printed.
		//  - "[INFO] In ja_send_ipchange_request(), Server Ip is up to date in host_name."

		linux_hostnames := [...]string{"oss.linux", "oss.linux", "oss.linux", "oss.linux", "oss.linux", "oss.linux"} // Multiple Hostnames of linux agents you want to run must be inputed here. Maximum supported up to 6 linux agents.
		var jobnet_id string = "Icon_1"                                                                              // This value must be Jobnet_ID that you want to run.
		var jobnet_name string = "jobicon_linux_" + strconv.Itoa(len(linux_hostnames)) + "_hosts"                    // This value must be Jobnet_Name that you want to run.
		var timeout_time int = 60                                                                                    // This is timeout value in seconds for waiting for the log to appear.
		var jobnet_run_manage_id string
		var result bool

		// 1.  Setup and initialize envs variable.
		// 2.  Run Jobarg_Cleanup_linux
		// 3.  Clear Linux Agent log
		// 4.  Clear Linux Server log
		// 5.  Enable Jobnet ID "Icon_1". Jobnet Name "jobicon_linux_" + (total hostname in array) + "_hosts"
		// 6.  Run Jobnet with Parallel Multiple Linux agents.
		// 7.  Sleep 10 seconds for all Jobicon to start running.
		// 8.  Removing all files under /var/lib/jobarranger/tmp/serverIPs/serverIPs.json
		// 9.  Restart linux Jobarg_Server
		// 10. Looping until timeout for the log to write.

		var args []string
		for i := 1; i <= len(linux_hostnames); i++ {
			ja_hostname := "JA_HOSTNAME" + strconv.Itoa(i)
			args = append(args, ja_hostname, linux_hostnames[i-1])
		}

		args = append(args, "JA_CMD", "sleep 120")
		envs, err := lib.Get_str_str_map(args...)

		if err != nil || len(linux_hostnames) >= 7 {
			fmt.Print("Environments values that is being set are wrong. Please reset them again.")
			return FAILED
		}

		if Run_Jobarg_cleanup_linux(tc_100) &&
			Run_Clear_Linux_Agent_log(tc_100) &&
			Run_Clear_Linux_Server_log(tc_100) &&
			Run_enable_jobnet(tc_100, jobnet_id, jobnet_name) &&
			func() bool {
				result, jobnet_run_manage_id = Run_Jobnet_Envs_And_Exec(tc_100, jobnet_id, envs)
				return result
			}() &&
			// Run_Job_process_count(tc_100, 1, 10) && // Process count is commented out since it's impossible/hard to count processes of multiple agents
			Run_sleep() && //Instead, a simple sleep is here.
			Run_Linux_Command(tc_100, "rm -rf /var/lib/jobarranger/tmp/serverIPs/serverIPs.json") &&
			Run_Restart_Linux_Jaz_server(tc_100) &&
			func() bool {
				var jobnet_run_info *common.Jobnet_run_info
				result, jobnet_run_info = Run_Jobarg_get_jobnet_run_info(tc_100, jobnet_run_manage_id)
				fmt.Print(tc_100.Info_log("Info: Jobnet Status: %s", jobnet_run_info.Jobnet_status))
				return result
			}() &&
			func() bool {
				var executeResult string
				timeoutDuration := time.Duration(timeout_time) * time.Second
				timeout := time.Now().Add(timeoutDuration)
				for time.Now().Before(timeout) {
					result, executeResult = Run_Linux_Command_Str(tc_99, "cat /var/log/jobarranger/jobarg_server.log | grep 'In ja_send_ipchange_request(), Server Ip is up to date in'")
					if executeResult != "" && executeResult != " " {
						break
					}
					time.Sleep(1 * time.Second)
				}
				if executeResult == "" || executeResult == " " {
					fmt.Print(tc_99.Err_log("Error: There's no 'In ja_send_ipchange_request(), Server Ip is up to date in' log jobarranger_server log."))
					return false
				}
				fmt.Println(tc_100.Info_log("Info: Execution result: %s", executeResult))
				return result
			}() {
			fmt.Println("All operations completed successfully")
			return PASSED
		}
		return FAILED
	}
	tc_100.Set_function(tc_func)
	t.Add_testcase(*tc_100)
}

func Run_sleep() bool {
	time.Sleep(time.Duration(10) * time.Second)
	return true
}
