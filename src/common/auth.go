package common

import (
	"golang.org/x/crypto/ssh"
)

// Deprecated: Use Host instead
type Auth struct {
	Username, Password, Hostname string
	Port                         int
}

// Deprecated: Use Host instead
func New_auth(username, password, hostname string, port int, ssh_client *ssh.Client) *Auth {
	return &Auth{
		Username: username,
		Password: password,
		Hostname: hostname,
		Port:     port,
	}
}
