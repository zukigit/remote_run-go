package tickets

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/zukigit/remote_run-go/src/common"
	"github.com/zukigit/remote_run-go/src/dao"
	"github.com/zukigit/remote_run-go/src/lib"
	"golang.org/x/crypto/ssh"
)

type Ticket_1234 struct {
	no          uint
	description string
	testcases   []dao.TestCase
}

func (t *Ticket_1234) New_testcase(testcase_id uint, testcase_description string) *dao.TestCase {
	return dao.New_testcase(testcase_id, testcase_description)
}

func (t *Ticket_1234) Get_no() uint {
	return t.no
}

func (t *Ticket_1234) Get_dsctn() string {
	return t.description
}

func (t *Ticket_1234) Add_testcase(tc dao.TestCase) {
	t.testcases = append(t.testcases, tc)
}

func (t *Ticket_1234) Get_testcases() []dao.TestCase {
	return t.testcases
}

// Enter your ticket information here
func (t *Ticket_1234) Set_values() {
	t.no = 1234 // Enter your ticket id
	t.description = "Run Jobnet WINRM and check std out."
}

// Add your test case here
func (t *Ticket_1234) Add_testcases() {
	// tc_1 := t.New_testcase(1, "Run WINRM Jobnet and check std out printed hostname.")
	// tc_func := func() common.Testcase_status {
	// 	if err := lib.Jobarg_enable_jobnet("Icon_1", "WINRM"); err != nil {
	// 		tc_1.Err_log("Failed to enable jobnet, Error: %s", err)
	// 		return FAILED
	// 	}
	// 	return RunJobnetWinRMJobnet("Icon_1", 1600, 80, tc_1, common.Client)
	// }
	// tc_1.Set_function(tc_func)
	// t.Add_testcase(*tc_1)

	tc_2 := t.New_testcase(2, "Run WINRM getHost Jobnet and check std out printed hostname.")
	tc_func := func() common.Testcase_status {
		if err := lib.Jobarg_enable_jobnet("Icon_1", "WINRM_SRV"); err != nil {
			tc_2.Err_log("Failed to enable jobnet, Error: %s", err)
			return FAILED
		}
		return RunJobnetWinRMGetHostJobnet("Icon_1", 1600, 80, tc_2, common.Client)
	}
	tc_2.Set_function(tc_func)
	t.Add_testcase(*tc_2)
}

// Run the WINRM Jobnet and check the std out
func RunJobnetWinRMJobnet(jobnetId string, processCount int, processCheckTimeout int, testcase *dao.TestCase, client *ssh.Client) common.Testcase_status {
	/*
		Prepare process before execute the ext jobnet
		1. cleanup data from ja_run_jobnet_table
	*/
	lib.Jobarg_cleanup_linux()

	// Get hostname
	hostname, _err := lib.GetOutputStrFromSSHCommand(client, "hostname")
	if _err != nil {
		fmt.Println(testcase.Info_log("Error fetching hostname: %v", _err))
		return FAILED
	}

	// command := fmt.Sprintf("export WINRM_SRV=\"%s\"; echo $WINRM_SRV", hostname)
	env_hostname, err := lib.GetOutputStrFromSSHCommand(client, fmt.Sprintf("export WINRM_SRV=\"%s\"; echo $WINRM_SRV", hostname))
	// fmt.Println(testcase.Err_log("env hostname info: %s", env_hostname))
	if err != nil {
		fmt.Println("Error executing command:", err)
		return FAILED
	}
	if env_hostname != "" {
		envs, _ := lib.Get_str_str_map("JA_HOSTNAME", env_hostname, "JA_CMD", "exit -100000;", "WINRM_SRV", env_hostname)
		run_jobnet_id, err := lib.Jobarg_exec_E("Icon_1", envs)

		if err != nil {
			fmt.Println(testcase.Err_log("Error getting jobnet info: %s", _err))
		}

		if run_jobnet_id != "" {
			fmt.Println(testcase.Info_log("Running the jobnet info: %s", run_jobnet_id))
		}
	}

	return PASSED
}

