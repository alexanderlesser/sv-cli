package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/alexanderlesser/sv-cli/datastore"
	"github.com/alexanderlesser/sv-cli/types"
	"github.com/alexanderlesser/sv-cli/utils/encrypt"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

func readPassword(prompt string) (string, error) {
	fmt.Print(prompt)
	password, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println() // Print a newline after reading password
	if err != nil {
		return "", err
	}
	return string(password), nil
}

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Creates a new record",
	Long:  `The add command creates and stores a new record in the datastore file. Records are crutial for deploying files later.`,
	Run: func(cmd *cobra.Command, args []string) {
		record := types.Record{}
		prompt := &survey.Input{
			Message: "Enter username:",
		}
		survey.AskOne(prompt, &record.Username, survey.WithValidator(survey.Required))

		// Prompt for password
		password, err := readPassword("Enter password: ")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		encryptedPassword, err := encrypt.EncryptPassword(password)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		record.Password = encryptedPassword

		// Prompt for domain
		prompt = &survey.Input{
			Message: "Enter domain:",
		}
		survey.AskOne(prompt, &record.Domain, survey.WithValidator(survey.Required))

		// Prompt for name
		prompt = &survey.Input{
			Message: "Enter name:",
		}
		survey.AskOne(prompt, &record.Name, survey.WithValidator(survey.Required))

		currentPath, err := os.Getwd()
		if err != nil {
			log.Println(err)
			currentPath = ""
		}

		// Prompt for path with default value as the current directory
		pathPrompt := &survey.Input{
			Message: "Enter root path for project:",
			Default: currentPath,
		}
		survey.AskOne(pathPrompt, &record.Path, survey.WithValidator(survey.Required))

		err = datastore.Save(record)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	recordCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
