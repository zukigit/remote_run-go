/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// testcasesCmd represents the testcases command
var testcasesCmd = &cobra.Command{
	Use:   "testcases",
	Short: "List all avaliable tickets.",
	Run: func(cmd *cobra.Command, args []string) {
		tc_index := 0
		fmt.Println("Testcases:")

		for _, tk := range tkts {
			for _, tc := range tk.Get_testcases() {
				fmt.Printf("(%d) Testcase: %d, Ticket: %d, Testcase_description: %s\n", tc_index+1, tc.Get_id(), tk.Get_no(), tc.Get_dsctn())
				tc_index++
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(testcasesCmd)
}
