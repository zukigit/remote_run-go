package tickets

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/zukigit/remote_run-go/src/common"
	"github.com/zukigit/remote_run-go/src/dao"
	"github.com/zukigit/remote_run-go/src/lib"
	"golang.org/x/crypto/ssh"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

type Ticket_840 struct {
	no          uint
	description string
	testcases   []dao.TestCase
}

func (t *Ticket_840) New_testcase(testcase_id uint, testcase_description string) *dao.TestCase {
	return dao.New_testcase(testcase_id, testcase_description)
}

func (t *Ticket_840) Get_no() uint {
	return t.no
}

func (t *Ticket_840) Get_dsctn() string {
	return t.description
}

func (t *Ticket_840) Add_testcase(tc dao.TestCase) {
	t.testcases = append(t.testcases, tc)
}

func (t *Ticket_840) Get_testcases() []dao.TestCase {
	return t.testcases
}

// Enter your ticket information here
func (t *Ticket_840) Set_values() {
	t.no = 840 // Enter your ticket id
	t.description = "(Fujitsu) Fwait icon can't execute well on SJIS environment"
}

// OS Types
const (
	Windows = "Windows"
	Linux   = "Linux"
)

// Encoding Types
const (
	SJIS = "SJIS"
	UTF8 = "UTF-8"
)

// struct for each environment
type Environment struct {
	Encoding string
	OSType   string
	Hostname string
}

const otherLinuxHost = "oss.linux1" // this should be hostname of other linux server

// Add your test case here
func (t *Ticket_840) Add_testcases() {

	// initialize the environments
	// ***** Both windows for SJIS or UTF-8 should be local windows, which is running the automation tool *****
	envrionments := []Environment{
		{Encoding: SJIS, OSType: Linux, Hostname: "oss.linux"},
		// {Encoding: SJIS, OSType: Windows, Hostname: "oss.windows"},
		{Encoding: UTF8, OSType: Linux, Hostname: "oss.linux"},
		{Encoding: UTF8, OSType: Windows, Hostname: "oss.windows"},
	}

	var testcaseCounter uint = 0
	// run testcases for each environment
	for _, env := range envrionments {
		// Waiting File Creation (the file will be created)
		testcaseCounter++
		testcaseA := t.New_testcase(testcaseCounter, fmt.Sprintf("Can be successfully run on %s (%s) [FWait - Waiting File Creation]", env.Encoding, env.OSType))

		tcFunc := func() common.Testcase_status {
			return fileWaitTestcase(testcaseA, env)
		}
		testcaseA.Set_function(tcFunc)
		t.Add_testcase(*testcaseA)

		// Check File Existence (the file exists)
		testcaseCounter++
		testcaseB := t.New_testcase(testcaseCounter, fmt.Sprintf("Can be successfully run on %s (%s) [FWait - Check File Existence: file exists]", env.Encoding, env.OSType))
		tcFunc = func() common.Testcase_status {
			return fileCheckTestcase(testcaseB, env, true)
		}
		testcaseB.Set_function(tcFunc)
		t.Add_testcase(*testcaseB)

		// Check File Existence (the file does not exist)
		testcaseCounter++
		testcaseC := t.New_testcase(testcaseCounter, fmt.Sprintf("Can be successfully run on %s (%s) [FWait - Check File Existence: file does not exist]", env.Encoding, env.OSType))
		tcFunc = func() common.Testcase_status {
			return fileCheckTestcase(testcaseC, env, false)
		}
		testcaseC.Set_function(tcFunc)
		t.Add_testcase(*testcaseC)

		// Normal Job Icon with file copy command
		testcaseCounter++
		testcaseD := t.New_testcase(testcaseCounter, fmt.Sprintf("Can be successfully run on %s (%s) [Job Icon - Copy Command with japanese name]", env.Encoding, env.OSType))
		tcFunc = func() common.Testcase_status {
			return normalIconTestcase(testcaseD, env)
		}
		testcaseD.Set_function(tcFunc)
		t.Add_testcase(*testcaseD)

		if env.OSType == Linux {
			// Reboot job
			testcaseCounter++
			testcaseE := t.New_testcase(testcaseCounter, fmt.Sprintf("Can be successfully run on %s (%s) [Reboot Icon]", env.Encoding, env.OSType))
			tcFunc := func() common.Testcase_status {
				return rebootIconTestcase(testcaseE, env)
			}
			testcaseE.Set_function(tcFunc)
			t.Add_testcase(*testcaseE)
		}

		// FCopy Icon
		// Own host
		testcaseCounter++
		testcaseF := t.New_testcase(testcaseCounter, fmt.Sprintf("Can be successfully run on %s (%s) [FCopy Icon - own host]", env.Encoding, env.OSType))
		tcFunc = func() common.Testcase_status {
			envs := map[string]string{
				"JA_SRC_HOST":  env.Hostname,
				"JA_DEST_HOST": env.Hostname,
			}
			return fcopyIconTestcase(testcaseF, env, envs)
		}
		testcaseF.Set_function(tcFunc)
		t.Add_testcase(*testcaseF)

		// Different host (linux)
		if env.OSType == Linux {
			testcaseCounter++
			testcaseG := t.New_testcase(testcaseCounter, fmt.Sprintf("Can be successfully run on %s (%s) [FCopy Icon - different host, same OS]", env.Encoding, env.OSType))
			tcFunc = func() common.Testcase_status {
				envs := map[string]string{
					"JA_SRC_HOST":  env.Hostname,
					"JA_DEST_HOST": otherLinuxHost,
				}
				return fcopyIconTestcase(testcaseG, env, envs)
			}
			testcaseG.Set_function(tcFunc)
			t.Add_testcase(*testcaseG)
		}

	}

}

