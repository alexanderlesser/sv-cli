package testUtil

import (
	"github.com/alexanderlesser/sv-cli/internal/constants"
)

// "github.com/your-username/your-project/internal/constants"
// "github.com/your-username/your-project/internal/config"

func SetupTestEnvironment() func() {
	// Store the original values
	originalConfigName := constants.CONFIG_FILE_NAME
	originalDevelopment := constants.DEVELOPMENT
	originalDatastoreName := constants.DATASTORE_NAME

	// Modify the necessary constants for testing
	constants.CONFIG_FILE_NAME = "config_test.yaml"
	constants.DEVELOPMENT = true
	constants.DATASTORE_NAME = "data_test.json"

	// config.InitializeConfig()

	// Return a function to restore the original values
	return func() {
		constants.CONFIG_FILE_NAME = originalConfigName
		constants.DEVELOPMENT = originalDevelopment
		constants.DATASTORE_NAME = originalDatastoreName
	}
}
