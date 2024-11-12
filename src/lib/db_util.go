package lib

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"github.com/zukigit/remote_run-go/src/common"
	"golang.org/x/crypto/ssh"
)

type DBQuery string

const (
	DeleteRunJobnetQuery     DBQuery = "DELETE FROM ja_run_jobnet_table"
	AbortSingleFWaitJobQuery DBQuery = `UPDATE ja_run_job_table SET method_flag = 3 WHERE inner_job_id = (
		SELECT inner_job_id FROM ja_run_icon_fwait_table WHERE inner_jobnet_id = $1 limit 1
	)`

	CheckJobStatusCountQuery DBQuery = "SELECT count(*) FROM ja_run_job_table where status = 4 and job_type = 4 and inner_jobnet_main_id = $1"
	AbortExtJobQuery         DBQuery = `UPDATE ja_run_jobnet_summary_table SET jobnet_abort_flag = 1 WHERE inner_jobnet_id = ?`
	AbortJobnetQuery         DBQuery = "UPDATE ja_run_jobnet_summary_table SET jobnet_abort_flag = 1 WHERE inner_jobnet_id = $1"
	AbortSingleJOBIconQuery  DBQuery = "UPDATE ja_run_job_table SET method_flag = 3 WHERE inner_jobnet_id = $1"
	CheckJobnetDoneWithRed   DBQuery = "select * from ja_run_jobnet_table where status = 2 or status = 6"
	CheckAllRunCount         DBQuery = `SELECT (
		(SELECT COUNT(*) FROM ja_run_jobnet_table) +
		(SELECT COUNT(*) FROM ja_run_jobnet_summary_table) +
		(SELECT COUNT(*) FROM ja_run_flow_table) +
		(SELECT COUNT(*) FROM ja_value_after_jobnet_table) +
		(SELECT COUNT(*) FROM ja_value_before_jobnet_table) +
		(SELECT COUNT(*) FROM ja_run_value_before_table) +
		(SELECT COUNT(*) FROM ja_run_value_after_table) +
		(SELECT COUNT(*) FROM ja_run_icon_task_table) +
		(SELECT COUNT(*) FROM ja_run_icon_value_table) +
		(SELECT COUNT(*) FROM ja_run_icon_release_table) +
		(SELECT COUNT(*) FROM ja_run_icon_calc_table) +
		(SELECT COUNT(*) FROM ja_run_icon_reboot_table) +
		(SELECT COUNT(*) FROM ja_run_icon_fwait_table) +
		(SELECT COUNT(*) FROM ja_run_icon_info_table) +
		(SELECT COUNT(*) FROM ja_run_icon_zabbix_link_table) +
		(SELECT COUNT(*) FROM ja_run_icon_agentless_table) +
		(SELECT COUNT(*) FROM ja_run_icon_jobnet_table) +
		(SELECT COUNT(*) FROM ja_run_icon_end_table) +
		(SELECT COUNT(*) FROM ja_run_icon_extjob_table) +
		(SELECT COUNT(*) FROM ja_run_icon_job_table) +
		(SELECT COUNT(*) FROM ja_run_icon_if_table) +
		(SELECT COUNT(*) FROM ja_run_icon_fcopy_table) +
		(SELECT COUNT(*) FROM ja_run_job_table)
	) AS total_count;`
)

// Converts the parameter in postgresql query to a compatible version for mysql
func ConvertParamPostgresToMysql(query string) string {
	if common.DB_type == common.MYSQL {
		for i := 1; strings.Contains(query, fmt.Sprintf("$%d", i)); i++ {
			query = strings.ReplaceAll(query, fmt.Sprintf("$%d", i), "?")
		}
	}

	return query
}

// ConnectDB initializes the database connection
func ConnectDB(user, password, dbname string) {
	var dsn string
	switch common.DB_type {
	case common.PSQL:
		dsn = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			common.DB_hostname, common.DB_port, user, password, dbname)
	case common.MYSQL:
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", user, password, common.DB_hostname, common.DB_port, dbname)
	default:
		fmt.Printf("Err: unsupported db type: %s\n", common.DB_type)
		os.Exit(1)
	}

	db, err := sql.Open(string(common.DB_type), dsn)
	if err != nil {
		fmt.Println("Could not connect to database, Err: ", err.Error())
		os.Exit(1)
	}

	// Set maximum number of open connections
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(0 * time.Minute)

	// Check the connection
	err = db.Ping()
	if err != nil {
		fmt.Println("Could not connect to database, Err: ", err.Error())
		os.Exit(1)
	}

	common.DB = db
}

// ExecuteQuery that changes the state of the database
func ExecuteQuery(query DBQuery, args ...interface{}) (sql.Result, error) {
	queryStr := ConvertParamPostgresToMysql(string(query))
	stmt, err := common.DB.Prepare(queryStr)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	return stmt.Exec(args...)
}

