package config

import (
	"os"
	"testing"

	"github.com/alexanderlesser/sv-cli/internal/constants"
	"github.com/alexanderlesser/sv-cli/testUtil"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestInitializeConfig(t *testing.T) {
	// Clean up any existing config file before testing
	defer os.Remove(constants.CONFIG_FILE_NAME)

	// Setup
	defer func() {
		restore := testUtil.SetupTestEnvironment()
		defer restore()
	}()

	// Testing
	InitializeConfig()

	// Verify that the config file has been created
	assert.FileExists(t, constants.CONFIG_FILE_NAME)

	// Verify that the EncryptionKey has been set
	encryptionKey := viper.GetString(constants.ENCRYPTION_KEY_NAME)
	assert.NotEmpty(t, encryptionKey)

	// bytes exist
	encryptionBytesStr := viper.GetString(constants.ENCRYPTION_BYTES_NAME)
	assert.NotEmpty(t, encryptionBytesStr)

	defer os.Remove(constants.CONFIG_FILE_NAME)
}

// Run tests
func TestMain(m *testing.M) {

	exitCode := m.Run()

	os.Exit(exitCode)
}
