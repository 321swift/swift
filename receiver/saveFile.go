package receiver

import "os"

func saveFile(file []byte) error {
	// get file name
	// bytes.
	err := os.WriteFile("received", file, 0775)
	if err != nil {
		return err
	}
	return nil
}
