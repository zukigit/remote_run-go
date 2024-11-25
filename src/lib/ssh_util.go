package lib

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"time"

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

func GetSSHClientWithKey(hostIP string, port int, username string, keyfilepath string) *ssh.Client {
	// Load the private key
	key, err := os.ReadFile(keyfilepath)
	if err != nil {
		fmt.Println("unable to read private key, Error:", err.Error())
		fmt.Println("use the following command to generate private key, 'ssh-keygen -t rsa -b 4096'")
		os.Exit(1)
	}

	// Create the signer for the private key
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		fmt.Println("unable to parse private key Error:", err.Error())
		os.Exit(1)
	}

	clientConfig := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	address := fmt.Sprintf("%s:%d", hostIP, port)

	client, err := ssh.Dial("tcp", address, clientConfig)
	if err != nil {
		fmt.Println("Failed in getting ssh client, Error:", err.Error())
		os.Exit(1)
	}

	return client
}

func GetSSHClient(hostIP string, port int, username string, password string) *ssh.Client {
	clientConfig := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	address := fmt.Sprintf("%s:%d", hostIP, port)

	client, err := ssh.Dial("tcp", address, clientConfig)
	if err != nil {
		fmt.Println("Failed in getting ssh client, Error:", err.Error())
		os.Exit(1)
	}

	return client
}

// CheckSSHforRebootingAfterDelay checks if the system is up using SSH connection after waiting for a specified timeAmount.
func CheckSSHforRebootingAfterDelay(client *ssh.Client, timeAmount int) bool {
	// Wait for the specified timeAmount (in seconds) before checking SSH
	time.Sleep(time.Duration(timeAmount) * time.Second)

	// Try to create an SSH session
	session, err := client.NewSession()
	if err != nil {
		// If session creation fails, assume the system might be rebooting
		return true
	}

	// If session creation succeeds, close the session and return false
	defer session.Close()
	return false
}

func ConnectWithRetry(hostIP string, port int, username string, password string, maxRetries int) *ssh.Client {
	var client *ssh.Client

	for attempts := 0; attempts < maxRetries; attempts++ {
		client = GetSSHClient(hostIP, port, username, password)
		if client != nil {
			fmt.Println("SSH client connected successfully.")
			return client // Return on successful connection
		}

		// Log the retry attempt
		fmt.Printf("Retrying to establish SSH connection... Attempt %d of %d\n", attempts+1, maxRetries)
		// Sleep for 1 second before retrying
		time.Sleep(1 * time.Second)
	}

	// If we exhaust all attempts without success
	fmt.Println("Max retry limit reached. Exiting.")
	os.Exit(1)
	return nil // Not reached, just for clarity
}

func ExecuteSSHCommand(client *ssh.Client, command string) ([]byte, error) {
	session, err := client.NewSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()

	return session.Output(command)
}

func GetOutputStrFromSSHCommand(client *ssh.Client, command string) (string, error) {
	output, err := ExecuteSSHCommand(client, command)
	return string(output), err
}

func Set_common_client(username, passwd, hostname string, port int) {
	common.Client = GetSSHClient(hostname, port, username, passwd)
}

func Set_host_pool() {
	common.Host_pool = Get_hosts_from_jsonfile("hosts.json")

	current_user, err := user.Current()
	if err != nil {
		fmt.Printf("failed in getting run user, Error: %v", err)
		os.Exit(1)
	}
	ssh_private_filepath := filepath.Join(current_user.HomeDir, ".ssh", "id_rsa")

	for i := range *common.Host_pool {
		host := &(*common.Host_pool)[i] // Get a pointer to the actual host
		if host.Host_use_ip {
			host.Host_client = GetSSHClientWithKey(host.Host_ip, host.Host_port, host.Host_run_username, ssh_private_filepath)
		} else {
			host.Host_client = GetSSHClientWithKey(host.Host_dns, host.Host_port, host.Host_run_username, ssh_private_filepath)
		}
	}
}
