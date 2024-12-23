package common

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"syscall"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/crypto/ssh"
	"golang.org/x/term"
)

// type Host_temp struct {
// 	Host_name, Host_run_username string
// 	Host_ip, Host_dns            string      `json:"-"`
// 	Host_client                  *ssh.Client `json:"-"`
// 	Host_use_ip                  bool        `json:"-"`
// 	Host_port                    int
// }

type Testcase_status string
type Database string
type Doc_data_type string
type Host_type string

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

	LINUX_SERVER  Host_type = "auto.linux.server."
	LINUX_AGENT   Host_type = "auto.linux.agent."
	WINDOWS_AGENT Host_type = "auto.windows.agent."

	LOG_LEVEL_INFO = 1
	LOG_LEVEL_ERR  = 2
)

var Log_filepath, DB_hostname, DB_user, DB_passwd, DB_name string
var Specific_ticket_no, Specific_testcase_no, DB_port, Timeout, Current_tk_no, Current_tc_no uint
var Client *ssh.Client // server's ssh client
var Login_info Auth
var Is_mysql, Is_psql bool
var DB_type Database
var DB *sql.DB
var Sugar *zap.SugaredLogger
var Host_pool []Host
var Linux_server_host Host // only supports single server.

func Set_linux_server_host() {
	for _, host := range Host_pool {
		if host.Get_Host_type() == LINUX_SERVER {
			Linux_server_host = host
		}
	}
}

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
	if DB_hostname == "" {
		DB_hostname = Login_info.Hostname
	}
}

func Set_default_db_port() {
	if Is_mysql && DB_port == 0 {
		DB_port = 3306
	} else if Is_psql && DB_port == 0 {
		DB_port = 5432
	}
}

func Set_db_type() error {
	if !Is_mysql && !Is_psql {
		return fmt.Errorf("choose db type using --with-mysql or --with-postgresql flag")
	}

	if Is_mysql {
		DB_type = MYSQL
	} else {
		DB_type = PSQL
	}

	return nil
}

func Set_usr_hst(args []string) error {
	parts := strings.Split(args[0], "@")
	if len(parts) != 2 {
		fmt.Println("Usage: remote_run user@host <port>")
		os.Exit(1)
	}

	Login_info.Username = parts[0]
	Login_info.Hostname = parts[1]

	return nil
}

func Set_passwd() {
	fmt.Printf("%s@%s's password:", Login_info.Username, Login_info.Hostname)
	bytepw, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println() //new line
	if err != nil {
		fmt.Println("Failed in getting password, Error:", err.Error())
		os.Exit(1)
	}
	Login_info.Password = string(bytepw)
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
