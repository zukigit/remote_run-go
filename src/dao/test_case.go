package dao

import "golang.org/x/crypto/ssh"

type Test_case interface {
	Run() bool
	Set_tc_values(session *ssh.Session)
	Get_is_passed() bool
	Set_is_passed(is_passed bool)
	Get_tc_id() string
	Get_tc_dsctn() string
	Set_tc_log(tc_log string)
	Get_tc_log() string
}
