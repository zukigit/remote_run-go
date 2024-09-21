package lib

import (
	"fmt"

	"golang.org/x/crypto/ssh"
)

func Ssh_exec(command string, session *ssh.Session) ([]byte, error) {
	return session.Output(command)
}

func Ssh_exec_to_str(command string, session *ssh.Session) (string, error) {
	output, err := Ssh_exec(command, session)
	fmt.Println("string(output)", string(output))

	return string(output), err
}
