package lib

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
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

func Generate_sshkeys(key_filepath string) {
	private_key_filepath := filepath.Join(key_filepath, "id_rsa")
	public_key_filepath := filepath.Join(key_filepath, "id_rsa.pub")

	private_key, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		fmt.Println("unable to generate private key, Error:", err.Error())
		os.Exit(1)
	}

	private_key_pem := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(private_key),
	})

	private_key_file, err := os.Create(private_key_filepath)
	if err != nil {
		fmt.Println("unable to create private key file, Error:", err.Error())
		os.Exit(1)
	}
	defer private_key_file.Close()

	_, err = private_key_file.Write(private_key_pem)
	if err != nil {
		fmt.Println("unable to write private_key_pem, Error:", err.Error())
		os.Exit(1)
	}

	public_key, err := ssh.NewPublicKey(&private_key.PublicKey)
	if err != nil {
		fmt.Println("unable to get public key, Error:", err.Error())
		os.Exit(1)
	}

	public_key_bytes := ssh.MarshalAuthorizedKey(public_key)

	public_key_file, err := os.Create(public_key_filepath)
	if err != nil {
		fmt.Println("unable to create public key file, Error:", err.Error())
		os.Exit(1)
	}
	defer public_key_file.Close()

	_, err = public_key_file.Write(public_key_bytes)
	if err != nil {
		fmt.Println("unable to write public_key_pem, Error:", err.Error())
		os.Exit(1)
	}
}

func GetSSHClientWithKey(hostIP string, port int, username string, keyfilepath string) (*ssh.Client, error) {
	private_key_filepath := filepath.Join(keyfilepath, "id_rsa")

	// Load the private key
	key, err := os.ReadFile(private_key_filepath)
	if err != nil {
		return nil, fmt.Errorf("unable to read private key, Error: %s", err.Error())
	}

	// Create the signer for the private key
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, fmt.Errorf("unable to parse private key, Error: %s", err.Error())
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
		return nil, fmt.Errorf("failed in getting ssh client, Error: %s", err.Error())
	}

	return client, err
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

func update_hosts_ips() {
	var host_dns, host_ip string
	var host_use_ip bool

	for i := range common.Host_pool {
		host := common.Host_pool[i]

		rows, err := GetData(`select i.useip, i.dns, i.ip from hosts h, interface i
							where h.hostid = i.hostid and i.main = 1 and i.type = 1 and h.host = $1`, host.Get_Host_name())
		if err != nil {
			fmt.Println("Failed in quering hosts, Error:", err.Error())
			os.Exit(1)
		}
		defer rows.Close()

		i := 0
		for rows.Next() {
			if err := rows.Scan(&host_use_ip, &host_dns, &host_ip); err != nil {
				fmt.Println("Failed in scanning hosts, Error:", err.Error())
				os.Exit(1)
			}

			host.Set_Host_use_ip(host_use_ip)
			host.Set_Host_dns(host_dns)
			host.Set_Host_ip(host_ip)

			i++
		}
	}
}

func Set_host_pool() {
	Get_hosts_from_jsonfile("hosts.json")

	if len(common.Host_pool) <= 0 {
		fmt.Println("error: no hosts to run, use 'register_hosts' command to register.")
		os.Exit(1)
	}

	update_hosts_ips()

	current_user, err := user.Current()
	if err != nil {
		fmt.Printf("failed in getting run user, Error: %s\n", err.Error())
		os.Exit(1)
	}
	ssh_key_filepath := filepath.Join(current_user.HomeDir, ".ssh")

	for i := range common.Host_pool {
		host := common.Host_pool[i]

		if host.Get_Host_use_ip() {
			fmt.Printf("Getting client from %s:%d ...", host.Get_Host_ip(), host.Get_Host_connect_port())

			client, err := GetSSHClientWithKey(host.Get_Host_ip(), host.Get_Host_connect_port(), host.Get_Host_run_username(), ssh_key_filepath)
			if err != nil {
				fmt.Printf("GetSSHClientWithKey failed, Error: %s\n", err.Error())
				os.Exit(1)
			}

			host.Set_Host_ssh_client(client)
		} else {
			fmt.Printf("Getting client from %s:%d ...", host.Get_Host_dns(), host.Get_Host_connect_port())

			client, err := GetSSHClientWithKey(host.Get_Host_dns(), host.Get_Host_connect_port(), host.Get_Host_run_username(), ssh_key_filepath)
			if err != nil {
				fmt.Printf("GetSSHClientWithKey failed, Error: %s\n", err.Error())
				os.Exit(1)
			}

			host.Set_Host_ip(host.Get_Host_dns())
			host.Set_Host_ssh_client(client)
		}
		fmt.Println(" done.")
	}

	common.Set_linux_server_host()

	// to delete later
	if common.Server_host == nil {
		fmt.Println("error: no server host to run, use 'register_hosts' command to register.")
		os.Exit(1)
	} else {
		common.Client = common.Server_host.Get_Host_ssh_client()
		common.Login_info.Hostname = common.Server_host.Get_Host_ip()
		common.Login_info.Username = common.Server_host.Get_Host_run_username()
		common.Login_info.Port = common.Server_host.Get_Host_connect_port()
	}
}
