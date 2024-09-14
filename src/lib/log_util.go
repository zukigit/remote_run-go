package lib

import (
	"fmt"
	"time"
)

func logi(log string) {
	currentTime := time.Now()
	formattedTime := currentTime.Format("2006-01-02 15:04:05")

	tem_log := test_case.Get_tc_log() + "\t"
	test_case.Set_tc_log(tem_log + formattedTime + " " + log)
}

func Err_log(unfmt string, arg ...any) {
	log := fmt.Sprintf(unfmt, arg...)
	logi("Err: " + log)
}

func Info_log(unfmt string, arg ...any) {
	log := fmt.Sprintf(unfmt, arg...)
	logi("Info: " + log)
}
