package common

import (
	"golang.org/x/crypto/ssh"
)

// Linux_host represents the configuration and state of a Linux host.
type Linux_host struct {
	Host_name, Host_run_username *string     // The host name and username for running commands.
	Host_ip, Host_dns            *string     `json:"-"` // IP and DNS of the host (excluded from JSON serialization).
	Host_ssh_client              *ssh.Client `json:"-"` // SSH client for connecting to the host (excluded from JSON).
	Host_use_ip                  *bool       `json:"-"` // Whether to use the IP address for connection (excluded from JSON).
	Host_connect_port            *int        // The port to use for SSH connection.
	Host_type                    *Host_type  // The type of host (custom type).
}

// Get_Host_name retrieves the host name.
// Returns an empty string if the field is nil to ensure safe dereferencing.
func (host *Linux_host) Get_Host_name() string {
	if host.Host_name == nil {
		return ""
	}
	return *host.Host_name
}

// Set_Host_name sets the host name directly.
func (host *Linux_host) Set_Host_name(name string) {
	*host.Host_name = name
}

// Get_Host_run_username retrieves the username for running commands on the host.
// Returns an empty string if the field is nil to ensure safe dereferencing.
func (host *Linux_host) Get_Host_run_username() string {
	if host.Host_run_username == nil {
		return ""
	}
	return *host.Host_run_username
}

// Set_Host_run_username sets the username for running commands directly.
func (host *Linux_host) Set_Host_run_username(username string) {
	*host.Host_run_username = username
}

// Get_Host_ip retrieves the IP address of the host.
// Returns an empty string if the field is nil to ensure safe dereferencing.
func (host *Linux_host) Get_Host_ip() string {
	if host.Host_ip == nil {
		return ""
	}
	return *host.Host_ip
}

// Set_Host_ip sets the IP address of the host directly.
func (host *Linux_host) Set_Host_ip(ip string) {
	*host.Host_ip = ip
}

// Get_Host_dns retrieves the DNS of the host.
// Returns an empty string if the field is nil to ensure safe dereferencing.
func (host *Linux_host) Get_Host_dns() string {
	if host.Host_dns == nil {
		return ""
	}
	return *host.Host_dns
}

// Set_Host_dns sets the DNS of the host directly.
func (host *Linux_host) Set_Host_dns(dns string) {
	*host.Host_dns = dns
}

// Get_Host_ssh_client retrieves the SSH client for connecting to the host.
// Returns nil if the SSH client is not set.
func (host *Linux_host) Get_Host_ssh_client() *ssh.Client {
	return host.Host_ssh_client
}

// Set_Host_ssh_client sets the SSH client for connecting to the host directly.
func (host *Linux_host) Set_Host_ssh_client(client *ssh.Client) {
	host.Host_ssh_client = client
}

// Get_Host_use_ip checks whether the IP address is used for connections.
// Returns false if the field is nil to ensure safe dereferencing.
func (host *Linux_host) Get_Host_use_ip() bool {
	if host.Host_use_ip == nil {
		return false
	}
	return *host.Host_use_ip
}

// Set_Host_use_ip sets whether the IP address is used for connections directly.
func (host *Linux_host) Set_Host_use_ip(useIP bool) {
	*host.Host_use_ip = useIP
}

// Get_Host_connect_port retrieves the SSH connection port.
// Returns 0 if the field is nil to ensure safe dereferencing.
func (host *Linux_host) Get_Host_connect_port() int {
	if host.Host_connect_port == nil {
		return 0
	}
	return *host.Host_connect_port
}

// Set_Host_connect_port sets the SSH connection port directly.
func (host *Linux_host) Set_Host_connect_port(port int) {
	*host.Host_connect_port = port
}

// Get_Host_type retrieves the host type.
// Returns nil if the Host_type field is not set.
func (host *Linux_host) Get_Host_type() Host_type {
	return *host.Host_type
}

// Set_Host_type sets the host type directly.
func (host *Linux_host) Set_Host_type(hostType Host_type) {
	*host.Host_type = hostType
}

// Get_Linux_host creates and returns a new instance of Linux_host with default values.
// Fields with pointer types are initialized to avoid nil pointer dereferences.
func New_linux_host() Host {
	default_Host_name := "" // Default empty string for Host_name and Host_run_username.
	default_run_username := ""
	default_Host_ip := ""
	default_Host_dns := ""
	default_Host_use_ip := false       // Default value for Host_use_ip.
	default_Host_connect_port := 22    // Default SSH port.
	default_Host_type := Host_type("") // Default Host_type value (customize if necessary).

	return &Linux_host{
		Host_name:         &default_Host_name,
		Host_run_username: &default_run_username,
		Host_ip:           &default_Host_ip,
		Host_dns:          &default_Host_dns,
		Host_ssh_client:   nil, // Default to nil; SSH client not initialized.
		Host_use_ip:       &default_Host_use_ip,
		Host_connect_port: &default_Host_connect_port,
		Host_type:         &default_Host_type,
	}
}
