package jalibs_t

/*
#include "../../../jainclude/jalockutil.h"
*/
import "C"
import "fmt"

type Jalockutil_ja_test struct{}

func init_session_dbc_locks_test() bool {
	const __function_name = "init_session_dbc_locks_test"
	var rtn bool = false

	rtn = int(C.init_session_dbc_locks()) == 1
	if !rtn {
		fmt.Printf("%s() got failed!\n", __function_name)
	}

	return rtn
}

func get_jaz_folder_path() bool {
	return true
}

func (ja *Jalockutil_ja_test) TestJaz() bool {
	return init_session_dbc_locks_test() && get_jaz_folder_path()
}
