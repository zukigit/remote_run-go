package dao

import (
	"fmt"

	"github.com/zukigit/remote_run-go/src/common"
	"github.com/zukigit/remote_run-go/src/lib"
)

type TestCase struct {
	testcase_no uint
	ticket_no   uint
	description string
	status      *common.Testcase_status
	function    func() common.Testcase_status
}

func New_testcase(testcase_id uint, testcase_description string) *TestCase {
	status := FAILED
	return &TestCase{
		testcase_no: testcase_id,
		description: testcase_description,
		status:      &status,
	}
}

func (t *TestCase) Get_no() uint {
	return t.testcase_no
}

func (t *TestCase) Get_ticket_no() uint {
	return t.ticket_no
}

func (t *TestCase) Set_ticket_no(ticket_no uint) {
	t.ticket_no = ticket_no
}

func (t *TestCase) Get_dsctn() string {
	return t.description
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
