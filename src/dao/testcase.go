package dao

import (
	"fmt"
	"time"
	"zukigit/remote_run-go/src/lib"
)

const INFO = 1
const ERR = 2

type TestCase struct {
	id          uint
	description string
	logs        *[]string
	is_passed   *bool
	function    func() bool
	auth        *Auth
}

func New_testcase(testcase_id uint, testcase_description string, auth *Auth) *TestCase {
	status := false
	logs := []string{}
	return &TestCase{
		id:          testcase_id,
		description: testcase_description,
		auth:        auth,
		is_passed:   &status,
		logs:        &logs,
	}
}

func (t *TestCase) Set_log(log string) {
	*t.logs = append(*t.logs, log)
}

func (t *TestCase) Get_id() uint {
	return t.id
}

func (t *TestCase) Get_auth() Auth {
	return *t.auth
}

func (t *TestCase) Get_dsctn() string {
	return t.description
}

func (t *TestCase) Get_logs() []string {
	return *t.logs
}

func (t *TestCase) Set_is_passed(is_passed bool) {
	*t.is_passed = is_passed
}

func (t *TestCase) Get_is_passed() bool {
	return *t.is_passed
}

func (t *TestCase) Set_function(function func() bool) {
	t.function = function
}

func (t *TestCase) Run_function() bool {
	return t.function()
}

func (t *TestCase) Is_function_nil() bool {
	return t.function == nil
}

// From here is test case util functions

func (t *TestCase) Logi(level int, log string) string {
	log = fmt.Sprintf("Test_case no[%d], %s", t.Get_id(), log)
	log = lib.Formatted_log(level, log)
	t.Set_log(log)

	return log
}

func (t *TestCase) Err_log(unfmt string, arg ...any) string {
	log := fmt.Sprintf(unfmt, arg...)
	return t.Logi(ERR, log)
}

func (t *TestCase) Info_log(unfmt string, arg ...any) string {
	log := fmt.Sprintf(unfmt, arg...)
	return t.Logi(INFO, log)
}

func (t *TestCase) Ssh_exec(command string) ([]byte, error) {
	return lib.Ssh_exec(command, t.auth.Ssh_client)
}

func (t *TestCase) Ssh_exec_to_str(command string) (string, error) {
	return lib.Ssh_exec_to_str(command, t.auth.Ssh_client)
}

func (t *TestCase) Jobarg_exec(jobid string) (string, error) {
	cmd := fmt.Sprintf("jobarg_exec -z %s -U Admin -P zabbix -j %s &> /tmp/moon_jobarg_exec_result", t.Get_auth().Hostname, jobid)

	_, err := lib.Ssh_exec_to_str(cmd, t.auth.Ssh_client)
	if err != nil {
		return "", err
	}

	cmd = "cat /tmp/moon_jobarg_exec_result"
	result, err := lib.Ssh_exec_to_str(cmd, t.auth.Ssh_client)
	if err != nil {
		return "", err
	}

	return lib.Get_res_no(result)
}

func (t *TestCase) Jobarg_get_JA_JOBNETSTATUS(registry_number string) (string, error) {
	cmd := fmt.Sprintf("eval $(jobarg_get -z %s -U Admin -P zabbix -r %s -e) && echo -n $JA_JOBNETSTATUS", t.Get_auth().Hostname, registry_number)
	return lib.Ssh_exec_to_str(cmd, t.auth.Ssh_client)
}

func (t *TestCase) Jobarg_get_JA_JOBSTATUS(registry_number string) (string, error) {
	cmd := fmt.Sprintf("eval $(jobarg_get -z %s -U Admin -P zabbix -r %s -e) && echo -n $JA_JOBSTATUS", t.Get_auth().Hostname, registry_number)
	return lib.Ssh_exec_to_str(cmd, t.auth.Ssh_client)
}

func (t *TestCase) Jobarg_get_LASTEXITCD(registry_number string) (string, error) {
	cmd := fmt.Sprintf("eval $(jobarg_get -z %s -U Admin -P zabbix -r %s -e) && echo -n $LASTEXITCD", t.Get_auth().Hostname, registry_number)
	return lib.Ssh_exec_to_str(cmd, t.auth.Ssh_client)
}

func (t *TestCase) Jobarg_get_LASTSTDOUT(registry_number string) (string, error) {
	cmd := fmt.Sprintf("eval $(jobarg_get -z %s -U Admin -P zabbix -r %s -e) && echo -n $LASTSTDOUT", t.Get_auth().Hostname, registry_number)
	return lib.Ssh_exec_to_str(cmd, t.auth.Ssh_client)
}

func (t *TestCase) Jobarg_get_LASTSTDERR(registry_number string) (string, error) {
	cmd := fmt.Sprintf("eval $(jobarg_get -z %s -U Admin -P zabbix -r %s -e) && echo -n $LASTSTDERR", t.Get_auth().Hostname, registry_number)
	return lib.Ssh_exec_to_str(cmd, t.auth.Ssh_client)
}

// Jobarg_get_jobnet_run_info waits util the jobnet is done or get error and returns Jobnet run info.
func (t *TestCase) Jobarg_get_jobnet_run_info(registry_number string) (*Jobnet, error) {
	var status, job_status, exit_cd, std_out, std_error string
	var err error
	var index int

	for {
		status, err = t.Jobarg_get_JA_JOBNETSTATUS(registry_number)
		if err != nil {
			lib.Formatted_log(INFO, "Error:%s", err.Error())
		}

		job_status, err = t.Jobarg_get_JA_JOBSTATUS(registry_number)
		if err != nil {
			lib.Formatted_log(INFO, "Error:%s", err.Error())
		}

		if status == "END" || (status == "RUN" && job_status == "ERROR") {
			break
		}
		lib.Spinner_log(index, lib.Formatted_log(INFO, "Getting jobnet run info but jobnet is not finished yet"))
		time.Sleep(1 * time.Second)
		index++
	}

	fmt.Println()
	return New_Jobnet(status, job_status, exit_cd, std_out, std_error), nil
}
