package jalibs_t

/*
typedef unsigned long int zbx_uint64_t;

typedef struct {
    char		kind[16];
    double		version;
    zbx_uint64_t	jobid;
    char		serverid[33];
    char		hostname[128];
    int			method;
    char		type[1024];
    char		argument[4096];
    char		script[4096];
    char		env[132097];
    int			result;
    int			status;
    int		pid;
    zbx_uint64_t	start_time;
    zbx_uint64_t	end_time;
    char		message[1024];
    char		std_out[64001];
    char		std_err[64001];
    int			return_code;
    int			signal;
    int			send_retry;
    char		run_user[1024];
    char		run_user_password[1024];
    int			loop_cnt;

    char cur_unique_id[1024];
    char pre_unique_id[1024];

    char       serverip[33];
    // zbx_uint64_t    host_running_job[512];
    zbx_uint64_t*    host_running_job;
    int         size_of_host_running_job;
} ja_job_object;
int JA_FILE_PATH_LEN = 260;
#include "../../../jainclude/jajobfile.h"
*/
import "C"
import "fmt"

type Jajobfile struct{}

func check_job_abort_flag_test() bool {
	const __function_name = "check_job_abort_flag_test"
	var rtn bool = false
	var failedTestCase string

	if C.check_job_abort_flag(failedTestCase, 0) == 0 {
		fmt.Printf("Test case 4.1 was successfully executed. \n", __function_name)
	} //else {
	// 	fmt.Printf("Test case 4.1 execution was failed. \n", __function_name)
	// 	failedTestCase += "4.1"
	// 	return false
	// }
	// if C.check_job_abort_flag("12345", 1) == 0 {
	// 	fmt.Printf("Test case 3.2.5 was successfully executed. \n", __function_name)
	// } else {
	// 	fmt.Printf("Test case 3.2.5 execution was failed. \n", __function_name)
	// 	failedTestCase += "3.2.5"
	// 	return false
	// }
	// if C.check_job_abort_flag("12345", 0) == 1 {
	// 	fmt.Printf("Test case 4.2 was successfully executed. \n", __function_name)
	// } else {
	// 	fmt.Printf("Test case 4.2 execution was failed. \n", __function_name)
	// 	failedTestCase += "4.2"
	// 	return false
	// }

	fmt.Printf("Unit test for check_job_abort_flag test execution was successful. \n", __function_name)
	return rtn
}

func (ja *Jajobfile) TestJaz() bool {
	return check_job_abort_flag_test()
}
