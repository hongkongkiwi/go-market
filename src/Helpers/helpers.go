package Helpers

import (
	"os"
)

func AppendToFile(filename, data string) error {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.WriteString(data); err != nil {
		return err
	}
	return nil
}

func CreateDirIfNotExist(dirPath string) error {
	// Check if the directory already exists
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		// Directory does not exist, so create it
		err := os.MkdirAll(dirPath, 0755) // 0755 sets the directory permissions
		if err != nil {
			return err
		}
	}
	return nil
}

func CreateTempFile(namePrefix string) (*os.File, error) {
	// Create a temporary file
	tmpDBFile, err := os.CreateTemp("", namePrefix+"-*.db")
	if err != nil {
		return nil, err
	}
	return tmpDBFile, nil
}
