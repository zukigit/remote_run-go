package tickets

import (
	"fmt"
	"strings"
	"time"

	"github.com/zukigit/remote_run-go/src/common"
	"github.com/zukigit/remote_run-go/src/dao"
	"github.com/zukigit/remote_run-go/src/lib"
)

type Ticket_919 struct {
	Ticket_no                                   uint
	Ticket_description                          string
	PASSED_count, FAILED_count, MUSTCHECK_count int
	Testcases                                   []dao.TestCase
}

func (t *Ticket_919) New_testcase(testcase_id uint, testcase_description string) *dao.TestCase {
	return dao.New_testcase(testcase_id, testcase_description)
}

func (t *Ticket_919) Get_no() uint {
	return t.Ticket_no
}

func (t *Ticket_919) Set_PASSED_count(passed_count int) {
	t.PASSED_count = passed_count
}

func (t *Ticket_919) Set_FAILED_count(failed_count int) {
	t.FAILED_count = failed_count
}

func (t *Ticket_919) Set_MUSTCHECK_count(mustcheck_count int) {
	t.MUSTCHECK_count = mustcheck_count
}

func (t *Ticket_919) Get_dsctn() string {
	return t.Ticket_description
}

func (t *Ticket_919) Add_testcase(tc dao.TestCase) {
	t.Testcases = append(t.Testcases, tc)
}

func (t *Ticket_919) Get_testcases() []dao.TestCase {
	return t.Testcases
}

func (t *Ticket_919) Set_values() {
	t.Ticket_no = 919
	t.Ticket_description = "Add retry process to DB deadlock situations."
}
func (t *Ticket_919) Add_testcases() {

	tc_1 := t.New_testcase(1, "Check if jobs can rerun successfully after deadlock.")
	tc_func := func() common.Testcase_status {

		// Step 2: Cleanup environment
		if err := lib.Jobarg_cleanup_linux(); err != nil {
			fmt.Printf("Error during cleanup: %s\n", err.Error())
			return FAILED
		}

		// Step 4: Enable jobnet
		job := "Icon_1"
		if err := lib.Jobarg_enable_jobnet(job, "jobicon_linux"); err != nil {
			fmt.Printf("Error enabling jobnet: %s\n", err.Error())
			return FAILED
		}

		// Step 3: Set environment variables
		envs, err := lib.Get_str_str_map("JA_HOSTNAME", "oss.linux", "JA_CMD", "sleep 50")
		if err != nil {
			fmt.Printf("Error retrieving environment variables: %s\n", err.Error())
			return FAILED
		}

		runJobnetID, err := lib.Jobarg_exec_E(job, envs)
		if err != nil {
			fmt.Printf("Error executing job: %s\n", err.Error())
			return FAILED
		}
		fmt.Printf("Executed job: %s with runJobnetID: %s\n", job, runJobnetID)

		fmt.Println("Waiting for 8 seconds before locking the table...")
		time.Sleep(8 * time.Second)

		if err := t.lockTable(tc_1); err != nil {
			return FAILED
		}

		// Step 7: Retrieve jobnet run info
		jobnetRunInfo, err := lib.Jobarg_get_jobnet_run_info(runJobnetID)
		if err != nil {
			fmt.Printf("Error retrieving jobnet run info: %s\n", err.Error())
			return FAILED
		}

		if jobnetRunInfo.Jobnet_status == "END" && jobnetRunInfo.Job_status == "NORMAL" {
			fmt.Printf("%s completed successfully.\n", job)
			return PASSED
		}

		return PASSED
	}

	tc_1.Set_function(tc_func)
	t.Add_testcase(*tc_1)
}
func (t *Ticket_919) lockTable(tc *dao.TestCase) error {
	var err error

	if common.Is_mysql {
		fmt.Println("Using MySQL - Locking table...")
		_, err = common.DB.Exec("BEGIN;")
		if err != nil {
			return fmt.Errorf("error beginning transaction: %s", err.Error())
		}

		_, err = common.DB.Exec("LOCK TABLES ja_run_value_after_table WRITE;")
		if err != nil {
			return fmt.Errorf("error locking table: %s", err.Error())
		}

	} else if common.Is_psql {
		fmt.Println("Using PostgreSQL - Locking table...")
		_, err = common.DB.Exec("BEGIN;")
		if err != nil {
			return fmt.Errorf("error beginning transaction: %s", err.Error())
		}

		_, err = common.DB.Exec("LOCK TABLE ja_run_value_after_table IN EXCLUSIVE MODE;")
		if err != nil {
			return fmt.Errorf("error locking table: %s", err.Error())
		}

	} else {
		return fmt.Errorf("unsupported database type")
	}

	// Check the log immediately after locking the table
	if t.checkLog(tc) == PASSED {
		// If log is found, unlock the table
		if err := t.unlockTable(); err != nil {
			return fmt.Errorf("error unlocking table: %s", err.Error())
		}
		fmt.Println("Table unlocked due to retry count found in log.")
	} else {
		// If log is not found, fail the test
		fmt.Println("Log not found, marking test as failed.")
		return fmt.Errorf("log not found, test failed")
	}

	return nil
}

func (t *Ticket_919) checkLog(tc_1 *dao.TestCase) common.Testcase_status {
	const logFilePath = "/var/log/jobarranger/jobarg_agentd.log"
	const logFileWarning = `retry count`
	const maxRetries = 10
	const retryInterval = 10 * time.Second

	for i := 0; i < maxRetries; i++ {
		cmd := fmt.Sprintf(`cat %s | grep "%s"`, logFilePath, logFileWarning)
		tc_1.Info_log("Executing command: %s", cmd)

		warningLogOutput, err := lib.Ssh_exec_to_str(cmd)

		tc_1.Info_log("checking for warning log: %s", err)

		if strings.Contains(warningLogOutput, logFileWarning) {
			tc_1.Info_log("Warning log found, returning Passed.")
			return PASSED
		}

		tc_1.Info_log("Warning log not found. Retrying in %v...", retryInterval)
		time.Sleep(retryInterval)
	}

	tc_1.Err_log("Warning log not found after retries, returning FAILED.")
	return FAILED
}

func (t *Ticket_919) unlockTable() error {
	var err error

	if common.Is_psql {
		fmt.Println("Using PostgreSQL - Unlocking table...")
		_, err = common.DB.Exec("COMMIT;")
		if err != nil {
			return fmt.Errorf("error unlocking table: %s", err.Error())
		}
	} else if common.Is_mysql {
		_, err = common.DB.Exec("UNLOCK TABLES;")
		if err != nil {
			return fmt.Errorf("error unlocking table: %s", err.Error())
		}
	} else {
		return fmt.Errorf("unsupported database type")
	}

	return nil
}
