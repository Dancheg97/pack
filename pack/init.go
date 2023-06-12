// 2023 FMNX team.
// Use of this code is governed by GNU General Public License.
// Official web page: https://fmnx.su/core/pack
// Contact email: help@fmnx.su

package pack

import (
	"bytes"
	"errors"
	"os/exec"
)

func formOptions[Opts any](arr []Opts, getdefault func() *Opts) *Opts {
	if len(arr) != 1 {
		return getdefault()
	}
	return &arr[0]
}

func call(cmd *exec.Cmd) error {
	var buf bytes.Buffer
	cmd.Stderr = &buf
	err := cmd.Run()
	if err != nil {
		return errors.New(buf.String())
	}
	return nil
}
