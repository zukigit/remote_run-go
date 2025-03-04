/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// ticketsCmd represents the tickets command
var ticketsCmd = &cobra.Command{
	Use:   "tickets",
	Short: "List all avaliable tickets.",
	Run: func(cmd *cobra.Command, args []string) {
		add_tickets(&tkts)
		check_duplicated_ticket()

		fmt.Println("Tickets:")
		for index, t := range tkts {
			fmt.Printf("(%d) %d - %s\n", index+1, t.Get_ticket_no(), t.Get_ticket_description())
		}
	},
}

func init() {
	rootCmd.AddCommand(ticketsCmd)
}
