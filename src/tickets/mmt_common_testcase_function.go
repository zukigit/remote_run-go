package tickets

import (
	"database/sql"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/sftp"
	"github.com/zukigit/remote_run-go/src/common"
	"github.com/zukigit/remote_run-go/src/dao"
	"github.com/zukigit/remote_run-go/src/lib"
	"golang.org/x/crypto/ssh"
)

// A Collection of all common testcase functions that write logs and perform tasks.
// Return True and False and some certain values depending on fail or pass.
//
// Most of these functions receive testcase parameter
// which means you must use them in testcase functions.
//
// Purpose:
//      - Don't want to write the same code that write logs statements over again and again.
//      - I hate copying and pasting code.
//      - Perhaps this could be used as common library as testcase functions but I'm not still sure.
//      - For now, I'm using it.
//
// Author:
//      - Myint Myat Thu
//      - ChatGPT
//

// To enable common jobnet
//
// Returns
//   - True if worked.
//   - False if failed.
func Run_enable_jobnet(testcase *dao.TestCase, jobnet_id string, jobnet_name string) bool {
	// "TICKET844_TESTCASE75"
	if err := lib.Jobarg_enable_jobnet(jobnet_id, jobnet_name); err != nil {
		fmt.Println(testcase.Err_log("Error: Failed to enable jobnet. Jobnet Id: %s, Jobnet Name: %s Error: %s", jobnet_id, jobnet_name, err))
		return false
	} else {
		fmt.Println(testcase.Info_log("Info: Jobnet Enabled successfully. Jobnet Id: %s, Jobnet Name: %s", jobnet_id, jobnet_name))
		return true
	}
}

// To run jobnet
//
// Returns
//   - True and jobnet_run_manage_id if worked.
//   - False and jobnet_run_manage_id if failed.
func Run_Jobnet(testcase *dao.TestCase, jobnet_id string) (bool, string) {
	envs, _ := lib.Get_str_str_map("JA_HOSTNAME", "oss.linux", "JA_CMD", "")
	jobnet_run_manage_id, err := lib.Jobarg_exec_E(jobnet_id, envs)
	if err != nil {
		fmt.Println(testcase.Err_log("Error: Failed when trying to run the job. %s", err.Error()))
		fmt.Println(testcase.Err_log("Error: Did you imported the necessary xml files?(%s)", jobnet_id))
		return false, jobnet_run_manage_id
	} else {
		fmt.Println(testcase.Info_log("Info: Job has start running."))
		return true, jobnet_run_manage_id
	}
}

// To run jobnet with Execute Command.
//
// Returns
//   - True and jobnet_run_manage_id if worked.
//   - False and jobnet_run_manage_id if failed.
func Run_Jobnet_Exec(testcase *dao.TestCase, jobnet_id string, execute_command string) (bool, string) {
	envs, _ := lib.Get_str_str_map("JA_HOSTNAME", "oss.linux", "JA_CMD", execute_command)
	jobnet_run_manage_id, err := lib.Jobarg_exec_E(jobnet_id, envs)
	if err != nil {
		fmt.Println(testcase.Err_log("Error: Failed when trying to run the job. %s", err.Error()))

		//fmt.Println(testcase.Err_log("Error: Did you imported the necessary xml files?(%s)", jobnet_id))
		return false, jobnet_run_manage_id
	} else {
		fmt.Println(testcase.Info_log("Info: Job has start running."))
		return true, jobnet_run_manage_id
	}
}

// To run jobnet with Multiple Environments Variable and Execute Command.
//
// Returns
//   - True and jobnet_run_manage_id if worked.
//   - False and jobnet_run_manage_id if failed.
func Run_Jobnet_Envs_And_Exec(testcase *dao.TestCase, jobnet_id string, envs map[string]string) (bool, string) {
	jobnet_run_manage_id, err := lib.Jobarg_exec_E(jobnet_id, envs)
	if err != nil {
		fmt.Println(testcase.Err_log("Error: Failed when trying to run the job. %s", err.Error()))
		//fmt.Println(testcase.Err_log("Error: Did you imported the necessary xml files?(%s)", jobnet_id))
		return false, jobnet_run_manage_id
	} else {
		fmt.Println(testcase.Info_log("Info: Job has start running."))
		return true, jobnet_run_manage_id
	}
}

