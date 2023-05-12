// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Additional information can be found on official web page: https://fmnx.io/
// Contact email: help@fmnx.io

package system

// This package contains functions for execution system calls and other
// utilities for different system IO operations.

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

// Prepare directories and all it's subdirs.
func MkDir(filePath string) error {
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

// Get directory for current process.
func Pwd() string {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("Unable to find curr dir")
		os.Exit(1)
	}
	return dir
}

// Move files with provided extension from one directory to another.
func MvExt(src string, dst string, ext string) error {
	const command = "sudo mv %s/*%s %s"
	o, err := Callf(command, src, ext, dst)
	if err != nil {
		return errors.New("unable to move files:\n" + o)
	}
	return nil
}

// Simply overwrite file also creating missing directories.
func WriteFile(file string, content string) error {
	err := MkDir(file)
	if err != nil {
		return err
	}
	err = os.WriteFile(file, []byte(content), 0o600)
	if err != nil {
		return err
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
