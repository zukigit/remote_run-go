/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/zukigit/remote_run-go/src/common"
	"github.com/zukigit/remote_run-go/src/lib"
	"golang.org/x/crypto/ssh"
)

var found_hosts []common.Host
var rsa_pub_key string
var chosen_hosts_index int

func check_id_rsa() error {
	current_user, err := user.Current()
	if err != nil {
		return fmt.Errorf("failed in getting run user, Error: %v", err)
	}

	key_file_path := filepath.Join(current_user.HomeDir, ".ssh")
	ssh_pub_filepath := filepath.Join(key_file_path, "id_rsa.pub")
	ssh_private_filepath := filepath.Join(key_file_path, "id_rsa")

	if _, err := os.Stat(ssh_pub_filepath); err != nil {
		fmt.Printf("expected key[%s] not found, generating new key files...\n", ssh_pub_filepath)
		lib.Generate_sshkeys(key_file_path)
		return nil
	}

	if _, err := os.Stat(ssh_private_filepath); err != nil {
		fmt.Printf("expected key[%s] not found, generating new key files...\n", ssh_private_filepath)
		lib.Generate_sshkeys(key_file_path)
		return nil
	}

	content, err := os.ReadFile(ssh_pub_filepath)
	if err != nil {
		return fmt.Errorf("failed to read file: %v", err)
	}

	rsa_pub_key = string(content)

	return nil
}

func get_hosts() {
	var host_name, host_dns, host_ip string
	var host_use_ip bool

	// rows, err := common.DB.Query(`select h.host, i.useip, i.dns, i.ip from hosts h, interface i
	// 	where h.hostid = i.hostid and i.main = 1 and i.type = 1 and h.host LIKE 'auto.linux.agent.%'`)
	rows, err := common.DB.Query(`select h.host, i.useip, i.dns, i.ip from hosts h, interface i
		where h.hostid = i.hostid and i.main = 1 and i.type = 1 and h.host LIKE 'auto.%'`)
	if err != nil {
		fmt.Println("Failed in quering hosts, Error:", err.Error())
		os.Exit(1)
	}
	defer rows.Close()

	for rows.Next() {
		var host common.Host

		if err := rows.Scan(&host_name, &host_use_ip, &host_dns, &host_ip); err != nil {
			fmt.Println("Failed in scanning hosts, Error:", err.Error())
			os.Exit(1)
		}

		switch {
		case strings.HasPrefix(host_name, string(common.LINUX_SERVER)):
			host = common.New_linux_host()
			host.Set_Host_type(common.LINUX_SERVER)
		case strings.HasPrefix(host_name, string(common.LINUX_AGENT)):
			host = common.New_linux_host()
			host.Set_Host_type(common.LINUX_AGENT)
		case strings.HasPrefix(host_name, string(common.WINDOWS_AGENT)):
			host = common.New_windows_host()
			host.Set_Host_type(common.WINDOWS_AGENT)
		default:
			fmt.Printf("Host_name: %s does not match formats('%s' or '%s' or '%s') and will not be registered.", host_name, common.LINUX_SERVER, common.LINUX_AGENT, common.WINDOWS_AGENT)
			continue
		}

		host.Set_Host_name(host_name)
		host.Set_Host_use_ip(host_use_ip)
		host.Set_Host_dns(host_dns)
		host.Set_Host_ip(host_ip)

		found_hosts = append(found_hosts, host)
	}
}

func get_host(hostname string) common.Host {
	if hostname != "" {
		for index, host := range found_hosts {
			if hostname == host.Get_Host_name() {
				chosen_hosts_index = index
				return host
			}
		}
	}

	return nil
}

func check_duplicated_hosts(temp_hosts *[]common.Host, temp_host common.Host) {
	// Iterate through the slice to check for duplicates
	for index, host := range *temp_hosts {
		if host.Get_Host_name() == temp_host.Get_Host_name() || (host.Get_Host_type() == common.LINUX_SERVER && temp_host.Get_Host_type() == common.LINUX_SERVER) {
			// If a duplicate is found, update the existing host
			(*temp_hosts)[index] = temp_host
			return
		}
	}

	// If no duplicate is found, append the new host
	*temp_hosts = append(*temp_hosts, temp_host)
}

