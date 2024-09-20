package dao

import (
	"fmt"
	"zukigit/remote_run-go/src/lib"
)

const INFO = 1
const ERR = 2

type TestCase struct {
	id          uint
	description string
	logs        []string
	is_passed   bool
	function    func() bool
	auth        *Auth
}

func New_testcase(testcase_id uint, testcase_description string, auth *Auth) *TestCase {
	return &TestCase{
		id:          testcase_id,
		description: testcase_description,
		auth:        auth,
	}
}

func (t *TestCase) Set_log(log string) {
	t.logs = append(t.logs, log)
}

func (t *TestCase) Get_id() uint {
	return t.id
}

func (t *TestCase) Get_logs() []string {
	return t.logs
}

func (t *TestCase) Set_is_passed(is_passed bool) {
	t.is_passed = is_passed
}

func (t *TestCase) Set_function(function func() bool) {
	t.function = function
}

func (t *TestCase) Run_function() bool {
	return t.function()
}

func (t *TestCase) Is_function_nil() bool {
	return t.function == nil
}

func (t *TestCase) Logi(level int, log string) string {
	log = lib.Formatted_log(level, log)
	t.Set_log(log)

	return log
}

func (t *TestCase) Err_log(unfmt string, arg ...any) string {
	log := fmt.Sprintf(unfmt, arg...)
	return t.Logi(ERR, log)
}

func (t *TestCase) Info_log(unfmt string, arg ...any) string {
	log := fmt.Sprintf(unfmt, arg...)
	return t.Logi(INFO, log)
}

func (t *TestCase) Ssh_exec(command string) ([]byte, error) {
	return lib.Ssh_exec(command, t.auth.Session)
}

func (t *TestCase) Ssh_exec_to_str(command string) (string, error) {
	return lib.Ssh_exec_to_str(command, t.auth.Session)
}
