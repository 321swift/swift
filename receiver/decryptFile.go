package receiver

import "encoding/hex"

func decryptFile(encFile []byte) ([]byte, error) {
	var decFile = make([]byte, len(encFile))
	_, err := hex.Decode(decFile, encFile)
	if err != nil {
		return nil, err
	}

	return decFile, nil
}