// To restart linux jobarg agentd
//
// Returns
//   - True if passed.
//   - False if error.
func Run_Restart_Linux_Jaz_agent(testcase *dao.TestCase) bool {
	// Restartng Jobarranger linux Agent.
	fmt.Println(testcase.Info_log("Info: Restarting Linux Jobarranger Agent now."))
	err := lib.Restart_jaz_agent_linux()
	if err != nil {
		fmt.Println(testcase.Err_log("Error: Failted at restarting Linux Jobarranger Agent. %s", err.Error()))
		return false
	} else {
		fmt.Println(testcase.Info_log("Info: Linux Jobarranger Agent has successfully restarted."))
		return true
	}
}

// To restart window jobarg agentd
//
// Returns
//   - True if passed.
//   - False if error.
func Run_Restart_Window_Jaz_agent(testcase *dao.TestCase) bool {
	// Restartng Jobarranger linux Agent.
	fmt.Println(testcase.Info_log("Info: Restarting Window Jobarranger Agent now."))
	err := lib.Restart_jaz_agent_windows()
	if err != nil {
		fmt.Println(testcase.Err_log("Error: Failted at restarting Window Jobarranger Agent. %s", err.Error()))
		return false
	} else {
		fmt.Println(testcase.Info_log("Info: Window Jobarranger Agent has successfully restarted."))
		return true
	}
}

// To restart linux jobarg server
//
// Returns
//   - True if passed.
//   - False if error.
func Run_Restart_Linux_Jaz_server(testcase *dao.TestCase) bool {
	// Restartng Jobarranger linux Server.
	fmt.Println(testcase.Info_log("Info: Restarting Jobarranger Server now."))
	err := lib.Restart_jaz_server()
	if err != nil {
		fmt.Println(testcase.Err_log("Error: Failted at restarting Jobarranger Server. %s", err.Error()))
		return false
	} else {
		fmt.Println(testcase.Info_log("Info: Jobarranger Server has successfully restarted."))
		return true
	}
}

// To count jobnet process
//
// Returns
//   - True if met with desired process count.
//   - False if failed to met desired process count.
func Run_Job_process_count(testcase *dao.TestCase, job_process_count_amt int, timeout_minute int) bool {
	// Counting Job Process Count
	fmt.Println(testcase.Info_log("Info: Counting Job processes."))
	err := JobProcessCountCheck_with_process_counter(job_process_count_amt, timeout_minute, common.Client)
	if err != nil {
		fmt.Println(testcase.Err_log("Error: Failed at counting or failed to meet required process count amount. %s", err.Error()))
		return false
	} else {
		fmt.Println(testcase.Info_log("Info: Job process count met."))
		return true
	}
}

// To clean up Jobarg_linux.
//
// Returns
//   - True if worked.
//   - False if failed.
func Run_Jobarg_cleanup_linux(testcase *dao.TestCase) bool {
	// Initiating Jobarg_cleanup
	fmt.Println(testcase.Info_log("Info: Starting Jobarg_cleanup. (Linux)"))
	err := lib.Jobarg_cleanup_linux()
	if err != nil {
		fmt.Println(testcase.Err_log("Error: %s", err.Error()))
		return false
	}
	fmt.Println(testcase.Info_log("Info: Jobarg_cleanup Success."))
	return true
}

// To clean up Jobarg_window.
//
// Returns
//   - True if worked.
//   - False if failed.
func Run_Jobarg_cleanup_window(testcase *dao.TestCase) bool {
	// Initiating Jobarg_cleanup
	fmt.Println(testcase.Info_log("Info: Starting Jobarg_cleanup. (Window)"))
	err := lib.Jobarg_cleanup_windows()
	if err != nil {
		fmt.Println(testcase.Err_log("Error: %s", err.Error()))
		return false
	}
	fmt.Println(testcase.Info_log("Info: Jobarg_cleanup Success."))
	return true
}