// this function is for the testcase of rebooting the server
func rebootIconTestcase(testcase *dao.TestCase, env Environment) common.Testcase_status {
	fmt.Printf("====== %d - %s (%s) : %s ======\n", testcase.Get_id(), env.Encoding, env.OSType, testcase.Get_dsctn())

	// set AllowRoot to 1 for reboot
	err := lib.Ja_set_agent_config_linux("AllowRoot", "1")
	if err != nil {
		fmt.Println(testcase.Err_log("failed to set agent config: %s", err.Error()))
		return FAILED
	}

	// restart the jaz agent server
	err = lib.Restart_jaz_agent_linux()
	if err != nil {
		fmt.Println(testcase.Err_log("failed to restart jaz agent: %s", err.Error()))
		return FAILED
	}

	// convert the encoding of the Linux server
	srvEncoding := "ja_JP." + env.Encoding
	encodingCmd := fmt.Sprintf("echo LANG=%s > /etc/locale.conf", srvEncoding)
	_, err = lib.Ssh_exec(encodingCmd)
	if err != nil {
		fmt.Println(testcase.Err_log("failed to convert the encoding: %s", err.Error()))
		return common.ERROR
	}

	fmt.Println(testcase.Info_log("Encoding is set as: %s", srvEncoding))

	jobnetId := "TICKET840"
	jobnetName := "reboot_icon"

	// Enabling the jobnet
	err = lib.Jobarg_enable_jobnet(jobnetId, jobnetName)
	if err != nil {
		fmt.Println("err in enable jobnet", err.Error())
		return FAILED
	}

	envs := map[string]string{
		"JA_HOSTNAME": env.Hostname,
	}

	registryNumber, err := lib.Jobarg_exec_E(jobnetId, envs)

	if err != nil {
		fmt.Println(testcase.Err_log("Error running the jobnet: %s, std_out: %s", err.Error(), registryNumber))
		return FAILED
	}

	fmt.Println(testcase.Info_log("Jobnet [%s:%s] is run successfully with registry number: %s", jobnetId, jobnetName, registryNumber))

	// Check if the server is rebooted and up again
	fmt.Println(testcase.Info_log("Waiting for server to be down..."))
	downTimeout := 1 // minutes
	if waitForServerDown(common.Login_info.Hostname, strconv.Itoa(common.Login_info.Port), 3*time.Second, time.Duration(downTimeout)*time.Minute) {
		fmt.Println(testcase.Info_log("the server is down"))
	} else {
		fmt.Println(testcase.Err_log("the server is not down in the expected time of %d minutes", downTimeout))
		return FAILED
	}

	fmt.Println(testcase.Info_log("Waiting for server to be up again..."))
	upTimeout := 5 // minutes
	if waitForServerUp(common.Login_info.Hostname, strconv.Itoa(common.Login_info.Port), 3*time.Second, time.Duration(upTimeout)*time.Minute) {
		fmt.Println(testcase.Info_log("the server is now up and running"))
	} else {
		fmt.Println(testcase.Err_log("the server is not up in the expected time of %d", upTimeout))
		return FAILED
	}

	// Re-create the ssh client
	common.Client, err = ssh.Dial("tcp", fmt.Sprintf("%s:%d", common.Login_info.Hostname, common.Login_info.Port), &ssh.ClientConfig{
		User: common.Login_info.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(common.Login_info.Username),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	})

	if err != nil {
		fmt.Println(testcase.Err_log("failed to re-create the ssh client: %s", err.Error()))
		return FAILED
	}

	// Wait jobnet finishes and get jobnet run info.
	targetJobnetStatus := "END"
	targetJobStatus := "NORMAL"
	jobnetRunInfo, err := lib.Jobarg_get_jobnet_info(registryNumber, targetJobnetStatus, targetJobStatus)
	if err != nil {
		fmt.Println(testcase.Err_log("Error getting jobnet info: %s", err.Error()))
		return FAILED
	}

	// Success (obtain the expected status, message, or exit code)
	fmt.Println(testcase.Info_log("Jobnet_status: %s, Job_status: %s, std_out: %s, std_err: %s", jobnetRunInfo.Jobnet_status, jobnetRunInfo.Job_status, jobnetRunInfo.Std_out, jobnetRunInfo.Std_error))

	// Check the stderr
	stderr := strings.TrimSpace(jobnetRunInfo.Std_error)
	if stderr != "" {
		fmt.Println(testcase.Err_log("stderr should be empty, but it is not: '%s'", stderr))
		return FAILED
	}

	fmt.Println(testcase.Info_log("%s with registry number %s is completed with expected status and output.", jobnetId, registryNumber))

	return PASSED
}

