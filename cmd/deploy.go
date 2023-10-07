/*
Copyright © 2023 github.com/alexanderlesser
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/alexanderlesser/sv-cli/datastore"
	"github.com/alexanderlesser/sv-cli/internal/components"
	"github.com/alexanderlesser/sv-cli/types"
	"github.com/alexanderlesser/sv-cli/utils/helpers"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/spf13/cobra"
)

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		deployJs, _ := cmd.Flags().GetBool("js")

		var startGulp bool

		prompt := &survey.Confirm{
			Message: "Do you want to start Gulp?",
			Default: false,
		}

		err := survey.AskOne(prompt, &startGulp)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		records, err := datastore.Load()
		if err != nil {
			fmt.Println("Cannot load records: ", err)
			os.Exit(1)
		}

		app := tview.NewApplication()

		pages := tview.NewPages()
		grid := tview.NewGrid()

		table := components.DisplayRecordsTable(records, func(k tcell.Key) {
			if k == tcell.KeyEscape {
				app.Stop()
			}
		}, func(row, column int) {
			record := records[row-1]

			var files []types.File
			if deployJs {

				f, err := helpers.GetJSFiles(record.Path)
				files = f
				if err != nil {
					app.Stop()
					fmt.Println("Cannot fetch files: ", err)
					os.Exit(1)
				}

			} else {
				f, err := helpers.GetCSSFiles(record.Path)
				files = f

				if err != nil {
					app.Stop()
					fmt.Println("Cannot fetch files: ", err)
					os.Exit(1)
				}
			}

			components.GenerateGrid(app, grid, record, files, deployJs, startGulp)
			pages.AddPage("grid", grid, true, false)
			pages.SwitchToPage("grid")

		})

		pages.AddPage("grid", grid, true, true)
		pages.AddPage("table", table, true, true)

		if err := app.SetRoot(pages, true).Run(); err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(deployCmd)
	deployCmd.Flags().BoolP("js", "j", false, "Deploy javascript files")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deployCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deployCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
