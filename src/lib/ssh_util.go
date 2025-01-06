package lib

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
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

// This is new function of GetSSHClient that does not exit on error.
func GetSSHClient_(hostIP string, port int, username string, password string) (*ssh.Client, error) {
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

func Set_host_pool(jsonfilepath string) {
	var temp_hosts []common.Host_struct
	var host common.Host
	var host_dns, host_ip string
	var host_use_ip bool

	// Get ssh key filepath
	current_user, err := user.Current()
	if err != nil {
		fmt.Printf("failed in getting run user, Error: %s\n", err.Error())
		os.Exit(1)
	}
	ssh_key_filepath := filepath.Join(current_user.HomeDir, ".ssh")

	// Open the JSON file
	host_jsonfile := Get_file_trunc(jsonfilepath, os.O_CREATE|os.O_RDONLY, 0644)
	defer host_jsonfile.Close()

	// Decode the JSON file into the temp_hosts slice
	decoder := json.NewDecoder(host_jsonfile)
	if err := decoder.Decode(&temp_hosts); err != nil {
		if err == io.EOF {
			fmt.Println("error: no hosts to run, use 'register_hosts' command to register.")
			os.Exit(1)
		}
		fmt.Printf("Failed to decode hosts.json file, Error: %s\n", err.Error())
		os.Exit(1)
	}

	if len(temp_hosts) <= 0 {
		fmt.Println("error: no hosts to run, use 'register_hosts' command to register.")
		os.Exit(1)
	}

	// Iterate through temp_hosts and create appropriate host type (Linux_host or Windows_host)
	for _, temp_host := range temp_hosts {
		if temp_host.Host_type == nil {
			fmt.Println("error: Host_type is nil and skipping hostname:", temp_host.Host_name)
			continue
		}

		switch *temp_host.Host_type {
		case common.LS_HOST_TYPE:
			host = &common.Linux_host{
				Host_name:         temp_host.Host_name,
				Host_run_username: temp_host.Host_run_username,
				Host_ip:           temp_host.Host_ip,
				Host_dns:          temp_host.Host_dns,
				Host_connect_port: temp_host.Host_connect_port,
				Host_use_ip:       temp_host.Host_use_ip,
				Host_type:         temp_host.Host_type,
			}

			common.Server_host = host
		case common.LA_HOST_TYPE:
			host = &common.Linux_host{
				Host_name:         temp_host.Host_name,
				Host_run_username: temp_host.Host_run_username,
				Host_ip:           temp_host.Host_ip,
				Host_dns:          temp_host.Host_dns,
				Host_connect_port: temp_host.Host_connect_port,
				Host_use_ip:       temp_host.Host_use_ip,
				Host_type:         temp_host.Host_type,
			}
		case common.WA_HOST_TYPE:
			host = &common.Windows_host{
				Host_name:         temp_host.Host_name,
				Host_run_username: temp_host.Host_run_username,
				Host_ip:           temp_host.Host_ip,
				Host_dns:          temp_host.Host_dns,
				Host_connect_port: temp_host.Host_connect_port,
				Host_use_ip:       temp_host.Host_use_ip,
				Host_type:         temp_host.Host_type,
			}
		default:
			fmt.Println("error: Host_type is unknown and skipping hostname:", temp_host.Host_name)
			continue
		}

		// Get ips and dns from database
		rows, err := GetData(`select i.useip, i.dns, i.ip from hosts h, interface i
							where h.hostid = i.hostid and i.main = 1 and i.type = 1 and h.host = $1`, host.Get_Host_name())
		if err != nil {
			fmt.Print("Failed in quering hosts, Error:", err.Error())
			fmt.Println(" and skipping hostname:", temp_host.Host_name)
			continue
		}
		defer rows.Close()

		err_in_scan := false
		for rows.Next() {
			if err := rows.Scan(&host_use_ip, &host_dns, &host_ip); err != nil {
				fmt.Println("Failed in scanning hosts, Error:", err.Error())
				err_in_scan = true
			}

			host.Set_Host_use_ip(host_use_ip)
			host.Set_Host_dns(host_dns)
			host.Set_Host_ip(host_ip)
		}

		if err_in_scan {
			fmt.Println("error: can not get ips data from database and skipping hostname:", temp_host.Host_name)
			continue
		}

		if host.Get_Host_use_ip() {
			fmt.Printf("Getting client from %s:%d ...", host.Get_Host_ip(), host.Get_Host_connect_port())

			client, err := GetSSHClientWithKey(host.Get_Host_ip(), host.Get_Host_connect_port(), host.Get_Host_run_username(), ssh_key_filepath)
			if err != nil {
				fmt.Printf("GetSSHClientWithKey failed, Error: %s", err.Error())
				fmt.Println(" and skipping hostname:", temp_host.Host_name)
				continue
			}

			host.Set_Host_ssh_client(client)
		} else {
			fmt.Printf("Getting client from %s:%d ...", host.Get_Host_dns(), host.Get_Host_connect_port())

			client, err := GetSSHClientWithKey(host.Get_Host_dns(), host.Get_Host_connect_port(), host.Get_Host_run_username(), ssh_key_filepath)
			if err != nil {
				fmt.Printf("GetSSHClientWithKey failed, Error: %s", err.Error())
				fmt.Println(" and skipping hostname:", temp_host.Host_name)
				continue
			}

			host.Set_Host_ip(host.Get_Host_dns())
			host.Set_Host_ssh_client(client)
		}
		fmt.Println(" done.")

		// Append the created host to the Host_pool
		common.Host_pool = append(common.Host_pool, host)
	}

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
