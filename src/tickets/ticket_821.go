package tickets

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/zukigit/remote_run-go/src/common"
	"github.com/zukigit/remote_run-go/src/lib"
	"golang.org/x/crypto/ssh"
)

type Ticket_821 struct {
	Ticket_no                                   uint
	Ticket_description                          string
	PASSED_count, FAILED_count, MUSTCHECK_count int
	Testcases                                   []common.TestCase
}

func (t *Ticket_821) New_testcase(testcase_id uint, testcase_description string) *common.TestCase {
	return common.New_testcase(testcase_id, testcase_description)
}

func (t *Ticket_821) Get_no() uint {
	return t.Ticket_no
}

func (t *Ticket_821) Set_PASSED_count(passed_count int) {
	t.PASSED_count = passed_count
}

func (t *Ticket_821) Set_FAILED_count(failed_count int) {
	t.FAILED_count = failed_count
}

func (t *Ticket_821) Set_MUSTCHECK_count(mustcheck_count int) {
	t.MUSTCHECK_count = mustcheck_count
}

func (t *Ticket_821) Get_dsctn() string {
	return t.Ticket_description
}

func (t *Ticket_821) Add_testcase(tc common.TestCase) {
	t.Testcases = append(t.Testcases, tc)
}

func (t *Ticket_821) Get_testcases() []common.TestCase {
	return t.Testcases
}

// Enter your ticket information here
func (t *Ticket_821) Set_values() {
	t.Ticket_no = 821 // Enter your ticket id
	t.Ticket_description = "[Not Customer] Skip設定が入っているときの変数処理"
}

// Add your test case here
func (t *Ticket_821) Add_testcases() {
	// TESTCASE 001
	tc_1 := t.New_testcase(1, "Check the Jobnet is successully ran or not")
	tc_func := func() common.Testcase_status {
		// enable common jobnet
		if err := lib.Jobarg_enable_jobnet("Icon_1", "821Loop"); err != nil {
			lib.Logi(common.LOG_LEVEL_ERR, "Failed to enable jobnet, Error: %s", err)
			return FAILED
		}
		return CheckJobnetSuccess("Icon_1", tc_1, common.Client)

	}

	tc_1.Set_function(tc_func)
	t.Add_testcase(*tc_1)

	// TESTCASE 002
	tc_2 := t.New_testcase(2, "Jobarranger Agent folder and file creation Process [Agent]")
	tc_func = func() common.Testcase_status {
		// enable common jobnet
		if err := lib.Jobarg_enable_jobnet("Icon_1", "821Loop"); err != nil {
			lib.Logi(common.LOG_LEVEL_ERR, "Failed to enable jobnet, Error: %s", err)
			return FAILED
		}
		return CheckJobarrangerAgentfolderandFileCreationProcess("Icon_1", tc_2, common.Client)

	}

	tc_2.Set_function(tc_func)
	t.Add_testcase(*tc_2)

	// TESTCASE 003
	tc_3 := t.New_testcase(3, "Jobarranger server Process [Server]")
	tc_func = func() common.Testcase_status {
		// enable common jobnet
		if err := lib.Jobarg_enable_jobnet("Icon_1", "821Loop"); err != nil {
			lib.Logi(common.LOG_LEVEL_ERR, "Failed to enable jobnet, Error: %s", err)
			return FAILED
		}
		return CheckJobarrangerserverProcess("Icon_1", tc_3, common.Client)

	}

	tc_3.Set_function(tc_func)
	t.Add_testcase(*tc_3)

	// TESTCASE 004
	tc_4 := t.New_testcase(4, "Check the data recovery process is properly runnning or not when the index files are accidently deleted by user or corrupted.")
	tc_func = func() common.Testcase_status {
		// enable common jobnet
		if err := lib.Jobarg_enable_jobnet("Icon_1", "821Loop"); err != nil {
			lib.Logi(common.LOG_LEVEL_ERR, "Failed to enable jobnet, Error: %s", err)
			return FAILED
		}
		return ChecktheDatarecoveryProcess("Icon_1", tc_4, common.Client)

	}

	tc_4.Set_function(tc_func)
	t.Add_testcase(*tc_4)

	// TESTCASE 004
	tc_5 := t.New_testcase(5, "Abnormal case (DB down retry count 20)")
	tc_func = func() common.Testcase_status {
		// enable common jobnet
		if err := lib.Jobarg_enable_jobnet("Icon_1", "821Loop"); err != nil {
			lib.Logi(common.LOG_LEVEL_ERR, "Failed to enable jobnet, Error: %s", err)
			return FAILED
		}
		return AbnormalcaseDBdownretrycount20("Icon_1", tc_5, common.Client)

	}

	tc_5.Set_function(tc_func)
	t.Add_testcase(*tc_5)
}

