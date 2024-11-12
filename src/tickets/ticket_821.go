package tickets

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/zukigit/remote_run-go/src/common"
	"github.com/zukigit/remote_run-go/src/dao"
	"github.com/zukigit/remote_run-go/src/lib"
	"golang.org/x/crypto/ssh"
)

type Ticket_821 struct {
	no          uint
	description string
	testcases   []dao.TestCase
}

func (t *Ticket_821) New_testcase(testcase_id uint, testcase_description string) *dao.TestCase {
	return dao.New_testcase(testcase_id, testcase_description)
}

func (t *Ticket_821) Get_no() uint {
	return t.no
}

func (t *Ticket_821) Get_dsctn() string {
	return t.description
}

func (t *Ticket_821) Add_testcase(tc dao.TestCase) {
	t.testcases = append(t.testcases, tc)
}

func (t *Ticket_821) Get_testcases() []dao.TestCase {
	return t.testcases
}

// Enter your ticket information here
func (t *Ticket_821) Set_values() {
	t.no = 821 // Enter your ticket id
	t.description = "[Not Customer] Skip設定が入っているときの変数処理"
}

// Add your test case here
func (t *Ticket_821) Add_testcases() {
	// TESTCASE 001
	tc_1 := t.New_testcase(1, "Check the Jobnet is successully ran or not")
	tc_func := func() common.Testcase_status {
		// enable common jobnet
		if err := lib.Jobarg_enable_jobnet("Icon_1", "821Loop"); err != nil {
			tc_1.Err_log("Failed to enable jobnet, Error: %s", err)
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
			tc_2.Err_log("Failed to enable jobnet, Error: %s", err)
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
			tc_3.Err_log("Failed to enable jobnet, Error: %s", err)
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
			tc_4.Err_log("Failed to enable jobnet, Error: %s", err)
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
			tc_5.Err_log("Failed to enable jobnet, Error: %s", err)
			return FAILED
		}
		return AbnormalcaseDBdownretrycount20("Icon_1", tc_5, common.Client)

	}

	tc_5.Set_function(tc_func)
	t.Add_testcase(*tc_5)
}

func CheckJobnetSuccess(jobnetId string, testcase *dao.TestCase, sshClient *ssh.Client) common.Testcase_status {

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
		fmt.Println(testcase.Err_log("Error: %s, std_out: %s", error.Error(), run_jobnet_id))
		return FAILED
	}
	fmt.Println(testcase.Info_log("Job Icon %s has been successfully run with registry number: %s", jobnetId, run_jobnet_id))
	fmt.Println(testcase.Info_log("Wating for the job complete..."))
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
		fmt.Println(testcase.Err_log("Error : %s", err.Error()))
		return FAILED
	}

	// /***************
	// Expected Results
	// ****************/
	// Check if we have any results
	if len(results) > 0 {
		// Print only the first result
		result := results[0] // Get the first result
		fmt.Println(testcase.Info_log("The value of LOOPCNT: %d, and  STD_OUT: %d", result.AfterValueLoopcnt, result.AfterValueStdOut))
	} else {
		fmt.Println(testcase.Err_log("STD_OUT and LOOPCNT data not found."))
		return FAILED
	}
	return PASSED
}