// To get jobnet_run_info
//
// Returns
//   - True and jobnet_run_info if worked.
//   - False and jobnet_run_info if failed.
func Run_Jobarg_get_jobnet_run_info(testcase *dao.TestCase, jobnet_run_manage_id string) (bool, *common.Jobnet_run_info) {
	//In case if getting jobnet info failed
	fmt.Println(testcase.Info_log("Info: Getting Jobnet Info."))
	jobnet_run_info, err := lib.Jobarg_get_jobnet_run_info(jobnet_run_manage_id)
	if err != nil {
		fmt.Println(testcase.Err_log("Error: Failted at getting jobnet run info. %s. Jobnet Management Id: %s ", err.Error(), jobnet_run_manage_id)) //Possible error: Jobnet process killed, Jobnet force-stopped, Wrong Hostname in Job Icon, Job Icon failed, Job Icon force-stopped."
		return false, jobnet_run_info
	} else {
		fmt.Println(testcase.Info_log("Info: Successfully received getting Jobnet Info."))
		fmt.Println(testcase.Info_log("Info: Jobnet Status: %s, Job Status: %s, STD ERR : %s. STD OUT : %s.", jobnet_run_info.Jobnet_status, jobnet_run_info.Job_status, jobnet_run_info.Std_error, jobnet_run_info.Std_out))
		return true, jobnet_run_info
	}
}

// To check if jobnet finished with no zombie process
//
// Returns
//   - True if worked.
//   - False if failed.
func Run_Check_Jobnet_Finish_With_No_Zombie_Process(testcase *dao.TestCase, jobnet_run_info *common.Jobnet_run_info) bool {

	//Checking whether Jobnet is finished or not
	fmt.Println(testcase.Info_log("Info: Checking for Job Status, Jobnet Status and Exit CD."))
	if jobnet_run_info.Job_status == "NORMAL" && jobnet_run_info.Jobnet_status == "END" && jobnet_run_info.Exit_cd == 0 {

		// Checking whether "checking at zombie process" process failed or not.
		fmt.Println(testcase.Info_log("Info: Success. Job Status = %s, Jobnet Status = %s, Exit_CD = %d", jobnet_run_info.Job_status, jobnet_run_info.Jobnet_status, jobnet_run_info.Exit_cd))
		fmt.Println(testcase.Info_log("Info: Checking for Zombie Process."))
		zombieProcessCount, err := lib.CheckZombieProcess(1, common.Client)
		if err != nil {
			fmt.Println(testcase.Err_log("Error: Failed at checking zombie Process."))
			fmt.Println(testcase.Err_log(err.Error()))
			// Checking whether zombie process exist or not.
			if zombieProcessCount > 0 {
				fmt.Println(testcase.Err_log("Error: There is zombie Process left."))
				fmt.Println(testcase.Err_log(err.Error()))
			}

			return false
		}
		fmt.Println(testcase.Info_log("Info: There's no Zombie Processes. Success."))
		return true
	}

	fmt.Println(testcase.Err_log("Error: Either Job Status or Jobnet Status or Exit CD doesn't meet required contidions.(Normal, End, 0)"))
	fmt.Println(testcase.Err_log("Jobnet_status: %s, Job_status: %s, Exit_cd: %d", jobnet_run_info.Jobnet_status, jobnet_run_info.Job_status, jobnet_run_info.Exit_cd))
	return false
}

// To retrieve value from Jobarg_get_LASTSTDOUT
//
// Returns
//   - True and LASTSTDOUT if worked.
//   - False and nil if failed.
func Run_Jobarg_Get_LastSTDOUT(testcase *dao.TestCase, jobnet_run_manage_id string) (bool, string) {
	last_std_out, err := lib.Jobarg_get_LASTSTDOUT(jobnet_run_manage_id)
	if err != nil {
		fmt.Println(testcase.Err_log("Error: Failed to get LASTSTDOUT value. %s", err))
		return false, ""
	}
	fmt.Println(testcase.Info_log("Info: LASTSTDOUT value: %s", last_std_out))
	return true, last_std_out
}