func CheckJobnetSuccess(jobnetId string, testcase *common.TestCase, sshClient *ssh.Client) common.Testcase_status {

	/******************
	Pre-Operation State
	******************/
	lib.Jobarg_cleanup_linux()

	/**************
	Operation State
	***************/
	// Run jobnet
	run_jobnet_id, error := lib.Jobarg_exec(jobnetId)
	if error != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error: %s, std_out: %s", error.Error(), run_jobnet_id))
		return FAILED
	}
	fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "Job Icon %s has been successfully run with registry number: %s", jobnetId, run_jobnet_id))
	fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "Waiting for the job complete..."))
	sleepDuration := 3 * time.Minute
	time.Sleep(sleepDuration)

	defaulQuery := `SELECT a.after_value AS after_value_LOOPCNT,
		   b.after_value AS after_value_STD_OUT
	FROM ja_run_value_after_table AS a
	JOIN ja_run_value_after_table AS b
	  ON a.inner_job_id = b.inner_job_id
	  AND a.inner_jobnet_id = b.inner_jobnet_id
	WHERE a.after_value = 20
	  AND b.after_value = 20
	  AND a.value_name = 'LOOPCNT'
	  AND b.value_name = 'STD_OUT';`

	results, err := GetLoopCntAndStdOut(defaulQuery)
	if err != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error : %s", err.Error()))
		return FAILED
	}

	// /***************
	// Expected Results
	// ****************/
	// Check if we have any results
	if len(results) > 0 {
		// Print only the first result
		result := results[0] // Get the first result
		fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "The value of LOOPCNT: %d, and  STD_OUT: %d", result.AfterValueLoopcnt, result.AfterValueStdOut))
	} else {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "STD_OUT and LOOPCNT data not found."))
		return FAILED
	}
	return PASSED
}

func CheckJobarrangerAgentfolderandFileCreationProcess(jobnetId string, testcase *common.TestCase, sshClient *ssh.Client) common.Testcase_status {

	/******************
	Pre-Operation State
	******************/
	lib.Jobarg_cleanup_linux()
	// Restart the jobarg agent
	if err := lib.Restart_jaz_agent_linux(); err != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Faild to restart the JAZ agent, Error: %s", err.Error()))
		return FAILED
	}
	fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "JAZ agent has been restarted."))

	/**************
	Operation State
	***************/
	dirPath := "/var/lib/jobarranger/tmp/jobs"
	common.Client = lib.GetSSHClient(common.Login_info.Hostname, common.Login_info.Port, common.Login_info.Username, common.Login_info.Password)
	exists, err := lib.CheckRemoteDirectoryExists(common.Client, dirPath)
	if err != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error checking directory: %v\n", err))
		return FAILED
	}
	if exists {
		fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "Jobs Directory %s exists.", dirPath))
	} else {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Jobs Directory %s does not exist.", dirPath))
		return FAILED
	}

	err = lib.RemoveRemoteDirectory(common.Client, dirPath)
	if err != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error:%v", err))
	} else {
		fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "Jobs Directory removal successful."))
	}

	// Restart the jobarg agent
	if err := lib.Restart_jaz_agent_linux(); err != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Faild to restart the JAZ agent, Error: %s", err.Error()))
		return FAILED
	}
	fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "JAZ agent has been restarted."))

	exists, err = lib.CheckRemoteDirectoryExists(common.Client, dirPath)
	if err != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error checking directory: %v\n", err))
		return FAILED
	}
	if exists {
		fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "Jobs Directory %s exists.", dirPath))
	} else {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Jobs Directory %s does not exist.", dirPath))
		return FAILED
	}

	//Run jobnet
	run_jobnet_id, error := lib.Jobarg_exec(jobnetId)
	if error != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error: %s, std_out: %s", error.Error(), run_jobnet_id))
		return FAILED
	}
	fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "Job Icon %s has been successfully run with registry number: %s", jobnetId, run_jobnet_id))

	fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "Waiting for the .job file create in jobs folder..."))
	sleepDuration := 1 * time.Minute
	time.Sleep(sleepDuration)

	// Call the function to check if any .job files exist in the specified directory.
	exists, err = lib.CheckRemoteIndexFileExists(common.Client, dirPath)
	if err != nil {
		fmt.Println("Error:", err)
		return FAILED
	}

	if exists {
		fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "Index file is successfully created."))
	} else {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "No index file found in the jobs directory."))
		return FAILED
	}

	return PASSED
}

