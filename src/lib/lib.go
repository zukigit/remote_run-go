package lib

import "zukigit/remote_run-go/src/dao"

var test_case dao.Test_case

func Set_test_case(tc dao.Test_case) {
	test_case = tc
}
