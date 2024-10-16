package tickets

import (
	"database/sql"
	"fmt"

	"github.com/zukigit/remote_run-go/src/common"
	"github.com/zukigit/remote_run-go/src/dao"
	"github.com/zukigit/remote_run-go/src/lib"
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
		records_per_tx := 11
		var tx *sql.Tx
		var err error
		for i := 0; i < 12; i++ {
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
				fmt.Println("commited", i)
				record_counter = 0
			}
			record_counter += 1
			fmt.Println("inserted, counter", i, record_counter)
		}

		return FAILED
	}
	tc_79.Set_function(tc_func)
	t.Add_testcase(*tc_79)
}
