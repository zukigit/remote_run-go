package lib

import (
	"golang.org/x/crypto/ssh"
)

func Get_config(user, password string) *ssh.ClientConfig {
	return &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
}

// Connect to the SSH server
func Get_client(host string, config *ssh.ClientConfig) (*ssh.Client, error) {
	return ssh.Dial("tcp", host, config)
}

// Get new SSH session
func Get_session(client *ssh.Client) (*ssh.Session, error) {
	return client.NewSession()
}

// Run command on shell get output
func Get_output(shell_command string, session *ssh.Session) ([]byte, error) {
	return session.Output(shell_command)
}

func Ssh_exec(command string) ([]byte, error) {
	session, err := Get_session(Common_client)
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
