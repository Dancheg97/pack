// Copyright 2023 FMNX Linux team.
// This code is covered by GPL license, which can be found in LICENSE file.
// Additional information can be found on official web page: https://fmnx.io/
// Contact email: help@fmnx.io
package system

import (
	"os"
	"strings"
)

// Simply overwrite file also creating missing directories.
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

// Prepare directories and all it's subdirs.
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

// Add text contents to file.
func AppendToFile(path string, contents string) error {
	_, err := os.Stat(path)
	if err != nil {
		err = os.WriteFile(path, []byte(""), 0644)
		if err != nil {
			return err
		}
	}
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	if _, err := f.Write([]byte(contents)); err != nil {
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return nil
}
