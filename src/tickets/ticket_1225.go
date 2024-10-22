package tickets

import (
	"database/sql"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/zukigit/remote_run-go/src/common"
	"github.com/zukigit/remote_run-go/src/dao"
	"github.com/zukigit/remote_run-go/src/lib"
	"golang.org/x/crypto/ssh"
)

type Ticket_1225 struct {
	no          uint
	description string
	testcases   []dao.TestCase
}

func (t *Ticket_1225) New_testcase(testcase_id uint, testcase_description string) *dao.TestCase {
	return dao.New_testcase(testcase_id, testcase_description)
}

func (t *Ticket_1225) Get_no() uint {
	return t.no
}

func (t *Ticket_1225) Get_dsctn() string {
	return t.description
}

func (t *Ticket_1225) Add_testcase(tc dao.TestCase) {
	t.testcases = append(t.testcases, tc)
}

func (t *Ticket_1225) Get_testcases() []dao.TestCase {
	return t.testcases
}

// Enter your ticket information here
func (t *Ticket_1225) Set_values() {
	t.no = 1225 // Enter your ticket id
	t.description = "Japurge delete records with limit 500."
}

// Add your test case here
func (t *Ticket_1225) Add_testcases() {
	// TESTCASE 001
	tc_79 := t.New_testcase(79, "Test with 1 million records.")
	tc_func := func() common.Testcase_status {
		// Enter your test case logic here
		if err := lib.Stop_jaz_server(); err != nil {
			tc_79.Err_log("Failed to stop jobarg-server, Error: %s", err.Error())
			return FAILED
		}
		fmt.Println(tc_79.Info_log("JAZ Server has been stopped."))

		// Set log level
		log_level := "4"
		if err := lib.Ja_set_server_config_linux("DebugLevel", log_level); err != nil {
			tc_79.Err_log("Failed to update server config for log level, Error: %s", err.Error())
			return FAILED
		}
		fmt.Println(tc_79.Info_log("Log level has been set to %s.", log_level))

		// Set JaPurgeLimit
		ja_purge_limit := 500
		if err := lib.Ja_set_server_config_linux("JaPurgeLimit", strconv.Itoa(ja_purge_limit)); err != nil {
			tc_79.Err_log("Failed to update server config for JaPurgeLimit, Error: %s", err.Error())
			return FAILED
		}
		fmt.Println(tc_79.Info_log("JaPurgeLimit has been set to %d.", ja_purge_limit))

		if _, err := lib.DBexec("delete from ja_run_jobnet_table;"); err != nil {
			tc_79.Err_log("Failed to delete jobnet, Error: %s", err.Error())
			return FAILED
		}

		if _, err := lib.DBexec(`INSERT INTO ja_run_jobnet_table
						(inner_jobnet_id, inner_jobnet_main_id) VALUES (%d, %d);`,
			1500000000000000001, 1500000000000000001); err != nil {
			tc_79.Err_log("Failed to insert jobnet, Error: %s", err.Error())
			return FAILED
		}

		if _, err := lib.DBexec(`INSERT INTO ja_run_jobnet_summary_table (
    					inner_jobnet_id, update_date, created_date, invo_flag, run_type, status, job_status, 
						jobnet_abort_flag, load_status, scheduled_time, start_time, end_time, public_flag) 
						VALUES ('1500000000000000001', '20240704120626', '2024-07-04 07:50:40.064', 1, 1, 3, 0,
    					0, 0, 0, '20240704035042', '20240704035044', 0);`); err != nil {
			tc_79.Err_log("Failed to insert into jobnet summary table, Error: %s", err.Error())
			return FAILED
		}

		record_counter := 0
		records_per_tx := 499
		maximum_records := 1000000
		var tx *sql.Tx
		var err error
		for i := 0; i < maximum_records; i++ {
			if record_counter == 0 {
				tx, err = common.DB.Begin()
				if err != nil {
					tc_79.Err_log("Faild to start the transaction, Error: %s", err.Error())
					return FAILED
				}
			}

			query := lib.ConvertParamPostgresToMysql(`INSERT INTO ja_run_job_table  (
						inner_job_id, inner_jobnet_id, inner_jobnet_main_id, inner_job_id_fs_link, created_date,
						invo_flag, job_type, test_flag, method_flag, force_flag, timeout_flag, status, boot_count,
						end_count, start_time, end_time, point_x, point_y, job_id, job_name, continue_flag, run_user, run_user_password) 
						VALUES ($1, '1500000000000000001', '1500000000000000001', 0, '2024-07-04 08:02:01.434', 0, 0, 0, 0, 0, 0, 3, 0, 0,
						'20240704040203', '20240704040203', 354, 19, 'START', '', 0, '', '');`)
			if _, err := tx.Exec(query, i); err != nil {
				tc_79.Err_log("Failed to insert data, Error: %s", err.Error())
				tx.Rollback()
				return FAILED
			}

			query = lib.ConvertParamPostgresToMysql(`INSERT INTO ja_run_flow_table (
						inner_flow_id, inner_jobnet_id, start_inner_job_id, end_inner_job_id, flow_type, flow_width) 
						VALUES ($1, '1500000000000000001', 14, 13, 0, 0);`)
			if _, err := tx.Exec(query, i); err != nil {
				tc_79.Err_log("Failed to insert data, Error: %s", err.Error())
				tx.Rollback()
				return FAILED
			}

			query = lib.ConvertParamPostgresToMysql(`INSERT INTO ja_value_after_jobnet_table  (
						inner_jobnet_id, value_name, after_value , seq_no) 
						VALUES ('1500000000000000001', 'JOBNET_ID', 'JOBNET_1', $1);`)
			if _, err := tx.Exec(query, i); err != nil {
				tc_79.Err_log("Failed to insert data, Error: %s", err.Error())
				tx.Rollback()
				return FAILED
			}

			query = lib.ConvertParamPostgresToMysql(`INSERT INTO ja_value_before_jobnet_table  (
						inner_jobnet_id, value_name, before_value , seq_no) 
						VALUES ('1500000000000000001', 'JOBNET_ID', 'JOBNET_1', $1);`)
			if _, err := tx.Exec(query, i); err != nil {
				tc_79.Err_log("Failed to insert data, Error: %s", err.Error())
				tx.Rollback()
				return FAILED
			}

			query = lib.ConvertParamPostgresToMysql(`INSERT INTO ja_run_value_before_table (
						inner_job_id, inner_jobnet_id, value_name, before_value, seq_no) 
						VALUES ($1, '1500000000000000001', 'JOBNET_ID', 'JOBNET_1', $1);`)
			if _, err := tx.Exec(query, i); err != nil {
				tc_79.Err_log("Failed to insert data, Error: %s", err.Error())
				tx.Rollback()
				return FAILED
			}

			query = lib.ConvertParamPostgresToMysql(`INSERT INTO ja_run_value_after_table (
						inner_job_id, inner_jobnet_id, value_name, after_value, seq_no) 
						VALUES ($1, '1500000000000000001', 'JOBNET_ID', 'JOBNET_1', $1);`)
			if _, err := tx.Exec(query, i); err != nil {
				tc_79.Err_log("Failed to insert data, Error: %s", err.Error())
				tx.Rollback()
				return FAILED
			}

			query = lib.ConvertParamPostgresToMysql(`INSERT INTO ja_run_icon_job_table (
							inner_job_id, inner_jobnet_id, host_flag, stop_flag, command_type, timeout, host_name, stop_code, timeout_run_type) 
							VALUES ($1, '1500000000000000001', 0, 0, 0, 0, 'moon8.7', NULL, 0);`)
			if _, err := tx.Exec(query, i); err != nil {
				tc_79.Err_log("Failed to insert data, Error: %s", err.Error())
				tx.Rollback()
				return FAILED
			}

			query = lib.ConvertParamPostgresToMysql(`INSERT INTO ja_run_icon_task_table (
						inner_job_id, inner_jobnet_id) 
						VALUES ($1, '1500000000000000001');`)
			if _, err := tx.Exec(query, i); err != nil {
				tc_79.Err_log("Failed to insert data, Error: %s", err.Error())
				tx.Rollback()
				return FAILED
			}

			query = lib.ConvertParamPostgresToMysql(`INSERT INTO ja_run_icon_value_table (
					inner_job_id, inner_jobnet_id) 
					VALUES ($1, '1500000000000000001');`)
			if _, err := tx.Exec(query, i); err != nil {
				tc_79.Err_log("Failed to insert data, Error: %s", err.Error())
				tx.Rollback()
				return FAILED
			}

			query = lib.ConvertParamPostgresToMysql(`INSERT INTO ja_run_icon_release_table (
					inner_job_id, inner_jobnet_id) 
					VALUES ($1, '1500000000000000001');`)
			if _, err := tx.Exec(query, i); err != nil {
				tc_79.Err_log("Failed to insert data, Error: %s", err.Error())
				tx.Rollback()
				return FAILED
			}

			query = lib.ConvertParamPostgresToMysql(`INSERT INTO ja_run_icon_calc_table (
				inner_job_id, inner_jobnet_id
			) VALUES ($1, '1500000000000000001');`)
			if _, err := tx.Exec(query, i); err != nil {
				tc_79.Err_log("Failed to insert data, Error: %s", err.Error())
				tx.Rollback()
				return FAILED
			}

			query = lib.ConvertParamPostgresToMysql(`INSERT INTO ja_run_icon_reboot_table (
				inner_job_id, inner_jobnet_id
			) VALUES ($1, '1500000000000000001');`)
			if _, err := tx.Exec(query, i); err != nil {
				tc_79.Err_log("Failed to insert data, Error: %s", err.Error())
				tx.Rollback()
				return FAILED
			}

			query = lib.ConvertParamPostgresToMysql(`INSERT INTO ja_run_icon_fwait_table (
				inner_job_id, inner_jobnet_id
			) VALUES ($1, '1500000000000000001');`)
			if _, err := tx.Exec(query, i); err != nil {
				tc_79.Err_log("Failed to insert data, Error: %s", err.Error())
				tx.Rollback()
				return FAILED
			}

			query = lib.ConvertParamPostgresToMysql(`INSERT INTO ja_run_icon_info_table (
				inner_job_id, inner_jobnet_id
			) VALUES ($1, '1500000000000000001');`)
			if _, err := tx.Exec(query, i); err != nil {
				tc_79.Err_log("Failed to insert data, Error: %s", err.Error())
				tx.Rollback()
				return FAILED
			}

			query = lib.ConvertParamPostgresToMysql(`INSERT INTO ja_run_icon_zabbix_link_table (
				inner_job_id, inner_jobnet_id
			) VALUES ($1, '1500000000000000001');`)
			if _, err := tx.Exec(query, i); err != nil {
				tc_79.Err_log("Failed to insert data, Error: %s", err.Error())
				tx.Rollback()
				return FAILED
			}

			query = lib.ConvertParamPostgresToMysql(`INSERT INTO ja_run_icon_agentless_table (
				inner_job_id, inner_jobnet_id
			) VALUES ($1, '1500000000000000001');`)
			if _, err := tx.Exec(query, i); err != nil {
				tc_79.Err_log("Failed to insert data, Error: %s", err.Error())
				tx.Rollback()
				return FAILED
			}

			query = lib.ConvertParamPostgresToMysql(`INSERT INTO ja_run_icon_jobnet_table (
				inner_job_id, inner_jobnet_id
			) VALUES ($1, '1500000000000000001');`)
			if _, err := tx.Exec(query, i); err != nil {
				tc_79.Err_log("Failed to insert data, Error: %s", err.Error())
				tx.Rollback()
				return FAILED
			}

			query = lib.ConvertParamPostgresToMysql(`INSERT INTO ja_run_icon_end_table (
				inner_job_id, inner_jobnet_id
			) VALUES ($1, '1500000000000000001');`)
			if _, err := tx.Exec(query, i); err != nil {
				tc_79.Err_log("Failed to insert data, Error: %s", err.Error())
				tx.Rollback()
				return FAILED
			}

			query = lib.ConvertParamPostgresToMysql(`INSERT INTO ja_run_icon_extjob_table (
				inner_job_id, inner_jobnet_id
			) VALUES ($1, '1500000000000000001');`)
			if _, err := tx.Exec(query, i); err != nil {
				tc_79.Err_log("Failed to insert data, Error: %s", err.Error())
				tx.Rollback()
				return FAILED
			}

			query = lib.ConvertParamPostgresToMysql(`INSERT INTO ja_run_icon_if_table (
				inner_job_id, inner_jobnet_id
			) VALUES ($1, '1500000000000000001');`)
			if _, err := tx.Exec(query, i); err != nil {
				tc_79.Err_log("Failed to insert data, Error: %s", err.Error())
				tx.Rollback()
				return FAILED
			}

			query = lib.ConvertParamPostgresToMysql(`INSERT INTO ja_run_icon_fcopy_table (
				inner_job_id, inner_jobnet_id
			) VALUES ($1, '1500000000000000001');`)
			if _, err := tx.Exec(query, i); err != nil {
				tc_79.Err_log("Failed to insert data, Error: %s", err.Error())
				tx.Rollback()
				return FAILED
			}

			if record_counter >= records_per_tx {
				if err = tx.Commit(); err != nil {
					tc_79.Err_log("Faild to commit the transaction, Error: %s", err.Error())
					return FAILED
				}
				fmt.Println("commited", i+1)
				record_counter = 0
				continue
			}
			record_counter += 1
			fmt.Println("inserted, counter", i+1, record_counter)
		}
		fmt.Println(tc_79.Info_log("%d records have been inserted to DB.", maximum_records))

		// Backup and clean jobarg server log
		if err := moveAndBackupJAZServerLog(); err != nil {
			tc_79.Err_log("Faild to clean the JAZ server log, Error: %s", err.Error())
			return FAILED
		}

		// Restart the jobarg server
		if err := lib.Restart_jaz_server(); err != nil {
			tc_79.Err_log("Faild to restart the JAZ server, Error: %s", err.Error())
			return FAILED
		}
		fmt.Println(tc_79.Info_log("JAZ server has been restarted."))

		// Check the expected result; jobs are being purged at a rate of JaPurgeLimit
		var server_log string
		if err := checkPurgeLog(ja_purge_limit, 3, common.Client, &server_log); err != nil {
			tc_79.Err_log("Faild to check the JAZ server log for purge limit, Error: %s", err.Error())
			return FAILED
		}
		fmt.Println(tc_79.Info_log("The expected JAZ server logs are found. %s", server_log))

		// Check if all data inside the database gets deleted
		timeout_min := 240 // timeout in minutes
		if err := checkRunIconCount(timeout_min, 0); err != nil {
			tc_79.Err_log("Faild to check the run icon count, Error: %s", err.Error())
			return FAILED
		}
		fmt.Println(tc_79.Info_log("All run job icons are deleted from DB."))

		return PASSED
	}
	tc_79.Set_function(tc_func)
	t.Add_testcase(*tc_79)
}

