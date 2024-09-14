package testcases

import (
	"zukigit/remote_run-go/src/lib"

	"golang.org/x/crypto/ssh"
)

type Test_case_1 struct {
	test_case_id, test_case_description, test_case_log string
	session                                            *ssh.Session
	is_passed                                          bool
}

func (T *Test_case_1) Set_tc_values(session *ssh.Session) {
	T.test_case_id = "TestCase1"
	T.test_case_description = "Test 1 is goood lol"
	T.session = session
}

func (T *Test_case_1) Get_tc_id() string {
	return T.test_case_id
}

func (T *Test_case_1) Get_tc_dsctn() string {
	return T.test_case_description
}

func (T *Test_case_1) Set_tc_log(tc_log string) {
	T.test_case_log = tc_log
}

func (T *Test_case_1) Get_tc_log() string {
	return T.test_case_log
}

func (T *Test_case_1) Set_is_passed(is_passed bool) {
	T.is_passed = is_passed
}

func (T *Test_case_1) Get_is_passed() bool {
	return T.is_passed
}

func (T *Test_case_1) Get_session() *ssh.Session {
	return T.session
}

func (T *Test_case_1) Run() bool {
	lib.Set_test_case(T)

	// Your logic from here
	output, error := lib.Exec("pwd")
	if error != nil {
		lib.Err_log(error.Error())
		return false
	}

	lib.Info_log("output: %s", output)
	lib.Info_log("output: %s", output)

	return true
}
