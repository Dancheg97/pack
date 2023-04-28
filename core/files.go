package core

import (
	"os"
	"strings"
)

func WriteFile(file string, content string) error {
	err := PrepareDir(file)
	if err != nil {
		return err
	}
	err = os.WriteFile(file, []byte(content), 0o600)
	if err != nil {
		return err
	}
	return nil
}

func PrepareDir(filePath string) error {
	if len(strings.Split(filePath, `/`)) != 1 {
		splitted := strings.Split(filePath, `/`)
		path := strings.Join(splitted[0:len(splitted)-1], `/`)
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}
