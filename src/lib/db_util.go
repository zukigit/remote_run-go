package lib

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"github.com/zukigit/remote_run-go/src/common"
)

type DBQuery string

const (
	DeleteRunJobnetQuery     DBQuery = "DELETE FROM ja_run_jobnet_table"
	AbortSingleFWaitJobQuery DBQuery = `UPDATE ja_run_job_table SET method_flag = 3 WHERE inner_job_id = (
		SELECT inner_job_id FROM ja_run_icon_fwait_table WHERE inner_jobnet_id = $1 limit 1
	)`

	CheckJobStatusCountQuery DBQuery = "SELECT count(*) FROM ja_run_job_table where status = 4 and job_type = 4 and inner_jobnet_main_id = $1"
	AbortExtJobQuery        DBQuery = `UPDATE ja_run_jobnet_summary_table SET jobnet_abort_flag = 1 WHERE inner_jobnet_id = ?`
	AbortJobnetQuery        DBQuery = "UPDATE ja_run_jobnet_summary_table SET jobnet_abort_flag = 1 WHERE inner_jobnet_id = $1"
	AbortSingleJOBIconQuery DBQuery = "UPDATE ja_run_job_table SET method_flag = 3 WHERE inner_jobnet_id = $1"
	CheckJobnetDoneWithRed  DBQuery = "select * from ja_run_jobnet_table where status = 2 or status = 6"
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
			queryStr := convertParamPostgresToMysql(string(query))

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
