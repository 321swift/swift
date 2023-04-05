package receiver

import (
	"os"
	"path"
)

func saveFile(filename string, location string, file []byte) error {
	err := os.WriteFile(path.Join(location, filename), file, 0775)
	if err != nil {
		return err
	}
	return nil
}