// This function checks the testcase for waiting file creation with FWait Icon
func fileWaitTestcase(testcase *dao.TestCase, env Environment) common.Testcase_status {
	fmt.Printf("====== %d - %s (%s) : %s ======\n", testcase.Get_id(), env.Encoding, env.OSType, testcase.Get_dsctn())

	// File to be waited for creation or checked for existence (must be japanese name)
	fileName := "日本語ファイル.txt"
	var filePath string

	// Setup the file path according to the OS Type. Folder path: Linux -> /tmp, Windows -> C:\Users\{Username}\Documents
	if env.OSType == Linux {
		filePath = "/tmp/" + fileName

		// temporarily convert the encoding of the Linux server
		srvEncoding := "ja_JP." + env.Encoding
		encodingCmd := fmt.Sprintf("echo LANG=%s > /etc/locale.conf", srvEncoding)
		_, err := lib.Ssh_exec(encodingCmd)
		if err != nil {
			fmt.Println(testcase.Err_log("failed to convert the encoding: %s", err.Error()))
			return common.ERROR
		}

		fmt.Println(testcase.Info_log("Encoding is set as: %s", srvEncoding))

	} else if env.OSType == Windows {
		homeDir, err := os.UserHomeDir()

		if err != nil {
			fmt.Println(testcase.Err_log("failed to get home directory: %s", err.Error()))
			return common.ERROR
		}

		filePath = filepath.Join(homeDir, "Documents", fileName)
	} else {
		fmt.Println(testcase.Err_log("unsupported os type: %s", env.OSType))
		return common.ERROR
	}

	fmt.Println(testcase.Info_log("File to be waited: %s", filePath))

	jobnetId := "TICKET840"
	jobnetName := "fwait_waiting_file_creation"

	// Enabling the jobnet
	err := lib.Jobarg_enable_jobnet(jobnetId, jobnetName)
	if err != nil {
		fmt.Println("err in enable jobnet", err.Error())
		return FAILED
	}

	envs := map[string]string{
		"JA_HOSTNAME": env.Hostname,
		"FILE_PATH":   filePath,
	}
	registryNumber, err := lib.Jobarg_exec_E(jobnetId, envs)

	if err != nil {
		fmt.Println(testcase.Err_log("Error running the jobnet: %s, std_out: %s", err.Error(), registryNumber))
		return FAILED
	}

	fmt.Println(testcase.Info_log("Jobnet [%s:%s] is run successfully with registry number: %s", jobnetId, jobnetName, registryNumber))

	// Create the file if not exists
	if err = createFile(filePath, env.OSType); err != nil {
		fmt.Println(testcase.Err_log("failed to create the file: %s", err.Error()))
		return FAILED
	}

	fmt.Println(testcase.Info_log("File is created: %s", filePath))

	// Wait jobnet finishes and get jobnet run info.
	targetJobnetStatus := "END"
	targetJobStatus := "NORMAL"
	jobnetRunInfo, err := lib.Jobarg_get_jobnet_info(registryNumber, targetJobnetStatus, targetJobStatus)
	if err != nil {
		fmt.Println(testcase.Err_log("Error getting jobnet info: %s", err.Error()))
		return FAILED
	}

	// Success (obtain the expected status, message, or exit code)
	fmt.Println(testcase.Info_log("Jobnet_status: %s, Job_status: %s, std_out: %s, std_err: %s", jobnetRunInfo.Jobnet_status, jobnetRunInfo.Job_status, jobnetRunInfo.Std_out, jobnetRunInfo.Std_error))

	// Check the stderr
	stdout := strings.TrimSpace(jobnetRunInfo.Std_out)
	stderr := strings.TrimSpace(jobnetRunInfo.Std_error)
	if stderr != "" {
		fmt.Println(testcase.Err_log("stderr should be empty, but it is not: '%s'", stderr))
		return FAILED
	}

	// Check the stdout
	expectedStdOut := fmt.Sprintf("the file '%s' exists", filePath)
	if stdout != expectedStdOut {
		fmt.Println(testcase.Err_log("stdout is not as expected. Expected: %s, Actual: %s", expectedStdOut, stdout))
		return FAILED
	}

	fmt.Println(testcase.Info_log("%s with registry number %s is completed with expected status and output.", jobnetId, registryNumber))

	return PASSED
}

