package encrypt

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/alexanderlesser/sv-cli/internal/constants"
	"github.com/spf13/viper"
)

// Read the encryption key from config
func GetEncryptionKey() string {
	// Read encryptionKey from Viper
	encryptionKey := viper.GetString(constants.ENCRYPTION_KEY_NAME)
	return encryptionKey
}

// Read the byteslice from config
func GetEncryptionBytes() []byte {
	var encryptionBytes []byte
	encryptionBytesStr := viper.GetString(constants.ENCRYPTION_BYTES_NAME)
	if encryptionBytesStr != "" {
		bytes, err := decodeBytes(encryptionBytesStr)
		encryptionBytes = bytes
		if err != nil {
			fmt.Println("Error decoding encryption_bytes:", err)
		}
	}

	return encryptionBytes
}

// Decodes the byte slice
func decodeBytes(input string) ([]byte, error) {
	input = strings.ReplaceAll(input, "[", "")
	input = strings.ReplaceAll(input, "]", "")
	input = strings.ReplaceAll(input, " ", "")

	// Split the string by commas
	byteStrings := strings.Split(input, ",")

	// Convert strings to bytes
	var result []byte
	for _, str := range byteStrings {
		num, err := strconv.Atoi(str)
		if err != nil {
			return nil, err
		}
		result = append(result, byte(num))
	}

	return result, nil
}
