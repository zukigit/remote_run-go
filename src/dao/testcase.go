package dao

type testCase struct {
	id          uint
	description string
	logs        []string
	is_passed   bool
}

func New_TestCase(testcase_id uint, testcase_description string) *testCase {
	return &testCase{
		id:          testcase_id,
		description: testcase_description,
	}
}

func (t *testCase) Add_log(log string) {
	t.logs = append(t.logs, log)
}

func (t *testCase) Get_log() []string {
	return t.logs
}

func (t *testCase) Succeed() {
	t.is_passed = true
}

func (t *testCase) Fail() {
	t.is_passed = false
}
