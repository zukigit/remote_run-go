package testcases

import (
	"golang.org/x/crypto/ssh"
)

type Test_case_1 struct {
	test_case_id, test_case_description, test_case_log string
	session                                            *ssh.Session
}

func (T *Test_case_1) Set_tc_values(session *ssh.Session) {
	T.test_case_id = "Test case 1"
	T.test_case_description = "Test 1 is goood lol"
	T.session = session
}

func (T *Test_case_1) Get_tc_id() string {
	return T.test_case_id
}

func (T *Test_case_1) Get_tc_dsctn() string {
	return T.test_case_description
}

func (T *Test_case_1) Get_tc_log() string {
	return T.test_case_log
}

func (T *Test_case_1) Run() bool {
	T.test_case_log = "it ran successfully"

	return true
}
