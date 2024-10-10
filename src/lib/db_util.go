package lib

import (
	"database/sql"
	"fmt"
	"os"
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
	AbortJobnetQuery DBQuery = "UPDATE ja_run_jobnet_summary_table SET jobnet_abort_flag = 1 WHERE inner_jobnet_id = $1"
)

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
	db.SetConnMaxLifetime(5 * time.Minute)

	// Check the connection
	err = db.Ping()
	if err != nil {
		fmt.Println("Could not connect to database, Err: ", err.Error())
		os.Exit(1)
	}

	common.DB = db
}

// InsertData inserts data into a table
func InsertData(query DBQuery, args ...interface{}) (sql.Result, error) {
	stmt, err := common.DB.Prepare(string(query))
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	return stmt.Exec(args...)
}

// GetData fetches rows based on a query
func GetData(query DBQuery, args ...interface{}) (*sql.Rows, error) {
	rows, err := common.DB.Query(string(query), args...)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// UpdateData updates records in a table
func UpdateData(query DBQuery, args ...interface{}) (sql.Result, error) {
	return InsertData(query, args...)
}

// DeleteData deletes records from a table
func DeleteData(query DBQuery, args ...interface{}) (sql.Result, error) {
	return InsertData(query, args...)
}