func CheckJobarrangerAgentfolderandFileCreationProcess(jobnetId string, testcase *dao.TestCase, sshClient *ssh.Client) common.Testcase_status {

	/******************
	Pre-Operation State
	******************/
	lib.Jobarg_cleanup_linux()
	// Restart the jobarg agent
	if err := lib.Restart_jaz_agent_linux(); err != nil {
		fmt.Println(testcase.Err_log("Faild to restart the JAZ agent, Error: %s", err.Error()))
		return FAILED
	}
	fmt.Println(testcase.Info_log("JAZ agent has been restarted."))

	/**************
	Operation State
	***************/
	dirPath := "/var/lib/jobarranger/tmp/jobs"
	common.Client = lib.GetSSHClient(common.Login_info.Hostname, common.Login_info.Port, common.Login_info.Username, common.Login_info.Password)
	exists, err := lib.CheckRemoteDirectoryExists(common.Client, dirPath)
	if err != nil {
		fmt.Println(testcase.Err_log("Error checking directory: %v\n", err))
		return FAILED
	}
	if exists {
		fmt.Println(testcase.Info_log("Jobs Directory %s exists.", dirPath))
	} else {
		fmt.Println(testcase.Err_log("Jobs Directory %s does not exist.", dirPath))
		return FAILED
	}

	err = lib.RemoveRemoteDirectory(common.Client, dirPath)
	if err != nil {
		fmt.Println(testcase.Err_log("Error:%v", err))
	} else {
		fmt.Println(testcase.Info_log("Jobs Directory removal successful."))
	}

	// Restart the jobarg agent
	if err := lib.Restart_jaz_agent_linux(); err != nil {
		fmt.Println(testcase.Err_log("Faild to restart the JAZ agent, Error: %s", err.Error()))
		return FAILED
	}
	fmt.Println(testcase.Info_log("JAZ agent has been restarted."))

	exists, err = lib.CheckRemoteDirectoryExists(common.Client, dirPath)
	if err != nil {
		fmt.Println(testcase.Err_log("Error checking directory: %v\n", err))
		return FAILED
	}
	if exists {
		fmt.Println(testcase.Info_log("Jobs Directory %s exists.", dirPath))
	} else {
		fmt.Println(testcase.Err_log("Jobs Directory %s does not exist.", dirPath))
		return FAILED
	}

	//Run jobnet
	run_jobnet_id, error := lib.Jobarg_exec(jobnetId)
	if error != nil {
		fmt.Println(testcase.Err_log("Error: %s, std_out: %s", error.Error(), run_jobnet_id))
		return FAILED
	}
	fmt.Println(testcase.Info_log("Job Icon %s has been successfully run with registry number: %s", jobnetId, run_jobnet_id))

	fmt.Println(testcase.Info_log("Wating for the .job file create in jobs folder..."))
	sleepDuration := 1 * time.Minute
	time.Sleep(sleepDuration)

	// Call the function to check if any .job files exist in the specified directory.
	exists, err = lib.CheckRemoteIndexFileExists(common.Client, dirPath)
	if err != nil {
		fmt.Println("Error:", err)
		return FAILED
	}

	if exists {
		fmt.Println(testcase.Info_log("Index file is successfully created."))
	} else {
		fmt.Println(testcase.Err_log("No index file found in the jobs directory."))
		return FAILED
	}

	return PASSED
}

func CheckJobarrangerserverProcess(jobnetId string, testcase *dao.TestCase, sshClient *ssh.Client) common.Testcase_status {

	/******************
	Pre-Operation State
	******************/
	lib.Jobarg_cleanup_linux()
	// Restart the jobarg agent
	if err := lib.Restart_jaz_server(); err != nil {
		fmt.Println(testcase.Err_log("Faild to restart the JAZ server, Error: %s", err.Error()))
		return FAILED
	}
	fmt.Println(testcase.Info_log("JAZ server has been restarted."))

	/**************
	Operation State
	***************/
	dirPath := "/var/log/jobarranger/job"
	common.Client = lib.GetSSHClient(common.Login_info.Hostname, common.Login_info.Port, common.Login_info.Username, common.Login_info.Password)
	exists, err := lib.CheckRemoteDirectoryExists(common.Client, dirPath)
	if err != nil {
		fmt.Println(testcase.Err_log("Error checking directory: %v\n", err))
		return FAILED
	}
	if exists {
		fmt.Println(testcase.Info_log("Job Directory %s exists.", dirPath))
	} else {
		fmt.Println(testcase.Err_log("Job Directory %s does not exist.", dirPath))
		return FAILED
	}

	err = lib.RemoveRemoteDirectory(common.Client, dirPath)
	if err != nil {
		fmt.Println(testcase.Err_log("Error:%v", err))
	} else {
		fmt.Println(testcase.Info_log("Job Directory removal successful."))
	}

	// Restart the jobarg server
	if err := lib.Restart_jaz_server(); err != nil {
		fmt.Println(testcase.Err_log("Faild to restart the JAZ server, Error: %s", err.Error()))
		return FAILED
	}
	fmt.Println(testcase.Info_log("JAZ server has been restarted."))

	exists, err = lib.CheckRemoteDirectoryExists(common.Client, dirPath)
	if err != nil {
		fmt.Println(testcase.Err_log("Error checking directory: %v\n", err))
		return FAILED
	}
	if exists {
		fmt.Println(testcase.Info_log("Job Directory %s exists.", dirPath))
	} else {
		fmt.Println(testcase.Err_log("Job Directory %s does not exist.", dirPath))
		return FAILED
	}

	//Run jobnet
	run_jobnet_id, error := lib.Jobarg_exec(jobnetId)
	if error != nil {
		fmt.Println(testcase.Err_log("Error: %s, std_out: %s", error.Error(), run_jobnet_id))
		return FAILED
	}
	fmt.Println(testcase.Info_log("Job Icon %s has been successfully run with registry number: %s", jobnetId, run_jobnet_id))

	fmt.Println(testcase.Info_log("Wating for the .job file create in jobs folder..."))
	sleepDuration := 1 * time.Minute
	time.Sleep(sleepDuration)

	// Call the function to check if any .job files exist in the specified directory.
	exists, err = lib.CheckRemoteIndexFileExists(common.Client, dirPath)
	if err != nil {
		fmt.Println("Error:", err)
		return FAILED
	}

	if exists {
		fmt.Println(testcase.Info_log("Index file is successfully created."))
	} else {
		fmt.Println(testcase.Err_log("No index file found in the jobs directory."))
		return FAILED
	}

	return PASSED
}

