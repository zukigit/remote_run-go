package common

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"syscall"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/crypto/ssh"
	"golang.org/x/term"
)

type Testcase_status string
type Database string
type Doc_data_type string

const LOG_LEVEL_INFO = 1
const LOG_LEVEL_ERR = 2

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
)

var Log_filepath, DB_hostname, DB_user, DB_passwd, DB_name string
var Specific_ticket_no, Specific_testcase_no, DB_port, Timeout, Current_tk_no, Current_tc_no uint
var Client *ssh.Client
var Login_info Auth
var Is_mysql, Is_psql bool
var DB_type Database
var DB *sql.DB
var Sugar *zap.SugaredLogger

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
		return fmt.Errorf("please choose db type using --with-mysql --with-postgresql")
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

func Set_client() {
	config := &ssh.ClientConfig{
		User: Login_info.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(Login_info.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	hostname_with_port := fmt.Sprintf("%s:%d", Login_info.Hostname, Login_info.Port)

	client, err := ssh.Dial("tcp", hostname_with_port, config)
	if err != nil {
		fmt.Println("Failed in getting ssh client, Error:", err.Error())
		os.Exit(1)
	}

	Client = client
}
