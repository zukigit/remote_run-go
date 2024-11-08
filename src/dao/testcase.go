package dao

import (
	"fmt"

	"github.com/zukigit/remote_run-go/src/common"
	"github.com/zukigit/remote_run-go/src/lib"
)

type TestCase struct {
	Testcase_no          uint
	Testcase_description string
	Testcase_status      *common.Testcase_status
	function             func() common.Testcase_status
	ticket_no            uint
}

func New_testcase(testcase_id uint, testcase_description string) *TestCase {
	status := FAILED
	return &TestCase{
		Testcase_no:          testcase_id,
		Testcase_description: testcase_description,
		Testcase_status:      &status,
	}
}

func (t *TestCase) Get_no() uint {
	return t.Testcase_no
}

func (t *TestCase) Get_ticket_no() uint {
	return t.ticket_no
}

func (t *TestCase) Set_ticket_no(ticket_no uint) {
	t.ticket_no = ticket_no
}

func (t *TestCase) Get_dsctn() string {
	return t.Testcase_description
}

func (t *TestCase) Set_status(status common.Testcase_status) {
	*t.Testcase_status = status
}

func (t *TestCase) Get_status() common.Testcase_status {
	return *t.Testcase_status
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
	log = fmt.Sprintf("[%d] [%d] %s", t.Get_ticket_no(), t.Get_no(), log)
	log = lib.Formatted_log(level, log)

	common.Sugar.Infof(log)

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