// To clear jaz linux agent log.
//
// Returns
//   - True if worked.
//   - False if failed.
func Run_Clear_Linux_Agent_log(testcase *dao.TestCase) bool {
	err := lib.Clear_linux_jaz_agent_log()
	if err != nil {
		fmt.Println(testcase.Err_log("Error: Failed at clearing Jobarg Agentd log. Perhaps log file doesn't exist?"))
		return false
	}
	fmt.Println(testcase.Info_log("Info: Agent log is cleaned."))
	return true
}

// To clear jaz linux server log.
//
// Returns
//   - True if worked.
//   - False if failed.
func Run_Clear_Linux_Server_log(testcase *dao.TestCase) bool {
	err := lib.Clear_linux_jaz_server_log()
	if err != nil {
		fmt.Println(testcase.Err_log("Error: Failed at clearing Jobarg Server log. Perhaps log file doesn't exist?"))
		return false
	}
	fmt.Println(testcase.Info_log("Info: Server log is cleaned."))
	return true
}

// To set config files in linux.
//
// Parameter:
//   - Operations
//   - 1 for append
//   - 2 for replace.
//
// Returns
//   - True if worked.
//   - False if failed.
func Run_Set_Config_Linux(testcase *dao.TestCase, key string, value string, config_file_path string, operations int) bool {
	var err error
	switch operations {
	case 1:
		err = lib.Ja_set_config_linux(key, value, config_file_path)
	case 2:
		var command string
		command, err = lib.Ja_set_config_linux_str_replace(key, value, config_file_path)
		fmt.Println(testcase.Info_log("Info: Executed command: %s", command))
	default:
		fmt.Println("Invalid parameter. Try again.")
		return false
	}
	if err != nil {
		fmt.Println(testcase.Err_log("Error: Failed at setting config. %s", err.Error()))
		return false
	}
	fmt.Println(testcase.Info_log("Info: Config set successfully."))
	return true
}

// To run linux command
//
// Returns
//   - True if worked.
//   - False if failed.
func Run_Linux_Command(testcase *dao.TestCase, executeCommand string) bool {
	executeResult, err := lib.Ssh_exec_to_str(executeCommand)
	if err != nil {
		fmt.Println(testcase.Err_log("Error: Failed at executing linux command. Command executed: %s, Error: %s", executeCommand, err.Error()))
		executeResult = ""
		return false
	}
	fmt.Println(testcase.Info_log("Info: Command executed successfully: %s", executeCommand))
	if executeResult != "" && executeResult != " " {
		fmt.Print(testcase.Info_log("Info: Execution result: %s", strings.Trim(executeResult, "\n")))
	}
	return true
}

// To run linux command and return string
//
// Returns
//   - True and command result as string
//   - False and command result as string
func Run_Linux_Command_Str(testcase *dao.TestCase, executeCommand string) (bool, string) {
	result, err := lib.Ssh_exec_to_str(executeCommand)
	if err != nil {
		fmt.Println(testcase.Err_log("Error: Failed at executing linux command. Command executed: %s, Error: %s", executeCommand, err.Error()))
		result = ""
		return false, result
	}
	fmt.Println(testcase.Info_log("Info: Command executed successfully: %s", executeCommand))
	return true, result
}

// To run window command
//
// Returns
//   - True if worked.
//   - False if failed.
func Run_Window_Command(testcase *dao.TestCase, executeCommand string) bool {
	executeResult, err := lib.Execute_cmd_window_str(executeCommand)
	if err != nil {
		fmt.Println(testcase.Err_log("Error: Failed at executing window command. Command executed: %s, Error: %s", executeCommand, err.Error()))
		executeResult = ""
		return false
	}
	fmt.Println(testcase.Info_log("Info: Command executed successfully: %s", executeCommand))
	if executeResult != "" && executeResult != " " {
		fmt.Print(testcase.Info_log("Info: Execution result: %s", strings.Trim(executeResult, "\n")))
	}
	return true
}

