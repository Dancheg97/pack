// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Additional information can be found on official web page: https://fmnx.su/
// Contact email: help@fmnx.su

package system

// This package contains functions for execution system calls and other
// utilities for different system IO operations.

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"

	"fmnx.su/core/pack/config"
	"fmnx.su/core/pack/prnt"
)

// Execute external command with fmt like formatting.
func Callf(format string, a ...any) (string, error) {
	return Call(fmt.Sprintf(format, a...))
}

// Execute external command call in bash.
func Call(cmd string) (string, error) {
	execute := exec.Command("bash", "-c", cmd)
	var buf bytes.Buffer
	if config.DebugMode {
		prnt.Yellow("=> Executing system call: ", cmd)
		execute.Stdout = io.MultiWriter(&buf, os.Stdout)
		execute.Stderr = io.MultiWriter(&buf, os.Stderr)
	} else {
		execute.Stdout = &buf
		execute.Stderr = &buf
	}
	err := execute.Run()
	if err != nil {
		return buf.String(), errors.New(buf.String())
	}
	return buf.String(), nil
}