// This function checks the testcase for file existence with FWait Icon
// It can check for both conditions: file exisis, and file does not exist
func fileCheckTestcase(testcase *dao.TestCase, env Environment, shouldFileExist bool) common.Testcase_status {
	fmt.Printf("====== %d - %s (%s) : %s ======\n", testcase.Get_id(), env.Encoding, env.OSType, testcase.Get_dsctn())

	// File to be waited for creation or checked for existence (must be japanese name)
	fileName := "日本語ファイル.txt"
	var filePath string

	// Setup the file path according to the OS Type. Folder path: Linux -> /tmp, Windows -> C:\Users\{Username}\Documents
	if env.OSType == Linux {
		filePath = "/tmp/" + fileName

		// temporarily convert the encoding of the Linux server
		srvEncoding := "ja_JP." + env.Encoding
		encodingCmd := fmt.Sprintf("echo LANG=%s > /etc/locale.conf", srvEncoding)
		_, err := lib.Ssh_exec(encodingCmd)
		if err != nil {
			fmt.Println(testcase.Err_log("failed to convert the encoding: %s", err.Error()))
			return common.ERROR
		}

		fmt.Println(testcase.Info_log("Encoding is set as: %s", srvEncoding))

	} else if env.OSType == Windows {
		homeDir, err := os.UserHomeDir()

		if err != nil {
			fmt.Println(testcase.Err_log("failed to get home directory: %s", err.Error()))
			return common.ERROR
		}

		filePath = filepath.Join(homeDir, "Documents", fileName)
	} else {
		fmt.Println(testcase.Err_log("unsupported os type: %s", env.OSType))
		return common.ERROR
	}

	fmt.Println(testcase.Info_log("File to be checked: %s", filePath))

	// create or delete based on testcase
	if shouldFileExist {
		// Create the file
		if err := createFile(filePath, env.OSType); err != nil {
			fmt.Println(testcase.Err_log("failed to create the file: %s", err.Error()))
			return FAILED
		}

		fmt.Println(testcase.Info_log("File is created: %s", filePath))
	} else {
		// Delete the file if exists
		if err := deleteFileIfExists(filePath, env.OSType); err != nil {
			fmt.Println(testcase.Err_log("failed to delete the file: %s", err.Error()))
			return FAILED
		}

		fmt.Println(testcase.Info_log("File is deleted: %s", filePath))
	}

	jobnetId := "TICKET840"
	jobnetName := "fwait_check_file_existence"

	// Enabling the jobnet
	err := lib.Jobarg_enable_jobnet(jobnetId, jobnetName)
	if err != nil {
		fmt.Println("err in enable jobnet", err.Error())
		return FAILED
	}

	envs := map[string]string{
		"JA_HOSTNAME": env.Hostname,
		"FILE_PATH":   filePath,
	}
	registryNumber, err := lib.Jobarg_exec_E(jobnetId, envs)

	if err != nil {
		fmt.Println(testcase.Err_log("Error running the jobnet: %s, std_out: %s", err.Error(), registryNumber))
		return FAILED
	}

	fmt.Println(testcase.Info_log("Jobnet [%s:%s] is run successfully with registry number: %s", jobnetId, jobnetName, registryNumber))

	// Wait jobnet finishes and get jobnet run info.
	targetJobnetStatus := "END"
	targetJobStatus := "NORMAL"
	jobnetRunInfo, err := lib.Jobarg_get_jobnet_info(registryNumber, targetJobnetStatus, targetJobStatus)
	if err != nil {
		fmt.Println(testcase.Err_log("Error getting jobnet info: %s", err.Error()))
		return FAILED
	}

	// Success (obtain the expected status, message, or exit code)
	fmt.Println(testcase.Info_log("Jobnet_status: %s, Job_status: %s, std_out: %s, std_err: %s", jobnetRunInfo.Jobnet_status, jobnetRunInfo.Job_status, jobnetRunInfo.Std_out, jobnetRunInfo.Std_error))

	// Check the stderr
	stdout := strings.TrimSpace(jobnetRunInfo.Std_out)
	stderr := strings.TrimSpace(jobnetRunInfo.Std_error)
	if stderr != "" {
		fmt.Println(testcase.Err_log("stderr should be empty, but it is not: '%s'", stderr))
		return FAILED
	}

	// prepare expected stdout
	var expectedStdOutPattern string
	if shouldFileExist {
		expectedStdOutPattern = fmt.Sprintf(`(?i)the (?i)file '?%s'? exists`, regexp.QuoteMeta(filePath))
	} else {
		expectedStdOutPattern = fmt.Sprintf(`(?i)the (?i)file '?%s'? does not exist`, regexp.QuoteMeta(filePath))
	}

	re := regexp.MustCompile(expectedStdOutPattern)

	// Check the stdout
	if !re.MatchString(stdout) {
		fmt.Println(testcase.Err_log("stdout is not as expected. Expected: %s, Actual: %s", expectedStdOutPattern, stdout))
		return FAILED
	}

	fmt.Println(testcase.Info_log("%s with registry number %s is completed with expected status and output.", jobnetId, registryNumber))

	return PASSED
}

