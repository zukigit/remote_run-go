package lib

import (
	"fmt"
)

func Jobarg_enable_jobnet(jobnet_id string, jobnet_name string) error {
	_, err := DBexec("update ja_jobnet_control_table set valid_flag = 0 where jobnet_id = '%s' and valid_flag = 1;", jobnet_id)
	if err != nil {
		return err
	}
	res, err := DBexec("update ja_jobnet_control_table set valid_flag = 1 where jobnet_id = '%s' and jobnet_name = '%s'", jobnet_id, jobnet_name)
	if err != nil {
		return err
	}

	affected_rows, err := res.RowsAffected()
	if err != nil {
		return err
	} else if affected_rows > 1 {
		DBexec("update ja_jobnet_control_table set valid_flag = 0 where jobnet_id = '%s' and valid_flag = 1;", jobnet_id)
		return fmt.Errorf("this function does not supprt duplicated jobnet's version. jobnet_id: %s, jobnet_name: %s", jobnet_id, jobnet_name)
	}

	return nil
}

func Enable_common_jobnets() {
	if err := Jobarg_enable_jobnet("Icon_1", "jobicon_linux"); err != nil {
		fmt.Println("Failed to enable common jobnets, error: ", err.Error())
	}

	if err := Jobarg_enable_jobnet("Icon_2", "Icon_2"); err != nil {
		fmt.Println("Failed to enable common jobnets, error: ", err.Error())
	}

	if err := Jobarg_enable_jobnet("Icon_10", "Icon_10"); err != nil {
		fmt.Println("Failed to enable common jobnets, error: ", err.Error())
	}

	if err := Jobarg_enable_jobnet("Icon_100", "Icon_100"); err != nil {
		fmt.Println("Failed to enable common jobnets, error: ", err.Error())
	}

	if err := Jobarg_enable_jobnet("Icon_200", "Icon_200"); err != nil {
		fmt.Println("Failed to enable common jobnets, error: ", err.Error())
	}

	if err := Jobarg_enable_jobnet("Icon_400", "Icon_400"); err != nil {
		fmt.Println("Failed to enable common jobnets, error: ", err.Error())
	}

	if err := Jobarg_enable_jobnet("Icon_500", "Icon_500"); err != nil {
		fmt.Println("Failed to enable common jobnets, error: ", err.Error())
	}

	if err := Jobarg_enable_jobnet("Icon_510", "Icon_510"); err != nil {
		fmt.Println("Failed to enable common jobnets, error: ", err.Error())
	}

	if err := Jobarg_enable_jobnet("Icon_800", "Icon_800"); err != nil {
		fmt.Println("Failed to enable common jobnets, error: ", err.Error())
	}

	if err := Jobarg_enable_jobnet("Icon_1000", "Icon_1000"); err != nil {
		fmt.Println("Failed to enable common jobnets, error: ", err.Error())
	}

	if err := Jobarg_enable_jobnet("Icon_1020", "Icon_1020"); err != nil {
		fmt.Println("Failed to enable common jobnets, error: ", err.Error())
	}

	if err := Jobarg_enable_jobnet("Icon_2040", "Icon_2040"); err != nil {
		fmt.Println("Failed to enable common jobnets, error: ", err.Error())
	}

	if err := Jobarg_enable_jobnet("Icon_3000", "Icon_3000"); err != nil {
		fmt.Println("Failed to enable common jobnets, error: ", err.Error())
	}

}
