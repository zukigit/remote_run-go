/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var target_tc_no uint

// testcasesCmd represents the testcases command
var testcasesCmd = &cobra.Command{
	Use:   "testcases [testcase_number_to_search]",
	Short: "List all avaliable tickets.",
	Args: func(cmd *cobra.Command, args []string) error {
		if err := cobra.MaximumNArgs(1)(cmd, args); err != nil {
			return err
		}

		if len(args) > 0 {
			tc_int, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}

			target_tc_no = uint(tc_int)
		}

		return nil
	},

	Run: func(cmd *cobra.Command, args []string) {
		var tc_index int

		add_tickets(&tkts)
		check_duplicated_ticket()

		fmt.Println("Testcases:")
		for _, tk := range tkts {
			for _, tc := range tk.Get_testcases() {
				if target_tc_no == 0 || target_tc_no == tc.Get_ticket_no() {
					fmt.Printf("(%d) Ticket: %d, Testcase: %d, Testcase_description: %s\n",
						tc_index+1, tk.Get_ticket_no(), tc.Get_ticket_no(), tc.Get_ticket_description())

					if target_tc_no != 0 {
						break
					}
					tc_index++
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(testcasesCmd)
}
