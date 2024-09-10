package lib

import (
	"fmt"
	"os"

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
func Get_client(host string, config *ssh.ClientConfig) *ssh.Client {
	fmt.Print("Connecting...")
	client, err := ssh.Dial("tcp", host, config)
	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}
	fmt.Println("Connected!")
	return client
}

func Get_session(client *ssh.Client) *ssh.Session {
	session, err := client.NewSession()
	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}
	return session
}

func Get_output(shell_command string, session *ssh.Session) []byte {
	output, err := session.Output(shell_command)
	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}
	return output
}
