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

	rows, err := common.DB.Query(`select h.host, i.useip, i.dns, i.ip from hosts h, interface i
		where h.hostid = i.hostid and i.main = 1 and i.type = 1 and h.host LIKE 'auto.%'`)
	if err != nil {
		fmt.Println("Error: failed in quering hosts, Error:", err.Error())
		os.Exit(1)
	}
	defer rows.Close()

	for rows.Next() {
		var host common.Host

		if err := rows.Scan(&host_name, &host_use_ip, &host_dns, &host_ip); err != nil {
			fmt.Println("Error: failed in scanning hosts, Error:", err.Error())
			os.Exit(1)
		}

		switch {
		case strings.HasPrefix(host_name, string(common.LS_HOST_TYPE)):
			host = common.New_linux_host()
			host.Set_Host_type(common.LS_HOST_TYPE)
		case strings.HasPrefix(host_name, string(common.LA_HOST_TYPE)):
			host = common.New_linux_host()
			host.Set_Host_type(common.LA_HOST_TYPE)
		case strings.HasPrefix(host_name, string(common.WA_HOST_TYPE)):
			host = common.New_windows_host()
			host.Set_Host_type(common.WA_HOST_TYPE)
		default:
			fmt.Printf("Host_name: %s does not match formats('%s' or '%s' or '%s') and will not be registered.", host_name, common.LS_HOST_TYPE, common.LA_HOST_TYPE, common.WA_HOST_TYPE)
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
		if host.Get_Host_name() == temp_host.Get_Host_name() || (host.Get_Host_type() == common.LS_HOST_TYPE && temp_host.Get_Host_type() == common.LS_HOST_TYPE) {
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
			fmt.Printf("%d) %s\n", index+1, host.Get_Host_name())
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
			temp_sshcli, err = lib.GetSSHClient_(temp_host.Get_Host_ip(), temp_host.Get_Host_connect_port(), temp_host.Get_Host_run_username(), temp_passwd)
			if err != nil {
				fmt.Printf("Error: failed to register host[%s], Error: %s\n", temp_host.Get_Host_name(), err.Error())
				return
			}
		} else {
			already_registered = true
		}
	} else {
		fmt.Printf("host_dns: %s:%d...\n", temp_host.Get_Host_dns(), temp_host.Get_Host_connect_port())
		temp_sshcli, err = lib.GetSSHClientWithKey(temp_host.Get_Host_dns(), temp_host.Get_Host_connect_port(), temp_host.Get_Host_run_username(), ssh_key_filepath)
		if err != nil {
			temp_sshcli, err = lib.GetSSHClient_(temp_host.Get_Host_dns(), temp_host.Get_Host_connect_port(), temp_host.Get_Host_run_username(), temp_passwd)
			if err != nil {
				fmt.Printf("Error: failed to register host[%s], Error: %s\n", temp_host.Get_Host_name(), err.Error())
				return
			}
		} else {
			already_registered = true
		}
	}

	if !already_registered {
		temp_host.Set_Host_ssh_client(temp_sshcli)
		if err = temp_host.Register(rsa_pub_key); err != nil {
			fmt.Printf("Error: failed to register host[%s], Error: %s\n", temp_host.Get_Host_name(), err.Error())
			return
		}
	}

	// try to connect with ssh-keys after registration
	fmt.Print("Connecting to ")
	if temp_host.Get_Host_use_ip() {
		fmt.Printf("host_ip: %s:%d...", temp_host.Get_Host_ip(), temp_host.Get_Host_connect_port())

		_, err = lib.GetSSHClientWithKey(temp_host.Get_Host_ip(), temp_host.Get_Host_connect_port(), temp_host.Get_Host_run_username(), ssh_key_filepath)
		if err != nil {
			fmt.Printf("Error: failed to register host[%s], Error: %s\n", temp_host.Get_Host_name(), err.Error())
			return
		}
	} else {
		fmt.Printf("host_dns: %s:%d...\n", temp_host.Get_Host_dns(), temp_host.Get_Host_connect_port())

		_, err = lib.GetSSHClientWithKey(temp_host.Get_Host_dns(), temp_host.Get_Host_connect_port(), temp_host.Get_Host_run_username(), ssh_key_filepath)
		if err != nil {
			fmt.Printf("Error: failed to register host[%s], Error: %s\n", temp_host.Get_Host_name(), err.Error())
			return
		}
	}
	fmt.Println("connected!")

	lib.Get_hosts_from_jsonfile("hosts.json")
	check_duplicated_hosts(&common.Hosts, temp_host)

	lib.Set_hosts_to_jsonfile(&common.Hosts, "hosts.json")

	fmt.Printf("Registered host[%s]\n", temp_host.Get_Host_name())
	found_hosts = append(found_hosts[:chosen_hosts_index], found_hosts[chosen_hosts_index+1:]...)
}

// registerHostsCmd represents the registerHosts command
var registerHostsCmd = &cobra.Command{
	Use:   "register_hosts [-m YOUR_DB_HOSTNAME | -p YOUR_DB_HOSTNAME]",
	Short: "Register new hosts.",
	Long:  "This command will scan hosts that starts with 'auto.' from zabbix database and register it in hosts.json file.",
	Args: func(cmd *cobra.Command, args []string) error {
		if common.Temp_mysqlDB_hostname == "" && common.Temp_psqlDB_hostname == "" {
			return fmt.Errorf("specify database hostname using -m(for mysql) or -p(for psql) flags")
		}

		if common.Temp_mysqlDB_hostname != "" && common.Temp_psqlDB_hostname != "" {
			return fmt.Errorf("doesn't support for multiple databases yet")
		}

		return check_id_rsa()
	},
	Run: func(cmd *cobra.Command, args []string) {
		common.Set_db_hostname()

		// Initialize DB Connection
		fmt.Printf("Connecting to %s:%d ...", common.DB_hostname, common.DB_port)
		lib.ConnectDB(common.DB_user, common.DB_passwd, common.DB_name)
		defer common.DB.Close()
		fmt.Println("connected")

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
	registerHostsCmd.Flags().StringVarP(&common.Temp_mysqlDB_hostname, "mysql-hostname", "m", "", "Database specific hostname to connect.")
	registerHostsCmd.Flags().StringVarP(&common.Temp_psqlDB_hostname, "psql-hostname", "p", "", "Database specific hostname to connect.")
	registerHostsCmd.Flags().UintVar(&common.DB_port, "db-port", 0, "Database specific port to connect.")
}
