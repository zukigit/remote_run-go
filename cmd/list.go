/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all avaliable tickets.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Total ticket:", len(tkts))
		fmt.Println()
		for index, t := range tkts {
			fmt.Printf("(%d) %d - %s\n", index+1, t.Get_no(), t.Get_dsctn())
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
