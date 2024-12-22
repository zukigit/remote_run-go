package common

import "golang.org/x/crypto/ssh"

// Host defines the interface for host operations.
// Any struct implementing this interface must define these methods.
type Host interface {
	// Getters
	Get_Host_name() string
	Get_Host_run_username() string
	Get_Host_ip() string
	Get_Host_dns() string
	Get_Host_ssh_client() *ssh.Client
	Get_Host_use_ip() bool
	Get_Host_connect_port() int
	Get_Host_type() *Host_type

	// Setters
	Set_Host_name(name string)
	Set_Host_run_username(username string)
	Set_Host_ip(ip string)
	Set_Host_dns(dns string)
	Set_Host_ssh_client(client *ssh.Client)
	Set_Host_use_ip(useIP bool)
	Set_Host_connect_port(port int)
	Set_Host_type(hostType Host_type)
}
