package aosStreamTarget

import (
	"io"
	"os"
	"path/filepath"
)

const (
	keyLogFile = ".aosStream.keys"
)

func keyLogWriter() (io.Writer, error) {
	keyLogDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	keyLogFile := filepath.Join(keyLogDir, keyLogFile)

	err = os.MkdirAll(filepath.Dir(keyLogFile), os.FileMode(0644))
	if err != nil {
		return nil, err
	}

	return os.OpenFile(keyLogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
}
