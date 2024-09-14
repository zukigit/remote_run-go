package main

import (
	"fmt"
	"os"
	"zukigit/remote_run-go/src/dao"
	"zukigit/remote_run-go/src/lib"
	testcases "zukigit/remote_run-go/src/test_cases"

	"golang.org/x/crypto/ssh"
)

func set_tc_default_values(t []dao.Test_case, session *ssh.Session) {
	for _, test_case := range t {
		test_case.Set_tc_values(session)
	}
}

func run_tc(t []dao.Test_case) {
	for _, test_case := range t {
		test_case.Set_is_passed(test_case.Run())
		lib.Write_tc_log(test_case)
	}
}

func add_test_cases(test_cases *[]dao.Test_case) {
	// Add test cases here
	*test_cases = append(*test_cases, new(testcases.Test_case_1))
	// is it rainy day?
}

func main() {
	var test_cases []dao.Test_case
	add_test_cases(&test_cases)

	auth := dao.Get_auth() // Get login informations from user
	config := lib.Get_config(auth.Username, auth.Password)

	// Connect to the SSH server
	client, err := lib.Get_client(auth.Hostname, config)
	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}
	defer client.Close()

	// Create a session for the command
	session, err := lib.Get_session(client)
	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}
	defer session.Close()

	set_tc_default_values(test_cases, session)

	run_tc(test_cases) // run test cases
}