func ChecktheDatarecoveryProcess(jobnetId string, testcase *dao.TestCase, sshClient *ssh.Client) common.Testcase_status {

	/******************
	Pre-Operation State
	******************/
	lib.Jobarg_cleanup_linux()
	// Restart the jobarg agent
	if err := lib.Restart_jaz_agent_linux(); err != nil {
		fmt.Println(testcase.Err_log("Faild to restart the JAZ agent, Error: %s", err.Error()))
		return FAILED
	}
	fmt.Println(testcase.Info_log("JAZ agent has been restarted."))

	/**************
	Operation State
	***************/
	dirPath := "/var/lib/jobarranger/tmp/jobs"
	common.Client = lib.GetSSHClient(common.Login_info.Hostname, common.Login_info.Port, common.Login_info.Username, common.Login_info.Password)

	result := lib.RemoveAllFilesInDirectory(common.Client, dirPath)
	if result != nil {
		fmt.Println(testcase.Err_log("Error: Remove all file within directory"))
		return FAILED
	}

	fmt.Println(testcase.Info_log("All files within directory %s successfully removed.", dirPath))

	if err := lib.Jobarg_enable_jobnet("Icon_1", "821Loop"); err != nil {
		fmt.Println(testcase.Err_log("Failed to enable jobnet, Error: %s", err))
		return FAILED
	}

	//Run jobnet
	run_jobnet_id, error := lib.Jobarg_exec(jobnetId)
	if error != nil {
		fmt.Println(testcase.Err_log("Error: %s, std_out: %s", error.Error(), run_jobnet_id))
		return FAILED
	}
	fmt.Println(testcase.Info_log("Job Icon %s has been successfully run with registry number: %s", jobnetId, run_jobnet_id))

	fmt.Println(testcase.Info_log("Wating for the .job file create in agent jobs folder..."))
	sleepDuration := 1 * time.Minute
	time.Sleep(sleepDuration)

	result = lib.RemoveAllFilesInDirectory(common.Client, dirPath)
	if result != nil {
		fmt.Println(testcase.Err_log("Error: Remove index file within %s when running jobnet", dirPath))
		return FAILED
	} else {
		fmt.Println(testcase.Info_log("Successfully Remove index file within %s when running jobnet", dirPath))

	}

	pattern := "Create recovery file"
	timeout := 30 * time.Second // Timeout duration
	interval := 1 * time.Second // Polling interval
	logFilePath := "/var/log/jobarranger/jobarg_server.log"
	lib.ClearLogFile(common.Client, logFilePath)

	_, err := lib.WaitForPatternInLogFile(common.Client, logFilePath, pattern, timeout, interval)
	if err != nil {
		fmt.Println(testcase.Err_log("Error:%s", err))
	} else {
		fmt.Println(testcase.Info_log("'Create recovery' alert message found in server log."))
	}
	lib.Jobarg_cleanup_linux()

	dirPath = "/var/log/jobarranger/job"
	common.Client = lib.GetSSHClient(common.Login_info.Hostname, common.Login_info.Port, common.Login_info.Username, common.Login_info.Password)

	result = lib.RemoveAllFilesInDirectory(common.Client, dirPath)
	if result != nil {
		fmt.Println(testcase.Err_log("Error: Remove all file within directory"))
		return FAILED
	}

	fmt.Println(testcase.Info_log("All files within directory %s successfully removed.", dirPath))

	if err := lib.Jobarg_enable_jobnet("Icon_1", "821Loop"); err != nil {
		fmt.Println(testcase.Err_log("Failed to enable jobnet, Error: %s", err))
		return FAILED
	}

	//Run jobnet
	run_jobnet_id, error = lib.Jobarg_exec(jobnetId)
	if error != nil {
		fmt.Println(testcase.Err_log("Error: %s, std_out: %s", error.Error(), run_jobnet_id))
		return FAILED
	}
	fmt.Println(testcase.Info_log("Job Icon %s has been successfully run with registry number: %s", jobnetId, run_jobnet_id))

	fmt.Println(testcase.Info_log("Wating for the .job file create in server job folder..."))
	sleepDuration = 1 * time.Minute
	time.Sleep(sleepDuration)

	result = lib.RemoveAllFilesInDirectory(common.Client, dirPath)
	if result != nil {
		fmt.Println(testcase.Err_log("Error: Remove index file within %s when running jobnet", dirPath))
		return FAILED
	} else {
		fmt.Println(testcase.Info_log("Successfully Remove index file within %s when running jobnet", dirPath))
	}
	pattern = "Created backup file"
	timeout = 30 * time.Second // Timeout duration
	interval = 1 * time.Second // Polling interval
	logFilePath = "/var/log/jobarranger/jobarg_server.log"
	lib.ClearLogFile(common.Client, logFilePath)

	_, err = lib.WaitForPatternInLogFile(common.Client, logFilePath, pattern, timeout, interval)
	if err != nil {
		fmt.Println(testcase.Err_log("Error:%s", err))
	} else {
		fmt.Println(testcase.Info_log("'Created backup file' alert message found in server log."))
	}

	lib.Jobarg_cleanup_linux()

	dirPath1 := "/var/log/jobarranger/job"
	lib.RemoveAllFilesInDirectory(common.Client, dirPath1)
	lib.Restart_jaz_server()
	dirPath2 := " /var/lib/jobarranger/tmp/jobs"
	lib.RemoveAllFilesInDirectory(common.Client, dirPath1)
	lib.Restart_jaz_agent_linux()

	if err := lib.Jobarg_enable_jobnet("Icon_1", "821Loop"); err != nil {
		fmt.Println(testcase.Err_log("Failed to enable jobnet, Error: %s", err))
		return FAILED
	}

	//Run jobnet
	run_jobnet_id, error = lib.Jobarg_exec(jobnetId)
	if error != nil {
		fmt.Println(testcase.Err_log("Error: %s, std_out: %s", error.Error(), run_jobnet_id))
		return FAILED
	}
	fmt.Println(testcase.Info_log("Job Icon %s has been successfully run with registry number: %s", jobnetId, run_jobnet_id))

	fmt.Println(testcase.Info_log("Wating for the .job file create in both server and agent job folder..."))
	sleepDuration = 1 * time.Minute
	time.Sleep(sleepDuration)

	// Call the FindandTrimServerJobFilePrefix function
	ServerJobFileName, err := lib.FindandTrimServerJobFilePrefix(common.Client, dirPath1)
	if err != nil {
		fmt.Println(testcase.Err_log("Error finding job file in : %v", dirPath1))
		return FAILED
	}

	// Call the findJobFileName function
	AgentJobFileName, err := lib.FindandTrimAgentJobFilePrefix(common.Client, dirPath2)
	if err != nil {
		fmt.Println(testcase.Err_log("Error finding job file in : %v", dirPath2))
		return FAILED
	}

	if ServerJobFileName == AgentJobFileName {
		fmt.Println(testcase.Info_log("Both unique id is the same in server and agent job folder"))
	} else {
		fmt.Println(testcase.Err_log("No same unique id found"))
		return FAILED
	}

	return PASSED
}

