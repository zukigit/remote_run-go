package tickets

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/sftp"
	"github.com/zukigit/remote_run-go/src/common"
	"github.com/zukigit/remote_run-go/src/dao"
)

type Ticket_1264 struct {
	no          uint
	description string
	testcases   []dao.TestCase
}

func (t *Ticket_1264) New_testcase(testcase_id uint, testcase_description string) *dao.TestCase {
	return dao.New_testcase(testcase_id, testcase_description)
}

func (t *Ticket_1264) Get_no() uint {
	return t.no
}

func (t *Ticket_1264) Get_dsctn() string {
	return t.description
}

func (t *Ticket_1264) Add_testcase(tc dao.TestCase) {
	t.testcases = append(t.testcases, tc)
}

func (t *Ticket_1264) Get_testcases() []dao.TestCase {
	return t.testcases
}

// Enter your ticket information here
func (t *Ticket_1264) Set_values() {
	t.no = 1264 // Enter your ticket id
	t.description = "Issue 1264 - Automate zabbix-web, jaz-web and jobarg-server"
}

// Add your test case here
func (t *Ticket_1264) Add_testcases() {
	// TESTCASE 125 and TESTCASE 156

	// Issue 1264 - Automate zabbix-web, jaz-web and jobarg-server
	//
	// Default Case Check and Abnormal Case Check
	//
	// - Downd the jobarg_json_parse from the following link.
	// - Link: http://172.30.20.63/redmine/attachments/17496
	// - Move jobarg_json_parse to /usr/local/bin .
	//
	// - Change the config files.
	//
	//  - Run the following command to reastart jobarg-server and httpd
	//  - systemctl restart jobarg-server php-fpm httpd
	//  - Make sure all job icons are running.
	//  - Restart the jobarg-server
	//
	//  - Check jobarg-server, zabbix-web and jobarg-web like job creation, host creation and job running. Everything must work normally.

	tc_125 := t.New_testcase(125, "Default Case Check and Abnormal Case Check")
	tc_func := func() common.Testcase_status {

		var result bool
		var current_pwd, jobnet_manage_id string

		if Run_Jobarg_cleanup_linux(tc_125) &&
			Run_Linux_Command(tc_125, "cp /etc/jobarranger/jobarg_server.conf /etc/jobarranger/jobarg_server.conf.bk") &&
			Run_Linux_Command(tc_125, "cp /etc/jobarranger/web/jam.config.php /etc/jobarranger/web/jam.config.php.bk") &&
			Run_Linux_Command(tc_125, "cp /etc/zabbix/web/zabbix.conf.php /etc/zabbix/web/zabbix.conf.php.bk") &&
			Run_enable_jobnet(tc_125, "Icon_1", "jobicon_linux") &&
			func() bool {
				result, current_pwd = Run_Window_Command_Str(tc_125, "cd")
				fmt.Print(current_pwd)
				return result
			}() &&
			func() bool {

				//common.Set_passwd()
				//common.Set_client()

				// Create an SFTP client
				sftpClient, err := sftp.NewClient(common.Client)
				if err != nil {
					fmt.Println(tc_125.Err_log("Error: Failed to create SFTP client: %s", err))
					return false
				}
				defer sftpClient.Close()
				//defer common.Client.Close()

				// Specify the local and remote file paths
				localFilePath := filepath.ToSlash(filepath.Join(strings.TrimSpace(strings.Trim(current_pwd, "\n")), "exported_jobnets", "jobarg_json_parse")) // Replace with the desired full file path on the Window
				remoteFilePath := "/usr/local/bin/jobarg_json_parse"                                                                                          // Replace with the desired full file path on the Linux

				// Open the local file
				localFile, err := os.Open(localFilePath)
				if err != nil {
					fmt.Println(tc_125.Err_log("Error: Failed to open local file: %s", err))
					return false
				}
				defer localFile.Close()

				// Create the remote file
				remoteFile, err := sftpClient.Create(remoteFilePath)
				if err != nil {
					fmt.Println(tc_125.Err_log("Error: Failed to create remote file: %s", err))
					return false
				}
				defer remoteFile.Close()

				// Copy the file from local to remote
				bytes, err := io.ReadAll(localFile)
				if err != nil {
					fmt.Println(tc_125.Err_log("Error: Failed to read local file: %s", err))
					return false
				}
				_, err = remoteFile.Write(bytes)
				if err != nil {
					fmt.Println(tc_125.Err_log("Error: Failed to write to remote file: %s", err))
					return false
				}

				fmt.Println(tc_125.Info_log("Info: File transferred successfully!"))

				defer remoteFile.Close()

				return true
			}() &&
			Run_Linux_Command(tc_125, "chmod +x /usr/local/bin/jobarg_json_parse") &&
			func() bool {
				return (Run_Set_Config_Linux(tc_125, "# DBPasswordExternal=", "DBPasswordExternal=echo \"{\\\\\\\"password\\\\\\\":\\\\\\\"zabbix\\\\\\\"}\" | jobarg_json_parse password", "/etc/jobarranger/jobarg_server.conf", 2) ||
					Run_Set_Config_Linux(tc_125, "DBPasswordExternal=", "DBPasswordExternal=echo \"{\\\\\\\"password\\\\\\\":\\\\\\\"zabbix\\\\\\\"}\" | jobarg_json_parse password", "/etc/jobarranger/jobarg_server.conf", 2))
			}() &&
			func() bool {
				// I hate you shell script. sed command litreally makes no sense.
				//
				// Text we want to replace for:
				//   - define('DB_PASS' , `echo "{\"password\":\"zabbix\"}" | jobarg_json_parse password`);
				//   - $DB['PASSWORD']                 = `echo "{\"password\":\"zabbix\"}" | jobarg_json_parse password`;
				// __________________________________________
				//| Original Char |  In Script  | In String  |
				//| "             |  \"         | \\\"       |
				//| \             |  \\\\       | \\\\\\\\   |
				//| \"            |  \\\\\"     | \\\\\\\\\\"|
				//| `             |  \`         | \\`        |
				//|_______________|_____________|____________|
				cmd := "sed -i \"s/^#*\\(define('DB_PASS' , 'zabbix'); \\).*/define('DB_PASS' , \\`echo \\\"{\\\\\\\\\\\"password\\\\\\\\\\\":\\\\\\\\\\\"zabbix\\\\\\\\\\\"}\\\" | jobarg_json_parse password\\`); /\" /etc/jobarranger/web/jam.config.php"
				return Run_Linux_Command(tc_125, cmd)
			}() &&
			func() bool {
				cmd := "sed -i \"s/^#*\\(\\$DB\\['PASSWORD'\\]                 = 'zabbix';\\).*/\\$DB\\['PASSWORD'\\]                 = \\`echo \\\"{\\\\\\\\\\\"password\\\\\\\\\\\":\\\\\\\\\\\"zabbix\\\\\\\\\\\"}\\\" | jobarg_json_parse password\\`;/\" /etc/zabbix/web/zabbix.conf.php"
				return Run_Linux_Command(tc_125, cmd)
			}() &&
			func() bool {
				Run_Linux_Command(tc_125, "systemctl restart jobarg-server php-fpm httpd apache2")
				return true
			}() &&
			func() bool {
				result, jobnet_manage_id = Run_Jobnet_Exec(tc_125, "Icon_1", "hostname")
				return result
			}() &&
			func() bool {
				result, jobnet_run_info := Run_Jobarg_get_jobnet_run_info(tc_125, jobnet_manage_id)
				fmt.Println(tc_125.Info_log("Info: Jobnet Status: %s, Job Status: %s,Std Out: %s", jobnet_run_info.Jobnet_status, jobnet_run_info.Job_status, jobnet_run_info.Std_out))
				fmt.Println(tc_125.Info_log("Info: Everything works normally."))
				return result
			}() &&
			Run_Set_Config_Linux(tc_125, "DBPassword=zabbix", "DBPassword=YOUR_INCORRECT_PASSWORD", "/etc/jobarranger/jobarg_server.conf", 2) &&
			Run_Restart_Linux_Jaz_server(tc_125) &&
			func() bool {
				result, jobnet_manage_id = Run_Jobnet_Exec(tc_125, "Icon_1", "hostname")
				return result
			}() &&
			func() bool {
				result, jobnet_run_info := Run_Jobarg_get_jobnet_run_info(tc_125, jobnet_manage_id)
				fmt.Println(tc_125.Info_log("Info: Jobnet Status: %s, Job Status: %s,Std Out: %s", jobnet_run_info.Jobnet_status, jobnet_run_info.Job_status, jobnet_run_info.Std_out))
				fmt.Println(tc_125.Info_log("Info: Everything works normally."))
				return result
			}() &&
			Run_Set_Config_Linux(tc_125, "DBPassword=YOUR_INCORRECT_PASSWORD", "DBPassword=zabbix", "/etc/jobarranger/jobarg_server.conf", 2) &&
			Run_Set_Config_Linux(tc_125, "DBPasswordExternal=echo \"{\\\\\\\"password\\\\\\\":\\\\\\\"zabbix\\\\\\\"}\" | jobarg_json_parse password", "DBPasswordExternal=nsfkjnkjnsidufn", "/etc/jobarranger/jobarg_server.conf", 2) &&
			Run_Restart_Linux_Jaz_server(tc_125) &&
			func() bool {
				result, jobnet_manage_id = Run_Jobnet_Exec(tc_125, "Icon_1", "hostname")
				return result
			}() &&
			func() bool {
				result, jobnet_run_info := Run_Jobarg_get_jobnet_run_info(tc_125, jobnet_manage_id)
				fmt.Println(tc_125.Info_log("Info: Jobnet Status: %s, Job Status: %s,Std Out: %s", jobnet_run_info.Jobnet_status, jobnet_run_info.Job_status, jobnet_run_info.Std_out))
				fmt.Println(tc_125.Info_log("Info: Everything works normally."))
				return result
			}() {
			fmt.Println("All operations completed successfully")
			return PASSED
		}
		fmt.Println(tc_125.Info_log("Info: Resting config files back to normal."))
		if Run_Linux_Command(tc_125, "rm -rf /usr/local/bin/jobarg_json_parse") &&
			Run_Linux_Command(tc_125, "rm -rf /etc/jobarranger/jobarg_server.conf") &&
			Run_Linux_Command(tc_125, "rm -rf /etc/jobarranger/web/jam.config.php") &&
			Run_Linux_Command(tc_125, "rm -rf /etc/zabbix/web/zabbix.conf.php") &&
			Run_Linux_Command(tc_125, "mv /etc/jobarranger/jobarg_server.conf.bk /etc/jobarranger/jobarg_server.conf") &&
			Run_Linux_Command(tc_125, "mv /etc/jobarranger/web/jam.config.php.bk /etc/jobarranger/web/jam.config.php") &&
			Run_Linux_Command(tc_125, "mv /etc/zabbix/web/zabbix.conf.php.bk /etc/zabbix/web/zabbix.conf.php") &&
			Run_Linux_Command(tc_125, "chown zabbix:zabbix /etc/jobarranger/jobarg_server.conf") &&
			(Run_Linux_Command(tc_125, "chown apache:apache /etc/jobarranger/web/jam.config.php") ||
				Run_Linux_Command(tc_125, "chown www-data:www-data /etc/jobarranger/web/jam.config.php")) &&
			(Run_Linux_Command(tc_125, "chown apache:apache /etc/zabbix/web/zabbix.conf.php") ||
				Run_Linux_Command(tc_125, "chown www-data:www-data /etc/zabbix/web/zabbix.conf.php")) &&
			Run_Restart_Linux_Jaz_agent(tc_125) &&
			Run_Restart_Linux_Jaz_server(tc_125) {
			Run_Linux_Command(tc_125, "systemctl restart zabbix-server zabbix-agent jobarg-server php-fpm httpd apache2")
			Run_Jobarg_cleanup_linux(tc_125)
			return FAILED
		}
		fmt.Println(tc_125.Info_log("Error: Resting config files back to normal settings. Please reset the config files manually."))
		return MUST_CHECK
	}
	tc_125.Set_function(tc_func)
	t.Add_testcase(*tc_125)

}
