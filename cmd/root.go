/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/zukigit/remote_run-go/src/common"
	"github.com/zukigit/remote_run-go/src/dao"
	"github.com/zukigit/remote_run-go/src/lib"
	"github.com/zukigit/remote_run-go/src/tickets"

	"github.com/spf13/cobra"
)

var tkts, run_tickets []dao.Ticket

func set_ticket_values(t []dao.Ticket) {
	for _, ticket := range t {
		ticket.Set_values()
	}
}

func enable_common_jobnets() {
	if err := lib.Jobarg_enable_jobnet("Icon_1", "jobicon_linux"); err != nil {
		fmt.Println("Failed to enable common jobnets, error: ", err.Error())
	}

	if err := lib.Jobarg_enable_jobnet("Icon_2", "Icon_2"); err != nil {
		fmt.Println("Failed to enable common jobnets, error: ", err.Error())
	}

	if err := lib.Jobarg_enable_jobnet("Icon_10", "Icon_10"); err != nil {
		fmt.Println("Failed to enable common jobnets, error: ", err.Error())
	}

	if err := lib.Jobarg_enable_jobnet("Icon_100", "Icon_100"); err != nil {
		fmt.Println("Failed to enable common jobnets, error: ", err.Error())
	}

	if err := lib.Jobarg_enable_jobnet("Icon_200", "Icon_200"); err != nil {
		fmt.Println("Failed to enable common jobnets, error: ", err.Error())
	}

	if err := lib.Jobarg_enable_jobnet("Icon_400", "Icon_400"); err != nil {
		fmt.Println("Failed to enable common jobnets, error: ", err.Error())
	}

	if err := lib.Jobarg_enable_jobnet("Icon_500", "Icon_500"); err != nil {
		fmt.Println("Failed to enable common jobnets, error: ", err.Error())
	}

	if err := lib.Jobarg_enable_jobnet("Icon_510", "Icon_510"); err != nil {
		fmt.Println("Failed to enable common jobnets, error: ", err.Error())
	}

	if err := lib.Jobarg_enable_jobnet("Icon_800", "Icon_800"); err != nil {
		fmt.Println("Failed to enable common jobnets, error: ", err.Error())
	}

	if err := lib.Jobarg_enable_jobnet("Icon_1000", "Icon_1000"); err != nil {
		fmt.Println("Failed to enable common jobnets, error: ", err.Error())
	}

	if err := lib.Jobarg_enable_jobnet("Icon_1020", "Icon_1020"); err != nil {
		fmt.Println("Failed to enable common jobnets, error: ", err.Error())
	}

	if err := lib.Jobarg_enable_jobnet("Icon_2040", "Icon_2040"); err != nil {
		fmt.Println("Failed to enable common jobnets, error: ", err.Error())
	}

	if err := lib.Jobarg_enable_jobnet("Icon_3000", "Icon_3000"); err != nil {
		fmt.Println("Failed to enable common jobnets, error: ", err.Error())
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

// Add your tickets here
func add_testcases() {
	for _, tkt := range tkts {
		tkt.Add_testcases()
	}
}

func run_tc_(t dao.Ticket) {
	lib.Logi(fmt.Sprintf("Ticket[%d] %s\n", t.Get_no(), t.Get_dsctn()))
	lib.Logi("\nTest_cases:\n")
	lib.Logi(fmt.Sprintf("%s\n", common.Endtestcase_string))
	dao.Run_testcase(t)
	dao.Set_total_tc_results(t)
	lib.Logi("\nTotal_result:\n")

	if dao.Tc_unkown_cnt > 0 {
		lib.Logi(fmt.Sprintf("PASSED: %d, FAILED: %d, MUST_CHECK: %d, UNKNOWN: %d\n", dao.Tc_passed_cnt, dao.Tc_failed_cnt, dao.Tc_chk_cnt, dao.Tc_unkown_cnt))
	} else {
		lib.Logi(fmt.Sprintf("PASSED: %d, FAILED: %d, MUST_CHECK: %d\n", dao.Tc_passed_cnt, dao.Tc_failed_cnt, dao.Tc_chk_cnt))
	}
	lib.Logi(fmt.Sprintf("%s\n", common.Endticket_string))
}

func run_tc(t []dao.Ticket) {
	for _, ticket := range t {
		if common.Specific_testcase_no == 0 {
			run_tc_(ticket)
		} else {
			for _, tc := range ticket.Get_testcases() {
				if tc.Get_id() == common.Specific_testcase_no {
					run_tc_(ticket)
				}
			}
		}
	}
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "remote_run.exe user@host --with-mysql (or) --with-postgresql",
	Short: "Automated testing",
	Long:  "Automated testing",
	Args: func(cmd *cobra.Command, args []string) error {
		if err := cobra.MinimumNArgs(1)(cmd, args); err != nil {
			return err
		}

		if err := common.Set_db_type(); err != nil {
			return err
		}

		if common.Specific_testcase_no > 0 && common.Specific_ticket_no == 0 {
			return fmt.Errorf("specify the ticket number too by using --ticket")
		}

		return common.Set_usr_hst(args)
	},

	Run: func(cmd *cobra.Command, args []string) {
		common.Log_filename = lib.Get_log_filename()
		common.Set_passwd()
		common.Set_client()
		defer common.Client.Close()

		common.Set_db_hostname()
		common.Set_default_db_port()

		// Initialize DB Connection
		lib.ConnectDB("zabbix", "zabbix", "zabbix")
		defer common.DB.Close()

		common.Set_log_file(common.Log_filename)
		defer common.Log_file.Close()

		common.Set_ticket_logs_headers()

		add_tickets(&tkts)
		set_ticket_values(tkts)
		add_testcases()
		check_duplicated_ticket()
		add_run_tickets(common.Specific_ticket_no)
		enable_common_jobnets()
		run_tc(run_tickets) // run test cases

		if len(run_tickets) > 0 {
			fmt.Println(lib.Formatted_log(common.INFO, "Logged Filename: %s", common.Log_filename))
		} else {
			fmt.Println("There is no ticket to run.")
		}
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
	rootCmd.Flags().IntVarP(&common.Login_info.Port, "port", "p", 22, "Port")
	rootCmd.Flags().BoolVar(&common.Is_mysql, "with-mysql", false, "Use MySQL database")
	rootCmd.Flags().BoolVar(&common.Is_psql, "with-postgresql", false, "Use PostgreSQL database")
	rootCmd.Flags().UintVar(&common.Specific_ticket_no, "ticket", 0, "Ticket number to run specific ticket")
	rootCmd.Flags().UintVar(&common.Specific_testcase_no, "testcase", 0, "Testcase number to run specific testcase")
	rootCmd.Flags().StringVar(&common.DB_hostname, "db-hostname", "", "Database specific hostname to connect.")
	rootCmd.Flags().UintVar(&common.DB_port, "db-port", 0, "Database specific port to connect.")
	rootCmd.Flags().UintVar(&common.Timeout, "timeout", 300, "Common timeout in seconds. ")
}

// Add your tickets here
func add_tickets(t *[]dao.Ticket) {
	*t = append(*t, new(tickets.Ticket_000))
	*t = append(*t, new(tickets.Ticket_1318))
	*t = append(*t, new(tickets.Ticket_811))
	*t = append(*t, new(tickets.Ticket_800))
	*t = append(*t, new(tickets.Ticket_1225))
	*t = append(*t, new(tickets.Ticket_844))
	*t = append(*t, new(tickets.Ticket_794))
	*t = append(*t, new(tickets.Ticket_968))
	*t = append(*t, new(tickets.Ticket_940))
	*t = append(*t, new(tickets.Ticket_943))
	*t = append(*t, new(tickets.Ticket_952))

}
