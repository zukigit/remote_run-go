package lib

import (
	"golang.org/x/crypto/ssh"
)

func Ssh_exec(command string, client *ssh.Client) ([]byte, error) {
	session, err := Get_session(client)
	if err != nil {
		return nil, err
	}
	defer session.Close()

	return session.Output(command)
}

func Ssh_exec_to_str(command string, client *ssh.Client) (string, error) {
	output, err := Ssh_exec(command, client)
	return string(output), err
}