// To run window command
//
// Returns
//   - True and command result as string
//   - False and command result as string
func Run_Window_Command_Str(testcase *dao.TestCase, executeCommand string) (bool, string) {
	executeResult, err := lib.Execute_cmd_window_str(executeCommand)
	if err != nil {
		fmt.Println(testcase.Err_log("Error: Failed at executing window command. Command executed: %s, Error: %s", executeCommand, err.Error()))
		executeResult = ""
		return false, executeResult
	}
	fmt.Println(testcase.Info_log("Info: Command executed successfully: %s", executeCommand))
	if executeResult != "" && executeResult != " " {
		fmt.Print(testcase.Info_log("Info: Execution result: %s", strings.Trim(executeResult, "\n")))
	}
	return true, executeResult
}

// To execute Database Sql Script.
//
// Returns
//
//   - True if worked
//   - False if it didn't
func Run_Sql_Script(testcase *dao.TestCase, execute_Sql_Script string) bool {
	_, err := lib.ExecuteQuery(lib.DBQuery(execute_Sql_Script))
	if err != nil {
		fmt.Println(testcase.Err_log("Error: Failed at executing following Sql script: %s Error: %s", execute_Sql_Script, err.Error()))
		return false
	}
	fmt.Println(testcase.Info_log("Info: Sql script ran successfully: %s", execute_Sql_Script))
	return true
}

// To execute Database Sql Script and get rows effected.
//
// Returns
//
//   - True and returns sql.result if worked
//   - False and returns sql.result if worked
func Run_Sql_Script_Return_Rows(testcase *dao.TestCase, execute_Sql_Script string) (bool, *sql.Rows) {
	rows, err := lib.GetData(lib.DBQuery(execute_Sql_Script))
	if err != nil {
		fmt.Println(testcase.Err_log("Error: Failed at executing following Sql script: %s Error: %s", execute_Sql_Script, err.Error()))
		return false, rows
	}
	fmt.Println(testcase.Info_log("Info: Sql script ran successfully: %s", execute_Sql_Script))
	return true, rows
}

// To convert from string to int
//
// Returns
//   - True and number if worked.
//   - False and 0 if failed.
func Str_To_Int(testcase *dao.TestCase, numberStr string) (bool, int) {
	numberInt, err := strconv.Atoi(numberStr)
	if err != nil {
		fmt.Println(testcase.Err_log("Failed to conver String to integer. %s", err.Error()))
		return false, 0
	}
	return true, numberInt
}

// Checking 'Job_Status' to get the value 'Error' with 'Jobarg_get' command."
//
// Returns
//   - True if Passed.
//   - False if Not.
func Run_Job_Status_Check_For_Error(testcase *dao.TestCase, timeout_time int, jobnet_run_manage_id string) bool {
	// Checking 'Job_Status' to get the value 'Error' with 'Jobarg_get' command."
	fmt.Println(testcase.Info_log("Info: Total 0 Job process count met. Waiting for Jobnet to end with RUNERR.(Running Error)"))

	// Setting timeout duration
	timeoutDuration := time.Duration(timeout_time) * time.Minute
	timeout := time.Now().Add(timeoutDuration)

	// taking current time snapshot
	start := time.Now()

	for time.Now().Before(timeout) {

		// Calculating elapsed time.
		elapsed := time.Since(start)

		// Extract hours, minutes, and seconds from elapsed time
		hours := int(elapsed.Hours())
		minutes := int(elapsed.Minutes()) % 60
		seconds := int(elapsed.Seconds()) % 60

		jobstatus, err := lib.Jobarg_get_JA_JOBSTATUS(jobnet_run_manage_id)
		if err == nil && jobstatus == "ERROR" {
			fmt.Println(testcase.Info_log("Info: The job returns with Error. Success."))

			return true
		}
		// Print in HH:MM:SS format with \r to overwrite the line
		fmt.Printf("\rElapsed Time: %02d:%02d:%02d", hours, minutes, seconds)
		time.Sleep(1 * time.Millisecond)
	}
	return false
}