func normalIconTestcase(testcase *dao.TestCase, env Environment) common.Testcase_status {
	fmt.Printf("====== %d - %s (%s) : %s ======\n", testcase.Get_id(), env.Encoding, env.OSType, testcase.Get_dsctn())

	// File to be waited for creation or checked for existence (must be japanese name)
	fileName := "日本語ファイル.txt"
	var filePath string
	var homeDir string

	// Setup the file path according to the OS Type. Folder path: Linux -> /tmp, Windows -> C:\Users\{Username}\Documents
	if env.OSType == Linux {
		filePath = "/tmp/" + fileName

		// temporarily convert the encoding of the Linux server
		srvEncoding := "ja_JP." + env.Encoding
		encodingCmd := fmt.Sprintf("echo LANG=%s > /etc/locale.conf", srvEncoding)
		_, err := lib.Ssh_exec(encodingCmd)
		if err != nil {
			fmt.Println(testcase.Err_log("failed to convert the encoding: %s", err.Error()))
			return common.ERROR
		}

		fmt.Println(testcase.Info_log("Encoding is set as: %s", srvEncoding))

	} else if env.OSType == Windows {
		var err error
		homeDir, err = os.UserHomeDir()

		if err != nil {
			fmt.Println(testcase.Err_log("failed to get home directory: %s", err.Error()))
			return common.ERROR
		}

		filePath = filepath.Join(homeDir, "Documents", fileName)
	} else {
		fmt.Println(testcase.Err_log("unsupported os type: %s", env.OSType))
		return common.ERROR
	}

	fmt.Println(testcase.Info_log("File to be copied: %s", filePath))

	// Create the file
	if err := createFile(filePath, env.OSType); err != nil {
		fmt.Println(testcase.Err_log("failed to create the file: %s", err.Error()))
		return FAILED
	}

	fmt.Println(testcase.Info_log("File is created: %s", filePath))

	jobnetId := "TICKET840"
	var jobnetName string

	if env.OSType == Windows {
		jobnetName = "job_icon_windows"
	} else if env.OSType == Linux {
		jobnetName = "job_icon_linux"
	} else {
		fmt.Println(testcase.Err_log("unsupported os type: %s", env.OSType))
		return FAILED
	}

	// Enabling the jobnet
	err := lib.Jobarg_enable_jobnet(jobnetId, jobnetName)
	if err != nil {
		fmt.Println("err in enable jobnet", err.Error())
		return FAILED
	}

	envs := map[string]string{
		"JA_HOSTNAME": env.Hostname,
	}

	if env.OSType == Windows {
		envs["SOURCE_DIR"] = filepath.Join(homeDir, "Documents")
		envs["DEST_DIR"] = homeDir
	}

	registryNumber, err := lib.Jobarg_exec_E(jobnetId, envs)

	if err != nil {
		fmt.Println(testcase.Err_log("Error running the jobnet: %s, std_out: %s", err.Error(), registryNumber))
		return FAILED
	}

	fmt.Println(testcase.Info_log("Jobnet [%s:%s] is run successfully with registry number: %s", jobnetId, jobnetName, registryNumber))

	// Wait jobnet finishes and get jobnet run info.
	targetJobnetStatus := "END"
	targetJobStatus := "NORMAL"
	jobnetRunInfo, err := lib.Jobarg_get_jobnet_info(registryNumber, targetJobnetStatus, targetJobStatus)
	if err != nil {
		fmt.Println(testcase.Err_log("Error getting jobnet info: %s", err.Error()))
		return FAILED
	}

	// Success (obtain the expected status, message, or exit code)
	fmt.Println(testcase.Info_log("Jobnet_status: %s, Job_status: %s, std_out: %s, std_err: %s", jobnetRunInfo.Jobnet_status, jobnetRunInfo.Job_status, jobnetRunInfo.Std_out, jobnetRunInfo.Std_error))

	// Check the stderr
	stdout := strings.TrimSpace(jobnetRunInfo.Std_out)
	stderr := strings.TrimSpace(jobnetRunInfo.Std_error)
	if stderr != "" {
		fmt.Println(testcase.Err_log("stderr should be empty, but it is not: '%s'", stderr))
		return FAILED
	}

	// prepare expected stdout
	var expectedStdOutPattern string
	if env.OSType == Windows {
		expectedStdOutPattern = regexp.QuoteMeta("1 file(s) copied.")

		re := regexp.MustCompile(expectedStdOutPattern)

		// Check the stdout
		if !re.MatchString(stdout) {
			fmt.Println(testcase.Err_log("stdout is not as expected. Expected: %s, Actual: %s", expectedStdOutPattern, stdout))
			return FAILED
		}
	}

	fmt.Println(testcase.Info_log("%s with registry number %s is completed with expected status and output.", jobnetId, registryNumber))

	return PASSED
}