func CheckJobarrangerserverProcess(jobnetId string, testcase *common.TestCase, sshClient *ssh.Client) common.Testcase_status {

	/******************
	Pre-Operation State
	******************/
	lib.Jobarg_cleanup_linux()
	// Restart the jobarg agent
	if err := lib.Restart_jaz_server(); err != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Faild to restart the JAZ server, Error: %s", err.Error()))
		return FAILED
	}
	fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "JAZ server has been restarted."))

	/**************
	Operation State
	***************/
	dirPath := "/var/log/jobarranger/job"
	common.Client = lib.GetSSHClient(common.Login_info.Hostname, common.Login_info.Port, common.Login_info.Username, common.Login_info.Password)
	exists, err := lib.CheckRemoteDirectoryExists(common.Client, dirPath)
	if err != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error checking directory: %v\n", err))
		return FAILED
	}
	if exists {
		fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "Job Directory %s exists.", dirPath))
	} else {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Job Directory %s does not exist.", dirPath))
		return FAILED
	}

	err = lib.RemoveRemoteDirectory(common.Client, dirPath)
	if err != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error:%v", err))
	} else {
		fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "Job Directory removal successful."))
	}

	// Restart the jobarg server
	if err := lib.Restart_jaz_server(); err != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Faild to restart the JAZ server, Error: %s", err.Error()))
		return FAILED
	}
	fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "JAZ server has been restarted."))

	exists, err = lib.CheckRemoteDirectoryExists(common.Client, dirPath)
	if err != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error checking directory: %v\n", err))
		return FAILED
	}
	if exists {
		fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "Job Directory %s exists.", dirPath))
	} else {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Job Directory %s does not exist.", dirPath))
		return FAILED
	}

	//Run jobnet
	run_jobnet_id, error := lib.Jobarg_exec(jobnetId)
	if error != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error: %s, std_out: %s", error.Error(), run_jobnet_id))
		return FAILED
	}
	fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "Job Icon %s has been successfully run with registry number: %s", jobnetId, run_jobnet_id))

	fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "Waiting for the .job file create in jobs folder..."))
	sleepDuration := 1 * time.Minute
	time.Sleep(sleepDuration)

	// Call the function to check if any .job files exist in the specified directory.
	exists, err = lib.CheckRemoteIndexFileExists(common.Client, dirPath)
	if err != nil {
		fmt.Println("Error:", err)
		return FAILED
	}

	if exists {
		fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "Index file is successfully created."))
	} else {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "No index file found in the jobs directory."))
		return FAILED
	}

	return PASSED
}

