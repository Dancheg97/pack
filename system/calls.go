// Copyright 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Additional information can be found on official web page: https://fmnx.io/
// Contact email: help@fmnx.io

package system

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"

	"fmnx.io/core/pack/config"
)

// Execute external command with fmt like formatting.
func Callf(format string, a ...any) (string, error) {
	return Call(fmt.Sprintf(format, a...))
}

// Execute external command call in bash.
func Call(cmd string) (string, error) {
	if config.DebugMode {
		fmt.Println("Syscall - executing: ", cmd)
	}
	commad := exec.Command("bash", "-c", cmd)
	var buf bytes.Buffer
	if config.DebugMode {
		commad.Stdout = io.MultiWriter(&buf, os.Stdout)
		commad.Stderr = io.MultiWriter(&buf, os.Stderr)
	} else {
		commad.Stdout = &buf
		commad.Stderr = &buf
	}
	err := commad.Run()
	return buf.String(), err
}
