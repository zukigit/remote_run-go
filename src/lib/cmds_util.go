package lib

import "golang.org/x/crypto/ssh"

func Ssh_exec(command string, session *ssh.Session) ([]byte, error) {
	return session.Output(command)
}

func Ssh_exec_to_str(command string, session *ssh.Session) (string, error) {
	output, err := Ssh_exec(command, session)
	return string(output), err
}

// func Restart_jaz_server() ([]byte, error) {
// 	return Ssh_exec("systemctl restart jobarg-server")
// }

// func Restart_docker() ([]byte, error) {
// 	return Ssh_exec("systemctl restart docker")
// }

// func Stop_docker() ([]byte, error) {
// 	return Ssh_exec("systemctl stop docker")
// }

// func Stop_jaz_server() ([]byte, error) {
// 	return Ssh_exec("systemctl stop jobarg-server")
// }

// func Start_jaz_server() ([]byte, error) {
// 	return Ssh_exec("systemctl start jobarg-server")
// }
