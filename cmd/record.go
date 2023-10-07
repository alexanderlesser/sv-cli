/*
Copyright © 2023 github.com/alexanderlesser
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// recordCmd represents the record command
var recordCmd = &cobra.Command{
	Use:   "record",
	Short: "record [command] handles all tasks for records",
	Long: `
	record [command] handles all tasks for records such as add, list and delete.

	Examples:
	record add - To add a new record you run command "record add" and fill in the fields prompted.
	record list - to list all existing records you run command "record list".
	record delete - to delete an existing record you run command "record delete".
	
	For more information about a specific command run record [command] -h`,
	Run: func(cmd *cobra.Command, args []string) {
		list, _ := cmd.Flags().GetBool("list")
		add, _ := cmd.Flags().GetBool("add")
		delete, _ := cmd.Flags().GetBool("delete")

		if list {
			listCmd.Run(cmd, []string{""})
		}

		if add {
			addCmd.Run(cmd, []string{""})
		}

		if delete {
			deleteCmd.Run(cmd, []string{""})
		}

		if !list && !add && !delete {
			fmt.Println("record [command] handles all tasks for records.\nFor more information about this command run record -h")
		}
	},
}

func init() {
	rootCmd.AddCommand(recordCmd)
	recordCmd.Flags().BoolP("add", "a", false, "Run add command")
	recordCmd.Flags().BoolP("list", "l", false, "Run list command")
	recordCmd.Flags().BoolP("delete", "d", false, "Run delete command")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// recordCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// recordCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
