package tickets

import (
	"fmt"
	"strings"

	"github.com/zukigit/remote_run-go/src/common"
	"github.com/zukigit/remote_run-go/src/dao"
)

type Ticket_1021 struct {
	no          uint
	description string
	testcases   []dao.TestCase
}

func (t *Ticket_1021) New_testcase(testcase_id uint, testcase_description string) *dao.TestCase {
	return dao.New_testcase(testcase_id, testcase_description)
}

func (t *Ticket_1021) Get_no() uint {
	return t.no
}

func (t *Ticket_1021) Get_dsctn() string {
	return t.description
}

func (t *Ticket_1021) Add_testcase(tc dao.TestCase) {
	t.testcases = append(t.testcases, tc)
}

func (t *Ticket_1021) Get_testcases() []dao.TestCase {
	return t.testcases
}

// Enter your ticket information here
func (t *Ticket_1021) Set_values() {
	t.no = 1021 // Enter your ticket id
	t.description = "Issue 1021 - Converts the character code from UTF-8 and creates transaction files"
}

// Add your test case here
func (t *Ticket_1021) Add_testcases() {
	// TESTCASE 001
	tc_104 := t.New_testcase(1021, "Normal Case Check")

	tc_func := func() common.Testcase_status {

		// Issue 1021 - Converts the character code from UTF-8 and creates transaction files
		//
		// Normal Case Check
		//  - Check if Encoding is displayed correctly at jobarranger agent startup.
		//  - Check if transaction file(.json/.sh files) have correct encoding text.
		//  - Job can be executed correctly with full-width characters.
		//  - filewait can be executed correctly with full-width characters
		//  - file check can be executed correctly with full-width characters.

		var jobnet_id string = "Icon_1"                   // This value must be Jobnet_ID that you want to run.
		var jobnet_name string = "TICKET1021_TESTCASE104" // This value must be Jobnet_Name that you want to run.
		var systemEncode, agentEncode, jobnet_run_manage_id, close_folder_name, sh_file_encoding, json_file_encoding string
		var result bool

		// 1.  Clear Agent log
		// 2.  Restart JAZ Agent.
		// 3.  Remove All Files in Close Folder.
		// 4.  Get System Encoding Value.
		// 5.  Get Agent Log Encoding Value.
		// 6.  Compare Encoding value.
		// 7.  Run Jobarg_Cleanup_linux
		// 8.  Enable Jobnet jobicon_linux
		// 9.  Run Jobnet With Full-width Character
		// 10. Get Jobnet Run Info
		// 11. Run Jobarg_Cleanup_linux
		// 12. Enable Jobnet TICKET1021_TESTCASE104
		// 13. Make temporary Full-width character file for File Wait and File Check. (Allow Root must be 1)
		// 14. Run Jobnet With Full-width Character for File Wait and File Check
		// 15. Get Jobnet Run Info
		// 16. Remove the temporary file.
		// 17. Get File name of Close folder
		// 18. Get File encoding of .sh file in close folder with file --mime command.
		// 19. Get File encoding of .json file in close folder with file --mime command.
		// 20. Compare Encoding value.

		if Run_Clear_Linux_Agent_log(tc_104) &&
			Run_Restart_Linux_Jaz_agent(tc_104) &&
			Run_Linux_Command(tc_104, "rm -rf /var/lib/jobarranger/tmp/close/*") &&
			func() bool {
				result, executeResult := Run_Linux_Command_Str(tc_104, "locale | grep 'LC_CTYPE='")
				systemEncode = executeResult[strings.LastIndex(executeResult, "=")+1:]
				return result
			}() &&
			func() bool {
				result, executeResult := Run_Linux_Command_Str(tc_104, "cat /var/log/jobarranger/jobarg_agentd.log | grep 'LC_CTYPE :'")
				agentEncode = executeResult[strings.LastIndex(executeResult, ":")+1:]
				return result
			}() &&
			func() bool {
				systemEncode = strings.TrimSpace(systemEncode)
				agentEncode = strings.TrimSpace(agentEncode)
				systemEncode = strings.Trim(systemEncode, "\"\n")
				agentEncode = strings.Trim(agentEncode, "\"\n")

				fmt.Println(tc_104.Info_log("Info: System Encode: " + systemEncode))
				fmt.Println(tc_104.Info_log("Info: Agent Encode: " + systemEncode))

				if strings.TrimSpace(systemEncode) == strings.TrimSpace(agentEncode) {
					fmt.Println(tc_104.Info_log("Info: Both Encoding are same."))
					return true
				} else {
					fmt.Print(tc_104.Info_log("Info: Encoding are not same."))
					return false
				}
			}() &&
			Run_Jobarg_cleanup_linux(tc_104) &&
			Run_enable_jobnet(tc_104, jobnet_id, "jobicon_linux") &&
			func() bool {
				result, jobnet_run_manage_id = Run_Jobnet_Exec(tc_104, jobnet_id, "echo ａｂｃｄｅｆｇｈｉｊｋｌｍｎｏｐｑｒｓｔｕｖｗｘｙｚ｀１２３４５６７８９０－＝～！＠＃＄％＾＆＊（）＿＋，．／＜＞？；＇：＂［］｛｝＼｜ジョブの単一実行")
				return result
			}() &&
			func() bool {
				result, jobnet_run_info := Run_Jobarg_get_jobnet_run_info(tc_104, jobnet_run_manage_id)
				fmt.Println(tc_104.Info_log("Info: Jobnet Std_out: %s", jobnet_run_info.Std_out))
				fmt.Println(tc_104.Info_log("Info: Job can be executed correctly with full-width characters."))
				return result
			}() &&
			Run_Jobarg_cleanup_linux(tc_104) &&
			Run_enable_jobnet(tc_104, jobnet_id, jobnet_name) &&
			Run_Linux_Command(tc_104, "touch /home/ａｂｃｄｅｆｇｈｉｊｋｌｍｎｏｐｑｒｓｔｕｖｗｘｙｚ｀１２３４５６７８９０－＝～！＠＃＄％＾＆＊（）＿＋，．／＜＞？；＇：＂［］｛｝＼｜ジョブの単一実行．ｔｘｔ") &&
			func() bool {
				result, jobnet_run_manage_id = Run_Jobnet(tc_104, jobnet_id)
				return result
			}() &&
			func() bool {
				result, jobnet_run_info := Run_Jobarg_get_jobnet_run_info(tc_104, jobnet_run_manage_id)
				fmt.Println(tc_104.Info_log("Info: Jobnet Std_out: %s", jobnet_run_info.Std_out))
				fmt.Println(tc_104.Info_log("Info: Both File Wait and File Check can be executed correctly with full-width characters."))
				return result
			}() &&
			Run_Linux_Command(tc_104, "rm -rf /home/ａｂｃｄｅｆｇｈｉｊｋｌｍｎｏｐｑｒｓｔｕｖｗｘｙｚ｀１２３４５６７８９０－＝～！＠＃＄％＾＆＊（）＿＋，．／＜＞？；＇：＂［］｛｝＼｜ジョブの単一実行．ｔｘｔ") &&
			func() bool {
				result, executeResult := Run_Linux_Command_Str(tc_104, "ls /var/lib/jobarranger/tmp/close/")
				close_folder_name = strings.Split(executeResult, "\n")[0]
				fmt.Print(tc_104.Info_log("Info: Execution result: %s", strings.Trim(executeResult, "\n")))
				return result
			}() &&
			func() bool {
				result, executeResult := Run_Linux_Command_Str(tc_104, "file --mime-encoding /var/lib/jobarranger/tmp/close/"+close_folder_name+"/"+close_folder_name+".sh")
				fmt.Print(tc_104.Info_log("Info: Execution result: %s", strings.Trim(executeResult, "\n")))
				sh_file_encoding = executeResult
				return result
			}() &&
			func() bool {
				result, executeResult := Run_Linux_Command_Str(tc_104, "file --mime-encoding /var/lib/jobarranger/tmp/close/"+close_folder_name+"/"+close_folder_name[:strings.LastIndex(close_folder_name, "-")]+".json")
				fmt.Print(tc_104.Info_log("Info: Execution result: %s", strings.Trim(executeResult, "\n")))
				json_file_encoding = executeResult
				return result
			}() &&
			func() bool {
				//sh_file_encoding = strings.TrimSpace(strings.Trim(strings.Split(sh_file_encoding[strings.LastIndex(sh_file_encoding, ":")+1:], " ")[3], ","))
				sh_file_encoding = strings.TrimSpace(sh_file_encoding[strings.LastIndex(sh_file_encoding, ":")+1:])
				json_file_encoding = strings.TrimSpace(json_file_encoding[strings.LastIndex(json_file_encoding, ":")+1:])
				fmt.Println(tc_104.Info_log("Info: Sh file value = %s", sh_file_encoding))
				fmt.Println(tc_104.Info_log("Info: Json file value = %s", json_file_encoding))
				if strings.EqualFold(sh_file_encoding, json_file_encoding) {
					fmt.Println(tc_104.Info_log("Info: Both Encoding are same."))
					return true
				} else {
					fmt.Print(tc_104.Info_log("Info: Encoding are not same."))
					return false
				}
			}() {
			fmt.Println("All operations completed successfully")
			return PASSED
		} else {
			return FAILED
		}
	}
	tc_104.Set_function(tc_func)
	t.Add_testcase(*tc_104)
}