// Backup and clean JAZ server log
func moveAndBackupJAZServerLog() error {

	// Get the current time
	currentTime := time.Now()

	// Format the time as yyyymmddhhmmss
	formattedTime := currentTime.Format("20060102150405")

	_, err := lib.Ssh_exec("mv /var/log/jobarranger/jobarg_server.log /var/log/jobarranger/jobarg_server.log.bk.1225." + formattedTime)

	if err != nil {
		return fmt.Errorf("failed to clean jaz server log: %s", err.Error())
	}

	return nil
}

// This function fetches JAZ server purge log, and checks if the purge rate is correct; jobs need to be purged at a rate of JaPurgeLimit
func checkPurgeLog(jaPurgeLimit int, timeoutDuration int, client *ssh.Client, serverLog *string) error {
	// set timeout
	timeout := time.After(time.Duration(timeoutDuration) * time.Minute)

	for {
		select {
		case <-timeout:
			return fmt.Errorf("timeout after %d minutes while checking purge log", timeoutDuration)
		default:
			purgeLog, err := lib.GetOutputStrFromSSHCommand(client, "grep 'Delete not finished yet from ja_run_flow_table' /var/log/jobarranger/jobarg_server.log | head -3")

			if err != nil {
				return fmt.Errorf("failed to fetch purge log: %s", err.Error())
			}

			// trim the log
			purgeLog = strings.TrimSpace(purgeLog)

			// get array of lines from log
			lines := strings.Split(purgeLog, "\n")

			if len(lines) < 3 {
				time.Sleep(1 * time.Second)
				continue
			}

			*serverLog = fmt.Sprintf("Server log: \n%s", purgeLog)

			// Regular expression to extract the remaining record count from each line
			re := regexp.MustCompile(`left records: (\d+)`)
			nextNum := 0

			// Check if the left records have a difference of jaPurgeLimit
			for _, line := range lines {
				matches := re.FindStringSubmatch(line)
				if len(matches) > 1 {
					num, err := strconv.Atoi(matches[1])

					if err != nil {
						return fmt.Errorf("failed to convert left records to int: %s", err.Error())
					}

					// first encounter
					if nextNum == 0 {
						nextNum = num - jaPurgeLimit
						continue
					}

					if num != nextNum {
						return fmt.Errorf("the job purge rate is incorrect. Expected left records: %d, Actual left records: %d", nextNum, num)
					}

					nextNum = num - jaPurgeLimit
				}
			}

			return nil

		}

	}
}

// Check the run icon count in database until it reaches the target count
func checkRunIconCount(timeoutDuration int, targetCount int) error {
	// set timeout
	timeout := time.After(time.Duration(timeoutDuration) * time.Minute)

	for {
		select {
		case <-timeout:
			return fmt.Errorf("timeout after %d minutes while checking run icon count", timeoutDuration)
		default:
			var count int
			if err := lib.GetSingleRow(lib.CheckAllRunCount, nil, &count); err != nil {
				return fmt.Errorf("failed to get the count of running icons: %s", err.Error())
			}

			fmt.Printf("Run icon count: %d\n", count)

			if count == targetCount {
				return nil
			}

			time.Sleep(time.Second)

		}
	}

}
