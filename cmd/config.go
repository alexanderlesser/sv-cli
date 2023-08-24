/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/alexanderlesser/sv-cli/internal/constants"
	"github.com/rivo/tview"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Show and edit current config variables",
	Long:  `Command to display and edit current config variables`,
	Run: func(cmd *cobra.Command, args []string) {
		cssVal := viper.GetString(constants.CONFIG_CSS_NAME)
		jsVal := viper.GetString(constants.CONFIG_JS_NAME)
		minCss := viper.GetBool(constants.CONFIG_MINIFIED_CSS_NAME)
		minJs := viper.GetBool(constants.CONFIG_MINIFIED_JS_NAME)

		app := tview.NewApplication()
		form := tview.NewForm().
			AddInputField("CSS path", cssVal, 40, nil, nil).
			AddInputField("Javascript path", jsVal, 40, nil, nil).
			AddCheckbox("Only show minified css files", minCss, nil).
			AddCheckbox("Only show minified javascript files", minJs, nil)

		form.AddButton("Save", func() {
			cssVal = form.GetFormItem(0).(*tview.InputField).GetText()
			jsVal = form.GetFormItem(1).(*tview.InputField).GetText()
			minCss = form.GetFormItem(2).(*tview.Checkbox).IsChecked()
			minJs = form.GetFormItem(3).(*tview.Checkbox).IsChecked()

			if cssVal != "" && jsVal != "" {
				viper.Set(constants.CONFIG_CSS_NAME, cssVal)
				viper.Set(constants.CONFIG_JS_NAME, jsVal)
				viper.Set(constants.CONFIG_MINIFIED_CSS_NAME, minCss)
				viper.Set(constants.CONFIG_MINIFIED_JS_NAME, minJs)

				err := viper.WriteConfig()
				if err != nil {
					fmt.Println("Error writing config file:", err)
					app.Stop()
					os.Exit(1)
				}

				app.Stop()
				fmt.Println("Config updated")
			}
		}).
			AddButton("Quit", func() {
				app.Stop()
			})

		form.SetBorder(true).SetTitle("Config settings").SetTitleAlign(tview.AlignLeft)
		if err := app.SetRoot(form, true).EnableMouse(true).Run(); err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(configCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// configCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
