package common

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"golang.org/x/crypto/ssh"
	"golang.org/x/term"
)

type Testcase_status string
type Database string

const INFO = 1
const ERR = 2

const (
	MYSQL Database = "mysql"
	PSQL  Database = "postgres"
)

var Left_string, Right_string, Endticket_string, Endtestcase_string, Log_filename, DB_hostname string
var Specific_ticket_no, Specific_testcase_no, DB_port uint
var Client *ssh.Client
var Login_info Auth
var Log_file *os.File
var Is_mysql, Is_psql bool
var DB_type Database
var DB *sql.DB

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

func Set_log_file(file_name string) {
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}

	sub_dir := filepath.Join(currentDir, "logs")
	file_path := filepath.Join(sub_dir, file_name)

	if _, err := os.Stat(sub_dir); os.IsNotExist(err) {
		err = os.Mkdir(sub_dir, 0755) // Create the directory with read/write permissions
		if err != nil {
			fmt.Println("Error:", err.Error())
			os.Exit(1)
		}
	}

	file, err := os.OpenFile(file_path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}
	Log_file = file
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

func Set_ticket_logs_headers() {
	Left_string = strings.Repeat("/", 60)
	Right_string = strings.Repeat("/", 60)
	Endticket_string = fmt.Sprintf("%s%s", Left_string, Right_string)

	Left_string = strings.Repeat("-", 25)
	Right_string = strings.Repeat("-", 25)
	Endtestcase_string = fmt.Sprintf("%s%s", Left_string, Right_string)
}
