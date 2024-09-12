package dao

import "golang.org/x/crypto/ssh"

type Test_case interface {
	Run() bool
	Set_tc_values(session *ssh.Session)
	Get_tc_id() string
	Get_tc_dsctn() string
	Get_tc_log() string
}