// GetData fetches rows based on a query
func GetData(query DBQuery, args ...interface{}) (*sql.Rows, error) {
	queryStr := ConvertParamPostgresToMysql(string(query))
	rows, err := common.DB.Query(queryStr, args...)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// This function will execute the query that will get exactly one row
func GetSingleRow(query DBQuery, args []interface{}, dest ...interface{}) error {
	// Prepare the query with arguments, then scan the result into the provided destination variables
	return common.DB.QueryRow(string(query), args...).Scan(dest...)
}

func DBexec(unfmt string, arg ...any) (sql.Result, error) {
	query := fmt.Sprintf(unfmt, arg...)

	return common.DB.Exec(query)
}

// Check count of the query
func JobProcessDBCountCheck(targetProcessCount int, timeoutDuration int, inner_jobnet_main_id string, query DBQuery, args ...interface{}) error {

	// set timeout
	timeout := time.After(time.Duration(timeoutDuration) * time.Minute)
	c := 0
	for {
		select {
		case <-timeout:
			return fmt.Errorf("timeout after %d minutes", timeoutDuration)
		default:
			// Convert query parameters (assuming convertParamPostgresToMysql is a valid function)
			queryStr := ConvertParamPostgresToMysql(string(query))

			// Execute query
			rows, err := common.DB.Query(queryStr, inner_jobnet_main_id)
			if err != nil {
				return fmt.Errorf("query execution error: %w", err)
			}
			// Ensure rows are closed after processing to avoid resource leaks
			defer rows.Close()

			// Variable to hold the count
			var count int

			// Fetch the first row
			if rows.Next() {
				err = rows.Scan(&count)
				if err != nil {
					return fmt.Errorf("error scanning row: %w", err)
				}
			} else {
				// Handle case where no rows are returned
				return fmt.Errorf("no records found")
			}

			// Check if the count matches the targetProcessCount
			if count == targetProcessCount {
				fmt.Println("Count : ", count)
				return nil
			} else if count > targetProcessCount {
				fmt.Println("Actual count is greater than the target count ", c, ". Count :", count)
				c++
			} else if count < targetProcessCount {
				fmt.Println("Actual count is less than the target count ", c, ". Count :", count)
				c++
			}

			// Sleep for 30 second before retrying
			time.Sleep(30 * time.Second)

		}
	}
}

// StopDatabaseService stops the service on the remote Linux machine using SSH
func StopDatabaseService(client *ssh.Client, serviceName string) error {
	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}
	defer session.Close()

	// Command to stop the service (e.g., mysqld or postgresql)
	cmd := fmt.Sprintf("sudo systemctl stop %s", serviceName)
	err = session.Run(cmd)
	if err != nil {
		return fmt.Errorf("failed to stop service %s: %w", serviceName, err)
	}

	log.Printf("Successfully stopped service: %s", serviceName)
	return nil
}

// CheckAndStopDBService executes the query and stops the database service if conditions are met
func CheckAndStopDBService(db *sql.DB, client *ssh.Client, dbType string, query string, timeout, interval time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Timeout reached; no matching data found.")
			return nil

		case <-ticker.C:
			log.Println("Executing query...")

			// Define variables to capture the row's columns
			var jobID string
			var status, jobType int

			// Execute the query and scan the result into variables
			err := db.QueryRow(query).Scan(&jobID, &status, &jobType)

			// Handle no rows found by continuing to the next interval
			if err == sql.ErrNoRows {
				log.Println("No matching data found; continuing to next interval.")
				continue
			} else if err != nil {
				log.Printf("Error executing query: %v\n", err)
				return fmt.Errorf("error executing query: %w", err)
			}

			// Log the result for debugging
			log.Printf("Matching data found: stop the db...")

			// Stop the database service if a matching row was found
			if dbType == "mysql" {
				err := StopDatabaseService(client, "mysqld")
				if err != nil {
					return fmt.Errorf("failed to stop MySQL service: %w", err)
				}
			} else if dbType == "postgresql" {
				err := StopDatabaseService(client, "postgresql")
				if err != nil {
					return fmt.Errorf("failed to stop PostgreSQL service: %w", err)
				}
			}
			return nil
		}
	}
}

// StartDatabaseService starts the specified database service on the remote Linux machine using SSH
func StartDatabaseService(client *ssh.Client, serviceName string) error {
	// Create a new SSH session
	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}
	defer session.Close()

	// Prepare the command to start the service
	cmd := fmt.Sprintf("sudo systemctl start %s", serviceName)

	// Run the command to start the service
	err = session.Run(cmd)
	if err != nil {
		return fmt.Errorf("failed to start service %s: %w", serviceName, err)
	}

	// Log success
	log.Printf("Successfully started service: %s", serviceName)
	return nil
}
