package common

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	"golang.org/x/crypto/ssh"
	"golang.org/x/term"
)

const INFO = 1
const ERR = 2

type Testcase_status string

var Left_string, Right_string, Endticket_string, Endtestcase_string, Log_filename string
var Specific_ticket_no uint
var Client *ssh.Client
var Login_info Auth
var Log_file *os.File

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

func Set_specific_ticket_no(args []string) {
	if len(args) < 2 {
		Specific_ticket_no = 0
	} else {
		num, err := strconv.ParseUint(args[1], 10, 64)
		if err != nil {
			Specific_ticket_no = 0
		} else {
			Specific_ticket_no = uint(num)
		}
	}
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
