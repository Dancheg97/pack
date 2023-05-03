package system

import (
	"bytes"
	"fmt"
	"os/exec"
)

var Debug = false

func SystemCallf(format string, a ...any) (string, error) {
	return Call(fmt.Sprintf(format, a...))
}

func Call(cmd string) (string, error) {
	if Debug {
		fmt.Println("Syscall - executing: ", cmd)
	}
	commad := exec.Command("bash", "-c", cmd)
	var buf bytes.Buffer
	commad.Stdout = &buf
	commad.Stderr = &buf
	if Debug {
		fmt.Println("Syscall - output: ", buf.String())
	}
	err := commad.Run()
	return buf.String(), err
}
