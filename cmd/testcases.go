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
		add_tickets(&tkts)
		set_ticket_values(tkts)
		add_testcases()
		check_duplicated_ticket()

		tc_index := 0
		fmt.Println("Testcases:")

		for _, tk := range tkts {
			for _, tc := range tk.Get_testcases() {
				if target_tc_no == 0 {
					fmt.Printf("(%d) Ticket: %d, Testcase: %d, Testcase_description: %s\n",
						tc_index+1, tk.Get_no(), tc.Get_id(), tc.Get_dsctn())
					tc_index++
				} else if target_tc_no == tc.Get_id() {
					fmt.Printf("(%d) Testcase: %d, Ticket: %d, Testcase_description: %s\n",
						tc_index+1, tc.Get_id(), tk.Get_no(), tc.Get_dsctn())
					tc_index++
					// Break out of the inner loop since we found the specific testcase
					break
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(testcasesCmd)
}