func AbnormalcaseDBdownretrycount20(jobnetId string, testcase *dao.TestCase, sshClient *ssh.Client) common.Testcase_status {

	/******************
	Pre-Operation State
	******************/
	configFilePath := "/etc/jobarranger/jobarg_agentd.conf"
	err := lib.UpdateDebugLevel(common.Client, configFilePath, 3)
	if err != nil {
		fmt.Println(testcase.Err_log("Error updating DebugLevel: %v", err))
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
		fmt.Println(testcase.Err_log("Error: Remove all file within agend directory"))
		return FAILED
	}

	dirPath = "/var/log/jobarranger/job"
	common.Client = lib.GetSSHClient(common.Login_info.Hostname, common.Login_info.Port, common.Login_info.Username, common.Login_info.Password)

	result = lib.RemoveAllFilesInDirectory(common.Client, dirPath)
	if result != nil {
		fmt.Println(testcase.Err_log("Error: Remove all file within server directory"))
		return FAILED
	}

	/**************
	Operation State
	***************/

	if err := lib.Jobarg_enable_jobnet("Icon_1", "821Loop"); err != nil {
		fmt.Println(testcase.Err_log("Failed to enable jobnet, Error: %s", err))
		return FAILED
	}

	//Run jobnet
	run_jobnet_id, error := lib.Jobarg_exec(jobnetId)
	if error != nil {
		fmt.Println(testcase.Err_log("Error: %s, std_out: %s", error.Error(), run_jobnet_id))
		return FAILED
	}
	fmt.Println(testcase.Info_log("Job Icon %s has been successfully run with registry number: %s", jobnetId, run_jobnet_id))

	//defaulQuery := `select * from ja_run_job_table jrjt where jrjt.job_id ='NORMAL-JOB' and jrjt.status =2 and jrjt.job_type = 4`
	defaulQuery := `SELECT job_id, status, job_type FROM ja_run_job_table WHERE job_id='NORMAL-JOB' AND status=2 AND job_type=4
`
	timeout := 2 * time.Minute
	interval := 1 * time.Second

	if common.Is_psql {
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
		fmt.Println(testcase.Err_log("Error:%s", err))
	} else {
		fmt.Println(testcase.Info_log("%s show in server log", pattern))
	}

	fmt.Println(testcase.Info_log("Wait for the log info from agent log ..."))
	sleepDuration := 2 * time.Minute
	time.Sleep(sleepDuration)

	pattern = "[WARN] In ja_job_exec_close() agent close failed. retry count :[20]"
	logFilePath2 := "/var/log/jobarranger/jobarg_agentd.log"

	timeout = 1 * time.Minute  // Timeout duration
	interval = 1 * time.Second // Polling interval

	_, err = lib.WaitForPatternInLogFile(common.Client, logFilePath2, pattern, timeout, interval)
	if err != nil {
		fmt.Println(testcase.Err_log("Error:%s", err))
	} else {
		fmt.Println(testcase.Info_log("%s show in agent log", pattern))
	}

	fmt.Println(testcase.Info_log("start the database service and check the no index file remain from server side ..."))

	if common.Is_psql {
		lib.StartDatabaseService(common.Client, "postgresql")
	} else {
		lib.StartDatabaseService(common.Client, "mysqld")
	}
	fmt.Println(testcase.Info_log("Checking index file remain in server side ..."))
	sleepDuration = 3 * time.Minute
	time.Sleep(sleepDuration)

	dirPath = "/var/log/jobarranger/job"
	// Call the function to check if any .job files exist in the specified directory.
	exists, err := lib.CheckRemoteIndexFileExists(common.Client, dirPath)
	if err != nil {
		fmt.Println("Error:", err)
		return FAILED
	}

	if exists {
		fmt.Println(testcase.Err_log("Index file remain"))
		return FAILED
	} else {
		fmt.Println(testcase.Info_log("Index file does not remain in %s ", dirPath))
	}

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
		fmt.Println(testcase.Err_log("Error : %s", err.Error()))
		return FAILED
	}

	// /***************
	// Expected Results
	// ****************/
	// Check if we have any results
	if len(results) > 0 {
		// Print only the first result
		result := results[0] // Get the first result
		fmt.Println(testcase.Info_log("The value of LOOPCNT: %d, and  STD_OUT: %d", result.AfterValueLoopcnt, result.AfterValueStdOut))
	} else {
		fmt.Println(testcase.Err_log("STD_OUT and LOOPCNT data not found."))
		return FAILED
	}

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
