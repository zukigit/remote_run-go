package dao

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
