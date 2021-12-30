package seeder

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"github.com/marmorag/supateam/internal"
	"io"
	"os"
)

func WriteSecureData(filepath string) (string, error) {
	config := internal.GetConfig()

	text, _ := os.ReadFile(filepath)
	key := []byte(config.ApplicationAESPassphrase)

	encFilePath := fmt.Sprintf("%s.enc", filepath)

	// generate a new aes cipher using our 32 byte long key
	c, err := aes.NewCipher(key)
	// if there are any errors, handle them
	if err != nil {
		fmt.Println(err)
	}

	// gcm or Galois/Counter Mode, is a mode of operation
	// for symmetric key cryptographic block ciphers
	// - https://en.wikipedia.org/wiki/Galois/Counter_Mode
	gcm, err := cipher.NewGCM(c)
	// if any error generating new GCM
	// handle them
	if err != nil {
		fmt.Println(err)
	}

	// creates a new byte array the size of the nonce
	// which must be passed to Seal
	nonce := make([]byte, gcm.NonceSize())
	// populates our nonce with a cryptographically secure
	// random sequence
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		fmt.Println(err)
	}

	// here we encrypt our text using the Seal function
	// Seal encrypts and authenticates plaintext, authenticates the
	// additional data and appends the result to dst, returning the updated
	// slice. The nonce must be NonceSize() bytes long and unique for all
	// time, for a given key.
	enc := gcm.Seal(nonce, nonce, text, nil)
	err = os.WriteFile(encFilePath, enc, 0644)

	return encFilePath, err
}

func ReadSecureData(filePath string) (string, error) {
	config := internal.GetConfig()

	key := []byte(config.ApplicationAESPassphrase)
	enc, _ := os.ReadFile(filePath)

	c, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println(err)
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		fmt.Println(err)
	}

	nonceSize := gcm.NonceSize()
	if len(enc) < nonceSize {
		fmt.Println(err)
	}

	nonce, ciphertext := enc[:nonceSize], enc[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)

	return string(plaintext), err
}
