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
	t.description = "Enter your ticket description here."
}

// Add your test case here
func (t *Ticket_1264) Add_testcases() {
	// TESTCASE 125

	tc_125 := t.New_testcase(25, "Enter your test case description here.")
	tc_func := func() common.Testcase_status {

		var result bool
		var current_pwd string

		if func() bool {
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
			func() bool { // I hate you shell script.
				return (Run_Set_Config_Linux(tc_125, "# DBPasswordExternal=", "DBPasswordExternal=echo \"{\\\\\\\"password\\\\\\\":\\\\\\\"zabbix\\\\\\\"}\" | jobarg_json_parse password", "/etc/jobarranger/jobarg_server.conf", 2) ||
					Run_Set_Config_Linux(tc_125, "DBPasswordExternal=", "DBPasswordExternal=echo \"{\\\\\\\"password\\\\\\\":\\\\\\\"zabbix\\\\\\\"}\" | jobarg_json_parse password", "/etc/jobarranger/jobarg_server.conf", 2))

			}() &&
			func() bool {
				cmd := "sed -i \"s/^#*\\(define('DB_PASS' , 'zabbix'); \\).*/define('DB_PASS' , \\`echo \\\"{\\\\\\\\\\\"password\\\\\\\\\\\":\\\\\\\\\\\"zabbix\\\\\\\\\\\"}\\\" | jobarg_json_parse password\\`); /\" /etc/jobarranger/web/jam.config.php"
				return Run_Linux_Command(tc_125, cmd)
			}() &&
			func() bool {
				// sed -i 's/^#*\(%s\).*/%s/' %s`, key, value, config_file_path
				// sed -i "s/^#*\($DB['PASSWORD']                 = 'zabbix';\).*/$DB['PASSWORD']                 = 'labbix';/" /etc/zabbix/web/zabbix.conf.php"
				cmd := "sed -i \"s/^#*\\(\\$DB\\['PASSWORD'\\]                 = 'zabbix';\\).*/\\$DB['PASSWORD']                 = \\`echo \\\"{\\\\\\\\\\\"password\\\\\\\\\\\":\\\\\\\\\\\"zabbix\\\\\\\\\\\"}\\\" | jobarg_json_parse password\\`;/\" /etc/zabbix/web/zabbix.conf.php"
				return Run_Linux_Command(tc_125, cmd)
			}() {
			fmt.Println("All operations completed successfully")
			return PASSED
		}
		return FAILED
	}
	tc_125.Set_function(tc_func)
	t.Add_testcase(*tc_125)
}