// Wait until it reaches a specified process count
//
// Parameters:
//   - targetProcessCount: target process count to be reached
//   - timeoutDuration: timeout (minutes) for the process count checking
//   - client: ssh client
//
// Returns:
//
//	error message when the process does not reach the target count
//
// Ctl + C programmer. 99.99% copied from library. Don't wannt directly modify the library since I want to use it as template so I copied it and add the output line.
func JobProcessCountCheck_with_process_counter(targetProcessCount int, timeoutDuration int, client *ssh.Client) error {
	// set timeout
	timeout := time.After(time.Duration(timeoutDuration) * time.Minute)
	var index int
	for {
		select {
		case <-timeout:
			return fmt.Errorf("timeout after %d minutes", timeoutDuration)
		default:
			currentProcessCountStr, err := lib.GetOutputStrFromSSHCommand(client, "ps -aux | grep /etc/jobarranger/extendedjob/ | grep -v grep | wc -l")

			if err != nil {
				return fmt.Errorf("failed to obtain process count: %s", err.Error())
			}

			// Check the current job process count if it reaches the specified count
			currentProcessCountStr = strings.TrimSpace(currentProcessCountStr)

			currentProcessCount, err := strconv.Atoi(currentProcessCountStr)

			if err != nil {
				return fmt.Errorf("failed to convert the process count from string to int: %s", err.Error())
			}

			if currentProcessCount == targetProcessCount {
				return nil
			}
			fmt.Printf("\rProcess Count: %d", currentProcessCount)
			time.Sleep(1 * time.Second)
			index++
		}

	}
}

func Run_SFTP_File_Transfer(testcase *dao.TestCase, localFilePath string, remoteFilePath string) bool {
	//common.Set_passwd()
	//common.Set_client()

	// Create an SFTP client
	sftpClient, err := sftp.NewClient(common.Client)
	if err != nil {
		fmt.Println(testcase.Err_log("Error: Failed to create SFTP client: %s", err))
		return false
	}
	defer sftpClient.Close()
	//defer common.Client.Close()

	// Open the local file
	localFile, err := os.Open(localFilePath)
	if err != nil {
		fmt.Println(testcase.Err_log("Error: Failed to open local file: %s", err))
		return false
	}
	defer localFile.Close()

	// Create the remote file
	remoteFile, err := sftpClient.Create(remoteFilePath)
	if err != nil {
		fmt.Println(testcase.Err_log("Error: Failed to create remote file: %s", err))
		return false
	}
	defer remoteFile.Close()

	// Copy the file from local to remote
	bytes, err := io.ReadAll(localFile)
	if err != nil {
		fmt.Println(testcase.Err_log("Error: Failed to read local file: %s", err))
		return false
	}
	_, err = remoteFile.Write(bytes)
	if err != nil {
		fmt.Println(testcase.Err_log("Error: Failed to write to remote file: %s", err))
		return false
	}

	fmt.Println(testcase.Info_log("Info: File transferred successfully!"))

	defer remoteFile.Close()

	return true
}

// Put the Go program in timeout in seconds.
func Run_Timeout(testcase *dao.TestCase, timeout int) bool {
	var index int
	// taking current time snapshot
	start := time.Now()

	timeout_seconds := timeout % 60       // Remaining seconds
	timeout_minute := (timeout / 60) % 60 // Remaining minutes
	timeout_hour := timeout / 3600        // Total hours

	for i := 0; i < timeout; i++ {

		// Calculating elapsed time.
		elapsed := time.Since(start)

		// Extract hours, minutes, and seconds from elapsed time
		hours := int(elapsed.Hours())
		minutes := int(elapsed.Minutes()) % 60
		seconds := int(elapsed.Seconds()) % 60

		lib.Spinner_log(index, lib.Formatted_log(common.LOG_LEVEL_INFO, "Program in timeout. Elapsed Time: %02d:%02d:%02d|%02d:%02d:%02d", timeout_hour, timeout_minute, timeout_seconds, hours, minutes, seconds))

		time.Sleep(1 * time.Second)
		index++
	}
	return true
}
