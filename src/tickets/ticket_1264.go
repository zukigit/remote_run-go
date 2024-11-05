package tickets

import (
	"fmt"
	"io/ioutil"
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

	var result bool
	var current_pwd string

	tc_125 := t.New_testcase(1, "Enter your test case description here.")
	tc_func := func() common.Testcase_status {
		if func() bool {
			result, current_pwd = Run_Window_Command_Str(tc_125, "cd")
			fmt.Print(current_pwd)
			return result
		}() &&
			func() bool {
				// Create an SFTP client
				common.Set_passwd()
				common.Set_client()

				sftpClient, err := sftp.NewClient(common.Client)
				if err != nil {
					tc_125.Err_log("failed to create SFTP client: %s", err)
					return false
				}
				defer sftpClient.Close()
				defer common.Client.Close()

				// Specify the local and remote file paths
				localFilePath := filepath.ToSlash(filepath.Join(strings.TrimSpace(strings.Trim(current_pwd, "\n")))) // Replace with the path to your local file
				remoteFilePath := "/root/"                                                                           // Replace with the desired path on the Linux server

				// Open the local file
				localFile, err := os.Open(localFilePath)
				if err != nil {
					tc_125.Err_log("failed to open local file: %s", err)
					return false
				}
				defer localFile.Close()

				// Create the remote file
				remoteFile, err := sftpClient.Create(remoteFilePath)
				if err != nil {
					tc_125.Err_log("failed to create remote file: %s", err)
					return false
				}
				defer remoteFile.Close()

				// Copy the file from local to remote
				bytes, err := ioutil.ReadAll(localFile)
				if err != nil {
					tc_125.Err_log("failed to read local file: %s", err)
					return false
				}
				_, err = remoteFile.Write(bytes)
				if err != nil {
					tc_125.Err_log("failed to write to remote file: %s", err)
					return false
				}

				fmt.Println("File transferred successfully!")
				return true
			}() {
			fmt.Println("All operations completed successfully")
			return PASSED
		}
		return FAILED
	}
	tc_125.Set_function(tc_func)
	t.Add_testcase(*tc_125)
}
