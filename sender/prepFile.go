package sender

import (
	"encoding/hex"
	"log"
	"os"
)

// prepFile takes in the file path and reads the file to be processed.
// The function returns the encrypted file
func prepFile(file_path string) []byte {
	// read file
	file, err := os.ReadFile(file_path)
	if err != nil {
		log.Fatal("unable to read file: \n", err)
	}

	// encrypt file
	encoded := make([]byte, hex.EncodedLen(len(file)))
	hex.Encode(encoded, file)

	return encoded
}
