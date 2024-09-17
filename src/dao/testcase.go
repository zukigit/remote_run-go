package dao

import (
	"fmt"
	"time"
)

const INFO = 1
const ERR = 2

type TestCase struct {
	id          uint
	description string
	logs        []string
	is_passed   bool
	function    func() bool
}

func New_TestCase(testcase_id uint, testcase_description string) *TestCase {
	return &TestCase{
		id:          testcase_id,
		description: testcase_description,
	}
}

func (t *TestCase) Add_log(log string) {
	t.logs = append(t.logs, log)
}

func (t *TestCase) Get_id() uint {
	return t.id
}

func (t *TestCase) Get_log() []string {
	return t.logs
}

func (t *TestCase) Set_is_passed(is_passed bool) {
	t.is_passed = is_passed
}

func (t *TestCase) Add_function(function func() bool) {
	t.function = function
}

func (t *TestCase) Run_function() bool {
	return t.function()
}

func (t *TestCase) Is_function_nil() bool {
	return t.function == nil
}

func (t *TestCase) Logi(level int, log string) string {
	currentTime := time.Now()
	formattedTime := currentTime.Format("2001-02-10 15:04:05")

	switch level {
	case INFO:
		log = formattedTime + ", [INFO] " + log
		t.Add_log(log)
	case ERR:
		log = formattedTime + ", [ERR] " + log
		t.Add_log(formattedTime + ", [ERR] " + log)
	}

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
