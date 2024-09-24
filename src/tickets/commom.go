package tickets

import "zukigit/remote_run-go/src/dao"

const (
	PASSED     dao.Testcase_status = "PASSED"
	FAILED     dao.Testcase_status = "FAILED"
	MUST_CHECK dao.Testcase_status = "MUST_CHECK"

	END    string = "END"
	NORMAL string = "NORMAL"
)