func ChecktheDatarecoveryProcess(jobnetId string, testcase *common.TestCase, sshClient *ssh.Client) common.Testcase_status {

	/******************
	Pre-Operation State
	******************/
	lib.Jobarg_cleanup_linux()
	configFilePath := "/etc/jobarranger/jobarg_server.conf"
	err := lib.UpdateDebugLevel(common.Client, configFilePath, 3)
	if err != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error updating DebugLevel: %v", err))
	}
	lib.Restart_jaz_server()
	// Restart the jobarg agent
	if err := lib.Restart_jaz_agent_linux(); err != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Faild to restart the JAZ agent, Error: %s", err.Error()))
		return FAILED
	}
	fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "JAZ agent has been restarted."))

	/**************
	Operation State
	***************/
	dirPath := "/var/lib/jobarranger/tmp/jobs"
	common.Client = lib.GetSSHClient(common.Login_info.Hostname, common.Login_info.Port, common.Login_info.Username, common.Login_info.Password)

	result := lib.RemoveAllFilesInDirectory(common.Client, dirPath)
	if result != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error: Remove all file within directory"))
		return FAILED
	}

	fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "All files within directory %s successfully removed.", dirPath))

	if err := lib.Jobarg_enable_jobnet("Icon_1", "821Loop"); err != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Failed to enable jobnet, Error: %s", err))
		return FAILED
	}

	//Run jobnet
	run_jobnet_id, error := lib.Jobarg_exec(jobnetId)
	if error != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error: %s, std_out: %s", error.Error(), run_jobnet_id))
		return FAILED
	}
	fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "Job Icon %s has been successfully run with registry number: %s", jobnetId, run_jobnet_id))

	fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "Waiting for the .job file create in agent jobs folder..."))
	sleepDuration := 1 * time.Minute
	time.Sleep(sleepDuration)

	result = lib.RemoveAllFilesInDirectory(common.Client, dirPath)
	if result != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error: Remove index file within %s when running jobnet", dirPath))
		return FAILED
	} else {
		fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "Successfully Remove index file within %s when running jobnet", dirPath))

	}

	pattern := "Create recovery file"
	timeout := 30 * time.Second // Timeout duration
	interval := 1 * time.Second // Polling interval
	logFilePath := "/var/log/jobarranger/jobarg_server.log"
	lib.ClearLogFile(common.Client, logFilePath)

	_, err = lib.WaitForPatternInLogFile(common.Client, logFilePath, pattern, timeout, interval)
	if err != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error:%s", err))
	} else {
		fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "'Create recovery' alert message found in server log."))
	}
	lib.Jobarg_cleanup_linux()

	dirPath = "/var/log/jobarranger/job"
	common.Client = lib.GetSSHClient(common.Login_info.Hostname, common.Login_info.Port, common.Login_info.Username, common.Login_info.Password)

	result = lib.RemoveAllFilesInDirectory(common.Client, dirPath)
	if result != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error: Remove all file within directory"))
		return FAILED
	}

	fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "All files within directory %s successfully removed.", dirPath))

	if err := lib.Jobarg_enable_jobnet("Icon_1", "821Loop"); err != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Failed to enable jobnet, Error: %s", err))
		return FAILED
	}

	//Run jobnet
	run_jobnet_id, error = lib.Jobarg_exec(jobnetId)
	if error != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error: %s, std_out: %s", error.Error(), run_jobnet_id))
		return FAILED
	}
	fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "Job Icon %s has been successfully run with registry number: %s", jobnetId, run_jobnet_id))

	fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "Waiting for the .job file create in server job folder..."))
	sleepDuration = 1 * time.Minute
	time.Sleep(sleepDuration)

	result = lib.RemoveAllFilesInDirectory(common.Client, dirPath)
	if result != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error: Remove index file within %s when running jobnet", dirPath))
		return FAILED
	} else {
		fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "Successfully Remove index file within %s when running jobnet", dirPath))
	}
	pattern = "Created backup file"
	timeout = 30 * time.Second // Timeout duration
	interval = 1 * time.Second // Polling interval
	logFilePath = "/var/log/jobarranger/jobarg_server.log"
	lib.ClearLogFile(common.Client, logFilePath)

	_, err = lib.WaitForPatternInLogFile(common.Client, logFilePath, pattern, timeout, interval)
	if err != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error:%s", err))
	} else {
		fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "'Created backup file' alert message found in server log."))
	}

	lib.Jobarg_cleanup_linux()

	dirPath1 := "/var/log/jobarranger/job"
	lib.RemoveAllFilesInDirectory(common.Client, dirPath1)
	lib.Restart_jaz_server()
	dirPath2 := " /var/lib/jobarranger/tmp/jobs"
	lib.RemoveAllFilesInDirectory(common.Client, dirPath1)
	lib.Restart_jaz_agent_linux()

	if err := lib.Jobarg_enable_jobnet("Icon_1", "821Loop"); err != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Failed to enable jobnet, Error: %s", err))
		return FAILED
	}

	//Run jobnet
	run_jobnet_id, error = lib.Jobarg_exec(jobnetId)
	if error != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error: %s, std_out: %s", error.Error(), run_jobnet_id))
		return FAILED
	}
	fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "Job Icon %s has been successfully run with registry number: %s", jobnetId, run_jobnet_id))

	fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "Wating for the .job file create in both server and agent job folder..."))
	sleepDuration = 1 * time.Minute
	time.Sleep(sleepDuration)

	// Call the FindandTrimServerJobFilePrefix function
	ServerJobFileName, err := lib.FindandTrimServerJobFilePrefix(common.Client, dirPath1)
	if err != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error finding job file in : %v", dirPath1))
		return FAILED
	}

	// Call the findJobFileName function
	AgentJobFileName, err := lib.FindandTrimAgentJobFilePrefix(common.Client, dirPath2)
	if err != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error finding job file in : %v", dirPath2))
		return FAILED
	}

	if ServerJobFileName == AgentJobFileName {
		fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "Both unique id is the same in server and agent job folder"))
	} else {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "No same unique id found"))
		return FAILED
	}

	return PASSED
}

