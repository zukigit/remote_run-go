package dao

import (
	"fmt"

	"github.com/zukigit/remote_run-go/src/common"
	"github.com/zukigit/remote_run-go/src/lib"
)

type TestCase struct {
	id          uint
	description string
	logs        *[]string
	status      *common.Testcase_status
	function    func() common.Testcase_status
}

func New_testcase(testcase_id uint, testcase_description string) *TestCase {
	status := FAILED
	logs := []string{}
	return &TestCase{
		id:          testcase_id,
		description: testcase_description,
		status:      &status,
		logs:        &logs,
	}
}

func (t *TestCase) Set_log(log string) {
	*t.logs = append(*t.logs, log)
}

func (t *TestCase) Get_id() uint {
	return t.id
}

func (t *TestCase) Get_dsctn() string {
	return t.description
}

func (t *TestCase) Get_logs() []string {
	return *t.logs
}

func (t *TestCase) Set_status(status common.Testcase_status) {
	*t.status = status
}

func (t *TestCase) Get_status() common.Testcase_status {
	return *t.status
}

func (t *TestCase) Set_function(function func() common.Testcase_status) {
	t.function = function
}

func (t *TestCase) Run_function() common.Testcase_status {
	return t.function()
}

func (t *TestCase) Is_function_nil() bool {
	return t.function == nil
}

// From here is test case util functions

func (t *TestCase) Logi(level int, log string) string {
	log = fmt.Sprintf("Test_case[%d], %s", t.Get_id(), log)
	log = lib.Formatted_log(level, log)
	t.Set_log(log)

	return log
}

func (t *TestCase) Err_log(unfmt string, arg ...any) string {
	log := fmt.Sprintf(unfmt, arg...)
	return t.Logi(common.ERR, log)
}

func (t *TestCase) Info_log(unfmt string, arg ...any) string {
	log := fmt.Sprintf(unfmt, arg...)
	return t.Logi(common.INFO, log)
}
