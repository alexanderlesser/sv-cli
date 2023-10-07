package config

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"

	"github.com/alexanderlesser/sv-cli/internal/constants"
	"github.com/spf13/viper"
)

var EncryptionKey string

func getConfigFilePath() string {
	var configFilePath string

	if constants.DEVELOPMENT {
		configFilePath = constants.CONFIG_FILE_NAME
	} else {
		configDirPath := getConfigDirPath()
		configFilePath = filepath.Join(configDirPath, constants.CONFIG_FILE_NAME)
	}

	return configFilePath
}

func getConfigDirPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return filepath.Join(homeDir, constants.CONFIG_DIR_NAME)
}

func createConfigFile(filePath string) error {
	configContent := fmt.Sprintf("# %s\n%s: %s", constants.CONFIG_FILE_NAME, EncryptionKey, "")
	return os.WriteFile(filePath, []byte(configContent), 0600)
}

func generateRandomBytes(length int) ([]byte, error) {
	randomBytes := make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}
	return randomBytes, nil
}

func generateRandomKey(length int) string {
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()-_=+?/:;<>"
	var result strings.Builder
	for i := 0; i < length; i++ {
		result.WriteByte(chars[rand.Intn(len(chars))])
	}
	return result.String()
}

// Creates or load config file
func InitializeConfig() {
	done := make(chan bool)

	go func() {
		configFilePath := getConfigFilePath()
		if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
			err := createConfigFile(configFilePath)
			if err != nil {
				fmt.Println("Error creating config file:", err)
				done <- true
				return
			}

			viper.SetConfigFile(configFilePath)

			// Generate and store the encryption key
			encryptionKey := generateRandomKey(24)
			viper.Set(constants.ENCRYPTION_KEY_NAME, encryptionKey)

			viper.Set(constants.CONFIG_CSS_NAME, "assets/css")
			viper.Set(constants.CONFIG_JS_NAME, "assets/js")
			viper.Set(constants.CONFIG_REST_API_NAME, "rest-api/upload-css/upload")
			viper.Set(constants.CONFIG_MINIFIED_CSS_NAME, false)
			viper.Set(constants.CONFIG_MINIFIED_JS_NAME, false)

			// Generate and store the encryption bytes
			encryptionBytes, err := generateRandomBytes(16)
			encryptedBytesStr := strings.Join(strings.Fields(fmt.Sprint(encryptionBytes)), ", ")

			if err != nil {
				fmt.Println("Cannot generate random bytes")
				os.Exit(1)
			}

			viper.Set(constants.ENCRYPTION_BYTES_NAME, fmt.Sprintf("%v", encryptedBytesStr))

			// Save the configuration to the file
			err = viper.WriteConfig()
			if err != nil {
				fmt.Println("Error writing config file:", err)
				done <- true
				return
			}
		} else {
			viper.SetConfigFile(configFilePath)
			if err := viper.ReadInConfig(); err != nil {
				fmt.Println("Error reading config file:", err)
				done <- true
				return
			}
		}

		done <- true
	}()

	<-done
}
