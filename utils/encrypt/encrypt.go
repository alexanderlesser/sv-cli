package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

// Encodes string to base64
func encode(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

// Encrypt method is to encrypt or hide any classified text
func encrypt(text, MySecret string) (string, error) {
	bytes := GetEncryptionBytes()
	block, err := aes.NewCipher([]byte(MySecret))
	if err != nil {
		return "", err
	}
	plainText := []byte(text)
	cfb := cipher.NewCFBEncrypter(block, bytes)
	cipherText := make([]byte, len(plainText))
	cfb.XORKeyStream(cipherText, plainText)
	return encode(cipherText), nil
}

// Decodes base64 string
func decode(s string) []byte {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return data
}

// Decrypt method is to extract back the encrypted text
func decrypt(text, MySecret string) (string, error) {
	bytes := GetEncryptionBytes()

	block, err := aes.NewCipher([]byte(MySecret))
	if err != nil {
		return "", err
	}
	cipherText := decode(text)
	cfb := cipher.NewCFBDecrypter(block, bytes)
	plainText := make([]byte, len(cipherText))
	cfb.XORKeyStream(plainText, cipherText)
	return string(plainText), nil
}

// Encrypts password
func EncryptPassword(s string) (string, error) {
	secret := GetEncryptionKey()

	encText, err := encrypt(s, secret)
	if err != nil {
		return "", err
	}

	return encText, nil
}

// Decrypts password
func DecryptPassword(s string) (string, error) {
	secret := GetEncryptionKey()
	decText, err := decrypt(s, secret)
	if err != nil {
		return "", err
	}

	return decText, nil
}
