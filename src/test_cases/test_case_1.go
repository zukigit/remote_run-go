package testcases

import (
	"fmt"

	"golang.org/x/crypto/ssh"
)

type Test_case_1 struct {
	test_case_id, test_case_description string
	session                             *ssh.Session
}

func (T *Test_case_1) Set_tc_values(session *ssh.Session) {
	T.test_case_id = "Test case 1"
	T.test_case_description = "Test 1 is goood lol"
	T.session = session
}

func (T *Test_case_1) Run() {
	fmt.Println("Ran", T)
}
