/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/zukigit/remote_run-go/src/common"
	"github.com/zukigit/remote_run-go/src/lib"
	"golang.org/x/crypto/ssh"
)

var found_hosts []common.Host
var rsa_pub_key string

func check_id_rsa() error {
	current_user, err := user.Current()
	if err != nil {
		return fmt.Errorf("failed in getting run user, Error: %v", err)
	}

	ssh_pub_filepath := filepath.Join(current_user.HomeDir, ".ssh", "id_rsa.pub")
	ssh_private_filepath := filepath.Join(current_user.HomeDir, ".ssh", "id_rsa")

	if _, err := os.Stat(ssh_pub_filepath); err != nil {
		return fmt.Errorf("expected key[%s] not found, use the following command to generate keys 'ssh-keygen -t rsa -b 4096'", ssh_pub_filepath)
	}

	if _, err := os.Stat(ssh_private_filepath); err != nil {
		return fmt.Errorf("expected key[%s] not found, use the following command to generate keys 'ssh-keygen -t rsa -b 4096'", ssh_private_filepath)
	}

	content, err := os.ReadFile(ssh_pub_filepath)
	if err != nil {
		return fmt.Errorf("failed to read file: %v", err)
	}

	rsa_pub_key = string(content)
	fmt.Println(rsa_pub_key)

	return nil
}

func get_hosts() {
	rows, err := common.DB.Query(`select h.host, i.useip, i.dns, i.ip from hosts h, interface i
		where h.hostid = i.hostid and i.main = 1 and i.type = 1 and h.host LIKE 'auto.linux%'`)
	if err != nil {
		fmt.Println("Failed in quering hosts, Error:", err.Error())
		os.Exit(1)
	}
	defer rows.Close()

	for rows.Next() {
		var host common.Host

		if err := rows.Scan(&host.Host_name, &host.Host_use_ip, &host.Host_dns, &host.Host_ip); err != nil {
			fmt.Println("Failed in scanning hosts, Error:", err.Error())
			os.Exit(1)
		}
		found_hosts = append(found_hosts, host)
	}
}

func get_host(hostname string) *common.Host {
	for _, host := range found_hosts {
		if hostname == host.Host_name {
			return &host
		}
	}

	return nil
}

func ask_userinput_hostname() string {
	fmt.Println("Found hosts:")
	for index, host := range found_hosts {
		fmt.Printf("%d) %s\n", index+1, host.Host_name)
	}
	fmt.Println("--------------")

	return lib.Ask_usrinput_string("Enter hostname to register")
}

func register() {
	var temp_host *common.Host
	var temp_passwd string
	var err error
	var temp_sshcli *ssh.Client

	for {
		if temp_host = get_host(ask_userinput_hostname()); temp_host == nil {
			fmt.Print("you entered wrong hostname, ")
			continue
		}
		break
	}

	temp_host.Host_run_username = lib.Ask_usrinput_string("Enter ssh username to register")
	temp_passwd = lib.Ask_usrinput_passwd_string("Enter ssh password to register")
	for {
		temp_host.Host_port, err = lib.Ask_usrinput_int("Enter ssh port to register")
		if err != nil {
			fmt.Print("you entered wrong port, ")
			continue
		}
		break
	}

	if temp_host.Host_use_ip {
		temp_sshcli = lib.GetSSHClient(temp_host.Host_ip, temp_host.Host_port, temp_host.Host_run_username, temp_passwd)
	} else {
		temp_sshcli = lib.GetSSHClient(temp_host.Host_dns, temp_host.Host_port, temp_host.Host_run_username, temp_passwd)
	}

	cmd := fmt.Sprintf("echo '%s' >> ~/.ssh/authorized_keys", rsa_pub_key)
	lib.GetOutputStrFromSSHCommand(temp_sshcli, cmd)
}

// registerHostsCmd represents the registerHosts command
var registerHostsCmd = &cobra.Command{
	Use:   "register_hosts",
	Short: "Register new hosts.",
	Long:  "This command will scan the zabbix database to find hosts that has names started with 'auto.' and register in hosts.ini file.",
	Args: func(cmd *cobra.Command, args []string) error {
		if err := check_id_rsa(); err != nil {
			return err
		}
		if common.DB_hostname == "" {
			return fmt.Errorf("specify database hostname using --db-hostname flag")
		}
		return common.Set_db_type()
	},
	Run: func(cmd *cobra.Command, args []string) {
		// Initialize DB Connection
		common.Set_db_hostname()
		common.Set_default_db_port()
		lib.ConnectDB(common.DB_user, common.DB_passwd, common.DB_name)
		defer common.DB.Close()

		get_hosts()
		register()
	},
}

func init() {
	rootCmd.AddCommand(registerHostsCmd)
	registerHostsCmd.Flags().BoolVar(&common.Is_mysql, "with-mysql", false, "Use MySQL database")
	registerHostsCmd.Flags().BoolVar(&common.Is_psql, "with-postgresql", false, "Use PostgreSQL database")
	registerHostsCmd.Flags().StringVar(&common.DB_hostname, "db-hostname", "", "Database specific hostname to connect.")
	registerHostsCmd.Flags().StringVar(&common.DB_user, "db-user", "zabbix", "Database specific username to connect.")
	registerHostsCmd.Flags().StringVar(&common.DB_passwd, "db-password", "zabbix", "Database specific password to connect.")
	registerHostsCmd.Flags().StringVar(&common.DB_name, "db-name", "zabbix", "Database specific name to connect.")
	registerHostsCmd.Flags().UintVar(&common.DB_port, "db-port", 0, "Database specific port to connect.")
}
