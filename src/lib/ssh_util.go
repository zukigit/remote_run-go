package lib

import (
	"github.com/zukigit/remote_run-go/src/common"

	"golang.org/x/crypto/ssh"
)

// Get new SSH session
func Get_session(client *ssh.Client) (*ssh.Session, error) {
	return client.NewSession()
}

func Ssh_exec(command string) ([]byte, error) {
	session, err := Get_session(common.Client)
	if err != nil {
		return nil, err
	}
	defer session.Close()

	return session.Output(command)
}

func Ssh_exec_to_str(command string) (string, error) {
	output, err := Ssh_exec(command)
	return string(output), err
}
