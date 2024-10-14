package lib

import (
	"fmt"
	"strconv"
	"time"

	"github.com/zukigit/remote_run-go/src/common"
)

func Jobarg_exec(jobnet_id string) (string, error) {
	cmd := fmt.Sprintf("jobarg_exec -z %s -U Admin -P zabbix -j %s &> /tmp/moon_jobarg_exec_result", common.Login_info.Hostname, jobnet_id)

	_, err := Ssh_exec_to_str(cmd)

	cmd = "cat /tmp/moon_jobarg_exec_result"
	result, err1 := Ssh_exec_to_str(cmd)
	if err1 != nil {
		return result, err1
	}

	if err != nil {
		return result, err
	}

	return Get_res_no(result)
}

func Jobarg_get_JA_JOBNETSTATUS(registry_number string) (string, error) {
	cmd := fmt.Sprintf("eval $(jobarg_get -z %s -U Admin -P zabbix -r %s -e) && echo -n $JA_JOBNETSTATUS", common.Login_info.Hostname, registry_number)
	return Ssh_exec_to_str(cmd)
}

func Jobarg_get_JA_JOBSTATUS(registry_number string) (string, error) {
	cmd := fmt.Sprintf("eval $(jobarg_get -z %s -U Admin -P zabbix -r %s -e) && echo -n $JA_JOBSTATUS", common.Login_info.Hostname, registry_number)
	return Ssh_exec_to_str(cmd)
}

func Jobarg_get_LASTEXITCD(registry_number string) (int64, error) {
	cmd := fmt.Sprintf("eval $(jobarg_get -z %s -U Admin -P zabbix -r %s -e) && echo -n $JA_LASTEXITCD", common.Login_info.Hostname, registry_number)
	exit_cd, err := Ssh_exec_to_str(cmd)
	if err != nil {
		return -1, err
	}

	num, err := strconv.ParseInt(exit_cd, 10, 64) // Base 10, 64-bit integer
	if err != nil {
		return -1, err
	}

	return num, nil
}

func Jobarg_get_LASTSTDOUT(registry_number string) (string, error) {
	cmd := fmt.Sprintf("eval $(jobarg_get -z %s -U Admin -P zabbix -r %s -e) && echo -n $JA_LASTSTDOUT", common.Login_info.Hostname, registry_number)
	return Ssh_exec_to_str(cmd)
}

func Jobarg_get_LASTSTDERR(registry_number string) (string, error) {
	cmd := fmt.Sprintf("eval $(jobarg_get -z %s -U Admin -P zabbix -r %s -e) && echo -n $JA_LASTSTDERR", common.Login_info.Hostname, registry_number)
	return Ssh_exec_to_str(cmd)
}

// Jobarg_get_jobnet_run_info waits util the jobnet is done or get error and returns Jobnet run info.
func Jobarg_get_jobnet_run_info(registry_number string) (*common.Jobnet_run_info, error) {
	var jobnet_status, job_status, std_out, std_error string
	var err error
	var index int
	var exit_cd int64

	for {
		jobnet_status, err = Jobarg_get_JA_JOBNETSTATUS(registry_number)
		if err != nil {
			return nil, err
		}

		job_status, err = Jobarg_get_JA_JOBSTATUS(registry_number)
		if err != nil {
			return nil, err
		}

		if jobnet_status == "END" || (jobnet_status == "RUN" && job_status == "ERROR") {
			break
		}
		Spinner_log(index, Formatted_log(common.INFO, "Getting jobnet[%s] run info but jobnet is not finished yet", registry_number))
		time.Sleep(1 * time.Second)
		index++
	}

	exit_cd, err = Jobarg_get_LASTEXITCD(registry_number)
	if err != nil {
		return nil, err
	}

	std_out, err = Jobarg_get_LASTSTDOUT(registry_number)
	if err != nil {
		return nil, err
	}

	std_error, err = Jobarg_get_LASTSTDERR(registry_number)
	if err != nil {
		return nil, err
	}

	fmt.Println()
	return common.New_jobnet_run_info(jobnet_status, job_status, std_out, std_error, exit_cd), nil
}

// Jobarg_get_jobnet_run_info waits util the jobnet is done or get error and returns Jobnet run info.
func Jobarg_get_jobnet_info(registry_number string, targetJobnetStatus string, targetJobStatus string) (*common.Jobnet_run_info, error) {
	var jobnet_status, job_status, std_out, std_error string
	var err error
	var index int
	var exit_cd int64

	for {
		jobnet_status, err = Jobarg_get_JA_JOBNETSTATUS(registry_number)
		if err != nil {
			Formatted_log(common.INFO, "Error:%s", err.Error())
		}

		job_status, err = Jobarg_get_JA_JOBSTATUS(registry_number)
		if err != nil {
			Formatted_log(common.INFO, "Error:%s", err.Error())
		}

		if jobnet_status == targetJobnetStatus && job_status == targetJobStatus {
			break
		}
		Spinner_log(index, Formatted_log(common.INFO, "Getting jobnet[%s] run info but jobnet is not finished yet", registry_number))
		time.Sleep(1 * time.Second)
		index++
	}

	exit_cd, err = Jobarg_get_LASTEXITCD(registry_number)
	if err != nil {
		Formatted_log(common.INFO, "Error:%s", err.Error())
	}

	std_out, err = Jobarg_get_LASTSTDOUT(registry_number)
	if err != nil {
		Formatted_log(common.INFO, "Error:%s", err.Error())
	}

	std_error, err = Jobarg_get_LASTSTDERR(registry_number)
	if err != nil {
		Formatted_log(common.INFO, "Error:%s", err.Error())
	}

	fmt.Println()
	return common.New_jobnet_run_info(jobnet_status, job_status, std_out, std_error, exit_cd), nil
}

// To use this function, your jobarranger agent's TmpDir must be default(TmpDir=/var/lib/jobarranger/tmp and TmpDir=C:\Program Files\Job Arranger\Job Arranger Agent\temp).
//
// Jobarg_cleanup() cleans jobarg-server and jobarg-agentd(linux and windows) data.
// Since this is testcase utility funtion, you must use it in testcase function.
func Jobarg_cleanup() error {
	if err := Stop_jaz_server(); err != nil {
		return fmt.Errorf("failed to stop JAZ server: %w", err)
	}
	if err := Stop_jaz_agent_windows(); err != nil {
		return fmt.Errorf("failed to stop JAZ agent: %w", err)
	}
	if err := Stop_jaz_agent_linux(); err != nil {
		return fmt.Errorf("failed to stop JAZ agent: %w", err)
	}
	if _, err := DBexec("delete from ja_run_jobnet_table;"); err != nil {
		return fmt.Errorf("failed to execute DB command: %w", err)
	}
	if err := Cleanup_agent_windows(); err != nil {
		return fmt.Errorf("failed to clean up agent: %w", err)
	}
	if err := Cleanup_agent_linux(); err != nil {
		return fmt.Errorf("failed to clean up agent: %w", err)
	}
	if err := Restart_jaz_server(); err != nil {
		return fmt.Errorf("failed to stop JAZ server: %w", err)
	}
	if err := Restart_jaz_agent_windows(); err != nil {
		return fmt.Errorf("failed to stop JAZ server: %w", err)
	}

	return nil
}
