package lib

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

type DBType string

const (
	Postgres DBType = "postgres"
	MySQL    DBType = "mysql"
)

type DBQuery string

const (
	DeleteRunJobnetQuery     DBQuery = "DELETE FROM ja_run_jobnet_table"
	AbortSingleFWaitJobQuery DBQuery = `UPDATE ja_run_job_table SET method_flag = 3 WHERE inner_job_id = (
		SELECT inner_job_id FROM ja_run_icon_fwait_table WHERE inner_jobnet_id = $1 limit 1
	)`
	AbortJobnetQuery DBQuery = "UPDATE ja_run_jobnet_summary_table SET jobnet_abort_flag = 1 WHERE inner_jobnet_id = $1"
)

type DBUtil struct {
	DB *sql.DB
}

// ConnectDB initializes the database connection
func ConnectDB(dbType DBType, host, port, user, password, dbname string) (*DBUtil, error) {
	var dsn string
	switch dbType {
	case Postgres:
		dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			host, port, user, password, dbname)
	case MySQL:
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, port, dbname)
	default:
		return nil, fmt.Errorf("unsupported db type: %s", dbType)
	}

	db, err := sql.Open(string(dbType), dsn)
	if err != nil {
		return nil, err
	}

	// Set maximum number of open connections
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Check the connection
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	fmt.Printf("Successfully connected to the %s database!\n", dbType)
	return &DBUtil{DB: db}, nil
}

// CloseDB closes the database connection
func (dbu *DBUtil) CloseDB() {
	err := dbu.DB.Close()
	if err != nil {
		log.Fatal("Failed to close the database:", err)
	}
	fmt.Println("Database connection closed.")
}

// InsertData inserts data into a table
func (dbu *DBUtil) InsertData(query DBQuery, args ...interface{}) (sql.Result, error) {
	stmt, err := dbu.DB.Prepare(string(query))
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	return stmt.Exec(args...)
}

// GetData fetches rows based on a query
func (dbu *DBUtil) GetData(query DBQuery, args ...interface{}) (*sql.Rows, error) {
	rows, err := dbu.DB.Query(string(query), args...)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// UpdateData updates records in a table
func (dbu *DBUtil) UpdateData(query DBQuery, args ...interface{}) (sql.Result, error) {
	return dbu.InsertData(query, args...)
}

// DeleteData deletes records from a table
func (dbu *DBUtil) DeleteData(query DBQuery, args ...interface{}) (sql.Result, error) {
	return dbu.InsertData(query, args...)
}
