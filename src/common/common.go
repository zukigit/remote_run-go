package common

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/crypto/ssh"
)

type Testcase_status string
type Database string
type Doc_data_type string
type Host_type string // LS_HOST_TYPE, LA_HOST_TYPE and WA_HOST_TYPE

const (
	MYSQL    Database      = "mysql"
	PSQL     Database      = "postgres"
	NORMAL                 = "NORMAL"
	TIMEOUT                = "TIMEOUT"
	ERROR                  = "ERROR"
	RUN                    = "RUN"
	END                    = "END"
	RUNERR                 = "RUNERR"
	ENDERR                 = "ENDERR"
	PRE_OPT  Doc_data_type = "pre_operation"
	OPT      Doc_data_type = "operation"
	EXPT_RES Doc_data_type = "expected_results"

	PASSED     Testcase_status = "PASSED"
	FAILED     Testcase_status = "FAILED"
	MUST_CHECK Testcase_status = "MUST_CHECK"

	LS_HOST_TYPE Host_type = "auto.linux.server."  // Host_type: linux server
	LA_HOST_TYPE Host_type = "auto.linux.agent."   // Host_type: linux agent
	WA_HOST_TYPE Host_type = "auto.windows.agent." // Host_type: windows agent

	LOG_LEVEL_INFO = 1
	LOG_LEVEL_ERR  = 2
)

var Log_filepath, DB_hostname, Temp_mysqlDB_hostname, Temp_psqlDB_hostname, DB_user, DB_passwd, DB_name string
var Specific_ticket_no, Specific_testcase_no, DB_port, Timeout, Current_tk_no, Current_tc_no uint
var Run_alltc_flag bool

// Deprecated: Use Server_host or Hosts instead.
var Client *ssh.Client // will be deleted.
// Deprecated: Use Server_host instead.
var Login_info Auth
var DB_type Database
var DB *sql.DB
var Sugar *zap.SugaredLogger
var Hosts []Host // host pool
var Server_host Host

func Set_sugar(logfile_path string) {
	logger_conf := zap.NewProductionConfig()
	logger_conf.Encoding = "console"
	logger_conf.EncoderConfig = zapcore.EncoderConfig{
		MessageKey:       "msg",
		ConsoleSeparator: " ",
	}

	logger_conf.OutputPaths = []string{
		logfile_path,
	}

	logger, err := logger_conf.Build()
	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}

	Sugar = logger.Sugar()
}

func Set_db_hostname() {
	if Temp_mysqlDB_hostname != "" {
		DB_hostname = Temp_mysqlDB_hostname
		DB_type = MYSQL

		if DB_port == 0 {
			DB_port = 3306
		}
	} else if Temp_psqlDB_hostname != "" {
		DB_hostname = Temp_psqlDB_hostname
		DB_type = PSQL

		if DB_port == 0 {
			DB_port = 5432
		}
	}
}

func Run_testcase(tc TestCase) {
	// start time
	startTime := time.Now()

	tc.Set_status(tc.Run_function())

	// total elasped time or duration of testcase
	duration := time.Since(startTime)
	durationStr := fmt.Sprintf("%02d:%02d:%02d", int(duration/time.Hour), int(duration/time.Minute)%60, int(duration/time.Second)%60)

	tc.Set_duration(durationStr)
}

func Update_testcase_results_in_tickets(tks []Ticket) {
	var passed_count, failed_count, mustcheck_count int
	for _, ticket := range tks {
		for _, testcase := range ticket.Get_testcases() {
			switch testcase.Get_status() {
			case PASSED:
				passed_count++
			case FAILED:
				failed_count++
			default:
				mustcheck_count++
			}
		}
		ticket.Set_PASSED_count(passed_count)
		ticket.Set_FAILED_count(failed_count)
		ticket.Set_MUSTCHECK_count(mustcheck_count)
	}
}
