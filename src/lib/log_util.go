package lib

import "time"

func logi(log string) {
	currentTime := time.Now()
	test_case.Set_tc_log(currentTime.String() + " " + log)
}

func Err_log(log string) {
	logi("Err: " + log)
}

func Info_log(log string) {
	logi("Info: " + log)
}
