package utils

import "os"

func ReadFileContent(filePath string) (string, error) {
	var err error
	var content []byte

	content, err = os.ReadFile(filePath)
	return string(content), err
}
