// 2023 FMNX team.
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
	"fmnx.io/core/pack/prnt"
	"fmnx.io/core/pack/tmpl"
	"github.com/fatih/color"
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
		prnt.Yellow("Executing system call: ", cmd)
		execute.Stdout = io.MultiWriter(&buf, os.Stdout)
		execute.Stderr = io.MultiWriter(&buf, os.Stderr)
	} else {
		execute.Stdout = &buf
		execute.Stderr = &buf
	}
	err := execute.Run()
	if err != nil {
		if config.DisablePrettyPrint {
			return fmt.Sprintf(tmpl.SysCallErr, cmd, err, buf.String()), err
		}
		return fmt.Sprintf(
			tmpl.SysCallErr,
			color.RedString(cmd),
			err,
			buf.String(),
		), err
	}
	return buf.String(), nil
}