func fcopyIconTestcase(testcase *dao.TestCase, env Environment, envs map[string]string) common.Testcase_status {
	fmt.Printf("====== %d - %s (%s) : %s ======\n", testcase.Get_id(), env.Encoding, env.OSType, testcase.Get_dsctn())

	// File to be waited for creation or checked for existence (must be japanese name)
	fileName := "日本語ファイル.txt"
	var filePath string
	var homeDir string

	// Setup the file path according to the OS Type. Folder path: Linux -> /tmp, Windows -> C:\Users\{Username}\Documents
	if env.OSType == Linux {
		filePath = "/tmp/" + fileName

		// temporarily convert the encoding of the Linux server
		srvEncoding := "ja_JP." + env.Encoding
		encodingCmd := fmt.Sprintf("echo LANG=%s > /etc/locale.conf", srvEncoding)
		_, err := lib.Ssh_exec(encodingCmd)
		if err != nil {
			fmt.Println(testcase.Err_log("failed to convert the encoding: %s", err.Error()))
			return common.ERROR
		}

		fmt.Println(testcase.Info_log("Encoding is set as: %s", srvEncoding))

	} else if env.OSType == Windows {
		var err error
		homeDir, err = os.UserHomeDir()

		if err != nil {
			fmt.Println(testcase.Err_log("failed to get home directory: %s", err.Error()))
			return common.ERROR
		}

		filePath = filepath.Join(homeDir, "Documents", fileName)
	} else {
		fmt.Println(testcase.Err_log("unsupported os type: %s", env.OSType))
		return common.ERROR
	}

	fmt.Println(testcase.Info_log("File to be copied: %s", filePath))

	// Create the file
	if err := createFile(filePath, env.OSType); err != nil {
		fmt.Println(testcase.Err_log("failed to create the file: %s", err.Error()))
		return FAILED
	}

	fmt.Println(testcase.Info_log("File is created: %s", filePath))

	jobnetId := "TICKET840"
	jobnetName := "fcopy_icon"

	// Set up job icon value
	fromDir := map[bool]string{
		true:  "/tmp",
		false: filepath.Join(homeDir, "Documents"),
	}[strings.Contains(envs["JA_SRC_HOST"], "linux")]

	toDir := map[bool]string{
		true:  "/root",
		false: homeDir,
	}[strings.Contains(envs["JA_DEST_HOST"], "linux")]

	query := `UPDATE ja_icon_fcopy_table
SET 
    from_directory = $1,
    from_file_name = $2,
    to_directory = $3
WHERE 
    jobnet_id = $4;`

	if _, err := lib.ExecuteQuery(lib.DBQuery(query), fromDir, fileName, toDir, jobnetId); err != nil {
		fmt.Println(testcase.Err_log("failed to update values in fcopy icon: %s", err.Error()))
		return FAILED
	}

	// Enabling the jobnet
	err := lib.Jobarg_enable_jobnet(jobnetId, jobnetName)
	if err != nil {
		fmt.Println(testcase.Err_log("err in enable jobnet: %s", err.Error()))
		return FAILED
	}

	registryNumber, err := lib.Jobarg_exec_E(jobnetId, envs)

	if err != nil {
		fmt.Println(testcase.Err_log("Error running the jobnet: %s, std_out: %s", err.Error(), registryNumber))
		return FAILED
	}

	fmt.Println(testcase.Info_log("Jobnet [%s:%s] is run successfully with registry number: %s", jobnetId, jobnetName, registryNumber))

	// Wait jobnet finishes and get jobnet run info.
	targetJobnetStatus := "END"
	targetJobStatus := "NORMAL"
	jobnetRunInfo, err := lib.Jobarg_get_jobnet_info(registryNumber, targetJobnetStatus, targetJobStatus)
	if err != nil {
		fmt.Println(testcase.Err_log("Error getting jobnet info: %s", err.Error()))
		return FAILED
	}

	// Success (obtain the expected status, message, or exit code)
	fmt.Println(testcase.Info_log("Jobnet_status: %s, Job_status: %s, std_out: %s, std_err: %s", jobnetRunInfo.Jobnet_status, jobnetRunInfo.Job_status, jobnetRunInfo.Std_out, jobnetRunInfo.Std_error))

	// Check the stderr
	stderr := strings.TrimSpace(jobnetRunInfo.Std_error)
	if stderr != "" {
		fmt.Println(testcase.Err_log("stderr should be empty, but it is not: '%s'", stderr))
		return FAILED
	}

	fmt.Println(testcase.Info_log("%s with registry number %s is completed with expected status and output.", jobnetId, registryNumber))

	return PASSED
}

