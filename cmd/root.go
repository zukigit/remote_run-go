/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"zukigit/remote_run-go/src/common"
	"zukigit/remote_run-go/src/dao"
	"zukigit/remote_run-go/src/lib"
	"zukigit/remote_run-go/src/tickets"

	"github.com/spf13/cobra"
)

var tkts, run_tickets []dao.Ticket

func set_ticket_values(t []dao.Ticket) {
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

func run_tc(t []dao.Ticket) {
	for _, ticket := range t {
		ticket.Add_testcases()
		ticket.Run()
		lib.Logi(dao.Get_ticket_logs(ticket), common.Log_filename)
	}
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "remote_run.exe user@host",
	Short: "Automated testing",
	Long:  "Automated testing",
	Args: func(cmd *cobra.Command, args []string) error {
		if err := cobra.ExactArgs(1)(cmd, args); err != nil {
			return err
		}

		return common.Set_usr_hst(args)
	},

	Run: func(cmd *cobra.Command, args []string) {
		common.Log_filename = lib.Get_log_filename()
		common.Set_passwd()
		common.Set_client()
		defer common.Client.Close()

		common.Set_ticket_logs_headers()
		add_tickets(&tkts)
		set_ticket_values(tkts)
		check_duplicated_ticket()
		add_run_tickets(common.Specific_ticket_no)

		run_tc(run_tickets) // run test cases

		if len(run_tickets) > 0 {
			fmt.Println(lib.Formatted_log(common.INFO, "Logged Filename: %s", common.Log_filename))
		} else {
			fmt.Println("There is no ticket to run.")
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().UintVarP(&common.Specific_ticket_no, "ticket_no", "t", 0, "Ticket number to run")
	rootCmd.Flags().IntVarP(&common.Login_info.Port, "port", "p", 22, "Port")
}

// Add your tickets here
func add_tickets(t *[]dao.Ticket) {
	*t = append(*t, new(tickets.Ticket_1318))
	*t = append(*t, new(tickets.Ticket_000))
}