// Run the WINRM Jobnet and check the std out
func RunJobnetWinRMGetHostJobnet(jobnetId string, processCount int, processCheckTimeout int, testcase *dao.TestCase, client *ssh.Client) common.Testcase_status {
	/*
		Prepare process before execute the ext jobnet
		1. cleanup data from ja_run_jobnet_table
	*/
	lib.Jobarg_cleanup_linux()

	// Get hostname
	hostname, _err := lib.GetOutputStrFromSSHCommand(client, "hostname")
	if _err != nil {
		fmt.Println(testcase.Err_log("Error fetching hostname: %v", _err))
		return FAILED
	}

	fmt.Println(testcase.Info_log("hostname info: %s", hostname))

	cmd := fmt.Sprintf(`bash -c 'export WINRM_SRV="getHost(%s)"; 
	export HOSTNAME="%s"; 
	jobarg_exec -z %s -U Admin -P zabbix -j %s -E WINRM_SRV,HOSTNAME > /tmp/jobarg_output.txt 2>&1; 
	cat /tmp/jobarg_output.txt;'`, hostname, hostname, common.Login_info.Hostname, jobnetId)

	_run_jobnet, _err := lib.GetOutputStrFromSSHCommand(client, cmd)

	if _err != nil {
		fmt.Println(testcase.Err_log("Error getting jobnet info: %v", _err))
		fmt.Println("Output from command:", _run_jobnet) // Print output for debugging
		if exitErr, ok := _err.(*exec.ExitError); ok {
			fmt.Println("Command exited with status:", exitErr.ExitCode())
		}
	} else {
		fmt.Println(testcase.Info_log("Getting jobnet info: %s", _run_jobnet))

		// Adjusted regex pattern to account for spaces
		pattern := `Registry number\s*:\s*\[\s*(\d+)\s*\]`
		re := regexp.MustCompile(pattern)

		// Use FindStringSubmatch to get the ID
		matches := re.FindStringSubmatch(_run_jobnet)
		jobnetID := matches[1]
		fmt.Println("Extracted Jobnet ID:", jobnetID)
		time.Sleep(10 * time.Second) // Sleep for 2 seconds

		// cmd := fmt.Sprintf(`bash -c 'jobarg_get -z %s -U Admin -P zabbix -r %s -e 2>&1'`, common.Login_info.Hostname, jobnetID)
		cmd := fmt.Sprintf(`bash -c 'eval $(jobarg_get -z %s -U Admin -P zabbix -r %s -e) && echo -n $JA_LASTSTDOUT'`, common.Login_info.Hostname, jobnetID)
		fmt.Println("Executing command:", cmd)

		fmt.Println(cmd)
		lastSTDOut, _err_std_out := lib.GetOutputStrFromSSHCommand(client, cmd)
		fmt.Println(lastSTDOut)
		if lastSTDOut != "" {
			// Remove all spaces
			lastSTDOut = strings.ReplaceAll(lastSTDOut, " ", "")

			// Now replace with a comma
			lastSTDOut = strings.ReplaceAll(lastSTDOut, "(", ",(")
			lastSTDOut = strings.ReplaceAll(lastSTDOut, ")", "),")

			// Trim any leading/trailing whitespace
			lastSTDOut = strings.TrimSpace(lastSTDOut)

			// Split the output by commas
			values := strings.Split(lastSTDOut, ",")

			// Trim each value to remove whitespace
			for i := range values {
				values[i] = strings.TrimSpace(values[i])
			}

			// Print out the cleaned values
			for _, value := range values {
				if value != "" {
					// fmt.Println(testcase.Info_log("Value: %s", value))
					// Assuming 'hostname' is already set and contains the correct value
					hostnameTrimmed := strings.TrimSpace(hostname)
					valueTrimmed := strings.TrimSpace(value)
					ipTrimed := strings.TrimSpace(common.Login_info.Hostname)

					// Print for debugging
					fmt.Println(testcase.Info_log("Comparing Hostname: '%s' with Value: '%s'", hostnameTrimmed, valueTrimmed))

					if strings.EqualFold(valueTrimmed, hostnameTrimmed) && strings.EqualFold(valueTrimmed, ipTrimed) { // Case insensitive comparison
						return PASSED
					}
				}
			}
		} else {
			fmt.Println("No output for JA_LASTSTDOUT.")
		}

		if _err_std_out != nil {
			fmt.Println(testcase.Err_log("Error Std out log: %s", _err_std_out))
		}
	}

	return FAILED

}
