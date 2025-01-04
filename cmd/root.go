/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/zukigit/remote_run-go/src/common"
	"github.com/zukigit/remote_run-go/src/lib"
	"github.com/zukigit/remote_run-go/src/tickets"
	"gopkg.in/yaml.v3"

	"github.com/spf13/cobra"
)

var tkts, run_tickets []common.Ticket
var run_testcases []common.TestCase

func set_ticket_values(t []common.Ticket) {
	for _, ticket := range t {
		ticket.Set_values()
	}
}

func check_duplicated_ticket() {
	seen := make(map[uint]bool)

	for _, tkt := range tkts {
		if seen[tkt.Get_no()] {
			fmt.Printf("Error: ticket[%d] is duplicated\n", tkt.Get_no())
			os.Exit(1)
		}
		seen[tkt.Get_no()] = true
	}
}

func add_run_tickets(ticket_number uint) {
	if ticket_number == 0 {
		run_tickets = tkts
	} else {
		for _, t := range tkts {
			if ticket_number == t.Get_no() {
				run_tickets = append(run_tickets, t)
				break
			}
		}
	}
}

func add_run_testcases(testcase_number uint) {
	for _, ticket := range run_tickets {
		ticket.Add_testcases()

		for _, testcase := range ticket.Get_testcases() {
			if testcase_number == 0 || testcase_number == testcase.Get_no() {
				testcase.Set_ticket_no(ticket.Get_no())
				run_testcases = append(run_testcases, testcase)

				if testcase_number != 0 {
					return
				}
			}
		}
	}
}

func save_runtks_records() {
	yaml_data, err := yaml.Marshal(run_tickets)
	if err != nil {
		fmt.Println("Failed in Marshaling run_tickets, Error:", err.Error())
		os.Exit(1)
	}

	err = os.WriteFile(common.Log_filepath+".yml", yaml_data, 0644)
	if err != nil {
		fmt.Printf("Error writing YAML to file: %v\n", err)
		return
	}
}

func run_tc() {
	for _, testcase := range run_testcases {
		common.Current_tk_no = testcase.Get_ticket_no()
		common.Current_tc_no = testcase.Get_no()

		fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "running..."))
		if testcase.Is_function_nil() {
			testcase.Set_status(common.FAILED)
			fmt.Println(lib.Logi(common.LOG_LEVEL_INFO, "has no function. SKIPPED"))
		} else {
			common.Run_testcase(testcase)
		}
	}

	if len(run_testcases) > 0 {
		common.Update_testcase_results_in_tickets(run_tickets)
		save_runtks_records()

		fmt.Println(lib.Formatted_log(common.LOG_LEVEL_INFO, "Logged File: %s.log", common.Log_filepath))
		fmt.Println(lib.Formatted_log(common.LOG_LEVEL_INFO, "Yaml File: %s.yml", common.Log_filepath))
	} else {
		fmt.Println("There is no testcase to run.")
	}
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "remote_run.exe [-m YOUR_DB_HOSTNAME | -p YOUR_DB_HOSTNAME]",
	Short: "Run registered tickets on pre-defined hosts.",
	Long:  "Run registered tickets on pre-defined hosts.",
	Args: func(cmd *cobra.Command, args []string) error {
		if common.Temp_mysqlDB_hostname == "" && common.Temp_psqlDB_hostname == "" {
			return fmt.Errorf("specify database hostname using -m(for mysql) or -p(for psql) flags")
		}

		if common.Temp_mysqlDB_hostname != "" && common.Temp_psqlDB_hostname != "" {
			return fmt.Errorf("err: doesn't support for multiple databases yet")
		}

		if common.Specific_testcase_no > 0 && common.Specific_ticket_no == 0 {
			return fmt.Errorf("specify the ticket number too by using --ticket")
		}

		return nil
	},

	Run: func(cmd *cobra.Command, args []string) {
		common.Log_filepath = lib.Get_filepath()
		common.Set_sugar(common.Log_filepath + ".log")
		defer common.Sugar.Sync()

		// Initialize DB Connection
		common.Set_db_hostname()

		fmt.Printf("Connecting to %s:%d ...", common.DB_hostname, common.DB_port)
		lib.ConnectDB(common.DB_user, common.DB_passwd, common.DB_name)
		defer common.DB.Close()

		fmt.Println("connected")

		lib.Set_host_pool()

		lib.Enable_common_jobnets()

		add_tickets(&tkts)
		set_ticket_values(tkts)
		check_duplicated_ticket()

		add_run_tickets(common.Specific_ticket_no)
		add_run_testcases(common.Specific_testcase_no)

		run_tc() // run test cases
	},
}

func Execute() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().UintVar(&common.Specific_ticket_no, "ticket", 0, "Ticket number to run specific ticket")
	rootCmd.Flags().UintVar(&common.Specific_testcase_no, "testcase", 0, "Testcase number to run specific testcase")
	rootCmd.Flags().StringVarP(&common.Temp_mysqlDB_hostname, "mysql-hostname", "m", "", "Database specific hostname to connect.")
	rootCmd.Flags().StringVarP(&common.Temp_psqlDB_hostname, "psql-hostname", "p", "", "Database specific hostname to connect.")
	rootCmd.Flags().StringVar(&common.DB_user, "db-user", "zabbix", "Database specific username to connect.")
	rootCmd.Flags().StringVar(&common.DB_passwd, "db-password", "zabbix", "Database specific password to connect.")
	rootCmd.Flags().StringVar(&common.DB_name, "db-name", "zabbix", "Database specific name to connect.")
	rootCmd.Flags().UintVar(&common.DB_port, "db-port", 0, "Database specific port to connect.")
	rootCmd.Flags().UintVar(&common.Timeout, "timeout", 300, "Common timeout in seconds. ")
}

// Add your tickets here
func add_tickets(t *[]common.Ticket) {
	*t = append(*t, new(tickets.Ticket_000))
	*t = append(*t, new(tickets.Ticket_010))
	*t = append(*t, new(tickets.Ticket_698))
	*t = append(*t, new(tickets.Ticket_775))
	*t = append(*t, new(tickets.Ticket_794))
	*t = append(*t, new(tickets.Ticket_800))
	*t = append(*t, new(tickets.Ticket_811))
	*t = append(*t, new(tickets.Ticket_840))
	*t = append(*t, new(tickets.Ticket_844))
	*t = append(*t, new(tickets.Ticket_919))
	*t = append(*t, new(tickets.Ticket_923))
	*t = append(*t, new(tickets.Ticket_943))
	*t = append(*t, new(tickets.Ticket_952))
	*t = append(*t, new(tickets.Ticket_962))
	*t = append(*t, new(tickets.Ticket_1021))
	*t = append(*t, new(tickets.Ticket_1089))
	*t = append(*t, new(tickets.Ticket_1225))
	*t = append(*t, new(tickets.Ticket_1234))
	*t = append(*t, new(tickets.Ticket_1292))
	*t = append(*t, new(tickets.Ticket_1341))
	// *t = append(*t, new(tickets.Ticket_821))
	//*t = append(*t, new(tickets.Ticket_825))
	//*t = append(*t, new(tickets.Ticket_940))
	// *t = append(*t, new(tickets.Ticket_968))
	// *t = append(*t, new(tickets.Ticket_1281))
	// *t = append(*t, new(tickets.Ticket_1318))
}
