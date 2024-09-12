package lib

import (
	"fmt"
	"zukigit/remote_run-go/src/dao"
)

func prepare_data(t dao.Test_case) {
	data := "id: " + t.Get_tc_id() + "\n"
	data += "description: " + t.Get_tc_dsctn() + "\n"
	data += "status: "
	if t.Get_is_passed() {
		data += "passed\n"
	} else {
		data += "failed\n"
	}
	data += "log:\n\t" + t.Get_tc_log()

	t.Set_tc_log(data)
}

func Write_tc_log(t dao.Test_case) {
	prepare_data(t)
	fmt.Println(t.Get_tc_log())
}
