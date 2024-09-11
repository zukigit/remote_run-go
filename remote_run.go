package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"syscall"
	"zukigit/remote_run-go/src/lib"
	testcases "zukigit/remote_run-go/src/test_cases"

	"golang.org/x/term"
)

type auth struct {
	username, password, hostname string
}

func get_auth() *auth {
	arg_len := len(os.Args)
	var host string
	port := "22"

	// check args
	if arg_len == 2 {
		host = os.Args[1]
	} else if arg_len == 3 {
		host = os.Args[1]
		port = os.Args[2]

		// check port
		if _, err := strconv.Atoi(port); err != nil {
			fmt.Printf("Error: Bad port '%s'\n", port)
			os.Exit(1)
		}
	} else {
		fmt.Println("Usage: remote_run user@host <port>")
		os.Exit(1)
	}

	// split user and host
	parts := strings.Split(host, "@")
	if len(parts) != 2 {
		fmt.Println("Usage: remote_run user@host <port>")
		os.Exit(1)
	}

	user := parts[0]
	host = parts[1] + ":" + port

	// get password
	fmt.Printf("%s's password:", os.Args[1])
	bytepw, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println() //new line
	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}
	password := string(bytepw)

	return &auth{
		username: user,
		password: password,
		hostname: host,
	}
}

type Test_case interface {
	Run()
	Set_tc_values()
}

func set_tc_default_values(t []Test_case) {
	for _, test_case := range t {
		test_case.Set_tc_values()
	}
}

func run_tc(t []Test_case) {
	for _, test_case := range t {
		test_case.Run()
	}
}

func main() {
	var test_cases []Test_case

	auth := get_auth()
	config := lib.Get_config(auth.username, auth.password)

	// Connect to the SSH server
	client := lib.Get_client(auth.hostname, config)
	defer client.Close()

	// Create a session for the command
	session := lib.Get_session(client)
	defer session.Close()

	// add test cases from here
	test_cases = append(test_cases, new(testcases.Test_case_1))

	set_tc_default_values(test_cases)

	run_tc(test_cases) // run test cases
}
