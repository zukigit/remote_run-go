/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zukigit/remote_run-go/src/common"
	"github.com/zukigit/remote_run-go/src/lib"
)

// registerHostsCmd represents the registerHosts command
var registerHostsCmd = &cobra.Command{
	Use:   "register_hosts",
	Short: "Register new hosts.",
	Long:  "This command will scan the zabbix database to find hosts that has names started with 'auto.' and register in hosts.ini file.",
	Args: func(cmd *cobra.Command, args []string) error {
		if common.DB_hostname == "" {
			return fmt.Errorf("specify database hostname using --db-hostname flag")
		}
		return common.Set_db_type()
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("register_hosts called")
		// Initialize DB Connection
		common.Set_db_hostname()
		common.Set_default_db_port()
		lib.ConnectDB(common.DB_user, common.DB_passwd, common.DB_name)
		defer common.DB.Close()
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