func AbnormalcaseDBdownretrycount20(jobnetId string, testcase *common.TestCase, sshClient *ssh.Client) common.Testcase_status {

	/******************
	Pre-Operation State
	******************/
	configFilePath := "/etc/jobarranger/jobarg_server.conf"
	err := lib.UpdateDebugLevel(common.Client, configFilePath, 4)
	if err != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error updating DebugLevel: %v", err))
	}

	lib.Jobarg_cleanup_linux()
	logFilePath := "/var/log/jobarranger/jobarg_agentd.log"
	lib.ClearLogFile(common.Client, logFilePath)
	logFilePath = "/var/log/jobarranger/jobarg_server.log"
	lib.ClearLogFile(common.Client, logFilePath)

	dirPath := "/var/lib/jobarranger/tmp/jobs"
	common.Client = lib.GetSSHClient(common.Login_info.Hostname, common.Login_info.Port, common.Login_info.Username, common.Login_info.Password)

	result := lib.RemoveAllFilesInDirectory(common.Client, dirPath)
	if result != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error: Remove all file within agend directory"))
		return FAILED
	}

	dirPath = "/var/log/jobarranger/job"
	common.Client = lib.GetSSHClient(common.Login_info.Hostname, common.Login_info.Port, common.Login_info.Username, common.Login_info.Password)

	result = lib.RemoveAllFilesInDirectory(common.Client, dirPath)
	if result != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error: Remove all file within server directory"))
		return FAILED
	}

	/**************
	Operation State
	***************/

	if err := lib.Jobarg_enable_jobnet("Icon_1", "821Loop"); err != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Failed to enable jobnet, Error: %s", err))
		return FAILED
	}

	//Run jobnet
	run_jobnet_id, error := lib.Jobarg_exec(jobnetId)
	if error != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error: %s, std_out: %s", error.Error(), run_jobnet_id))
		return FAILED
	}
	fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "Job Icon %s has been successfully run with registry number: %s", jobnetId, run_jobnet_id))

	fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "Sleep the 1 minute to wait the next running job"))

	sleepDuration := 1 * time.Minute
	time.Sleep(sleepDuration)

	//defaulQuery := `select * from ja_run_job_table jrjt where jrjt.job_id ='NORMAL-JOB' and jrjt.status =2 and jrjt.job_type = 4`
	defaulQuery := `SELECT job_id, status, job_type FROM ja_run_job_table WHERE job_id='NORMAL-JOB' AND status=2 AND job_type=4`
	timeout := 2 * time.Minute
	interval := 1 * time.Second

	if common.DB_type == common.PSQL {
		lib.CheckAndStopDBService(common.DB, common.Client, "postgresql", defaulQuery, timeout, interval)
	} else {
		lib.CheckAndStopDBService(common.DB, common.Client, "mysql", defaulQuery, timeout, interval)
	}

	pattern := "[WARN] Database is down. Reconnecting in 10 seconds."
	logFilePath1 := "/var/log/jobarranger/jobarg_server.log"

	timeout = 1 * time.Minute  // Timeout duration
	interval = 1 * time.Second // Polling interval

	_, err = lib.WaitForPatternInLogFile(common.Client, logFilePath1, pattern, timeout, interval)
	if err != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error:%s", err))
	} else {
		fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "%s show in server log", pattern))
	}

	fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "Wait for the log info from agent log ..."))
	sleepDuration = 2 * time.Minute
	time.Sleep(sleepDuration)

	pattern = "[WARN] In ja_job_exec_close() agent close failed. retry count :[20]"
	logFilePath2 := "/var/log/jobarranger/jobarg_agentd.log"

	timeout = 1 * time.Minute  // Timeout duration
	interval = 1 * time.Second // Polling interval

	_, err = lib.WaitForPatternInLogFile(common.Client, logFilePath2, pattern, timeout, interval)
	if err != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error:%s", err))
	} else {
		fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "%s show in agent log", pattern))
	}

	fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "start the database service and check the no index file remain from server side ..."))

	if common.DB_type == common.PSQL {
		lib.StartDatabaseService(common.Client, "postgresql")
	} else {
		lib.StartDatabaseService(common.Client, "mysqld")
	}
	fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "Waiting and Checking index file remain in server side ..."))
	sleepDuration = 3 * time.Minute
	time.Sleep(sleepDuration)

	/////////////////////////////////////

	defaulQuery = `SELECT a.after_value AS after_value_LOOPCNT,
		   b.after_value AS after_value_STD_OUT
	FROM ja_run_value_after_table AS a
	JOIN ja_run_value_after_table AS b
	  ON a.inner_job_id = b.inner_job_id
	  AND a.inner_jobnet_id = b.inner_jobnet_id
	WHERE a.after_value = 20
	  AND b.after_value = 20
	  AND a.value_name = 'LOOPCNT'
	  AND b.value_name = 'STD_OUT';`

	results, err := GetLoopCntAndStdOut(defaulQuery)
	if err != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error : %s", err.Error()))
		return FAILED
	}

	// /***************
	// Expected Results
	// ****************/
	// Check if we have any results
	if len(results) > 0 {
		// Print only the first result
		result := results[0] // Get the first result
		fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "The value of LOOPCNT: %d, and  STD_OUT: %d", result.AfterValueLoopcnt, result.AfterValueStdOut))
	} else {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "STD_OUT and LOOPCNT data not found."))
		return FAILED
	}

	dbQuery := fmt.Sprintf(`SELECT inner_job_id, before_value, value_name, COUNT(*) as count, MIN(seq_no) as min_seq_no
FROM ja_run_value_before_table
WHERE inner_jobnet_id = '%s'
  AND (value_name = 'LOOPCNT' OR value_name = 'STD_OUT')
GROUP BY inner_job_id, before_value, value_name
HAVING COUNT(*) > 1
ORDER BY min_seq_no ASC;`, run_jobnet_id)

	rows, err := lib.GetData(lib.DBQuery(dbQuery))
	if err != nil {
		// If there was an error executing the query, print an error message
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error executing query: %v\n", err))
		configFilePath = "/etc/jobarranger/jobarg_server.conf"
		err = lib.UpdateDebugLevel(common.Client, configFilePath, 3)
		if err != nil {
			fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error updating DebugLevel: %v", err))
		}
		return FAILED
	}
	// Ensure the rows are closed after processing
	defer rows.Close()

	if !rows.Next() {
		fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "Data duplicate does not occur in the database"))
	} else {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Duplicate data found"))
		configFilePath = "/etc/jobarranger/jobarg_server.conf"
		err = lib.UpdateDebugLevel(common.Client, configFilePath, 3)
		if err != nil {
			fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error updating DebugLevel: %v", err))
		}
		return FAILED
	}

	fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "Wait for the duplicate data message in the serverlog"))

	pattern = "Duplicate Data"
	timeout = 30 * time.Second // Timeout duration
	interval = 1 * time.Second // Polling interval
	logFilePath = "/var/log/jobarranger/jobarg_server.log"

	_, err = lib.WaitForPatternInLogFile(common.Client, logFilePath, pattern, timeout, interval)
	if err != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error:%s", err))
	} else {
		fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "Duplicate Data found in server log"))
	}

	configFilePath = "/etc/jobarranger/jobarg_server.conf"
	err = lib.UpdateDebugLevel(common.Client, configFilePath, 3)
	if err != nil {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Error updating DebugLevel: %v", err))
	}

	dirPath = "/var/log/jobarranger/job"
	// Call the function to check if any .job files exist in the specified directory.
	exists, err := lib.CheckRemoteIndexFileExists(common.Client, dirPath)
	if err != nil {
		fmt.Println("Error:", err)
		return FAILED
	}

	if exists {
		fmt.Println(lib.Logi(common.LOG_LEVEL_ERR, "Index file remain"))
		return FAILED
	} else {
		fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "Index file does not remain in %s ", dirPath))
	}

	// lib.Jobarg_cleanup_linux()

	return PASSED
}

