package encrypt

import (
	"os"
	"testing"

	"github.com/alexanderlesser/sv-cli/internal/config"
	"github.com/alexanderlesser/sv-cli/internal/constants"
	"github.com/alexanderlesser/sv-cli/testUtil"
	"github.com/stretchr/testify/assert"
)

func TestEncryptDecryptPassword(t *testing.T) {
	// Clean up any existing config file before testing
	defer func() {
		restore := testUtil.SetupTestEnvironment()
		defer restore() // Restore the original environment after the tests are done
	}()

	// Initialize the config for testing
	// config.ToggleDevelopmentForTesting()
	config.InitializeConfig()

	// Sample password to test
	password := "mysecretpassword"

	// Encrypt the password
	encryptedPassword, err := EncryptPassword(password)
	assert.NoError(t, err)

	// Decrypt the password
	decryptedPassword, err := DecryptPassword(encryptedPassword)
	assert.NoError(t, err)

	// Check if the decrypted password matches the original
	assert.Equal(t, password, decryptedPassword)

	// Clean up the config file created during testing
	os.Remove(constants.CONFIG_FILE_NAME)
}
