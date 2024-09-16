package dao

type testCase struct {
	Id               uint
	Description, Log string
	IsPassed         bool
}

func New_TestCasep(testcase_id uint, testcase_description string) *testCase {
	return &testCase{
		Id:          testcase_id,
		Description: testcase_description,
	}
}