// Define a struct to hold the after_value_LOOPCNT and after_value_STD_OUT
type LoopCntAndStdOut struct {
	AfterValueLoopcnt int
	AfterValueStdOut  int
}

// Custom method to get both LOOPCNT and STD_OUT values in one method
func GetLoopCntAndStdOut(query string) ([]LoopCntAndStdOut, error) {
	dbQuery := lib.DBQuery(query) // Ensure to use the correct DBQuery from lib
	rows, err := lib.GetData(dbQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []LoopCntAndStdOut // Slice of structs to hold the results

	// Process the result rows
	for rows.Next() {
		var afterValueLoopcnt []byte // Scan after_value_LOOPCNT as []byte (raw data)
		var afterValueStdOut []byte  // Scan after_value_STD_OUT as []byte (raw data)

		// Scan the relevant columns
		err := rows.Scan(&afterValueLoopcnt, &afterValueStdOut)
		if err != nil {
			return nil, err
		}

		// Convert both after_value_LOOPCNT and after_value_STD_OUT from []byte to int in one method
		afterValueLoopcntInt, err := strconv.Atoi(strings.TrimSpace(string(afterValueLoopcnt)))
		if err != nil {
			return nil, fmt.Errorf("failed to convert after_value_LOOPCNT to int: %w", err)
		}

		afterValueStdOutInt, err := strconv.Atoi(strings.TrimSpace(string(afterValueStdOut)))
		if err != nil {
			return nil, fmt.Errorf("failed to convert after_value_STD_OUT to int: %w", err)
		}

		// Append the results as a struct
		results = append(results, LoopCntAndStdOut{
			AfterValueLoopcnt: afterValueLoopcntInt,
			AfterValueStdOut:  afterValueStdOutInt,
		})
	}

	// Check if any errors occurred during the iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}
