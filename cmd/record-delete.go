/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/alexanderlesser/sv-cli/datastore"
	"github.com/alexanderlesser/sv-cli/internal/components"
	"github.com/alexanderlesser/sv-cli/types"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/spf13/cobra"
)

var recordTable *tview.Table

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete record",
	Long:  `The delete command deletes a record in the datastore file.`,
	Run: func(cmd *cobra.Command, args []string) {
		id, _ := cmd.Flags().GetString("id")

		if id != "" {
			intVar, err := strconv.Atoi(id)

			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			err = datastore.DeleteRecord(int32(intVar))

			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			os.Exit(0)
		}

		data, err := datastore.Load()

		if err != nil {
			fmt.Println("Error loading data:", err)
			return
		}

		var app = tview.NewApplication()
		pages := tview.NewPages()

		modal := tview.NewModal().
			AddButtons([]string{"YES", "NO"})

		recordTable = components.DisplayRecordsTable(data, func(key tcell.Key) {
			if key == tcell.KeyEscape {
				app.Stop()
			}
		}, func(row int, column int) {
			selectedRecord := data[row-1]
			modal.SetText(fmt.Sprintf(`Are you sure you want to delete record "%s"?`, selectedRecord.Name))
			modal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
				if buttonLabel == "YES" {
					err := deleteRecord(selectedRecord)
					if err != nil {
						app.Stop()
						fmt.Println(err)
						os.Exit(1)
					}

					recordTable.RemoveRow(row)
					pages.SwitchToPage("table")
				}
				if buttonLabel == "NO" {
					pages.SwitchToPage("table")
				}
			})

			pages.SwitchToPage("modal")
		})

		pages.AddPage("modal", modal, false, false)
		pages.AddPage("table", recordTable, true, false)
		pages.SwitchToPage("table")

		if err := app.SetRoot(pages, true).Run(); err != nil {
			panic(err)
		}
	},
}

func init() {
	recordCmd.AddCommand(deleteCmd)
	deleteCmd.PersistentFlags().String("id", "", "ID of record to delete")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deleteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deleteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func deleteRecord(selectedRecord types.Record) error {
	err := datastore.DeleteRecord(selectedRecord.ID)
	if err != nil {
		// app.Stop()
		// fmt.Println(err)
		// os.Exit(1)
		return err
	}

	return nil
}