// Function to create a file if it doesn't exist. It works for both windows and linux
func createFile(filePath, osType string) error {

	if osType == Windows {
		file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			return fmt.Errorf("could not create file: %w", err)
		}
		defer file.Close() // Ensure the file is closed after creation

		// Encode content to Shift-JIS
		writer := transform.NewWriter(file, japanese.ShiftJIS.NewEncoder())
		_, err = writer.Write([]byte("これは日本語のファイルです。"))

		if err != nil {
			return fmt.Errorf("failed to write japanese content: %w", err)
		}

		fmt.Println("File created:", filePath)
		return nil
	} else if osType == Linux {
		// to be implemented.
		if _, err := lib.Ssh_exec("echo 'これは日本語のファイルです。' > " + filePath); err != nil {
			return fmt.Errorf("could not create file: %w", err)
		}
		return nil
	} else {
		return fmt.Errorf("unsupported os type: %s", osType)
	}

}

// This function deletes the file if exists. It works for both windows and linux
func deleteFileIfExists(filePath, osType string) error {
	if osType == Windows {
		// Check if the file already exists
		_, err := os.Stat(filePath)
		if err != nil {
			if os.IsNotExist(err) {
				// file does not exist, no need to delete it
				return nil
			}

			// some other error while checking the file
			return fmt.Errorf("could not check if file exists: %w", err)
		}

		err = os.Remove(filePath)
		if err != nil {
			return fmt.Errorf("could not delete file: %w", err)
		}

		fmt.Println("File is deleted:", filePath)
		return nil

	} else if osType == Linux {
		// to be implemented.
		if _, err := lib.Ssh_exec(fmt.Sprintf("rm -f '%s'", filePath)); err != nil {
			return fmt.Errorf("could not create file: %w", err)
		}
		return nil
	} else {
		return fmt.Errorf("unsupported os type: %s", osType)
	}
}

// waitForServerDown periodically checks if the server's SSH port is accessible
func waitForServerDown(server, port string, interval, timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)

	for {
		if time.Now().After(deadline) {
			return false // Timed out
		}

		conn, err := net.DialTimeout("tcp", net.JoinHostPort(server, port), interval)
		if err != nil {
			return true // Server is down
		}

		conn.Close()

		// Server is not down yet, wait and retry
		time.Sleep(interval)
	}
}

// waitForServerUp periodically checks if the server's SSH port is accessible
func waitForServerUp(server, port string, interval, timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)

	for {
		if time.Now().After(deadline) {
			return false // Timed out
		}

		conn, err := net.DialTimeout("tcp", net.JoinHostPort(server, port), interval)
		if err == nil {
			conn.Close()
			return true // Server is up
		}

		// Server is not up yet, wait and retry
		time.Sleep(interval)
	}
}
