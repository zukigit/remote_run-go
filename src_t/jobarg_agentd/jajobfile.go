package jalibs_t

/*
#include "../../../jainclude/jajobfile.h"
*/
import "C"
import "fmt"

type Jajobfile struct{}

func check_job_abort_flag_test() bool {
	const __function_name = "check_job_abort_flag_test"
	var rtn bool = false
	var failedTestCase string

	if C.check_job_abort_flag("12345", 0) == 0 {
		fmt.Printf("Test case 4.1 was successfully executed. \n", __function_name)
	} else {
		fmt.Printf("Test case 4.1 execution was failed. \n", __function_name)
		failedTestCase += "4.1"
		return false
	}
	if C.check_job_abort_flag("12345", 1) == 0 {
		fmt.Printf("Test case 3.2.5 was successfully executed. \n", __function_name)
	} else {
		fmt.Printf("Test case 3.2.5 execution was failed. \n", __function_name)
		failedTestCase += "3.2.5"
		return false
	}
	if C.check_job_abort_flag("12345", 0) == 1 {
		fmt.Printf("Test case 4.2 was successfully executed. \n", __function_name)
	} else {
		fmt.Printf("Test case 4.2 execution was failed. \n", __function_name)
		failedTestCase += "4.2"
		return false
	}

	fmt.Printf("Unit test for check_job_abort_flag test execution was successful. \n", __function_name)
	return rtn
}

func (ja *Jajobfile) TestJaz() bool {
	return check_job_abort_flag_test()
}
