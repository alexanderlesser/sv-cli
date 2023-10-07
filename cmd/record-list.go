/*
Copyright © 2023 github.com/alexanderlesser
*/
package cmd

import (
	"fmt"

	"github.com/alexanderlesser/sv-cli/datastore"
	"github.com/alexanderlesser/sv-cli/internal/components"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "lists all records",
	Long:  `The list command lists all records in the datastore file.`,
	Run: func(cmd *cobra.Command, args []string) {

		data, err := datastore.Load()
		if err != nil {
			fmt.Println("Error loading data:", err)
			return
		}

		app := tview.NewApplication()

		table := components.DisplayRecordsTable(data, func(key tcell.Key) {
			if key == tcell.KeyEscape {
				app.Stop()
			}
		}, func(row int, column int) {
			// Handle row selection but is not in use right now
		})

		if err := app.SetRoot(table, true).Run(); err != nil {
			panic(err)
		}

	},
}

func init() {
	recordCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
