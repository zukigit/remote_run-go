package common

import "golang.org/x/crypto/ssh"

// Host defines the interface for host operations.
type Host interface {
	// Getters
	Get_Host_name() string
	Get_Host_run_username() string
	Get_Host_ip() string
	Get_Host_dns() string
	Get_Host_ssh_client() *ssh.Client
	Get_Host_use_ip() bool
	Get_Host_connect_port() int
	Get_Host_type() Host_type

	// Setters
	Set_Host_name(name string)
	Set_Host_run_username(username string)
	Set_Host_ip(ip string)
	Set_Host_dns(dns string)
	Set_Host_ssh_client(client *ssh.Client)
	Set_Host_use_ip(useIP bool)
	Set_Host_connect_port(port int)
	Set_Host_type(hostType Host_type)

	// Host utility functions
	Register(public_key string) error       // Sets the public key for the specified host.
	Run_cmd(cmd string) ([]byte, error)     // Run cmd on host and get output as byte.
	Run_cmd_str(cmd string) (string, error) // Run cmd on host and getoutput as string.
}

type Host_struct struct {
	Host_name, Host_run_username *string     // The host name and username for running commands.
	Host_ip, Host_dns            *string     `json:"-"` // IP and DNS of the host (excluded from JSON).
	Host_ssh_client              *ssh.Client `json:"-"` // SSH client for connecting to the host (excluded from JSON).
	Host_use_ip                  *bool       `json:"-"` // Whether to use the IP address for connection (excluded from JSON).
	Host_connect_port            *int        // The port to use for SSH connection.
	Host_type                    *Host_type  // The type of host (custom type).
}
