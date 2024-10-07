package common

import (
	"golang.org/x/crypto/ssh"
)

type Auth struct {
	Username, Password, Hostname string
	Port                         int
}

func New_auth(username, password, hostname string, port int, ssh_client *ssh.Client) *Auth {
	return &Auth{
		Username: username,
		Password: password,
		Hostname: hostname,
		Port:     port,
	}
}