func register(ssh_key_filepath string) {
	var temp_host common.Host
	var temp_passwd string
	var err error
	var temp_sshcli *ssh.Client
	temp_port := 22
	var already_registered bool

	for {
		fmt.Println()

		if len(found_hosts) == 0 {
			fmt.Println("No hosts to register, exiting...")
			os.Exit(0)
		}

		fmt.Println("Found hosts:")
		for index, host := range found_hosts {
			switch host.(type) {
			case *common.Linux_host:
				fmt.Printf("%d) %s\n", index+1, host.Get_Host_name())
			default:
				fmt.Printf("%d) %s (not avaliable yet)\n", index+1, host.Get_Host_name())
			}
		}
		fmt.Println("--------------")

		temp_hostname := lib.Ask_usrinput_string("Enter hostname to register")

		parts := strings.Split(temp_hostname, ":")
		if len(parts) == 2 {
			temp_hostname = parts[0]
			temp_port, err = strconv.Atoi(parts[1])
			if err != nil {
				fmt.Println("Err: you entered wrong port")
				continue
			}
		}

		if temp_host = get_host(temp_hostname); temp_host == nil {
			fmt.Println("Err: you entered wrong hostname")
			continue
		}

		if temp_host.Get_Host_type() == common.WINDOWS_AGENT {
			fmt.Println("Err: WINDOWS_AGENT is not supported yet, skipping... ")
			continue
		}

		temp_host.Set_Host_connect_port(temp_port)
		break
	}

	temp_host.Set_Host_run_username(lib.Ask_usrinput_string("Enter ssh username to register"))
	temp_passwd = lib.Ask_usrinput_passwd_string("Enter ssh password to register")

	fmt.Print("Registering, ")
	if temp_host.Get_Host_use_ip() {
		fmt.Printf("host_ip: %s:%d...\n", temp_host.Get_Host_ip(), temp_host.Get_Host_connect_port())

		temp_sshcli, err = lib.GetSSHClientWithKey(temp_host.Get_Host_ip(), temp_host.Get_Host_connect_port(), temp_host.Get_Host_run_username(), ssh_key_filepath)
		if err != nil {
			temp_sshcli = lib.GetSSHClient(temp_host.Get_Host_ip(), temp_host.Get_Host_connect_port(), temp_host.Get_Host_run_username(), temp_passwd)
		} else {
			already_registered = true
		}
	} else {
		fmt.Printf("host_dns: %s:%d...\n", temp_host.Get_Host_dns(), temp_host.Get_Host_connect_port())
		temp_sshcli, err = lib.GetSSHClientWithKey(temp_host.Get_Host_dns(), temp_host.Get_Host_connect_port(), temp_host.Get_Host_run_username(), ssh_key_filepath)
		if err != nil {
			temp_sshcli = lib.GetSSHClient(temp_host.Get_Host_dns(), temp_host.Get_Host_connect_port(), temp_host.Get_Host_run_username(), temp_passwd)
		} else {
			already_registered = true
		}
	}

	if !already_registered {
		temp_host.Set_Host_ssh_client(temp_sshcli)
		if err = temp_host.Register(rsa_pub_key); err != nil {
			fmt.Printf("Failed to register host[%s], Error: %s\n", temp_host.Get_Host_name(), err.Error())
			os.Exit(1)
		}
	}

	lib.Get_hosts_from_jsonfile("hosts.json")
	check_duplicated_hosts(&common.Host_pool, temp_host)

	lib.Set_hosts_to_jsonfile(&common.Host_pool, "hosts.json")

	fmt.Printf("Registered host[%s]\n", temp_host.Get_Host_name())
	found_hosts = append(found_hosts[:chosen_hosts_index], found_hosts[chosen_hosts_index+1:]...)
}

// registerHostsCmd represents the registerHosts command
var registerHostsCmd = &cobra.Command{
	Use:   "register_hosts --db-hostname YOUR_DB_HOSTNAME [--with-postgresql | --with-mysql]",
	Short: "Register new hosts.",
	Long:  "This command will scan hosts that starts with 'auto.' from zabbix database and register it in hosts.json file.",
	Args: func(cmd *cobra.Command, args []string) error {
		if common.Temp_mysqlDB_hostname == "" || common.Temp_psqlDB_hostname == "" {
			return fmt.Errorf("specify database hostname using --db-hostname flag")
		}

		if common.Temp_mysqlDB_hostname != "" && common.Temp_psqlDB_hostname != "" {
			return fmt.Errorf("you can't use both -m and -p flags, just choose one")
		}

		return check_id_rsa()
	},
	Run: func(cmd *cobra.Command, args []string) {
		common.Set_db_hostname()

		// Initialize DB Connection
		lib.ConnectDB(common.DB_user, common.DB_passwd, common.DB_name)
		defer common.DB.Close()

		current_user, err := user.Current()
		if err != nil {
			fmt.Printf("Err: failed in getting run user, Error: %v", err.Error())
			os.Exit(1)
		}
		ssh_key_filepath := filepath.Join(current_user.HomeDir, ".ssh")

		get_hosts()
		for {
			register(ssh_key_filepath)
		}
	},
}

func init() {
	rootCmd.AddCommand(registerHostsCmd)
	registerHostsCmd.Flags().StringVar(&common.DB_user, "db-user", "zabbix", "Database specific username to connect.")
	registerHostsCmd.Flags().StringVar(&common.DB_passwd, "db-password", "zabbix", "Database specific password to connect.")
	registerHostsCmd.Flags().StringVar(&common.DB_name, "db-name", "zabbix", "Database specific name to connect.")
	registerHostsCmd.Flags().UintVar(&common.DB_port, "db-port", 0, "Database specific port to connect.")
}
