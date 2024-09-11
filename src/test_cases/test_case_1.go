package testcases

type Test_case_1 struct {
	test_case_description string
	jobnet_id             string
}

func (T *Test_case_1) Run() {

}

func get_test_case() *Test_case {
	return &Test_case_1{
		test_case_description: "Test 1 is goood lol",
	}
}
