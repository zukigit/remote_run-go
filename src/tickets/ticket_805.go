package tickets

import (
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/zukigit/remote_run-go/src/common"
	"github.com/zukigit/remote_run-go/src/dao"
	"github.com/zukigit/remote_run-go/src/lib"
)

type Ticket_805 struct {
	no          uint
	description string
	testcases   []dao.TestCase
}

func (t *Ticket_805) New_testcase(testcase_id uint, testcase_description string) *dao.TestCase {
	return dao.New_testcase(testcase_id, testcase_description)
}

func (t *Ticket_805) Get_no() uint {
	return t.no
}

func (t *Ticket_805) Get_dsctn() string {
	return t.description
}

func (t *Ticket_805) Add_testcase(tc dao.TestCase) {
	t.testcases = append(t.testcases, tc)
}

func (t *Ticket_805) Get_testcases() []dao.TestCase {
	return t.testcases
}

func (t *Ticket_805) Set_values() {
	t.no = 805 // Enter your ticket id
	t.description = "Run Jobnet"
}

func (t *Ticket_805) Add_testcases() {
	tc_1 := t.New_testcase(1, "Run Jobnet with External Command")
	tc_func := func() common.Testcase_status {
		// Enable the jobnet first
		if err := lib.Jobarg_enable_jobnet("Icon_1", "hold_job_icon"); err != nil {
			tc_1.Err_log("Failed to enable jobnet, Error: %s", err)
			return FAILED
		}

		// Now run the jobnet with automatic handling of held jobs
		if err := RunJobnetWithHeldJob("Icon_1", "JOB-1"); err != nil {
			tc_1.Err_log("Failed to run jobnet, Error: %s", err)
			return FAILED
		}

		// Continue with the next steps if needed
		// (e.g., running additional commands, assertions, etc.)
		return PASSED
	}
	tc_1.Set_function(tc_func)
	t.Add_testcase(*tc_1)
}

func RunJobnetWithHeldJob(jobnetId, jobId string) error {
	// Step 1: Start the jobnet execution
	startTime := time.Now().Format("20060102150405")

	// Set up environment variables as a map
	envs, err := lib.Get_str_str_map("JA_HOSTNAME", "oss.linux", "JA_CMD", "hostname")
	if err != nil {
		return fmt.Errorf("failed to get environment map: %w", err)
	}

	// Execute the jobnet with the specified environment variables
	output, err := lib.Jobarg_exec_E(jobnetId, envs)
	if err != nil {
		return fmt.Errorf("failed to execute jobnet: %w", err)
	}

	// Log the output
	logMessage := fmt.Sprintf("%s [INFO] Jobarg_exec output: %s", startTime, output)
	formattedLog := lib.Formatted_log(common.INFO, logMessage)
	lib.Logi(formattedLog)    // Write log to file
	fmt.Println(formattedLog) // Show in terminal as well

	// Record the job log date
	joblogDate := time.Now().Format("20060102")

	// Step 2: Log the job execution output using Jobarg_joblogput
	logputTime := time.Now().Format("20060102150405")
	joblogputOutput, err := Jobarg_joblogput(jobnetId, jobId, joblogDate, output)
	if err != nil {
		return fmt.Errorf("failed to log job execution output: %w", err)
	}

	// Format and log Jobarg_joblogput output
	logMessage = fmt.Sprintf("%s [INFO] Jobarg_joblogput output: \"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%d\"",
		logputTime,    // Log timestamp
		joblogDate,    // Log date
		output,        // Inner Jobnet Main ID (example placeholder)
		output,        // Inner Jobnet ID (example placeholder)
		"run_type",    // Run type (example placeholder)
		"public_flag", // Public flag (example placeholder)
		jobnetId,      // Jobnet ID
		jobId,         // Job ID
		"message_id",  // Message ID (example placeholder)
		"message",     // Message (example placeholder)
		"jobnet_name", // Jobnet name (example placeholder)
		"job_name",    // Job name (example placeholder)
		"user_name",   // User name (example placeholder)
		joblogDate,    // Update date
		0)             // Return code (example placeholder)

	formattedLog = lib.Formatted_log(common.INFO, logMessage, joblogputOutput)
	lib.Logi(formattedLog)    // Write log to file
	fmt.Println(formattedLog) // Show in terminal as well

	// Step 3: Wait to simulate an interruption
	time.Sleep(10 * time.Second)

	// Step 4: Release the held job
	releaseTime := time.Now().Format("20060102150405")
	releaseOutput, err := Jobarg_release(jobnetId, jobId) // Capture the output of the release command
	if err != nil {
		return fmt.Errorf("failed to release held job: %w", err)
	}

	// Log the result of Jobarg_release
	logMessage = fmt.Sprintf("%s [INFO] Jobarg_release output: %s", releaseTime, releaseOutput)
	formattedLog = lib.Formatted_log(common.INFO, logMessage)
	lib.Logi(formattedLog)    // Write log to file
	fmt.Println(formattedLog) // Show in terminal as well

	// Step 5: Get the jobnet status after release
	registryNumber := jobnetId // Adjust if registry number differs from jobnetId
	statusOutput, err := lib.Jobarg_get_JA_JOBNETSTATUS(registryNumber)
	if err != nil {
		return fmt.Errorf("failed to get JA_JOBNETSTATUS: %w", err)
	}

	// Log the result of Jobarg_get_JA_JOBNETSTATUS
	statusTime := time.Now().Format("20060102150405")
	statusLogMessage := fmt.Sprintf("%s [INFO] JA_JOBNETSTATUS output: %s", statusTime, statusOutput)
	formattedStatusLog := lib.Formatted_log(common.INFO, statusLogMessage)
	lib.Logi(formattedStatusLog)    // Write log to file
	fmt.Println(formattedStatusLog) // Show in terminal as well

	// Step 6: Return nil if all steps are successful
	return nil
}

// Jobarg_joblogput Cmd
func Jobarg_joblogput(jobnet_id, job_id, joblog_date, joblog_id string) (string, error) {
	// Construct the SSH command for jobarg_joblogput with the necessary parameters

	// Read the output from the temporary file created on the remote server
	cmd := fmt.Sprintf("jobarg_joblogput -z %s -U Admin -P zabbix -s %s -r %s", common.Login_info.Hostname, joblog_date, joblog_id)
	result, err := lib.Ssh_exec_to_str(cmd)
	if err != nil {
		return "", fmt.Errorf("failed to read jobarg_joblogput output: %w", err)
	}

	return result, nil
}

// Job Release Cmd
func Jobarg_release(jobnet_id, job_id string) (string, error) {
	// Construct the SSH command for jobarg_release with redirection to a temporary file
	cmd := fmt.Sprintf("jobarg_release -z %s -U Admin -P zabbix -j %s/%s &> /tmp/moon_jobarg_release_result",
		common.Login_info.Hostname, jobnet_id, job_id)

	// Execute the command over SSH and check for any initial errors
	_, err := lib.Ssh_exec_to_str(cmd)
	if err != nil {
		return "", fmt.Errorf("failed to execute jobarg_release command: %w", err)
	}

	// Read the output from the temporary file created on the remote server
	cmd = "cat /tmp/moon_jobarg_release_result"
	result, err := lib.Ssh_exec_to_str(cmd)
	if err != nil {
		return "", fmt.Errorf("failed to read jobarg_release output: %w", err)
	}

	return result, nil
}
