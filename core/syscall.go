package core

import (
	"bytes"
	"fmt"
	"os/exec"
)

func SystemCallf(format string, a ...any) (string, error) {
	return SystemCall(fmt.Sprintf(format, a...))
}

func SystemCall(cmd string) (string, error) {
	commad := exec.Command("bash", "-c", cmd)
	var buf bytes.Buffer
	commad.Stdout = &buf
	commad.Stderr = &buf
	err := commad.Run()
	return buf.String(), err
}
