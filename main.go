/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"fmt"
	"os"

	"github.com/alexanderlesser/sv-cli/cmd"
	"github.com/alexanderlesser/sv-cli/datastore"
	"github.com/alexanderlesser/sv-cli/internal/config"
)

func main() {
	config.InitializeConfig()

	err := datastore.InitDataStorage()
	if err != nil {
		fmt.Println("Error initializing data storage:", err)
		os.Exit(1)
	}

	cmd.Execute()
}
