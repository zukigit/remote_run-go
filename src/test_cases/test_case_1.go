package testcases

import "fmt"

type Test_case_1 struct {
	test_case_description string
	jobnet_id             string
}

func (T *Test_case_1) get_test_case() {
	T.test_case_description = "Test 1 is goood lol"
	T.jobnet_id = "JOBNET-1"
}

func (T *Test_case_1) Run() {
	fmt.Println("Ran", T)
}
